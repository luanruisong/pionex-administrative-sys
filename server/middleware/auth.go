package middleware

import (
	"net/http"
	"pionex-administrative-sys/db"
	"pionex-administrative-sys/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	ContextKeyClaims = "claims"
)

func r(c *gin.Context, code int, data string) {
	c.JSON(code, gin.H{
		"code": code,
		"msg":  data,
	})
	c.Abort()
}

// Auth JWT 认证中间件
func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Header 获取 token
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			r(c, http.StatusUnauthorized, "请求头中缺少 Authorization")
			return
		}

		// 解析 Bearer token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			r(c, http.StatusUnauthorized, "Authorization 格式错误")
			return
		}

		// 解析 token
		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			r(c, http.StatusUnauthorized, "无效的 token: "+err.Error())
			return
		}

		// 将用户信息存入上下文
		c.Set(ContextKeyClaims, claims)

		c.Next()
	}
}

// RequireRole 检查用户是否拥有指定权限
func RequireRole(role db.CommonRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		u := GetCurrentClaims(c)
		if !u.HasRole(role.Role) {
			r(c, http.StatusForbidden, "权限不足")
			return
		}
		c.Next()
	}
}

// GetCurrentClaims 从上下文获取 Claims
func GetCurrentClaims(c *gin.Context) *utils.Claims {
	if claims, exists := c.Get(ContextKeyClaims); exists {
		return claims.(*utils.Claims)
	}
	return nil
}
