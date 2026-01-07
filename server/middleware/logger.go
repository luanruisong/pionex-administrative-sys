package middleware

import (
	"pionex-administrative-sys/utils/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logger 使用项目 logger 的 Gin 日志中间件
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()

		fields := []zap.Field{
			zap.Int("status", status),
			zap.String("method", c.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", latency),
			zap.Int("size", c.Writer.Size()),
		}

		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("error", c.Errors.String()))
			logger.Error("request", fields...)
		} else if status >= 500 {
			logger.Error("request", fields...)
		} else if status >= 400 {
			logger.Warn("request", fields...)
		} else {
			logger.Info("request", fields...)
		}
	}
}

// Recovery 使用项目 logger 的 Gin 恢复中间件
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("panic recovered",
					zap.Any("error", err),
					zap.String("path", c.Request.URL.Path),
					zap.String("method", c.Request.Method),
				)
				c.AbortWithStatus(500)
			}
		}()
		c.Next()
	}
}
