package service

import (
	"errors"
	"main/IMChat/model"
	"main/IMChat/pkg/e"
	"main/IMChat/serializer"
	"main/ToDoList/pkg/utils"

	"gorm.io/gorm"
)

type UserRegisterService struct {
	UserName string `form:"user_name" json:"user_name" binding:"required,min=3,max=15" `
	Password string `form:"password" json:"password" binding:"required,min=5,max=16"`
}

type UserLoginService struct {
	UserName string `form:"user_name" json:"user_name" binding:"required,min=3,max=15" `
	Password string `form:"password" json:"password" binding:"required,min=5,max=16"`
}

// type FriendService struct {
// }

func (service *UserRegisterService) Register() serializer.Response {
	code := e.SUCCESS
	var user model.User
	var count int64

	model.DB.Model(&model.User{}).Where("user_name=?", service.UserName).First(&user).Count(&count)
	if count != 0 {
		code = e.ErrorExistUserName
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	user.UserName = service.UserName

	//密码加密
	if err := user.SetPassword(service.Password); err != nil {
		code = e.ErrorFailEncryption
		return serializer.Response{
			Status: code,
			Msg:    err.Error(),
			Error:  err.Error(),
		}
	}

	// 创建用户
	if err := model.DB.Create(&user).Error; err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}

}

func (service *UserLoginService) Login() serializer.Response {
	code := e.SUCCESS
	var user model.User

	//查找数据库中是否存在该用户
	if err := model.DB.Where("user_name=?", service.UserName).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			code = e.ErrorNotExistUser
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
				Error:  err.Error(),
			}
		}

		//如果是其他错误
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	if !user.CheckPassword(service.Password) {
		code = e.ErrorNotCompare
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	//发送token，为了其他功能需要身份验证所给前端存储
	token, err := utils.GenerateToken(user.ID, service.UserName, service.Password)
	if err != nil {
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	if err := model.DB.Where("user_name=?", service.UserName).First(&user).Update("status", 1).Error; err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}

	return serializer.Response{
		Status: code,
		Data: serializer.TokenData{
			User:  serializer.BuildUser(user),
			Token: token,
		},
		Msg: e.GetMsg(code),
	}
}
