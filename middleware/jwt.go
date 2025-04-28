package middleware

import (
	"time"
	"todo_list/package/utils"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := 200
		token := c.GetHeader("Authorization")
		// 判断是否是sse请求
		if token == "" {
			token = c.Query("token")
		}
		if token == "" {
			code = 400
		} else {
			claims, err := utils.ParseToken(token)
			if err != nil {
				code = 403 // 无权限
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = 401 // 过期
			}
			c.Set("claims", claims)
		}
		if code != 200 {
			c.JSON(200, gin.H{
				"code": code,
				"msg":  "token error!",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// SSE认证
func SSEAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从JWT中提取用户信息
		claims, exists := c.Get("claims")
		if !exists {
			c.AbortWithStatus(401)
			return
		}

		// 获取用户ID
		customClaims, ok := claims.(*utils.Claims)
		if !ok {
			c.AbortWithStatus(401)
			return
		}
		userID := customClaims.Id
		c.Set("currentUserID", userID)
		c.Next()
	}
}
