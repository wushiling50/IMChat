package router

import (
	"main/IMChat/api"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func UserRoutersInit(r *gin.Engine) {
	store := cookie.NewStore([]byte("something-very-secret"))
	r.Use(sessions.Sessions("mysession", store))

	v1 := r.Group("api/v1/user")
	{
		v1.GET("ping", func(c *gin.Context) {
			c.JSON(200, "SUCCESS") //ping
		})
		v1.POST("register", api.UserRegister) // 注册
		v1.POST("login", api.UserLogin)       //登录
		// v1.GET("friend", api.Friend)         //好友
	}
}
