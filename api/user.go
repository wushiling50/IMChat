package api

import (
	"main/IMChat/pkg/e"
	"main/IMChat/service"

	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
)

func UserRegister(c *gin.Context) {
	code := e.SUCCESS
	var userRegister service.UserRegisterService
	if err := c.ShouldBind(&userRegister); err == nil {
		res := userRegister.Register()
		c.JSON(code, res)
	} else {
		logging.Error(err)
		code = e.InvalidParams
		c.JSON(code, e.ErrorResponse(err))
	}
}

func UserLogin(c *gin.Context) {
	code := e.SUCCESS
	var userLogin service.UserLoginService
	if err := c.ShouldBind(&userLogin); err == nil {
		res := userLogin.Login()
		c.JSON(code, res)
	} else {
		logging.Error(err)
		code = e.InvalidParams
		c.JSON(code, e.ErrorResponse(err))
	}
}

// func Friend(c *gin.Context) {
// 	code := e.SUCCESS
// 	var friend service.FriendService
// 	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
// 	if err := c.ShouldBind(&friend); err != nil {
// 		logging.Error(err)
// 		code = e.InvalidParams
// 		c.JSON(code, e.ErrorResponse(err))
// 	} else {
// 		res := friend.Friend(claim.Id)
// 		c.JSON(code, res)
// 	}

// }
