package api

import (
	"main/IMChat/pkg/e"
	"main/IMChat/pkg/util"
	"main/IMChat/service"

	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
)

func Create(c *gin.Context) {
	code := e.SUCCESS
	var create service.CreateService
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&create); err != nil {
		logging.Error(err)
		code = e.InvalidParams
		c.JSON(code, e.ErrorResponse(err))
	} else {
		res := create.Create(claim.Id)
		c.JSON(code, res)
	}

}
