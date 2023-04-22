package router

import (
	"main/IMChat/api"
	"main/IMChat/middleware"
	"main/IMChat/service"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func ChatRoutersInit(r *gin.Engine) {
	store := cookie.NewStore([]byte("something-very-secret"))
	r.Use(sessions.Sessions("mysession", store))

	v1 := r.Group("api/v1/chat")
	v1.Use(middleware.JWT())
	{
		v1.GET("ws", service.Handle)      //单聊
		v1.POST("create", api.Create)     //建立群聊
		v1.GET("ws/group", service.Group) //群聊(只具备实时聊天功能)
	}
}
