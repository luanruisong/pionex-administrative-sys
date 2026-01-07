package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("pas-secret-key")

// Claims JWT claims
type Claims struct {
	UserId int64 `json:"user_id"`
	Role   int   `json:"role"` // 权限位: 1=admin, 2=login
	jwt.RegisteredClaims
}

// HasRole 检查是否拥有指定权限
func (c Claims) HasRole(role int) bool {
	return c.Role&role != 0
}

// GenerateToken 生成 JWT token
func GenerateToken(userId int64, role int, expireDuration time.Duration) (string, error) {
	claims := Claims{
		UserId: userId,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "pas",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// ParseToken 解析 JWT token
func ParseToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// SetJWTSecret 设置 JWT 密钥
func SetJWTSecret(secret string) {
	jwtSecret = []byte(secret)
}
