package handler

import (
	"pionex-administrative-sys/server/handler/coupon"
	my_coupon "pionex-administrative-sys/server/handler/my_coupon"
	"pionex-administrative-sys/server/handler/user"
	"pionex-administrative-sys/server/middleware"
	"pionex-administrative-sys/utils"

	"github.com/gin-gonic/gin"
)

func Register(r gin.IRouter) {
	r.GET("/health", healthHandler)
	r.Use(middleware.Logger(), middleware.Recovery())
	api := r.Group("/api/v1")
	{
		user.Register(api)
		coupon.Register(api)
		my_coupon.Register(api)
	}
}

// healthHandler 健康检查
func healthHandler(c *gin.Context) {
	utils.Resp(200, "OK", gin.H{
		"status": "ok",
	}).Success(c)
}
