package middleware

import (
	"MYshop/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未携带token",
			})
			return
		}
		parts := strings.SplitN(authHeader, "", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "token格式错误",
			})
			return
		}
		claims, err := util.ParseToken(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "token无效或过期",
			})
			return
		}
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username) //只在这一次http请求处理期间内有效
		c.Next()
	}
}
