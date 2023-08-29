package middleware

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JwtPrivateKey 定义 JWT 私钥
var JwtPrivateKey = []byte(os.Getenv("JWT_SECRET"))

// var JwtPrivateKey = []byte("secret")

// UserClaims 定义自定义的 Claims 结构体
type UserClaims struct {
	jwt.RegisteredClaims
	ID       uint   `json:"id"`
	UserType string `json:"userType"`
	UserName string `json:"username"`
	Role     []Role `json:"role"`
}

type Role string

// 角色常量
const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "User"
)

func (r Role) String() string {
	return string(r)
}

func RoleContains(roles []Role, role Role) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

// GenerateToken 生成 JWT
func GenerateToken(u UserClaims) (string, error) {
	claims := UserClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "",
			Subject:   "user token",
		},
		ID:       u.ID,
		UserType: u.UserType,
		UserName: u.UserName,
		Role:     u.Role,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的密钥签名并获得完整的编码后的字符串令牌
	return token.SignedString(JwtPrivateKey)
}

// LoginAuthMiddleware 登陆认证中间件
func LoginAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 HTTP Header 中获取 token
		tokenString := c.GetHeader("Authorization")
		// 如果 token 为空，返回错误信息
		if len(tokenString) <= 15 {
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Unauthorized",
			})
			return
		}
		// 去除 Bearer 前缀
		prefixLen := len("Bearer ")
		if len(tokenString) > prefixLen && tokenString[:prefixLen] == "Bearer " {
			tokenString = tokenString[prefixLen:]
			//log.Info("tokenString", tokenString)
		}
		// 校验 token 是否有效
		token, err := jwt.ParseWithClaims(tokenString, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return JwtPrivateKey, nil
		})

		if err == nil && token.Valid {
			// 用户已认证，继续执行请求
			// 设定UID
			claims := token.Claims.(*UserClaims)
			c.Set("user", claims)
			c.Next()
		} else {
			log.Printf("认证失败 %s", err)
			// 认证失败，返回错误信息
			c.AbortWithStatusJSON(401, gin.H{
				"error": "Unauthorized",
			})
		}
	}
}
