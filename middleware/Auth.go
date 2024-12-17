package middleware

import (
	"OuterChat/model"
	"OuterChat/util"
	"OuterChat/util/SError"
	"github.com/gin-gonic/gin"
)

func Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if !model.CheckTokenValid(token) {
			util.SendErrorResponse(c, SError.InvalidTokenError, "无效Token")
			c.Abort()
			return
		}
		c.Next()
	}
}
