package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// 用户列表
type User struct {
	gorm.Model
	UserName       string `gorm:"type:varchar(20);not null;unique"`
	PasswordDigest string //存储的是密文
	Email          string //`gorm:"unique"`
	Avatar         string `gorm:"size:1000"`
	Phone          string
	Status         uint
}

const PassWordCost = 12 //密码加密难度

// 加密
func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PassWordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

// 验证密码
func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil
}
