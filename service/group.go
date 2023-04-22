package service

import (
	"encoding/json"
	"errors"
	"log"
	"main/IMChat/model"
	"main/IMChat/pkg/e"
	"main/IMChat/pkg/util"
	"main/IMChat/serializer"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type CreateService struct {
	UserName string `form:"user_name" json:"user_name" binding:"required,min=3,max=15" `
	Password string `form:"password" json:"password" binding:"required,min=5,max=16"`
}

// 建立群聊
func (service *CreateService) Create(uid uint) serializer.Response {
	var user model.User
	code := e.SUCCESS
	//校验身份
	if err := model.DB.Where("id=?", uid).Where("user_name=?", service.UserName).
		First(&user).Error; err != nil {
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

	//生成群号
	rand.Seed(time.Now().Unix())
	gid := rand.Uint64()

	group := &model.Group{
		BuilderID: uid,
		GroupID:   gid,
	}

	if err := model.DB.Model(&group).Create(group).Error; err != nil {
		code = e.ErrorDatabase
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Error:  err.Error(),
		}
	}
	return serializer.Response{
		Status: code,
		Data:   gid,
		Msg:    e.GetMsg(code),
	}
}

func Group(c *gin.Context) {
	uid := c.Query("uid") // 自己的id
	gid := c.Query("gid") // 群聊的id

	//检验用户
	claim, err := util.ParseToken(c.GetHeader("Authorization"))
	if err != nil {
		log.Println("token解析错误")
		return
	}
	i1, err := strconv.ParseUint(uid, 10, 32)
	if uint(i1) != claim.Id || err != nil {
		log.Println("参数传输错误")
		return
	}
	var group model.Group
	i2, _ := strconv.ParseUint(gid, 10, 64)
	if err := model.DB.Model(&model.Group{}).Where("group_id=?", i2).First(&group).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("该群不存在", err)
			return
		}
	}

	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { // CheckOrigin解决跨域问题
			return true
		}}).Upgrade(c.Writer, c.Request, nil) // 升级成ws协议

	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	// 创建一个用户实例
	client := &Client{
		ID:     createId(uid, gid),
		ToID:   strconv.FormatUint(i2, 10),
		Socket: conn,
		Send:   make(chan []byte),
	}

	// 用户注册到用户管理上
	Manager.Register <- client
	go client.Read()
	go client.GWrite()
}

func (c *Client) GWrite() {
	defer func() {
		_ = c.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				_ = c.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			from, rmsg, _ := model.MsgSplit(string(message))
			log.Println(c.ID, "接受消息:", rmsg)
			replyMsg := ReplyMsg{
				From:    from,
				Code:    e.WebsocketSuccessMessage,
				Content: rmsg,
			}
			msg, _ := json.Marshal(replyMsg)
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
		}
	}

}
