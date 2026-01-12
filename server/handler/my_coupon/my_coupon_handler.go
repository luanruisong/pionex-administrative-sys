package my_coupon

import (
	"errors"
	"pionex-administrative-sys/db"
	"pionex-administrative-sys/server/middleware"
	"pionex-administrative-sys/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Register 注册路由
func Register(r gin.IRouter) {
	g := r.Group("/my-coupon")

	// 需要登录权限
	g.Use(middleware.Auth())

	g.GET("/list", listHandler)
	g.GET("/detail/:id", detailHandler)
	g.GET("/stock", stockHandler)

	// 申领卡券需要 RoleApplyCoupon 权限
	g.POST("/take", middleware.RequireRole(db.RoleApplyCoupon), takeHandler)
}

// MyCouponItem 我的卡券列表项
type MyCouponItem struct {
	Id        int64  `json:"id"`
	Type      int    `json:"type"`
	TypeName  string `json:"type_name"`
	TakenAt   int64  `json:"taken_at"` // 领取时间（使用 updated_at）
	CreatedAt int64  `json:"created_at"`
}

// MyCouponDetail 我的卡券详情
type MyCouponDetail struct {
	Id        int64  `json:"id"`
	Coupon    string `json:"coupon"` // 卡券码
	Type      int    `json:"type"`
	TypeName  string `json:"type_name"`
	TakenAt   int64  `json:"taken_at"`
	CreatedAt int64  `json:"created_at"`
}

func toMyCouponItem(c *db.Coupon) MyCouponItem {
	return MyCouponItem{
		Id:        c.Id,
		Type:      c.Type,
		TypeName:  db.GetCouponTypeName(c.Type),
		TakenAt:   c.UpdatedAt, // 领取时间用 updated_at
		CreatedAt: c.CreatedAt,
	}
}

func toMyCouponDetail(c *db.Coupon) MyCouponDetail {
	return MyCouponDetail{
		Id:        c.Id,
		Coupon:    c.Coupon,
		Type:      c.Type,
		TypeName:  db.GetCouponTypeName(c.Type),
		TakenAt:   c.UpdatedAt,
		CreatedAt: c.CreatedAt,
	}
}

// listHandler 我的卡券列表
func listHandler(c *gin.Context) {
	userId := middleware.GetCurrentClaims(c).UserId

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	// 筛选类型
	var typeFilter *int
	if typeStr := c.Query("type"); typeStr != "" {
		if t, err := strconv.Atoi(typeStr); err == nil {
			typeFilter = &t
		}
	}

	offset := (page - 1) * size
	coupons, err := db.GetCouponsByTaker(c.Request.Context(), userId, typeFilter, offset, size)
	if err != nil {
		utils.Resp(500, "查询失败", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	total, _ := db.CountCouponsByTaker(c.Request.Context(), userId, typeFilter)

	list := make([]MyCouponItem, 0, len(coupons))
	for _, cp := range coupons {
		list = append(list, toMyCouponItem(cp))
	}

	utils.Resp(0, "success", gin.H{
		"list":  list,
		"total": total,
		"page":  page,
		"size":  size,
	}).Success(c)
}

// detailHandler 我的卡券详情
func detailHandler(c *gin.Context) {
	userId := middleware.GetCurrentClaims(c).UserId

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

	// 验证是否是自己的卡券
	if coupon.Taker != userId {
		utils.Resp(403, "无权查看此卡券", gin.H{}).Fail(c)
		return
	}

	utils.Resp(0, "success", toMyCouponDetail(coupon)).Success(c)
}

// stockHandler 查询指定类型卡券库存
func stockHandler(c *gin.Context) {
	typeStr := c.Query("type")
	if typeStr == "" {
		utils.Resp(400, "参数错误", gin.H{"error": "缺少type参数"}).Fail(c)
		return
	}

	couponType, err := strconv.Atoi(typeStr)
	if err != nil {
		utils.Resp(400, "参数错误", gin.H{"error": "无效的卡券类型"}).Fail(c)
		return
	}

	// 校验卡券类型
	if !db.IsValidCouponType(couponType) {
		utils.Resp(400, "无效的卡券类型", gin.H{}).Fail(c)
		return
	}

	count, err := db.CountAvailableCouponsByType(c.Request.Context(), couponType)
	if err != nil {
		utils.Resp(500, "查询失败", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	utils.Resp(0, "success", gin.H{
		"type":     couponType,
		"typeName": db.GetCouponTypeName(couponType),
		"stock":    count,
	}).Success(c)
}

// TakeReq 申领卡券请求
type TakeReq struct {
	Type int `json:"type" binding:"required"`
}

// takeHandler 申领卡券
func takeHandler(c *gin.Context) {
	var req TakeReq
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

	// 校验12小时内是否已领取过该类型卡券
	lastCoupon, err := db.GetLastTakenCouponByTakerAndType(c.Request.Context(), userId, req.Type)
	if err == nil {
		// 找到了上次领取记录，检查时间间隔
		lastTakenTime := time.UnixMilli(lastCoupon.UpdatedAt)
		elapsed := time.Since(lastTakenTime)
		if elapsed < 12*time.Hour {
			remaining := 12*time.Hour - elapsed
			hours := int(remaining.Hours())
			minutes := int(remaining.Minutes()) % 60
			utils.Resp(400, "领取过于频繁，请稍后再试", gin.H{
				"message":           "同一类型卡券12小时内只能领取一张",
				"remaining_hours":   hours,
				"remaining_minutes": minutes,
			}).Fail(c)
			return
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// 查询出错（非"记录不存在"错误）
		utils.Resp(500, "查询失败", gin.H{"error": err.Error()}).Fail(c)
		return
	}

	// 获取一个可用的卡券
	coupon, err := db.GetOneAvailableCouponByType(c.Request.Context(), req.Type)
	if err != nil {
		utils.Resp(400, "该类型卡券库存不足", gin.H{}).Fail(c)
		return
	}

	// 领取卡券
	if err := db.TakeCoupon(c.Request.Context(), coupon.Id, userId); err != nil {
		if err == db.ErrCouponAlreadyTaken {
			utils.Resp(400, "卡券已被他人领取，请重试", gin.H{}).Fail(c)
		} else {
			utils.Resp(500, "领取失败", gin.H{"error": err.Error()}).Fail(c)
		}
		return
	}

	// 重新查询已领取的卡券信息
	takenCoupon, _ := db.GetCouponById(c.Request.Context(), coupon.Id)

	utils.Resp(0, "success", gin.H{
		"id":        takenCoupon.Id,
		"type":      takenCoupon.Type,
		"type_name": db.GetCouponTypeName(takenCoupon.Type),
		"coupon":    takenCoupon.Coupon,
	}).Success(c)
}
