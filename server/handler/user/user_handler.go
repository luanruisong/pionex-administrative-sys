package user

import (
	"pionex-administrative-sys/db"
	"pionex-administrative-sys/server/middleware"
	"pionex-administrative-sys/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Register 注册路由
func Register(r gin.IRouter) {
	g := r.Group("/user")

	// 公开接口
	g.POST("/login", loginHandler)
	g.POST("/register", registerHandler)

	// 需要登录权限的接口
	g.Use(middleware.Auth())
	g.GET("/profile", profileHandler)
	g.PUT("/profile", updateProfileHandler)

	// 需要管理员权限的接口
	g.Use(middleware.RequireRole(db.RoleAdmin))
	g.GET("/roles", allRoleHandler)
	g.POST("/add", addUserHandler)
	g.GET("/list", listHandler)
	g.PUT("/update", updateHandler)
	g.DELETE("/delete/:id", deleteHandler)
}

// RegisterReq 注册请求
type RegisterReq struct {
	Name     string `json:"name" binding:"required"`
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// registerHandler 用户注册
func registerHandler(c *gin.Context) {
	var req RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Resp(400, "参数错误", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	// 检查账号是否已存在
	if _, err := db.GetUserByAccount(c.Request.Context(), req.Account); err == nil {
		utils.Resp(400, "账号已存在", gin.H{}).Fail(c)
		return
	}

	// 创建用户
	user := &db.User{
		Name:    req.Name,
		Account: req.Account,
		Md5Pwd:  utils.MD5(req.Password),
	}
	if err := db.CreateUser(c.Request.Context(), user); err != nil {
		utils.Resp(500, "注册失败", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	utils.Resp(0, "success", gin.H{"account": req.Account}).Success(c)
}

// AddUserReq 管理员添加用户请求
type AddUserReq struct {
	Name     string `json:"name" binding:"required"`
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     *int   `json:"role"` // 权限位，不传则默认为 RoleLogin
}

// addUserHandler 管理员添加用户
func addUserHandler(c *gin.Context) {
	var req AddUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Resp(400, "参数错误", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	// 检查账号是否已存在
	if _, err := db.GetUserByAccount(c.Request.Context(), req.Account); err == nil {
		utils.Resp(400, "账号已存在", gin.H{}).Fail(c)
		return
	}

	// 设置权限，默认为登录权限
	role := db.RoleLogin.Role
	if req.Role != nil {
		role = *req.Role
	}

	// 创建用户
	user := &db.User{
		Name:    req.Name,
		Account: req.Account,
		Md5Pwd:  utils.MD5(req.Password),
		Role:    role,
	}
	if err := db.CreateUser(c.Request.Context(), user); err != nil {
		utils.Resp(500, "创建用户失败", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	utils.Resp(0, "success", gin.H{
		"id":      user.Id,
		"account": req.Account,
		"role":    role,
	}).Success(c)
}

// LoginReq 登录请求
type LoginReq struct {
	Account  string `json:"account" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResp 登录响应
type LoginResp struct {
	Token     string `json:"token"`
	ExpiresIn int64  `json:"expires_in"`
	Role      int    `json:"role"` // 权限位: 1=admin, 2=login
	Name      string `json:"name"`
}

// loginHandler 用户登录
func loginHandler(c *gin.Context) {
	var req LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Resp(400, "参数错误", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	// 查询用户
	user, err := db.GetUserByAccount(c.Request.Context(), req.Account)
	if err != nil {
		utils.Resp(401, "账号不存在", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	// 校验密码
	if err := user.CheckPwd(req.Password); err != nil {
		utils.Resp(401, "密码错误", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	// 校验登录权限
	if !user.HasRole(db.RoleLogin) {
		utils.Resp(403, "账号无登录权限", gin.H{}).Fail(c)
		return
	}

	// 签发 JWT token (24小时有效)
	expireDuration := 24 * time.Hour
	token, err := utils.GenerateToken(user.Id, user.Role, expireDuration)
	if err != nil {
		utils.Resp(500, "token生成失败", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	utils.Resp(0, "success", LoginResp{
		Token:     token,
		ExpiresIn: int64(expireDuration.Seconds()),
		Role:      user.Role,
		Name:      user.Name,
	}).Success(c)
}

// UserItem 用户列表项
type UserItem struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	Account   string `json:"account"`
	Role      int    `json:"role"` // 权限位: 1=admin, 2=login
	CreatedAt int64  `json:"created_at"`
}

// listHandler 用户列表
func listHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	offset := (page - 1) * size
	users, err := db.GetUserList(c.Request.Context(), offset, size)
	if err != nil {
		utils.Resp(500, "查询失败", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	total, _ := db.CountUsers(c.Request.Context())

	list := make([]UserItem, 0, len(users))
	for _, u := range users {
		list = append(list, UserItem{
			Id:        u.Id,
			Name:      u.Name,
			Account:   u.Account,
			Role:      u.Role,
			CreatedAt: u.CreatedAt,
		})
	}

	utils.Resp(0, "success", gin.H{
		"list":  list,
		"total": total,
		"page":  page,
		"size":  size,
	}).Success(c)
}

// UpdateReq 更新请求
type UpdateReq struct {
	Id       int64   `json:"id" binding:"required"`
	Name     *string `json:"name"`
	Account  *string `json:"account"`
	Password *string `json:"password"`
	Role     *int    `json:"role"` // 权限位: 1=admin, 2=login
}

// updateHandler 更新用户
func updateHandler(c *gin.Context) {
	var req UpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Resp(400, "参数错误", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	fields := make(map[string]interface{})
	if req.Name != nil && *req.Name != "" {
		fields["name"] = *req.Name
	}
	if req.Account != nil && *req.Account != "" {
		// 检查账号是否已被其他用户使用
		existUser, err := db.GetUserByAccount(c.Request.Context(), *req.Account)
		if err == nil && existUser.Id != req.Id {
			utils.Resp(400, "账号已被其他用户使用", gin.H{}).Fail(c)
			return
		}
		fields["account"] = *req.Account
	}
	if req.Password != nil && *req.Password != "" {
		fields["md5_pwd"] = utils.MD5(*req.Password)
	}
	if req.Role != nil {
		fields["role"] = *req.Role
	}

	if len(fields) == 0 {
		utils.Resp(400, "没有要更新的字段", gin.H{}).Fail(c)
		return
	}

	if err := db.UpdateUserFields(c.Request.Context(), req.Id, fields); err != nil {
		utils.Resp(500, "更新失败", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	utils.Resp(0, "success", gin.H{}).Success(c)
}

// deleteHandler 删除用户
func deleteHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Resp(400, "参数错误", gin.H{"error": "无效的用户ID"}).Fail(c)
		return
	}

	if err := db.DeleteUser(c.Request.Context(), id); err != nil {
		utils.Resp(500, "删除失败", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	utils.Resp(0, "success", gin.H{}).Success(c)
}

// ProfileResp 用户资料响应
type ProfileResp struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Account    string `json:"account"`
	Role       int    `json:"role"` // 权限位: 1=admin, 2=login
	PrivateKey string `json:"private_key"`
	CreatedAt  int64  `json:"created_at"`
}

// profileHandler 获取当前用户资料
func profileHandler(c *gin.Context) {
	userId := middleware.GetCurrentClaims(c).UserId

	user, err := db.GetUserById(c.Request.Context(), userId)
	if err != nil {
		utils.Resp(404, "用户不存在", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	utils.Resp(0, "success", ProfileResp{
		Id:         user.Id,
		Name:       user.Name,
		Account:    user.Account,
		Role:       user.Role,
		PrivateKey: user.PrivateKey,
		CreatedAt:  user.CreatedAt,
	}).Success(c)
}

// profileHandler 获取当前用户资料
func allRoleHandler(c *gin.Context) {
	utils.Resp(0, "success", gin.H{
		"list": db.AllRoles(),
	}).Success(c)
}

// UpdateProfileReq 更新资料请求
type UpdateProfileReq struct {
	Name       *string `json:"name"`
	Password   *string `json:"password"`
	PrivateKey *string `json:"private_key"`
}

// updateProfileHandler 更新当前用户资料
func updateProfileHandler(c *gin.Context) {
	userId := middleware.GetCurrentClaims(c).UserId

	var req UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Resp(400, "参数错误", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	fields := make(map[string]interface{})
	if req.Name != nil && *req.Name != "" {
		fields["name"] = *req.Name
	}
	if req.Password != nil && *req.Password != "" {
		fields["md5_pwd"] = utils.MD5(*req.Password)
	}
	if req.PrivateKey != nil {
		fields["private_key"] = *req.PrivateKey
	}

	if len(fields) == 0 {
		utils.Resp(400, "没有要更新的字段", gin.H{}).Fail(c)
		return
	}

	if err := db.UpdateUserFields(c.Request.Context(), userId, fields); err != nil {
		utils.Resp(500, "更新失败", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	utils.Resp(0, "success", gin.H{}).Success(c)
}
