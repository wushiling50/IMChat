package middleware

import (
	"main/IMChat/pkg/e"
	"main/IMChat/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int = e.SUCCESS

		token := c.GetHeader("Authorization")
		if token == "" {
			code = e.ErrorAuthToken
		} else {
			claim, err := util.ParseToken(token)
			if err != nil {
				code = e.ErrorAuthCheckTokenFail // 无权限，token是无权限的，是假的
			} else if time.Now().Unix() > claim.ExpiresAt {
				code = e.ErrorAuthCheckTokenTimeout //Token无效
			}
		}

		if code != e.SUCCESS {
			c.JSON(200, gin.H{
				"Status": code,
				"msg":    e.GetMsg(code),
			})
			c.Abort()
			return
		}
		c.Next()
	}

}
