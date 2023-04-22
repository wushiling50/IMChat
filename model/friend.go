package model

import (
	"fmt"
	"strconv"

	"gorm.io/gorm"
)

type Friend struct {
	gorm.Model
	Uid      uint
	FriendId uint
}

func AddFriends(id string) error {
	var friend Friend
	toid := string(id[len(id)-1])
	uid := string(id[0])
	id1, _ := strconv.ParseUint(toid, 10, 64)
	id2, _ := strconv.ParseUint(uid, 10, 64)
	friend1 := &Friend{
		Uid:      uint(id1),
		FriendId: uint(id2),
	}
	friend2 := &Friend{
		Uid:      uint(id2),
		FriendId: uint(id1),
	}

	err := DB.Model(&friend).Create(&friend1).Create(&friend2).Error

	return err
}

func IsFriends(id string) bool {
	var friend Friend
	var count int64
	toid := string(id[len(id)-1])
	uid := string(id[0])
	id1, _ := strconv.ParseUint(toid, 10, 64)
	id2, _ := strconv.ParseUint(uid, 10, 64)

	DB.Model(&friend).Where("uid=?", uint(id1)).Where("friend_id=?", uint(id2)).Count(&count).Find(&friend)
	DB.Model(&friend).Where("uid=?", uint(id2)).Where("friend_id=?", uint(id1)).Count(&count).Find(&friend)
	fmt.Println(count)

	return count != 0

}
