package model

import (
	"errors"
	"main/IMChat/pkg/util"

	"gorm.io/gorm"
)

type News struct {
	gorm.Model
	SendID    string //标记消息的用户发送方与接收方
	Msg       string
	ReadOrNot uint  //0为未读，1为已读，
	Time      int64 //返回发送时的时间
}

// 创建新的消息
func CreateNews(sendID, msg string, readOrNot uint, time int64) error {
	news := &News{
		SendID:    sendID,
		Msg:       msg,
		ReadOrNot: readOrNot,
		Time:      time,
	}

	err := DB.Model(&news).Create(news).Error

	return err
}

func FriendNews(msg []byte, toid string, id string) error {
	var news News
	var count int64
	umsg := util.FriendHash(msg)

	if string(msg) == "好友请求" {
		DB.Model(&news).Where("send_id=?", id).Where("msg=?", umsg).Find(&news).Count(&count)
		if count != 0 {
			return errors.New("等待对方应答")
		}
		DB.Model(&news).Where("send_id=?", toid).Where("msg=?", umsg).Find(&news).Count(&count)
		if count != 0 {
			return errors.New("尚未通过对方的好友请求")
		}
		return nil
	} else if string(msg) == "Y" || string(msg) == "N" {
		var num int64
		DB.Model(&news).Where("send_id=?", toid).Where("msg=?", util.FriendHash([]byte("好友请求"))).Count(&num).Find(&news)
		if num == 0 {
			return errors.New("对方尚未发送好友请求")
		}

		DB.Model(&news).Where("send_id=?", id).
			Where("msg in ?", []string{util.FriendHash([]byte("N")), util.FriendHash([]byte("Y"))}).
			Count(&count).Find(&news)
		if count != 0 {
			return errors.New("请勿重复发送验证消息")
		}
		return nil
	}

	return errors.New("错误的内容请求信息")

}
