package model

import (
	"log"
	"regexp"

	"gorm.io/gorm"
)

type Group struct {
	gorm.Model
	BuilderID uint   //群主id
	GroupID   uint64 //群号
}

func MsgSplit(msg string) (from string, rmsg string, err error) {
	s1 := "[0-9]+(->)[0-9]*"
	s2, err := regexp.Compile(s1)
	if err != nil {
		log.Println("正则表达式解析错误")
	}
	from = s2.FindString(msg) //from
	rmsg = msg[len(from)+1:]  //信息

	return from, rmsg, nil
}
