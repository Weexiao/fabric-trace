package middleware

import (
	"backend/pkg"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 基于JWT的认证中间件
func JWTAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 支持 Authorization: Bearer <token> 或直接传token
		raw := c.GetHeader("Authorization")
		if raw == "" {
			c.JSON(200, gin.H{
				"code": 401,
				"msg":  "请求未携带token，无权限访问1",
			})
			c.Abort()
			return
		}

		token := strings.TrimSpace(raw)
		const bearerPrefix = "Bearer "
		if strings.HasPrefix(token, bearerPrefix) {
			token = strings.TrimSpace(strings.TrimPrefix(token, bearerPrefix))
		}

		mc, err := pkg.ParseToken(token)
		if err != nil {
			c.JSON(200, gin.H{
				"code": 401,
				"msg":  "请求未携带token，无权限访问3",
				"data": err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("userID", mc.UserID)
		c.Next()
	}
}
