package main

import (
	"main/IMChat/conf"
	"main/IMChat/model"
	"main/IMChat/router"
	"main/IMChat/service"

	"github.com/gin-gonic/gin"
)

func main() {

	conf.SetUp()
	model.SetUp()

	go service.Manager.Start() //启动监听

	r := gin.Default()

	router.UserRoutersInit(r) //用户路由

	router.ChatRoutersInit(r) //ws路由

	r.Run(conf.HttpPort)

}
