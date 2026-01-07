package static

import "github.com/gin-gonic/gin"

func Register(engine *gin.Engine) {
	// 静态文件
	engine.Static("/static", "./static")

	// 根路径重定向到登录页
	engine.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/static/html/login.html")
	})
}
