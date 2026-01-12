package static

import (
	"embed"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

//go:embed html/** css/** js/** image/**
var StaticFS embed.FS

// fallback 路径，找不到文件时跳转到这里
const fallbackPath = "/static/html/main.html"

func toHome(c *gin.Context) {
	c.Redirect(302, fallbackPath)
}

func Register(engine *gin.Engine) {
	// 自定义静态文件处理，支持 fallback
	engine.GET("/static/*filepath", func(c *gin.Context) {
		filepath := c.Param("filepath")
		// 去掉开头的斜杠
		filepath = strings.TrimPrefix(filepath, "/")

		// 尝试打开文件
		f, err := StaticFS.Open(filepath)
		if err != nil {
			// 文件不存在，跳转到 fallback
			c.Redirect(302, fallbackPath)
			return
		}
		defer f.Close()

		// 检查是否是目录
		stat, err := f.Stat()
		if err != nil || stat.IsDir() {
			c.Redirect(302, fallbackPath)
			return
		}

		// 使用 http.ServeContent 提供文件
		http.ServeContent(c.Writer, c.Request, stat.Name(), stat.ModTime(), f.(io.ReadSeeker))
	})

	// 根路径重定向
	engine.GET("/", toHome)
}
