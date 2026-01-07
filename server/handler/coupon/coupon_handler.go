package coupon

import (
	"pionex-administrative-sys/db"
	"pionex-administrative-sys/server/middleware"
	"pionex-administrative-sys/utils"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// Register 注册路由
func Register(r gin.IRouter) {
	g := r.Group("/coupon")

	// 需要登录权限
	g.Use(middleware.Auth())

	// 获取卡券类型列表（不需要特殊权限）
	g.GET("/types", typesHandler)

	// 需要库存管理权限
	g.Use(middleware.RequireRole(db.RoleStock))

	g.POST("/add", addHandler)
	g.POST("/import", importHandler)
	g.GET("/list", listHandler)
	g.GET("/detail/:id", detailHandler)
	g.PUT("/update", updateHandler)
	g.DELETE("/delete/:id", deleteHandler)
}

// CouponItem 卡券列表项
type CouponItem struct {
	Id        int64  `json:"id"`
	Coupon    string `json:"coupon"`
	Type      int    `json:"type"`
	TypeName  string `json:"type_name"`
	Creator   int64  `json:"creator"`
	Taker     int64  `json:"taker"`
	IsTaken   bool   `json:"is_taken"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func toCouponItem(c *db.Coupon) CouponItem {
	return CouponItem{
		Id:        c.Id,
		Coupon:    c.Coupon,
		Type:      c.Type,
		TypeName:  db.GetCouponTypeName(c.Type),
		Creator:   c.Creator,
		Taker:     c.Taker,
		IsTaken:   c.IsTaken(),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

// typesHandler 获取卡券类型列表
func typesHandler(c *gin.Context) {
	utils.Resp(0, "success", gin.H{
		"list": db.AllCouponTypes(),
	}).Success(c)
}

// AddReq 添加卡券请求
type AddReq struct {
	Coupon string `json:"coupon" binding:"required"`
	Type   int    `json:"type" binding:"required"`
}

// addHandler 添加卡券
func addHandler(c *gin.Context) {
	var req AddReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Resp(400, "参数错误", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	// 校验卡券类型
	if !db.IsValidCouponType(req.Type) {
		utils.Resp(400, "无效的卡券类型", gin.H{}).Fail(c)
		return
	}

	userId := middleware.GetCurrentClaims(c).UserId

	// 检查卡券是否已存在
	if _, err := db.GetCouponByCode(c.Request.Context(), req.Coupon); err == nil {
		utils.Resp(400, "卡券已存在", gin.H{}).Fail(c)
		return
	}

	coupon := &db.Coupon{
		Coupon:  req.Coupon,
		Type:    req.Type,
		Creator: userId,
	}
	if err := db.CreateCoupon(c.Request.Context(), coupon); err != nil {
		utils.Resp(500, "创建卡券失败", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	utils.Resp(0, "success", gin.H{
		"id":     coupon.Id,
		"coupon": coupon.Coupon,
		"type":   coupon.Type,
	}).Success(c)
}

// ImportReq 导入卡券请求
type ImportReq struct {
	Coupons string `json:"coupons" binding:"required"` // 多个卡券用换行符分隔
	Type    int    `json:"type" binding:"required"`    // 卡券类型
}

// ImportResp 导入卡券响应
type ImportResp struct {
	Total      int      `json:"total"`      // 总数
	Success    int      `json:"success"`    // 成功数
	Failed     int      `json:"failed"`     // 失败数
	Duplicates []string `json:"duplicates"` // 重复的卡券
}

// importHandler 批量导入卡券
func importHandler(c *gin.Context) {
	var req ImportReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Resp(400, "参数错误", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	// 校验卡券类型
	if !db.IsValidCouponType(req.Type) {
		utils.Resp(400, "无效的卡券类型", gin.H{}).Fail(c)
		return
	}

	userId := middleware.GetCurrentClaims(c).UserId

	// 解析卡券列表（支持换行符、逗号分隔）
	lines := strings.Split(req.Coupons, "\n")
	var codes []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// 支持逗号分隔
		parts := strings.Split(line, ",")
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				codes = append(codes, part)
			}
		}
	}

	if len(codes) == 0 {
		utils.Resp(400, "没有有效的卡券", gin.H{}).Fail(c)
		return
	}

	// 去重
	codeMap := make(map[string]bool)
	var uniqueCodes []string
	for _, code := range codes {
		if !codeMap[code] {
			codeMap[code] = true
			uniqueCodes = append(uniqueCodes, code)
		}
	}

	var successCount int
	var duplicates []string

	// 逐个创建（检查重复）
	for _, code := range uniqueCodes {
		// 检查是否已存在
		if _, err := db.GetCouponByCode(c.Request.Context(), code); err == nil {
			duplicates = append(duplicates, code)
			continue
		}

		coupon := &db.Coupon{
			Coupon:  code,
			Type:    req.Type,
			Creator: userId,
		}
		if err := db.CreateCoupon(c.Request.Context(), coupon); err == nil {
			successCount++
		}
	}

	utils.Resp(0, "success", ImportResp{
		Total:      len(uniqueCodes),
		Success:    successCount,
		Failed:     len(duplicates),
		Duplicates: duplicates,
	}).Success(c)
}

// listHandler 卡券列表
func listHandler(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	// 解析筛选条件
	filter := db.CouponFilter{}
	if typeStr := c.Query("type"); typeStr != "" {
		if t, err := strconv.Atoi(typeStr); err == nil {
			filter.Type = &t
		}
	}
	if takenStr := c.Query("taken"); takenStr != "" {
		taken := takenStr == "1"
		filter.Taken = &taken
	}

	offset := (page - 1) * size
	coupons, err := db.GetCouponListWithFilter(c.Request.Context(), filter, offset, size)
	if err != nil {
		utils.Resp(500, "查询失败", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	total, _ := db.CountCouponsWithFilter(c.Request.Context(), filter)

	list := make([]CouponItem, 0, len(coupons))
	for _, cp := range coupons {
		list = append(list, toCouponItem(cp))
	}

	utils.Resp(0, "success", gin.H{
		"list":  list,
		"total": total,
		"page":  page,
		"size":  size,
	}).Success(c)
}

// detailHandler 卡券详情
func detailHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Resp(400, "参数错误", gin.H{"error": "无效的卡券ID"}).Fail(c)
		return
	}

	coupon, err := db.GetCouponById(c.Request.Context(), id)
	if err != nil {
		utils.Resp(404, "卡券不存在", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	utils.Resp(0, "success", toCouponItem(coupon)).Success(c)
}

// UpdateReq 更新卡券请求
type UpdateReq struct {
	Id     int64   `json:"id" binding:"required"`
	Coupon *string `json:"coupon"`
	Type   *int    `json:"type"`
}

// updateHandler 更新卡券
func updateHandler(c *gin.Context) {
	var req UpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Resp(400, "参数错误", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	// 检查卡券是否存在
	existing, err := db.GetCouponById(c.Request.Context(), req.Id)
	if err != nil {
		utils.Resp(404, "卡券不存在", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	// 已被领取的卡券不能修改
	if existing.IsTaken() {
		utils.Resp(400, "已被领取的卡券不能修改", gin.H{}).Fail(c)
		return
	}

	fields := make(map[string]interface{})
	if req.Coupon != nil && *req.Coupon != "" {
		// 检查新卡券码是否已被使用
		if existCoupon, err := db.GetCouponByCode(c.Request.Context(), *req.Coupon); err == nil && existCoupon.Id != req.Id {
			utils.Resp(400, "卡券码已被使用", gin.H{}).Fail(c)
			return
		}
		fields["coupon"] = *req.Coupon
	}
	if req.Type != nil {
		// 校验卡券类型
		if !db.IsValidCouponType(*req.Type) {
			utils.Resp(400, "无效的卡券类型", gin.H{}).Fail(c)
			return
		}
		fields["type"] = *req.Type
	}

	if len(fields) == 0 {
		utils.Resp(400, "没有要更新的字段", gin.H{}).Fail(c)
		return
	}

	if err := db.UpdateCouponFields(c.Request.Context(), req.Id, fields); err != nil {
		utils.Resp(500, "更新失败", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	utils.Resp(0, "success", gin.H{}).Success(c)
}

// deleteHandler 删除卡券
func deleteHandler(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		utils.Resp(400, "参数错误", gin.H{"error": "无效的卡券ID"}).Fail(c)
		return
	}

	// 检查卡券是否存在
	existing, err := db.GetCouponById(c.Request.Context(), id)
	if err != nil {
		utils.Resp(404, "卡券不存在", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	// 已被领取的卡券不能删除
	if existing.IsTaken() {
		utils.Resp(400, "已被领取的卡券不能删除", gin.H{}).Fail(c)
		return
	}

	if err := db.DeleteCoupon(c.Request.Context(), id); err != nil {
		utils.Resp(500, "删除失败", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	utils.Resp(0, "success", gin.H{}).Success(c)
}