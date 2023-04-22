package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"main/IMChat/cache"
	"main/IMChat/model"
	"main/IMChat/pkg/e"
	"main/IMChat/pkg/util"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

const passTime = time.Hour * 24 * 30 * 3 // 三个月

// 发送的消息
type SendMsg struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
}

// 回复的消息
type ReplyMsg struct {
	From    string      `json:"from"`
	Code    int         `json:"code"`
	Content interface{} `json:"content"`
}

// 用户类
type Client struct {
	ID     string
	ToID   string
	Socket *websocket.Conn //代表一个ws连接
	Send   chan []byte
}

// 广播类，包括广播内容和源用户，用于将消息发送给服务端
type Broadcast struct {
	Client  *Client
	Message []byte
	Type    int
}

// 群聊的广播
type GBroadcast struct {
	GBroadcast *Broadcast
}

// 好友验证的广播
type FBroadcast struct {
	FBroadcast *Broadcast
}

// 用户管理
type ClientManager struct {
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	GBroadcast chan *GBroadcast
	FBroadcast chan *FBroadcast
	Reply      chan *Client
	Register   chan *Client
	Unregister chan *Client
}

var Manager = ClientManager{
	Clients:    make(map[string]*Client), // 参与连接的用户
	Broadcast:  make(chan *Broadcast),
	GBroadcast: make(chan *GBroadcast),
	FBroadcast: make(chan *FBroadcast),
	Register:   make(chan *Client),
	Reply:      make(chan *Client),
	Unregister: make(chan *Client),
}

func createId(uid, toUid string) string {
	return uid + "->" + toUid
}

func Handle(c *gin.Context) {
	uid := c.Query("uid")     // 自己的id
	toUid := c.Query("toUid") // 对方的id

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
	var user model.User
	i2, _ := strconv.ParseUint(toUid, 10, 32)
	if err := model.DB.Model(&model.User{}).Where("id=?", i2).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("目标用户不存在", err)
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
		ID:     createId(uid, toUid),
		ToID:   createId(toUid, uid),
		Socket: conn,
		Send:   make(chan []byte),
	}

	// 用户注册到用户管理上
	Manager.Register <- client
	go client.Read()
	go client.Write()
}

func (c *Client) Read() {
	r := cache.RedisClient
	ctx := context.TODO()
	defer func() { // 避免忘记关闭，所以要加上close
		Manager.Unregister <- c
		_ = c.Socket.Close()
	}()
	for {
		c.Socket.PongHandler()
		sendMsg := new(SendMsg)
		err := c.Socket.ReadJSON(&sendMsg) // 读取json格式，如果不是json格式，会报错
		if err != nil {
			log.Println("数据传输错误", err)
			break
		}
		if sendMsg.Type == 1 { //单聊
			r.Incr(ctx, c.ID)                      //发消息的次数
			r.Expire(ctx, c.ID, passTime).Result() // 防止过快“分手”，建立连接三个月过期
			log.Println(c.ID, "发送消息", sendMsg.Content)
			Manager.Broadcast <- &Broadcast{
				Client:  c,
				Message: []byte(sendMsg.Content),
			}
		} else if sendMsg.Type == 2 { //查询
			timeT, err := strconv.Atoi(sendMsg.Content) // 传递时间戳
			if err != nil {
				timeT = 999999999
			}
			results, _ := FindHistory(c.ToID, c.ID, int64(timeT)) //获取10条历史消息
			if len(results) > 10 {
				results = results[:10]
			} else if len(results) == 0 {
				replyMsg := ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "暂无消息",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			}
			for _, result := range results {
				replyMsg := ReplyMsg{
					From:    result.From,
					Content: result.Content,
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
			}
		} else if sendMsg.Type == 3 { //群聊
			r.Incr(ctx, c.ID)                      //发消息的次数
			r.Expire(ctx, c.ID, passTime).Result() // 防止过快“分手”，建立连接三个月过期
			log.Println(c.ID, "发送消息", sendMsg.Content)
			Manager.GBroadcast <- &GBroadcast{
				GBroadcast: &Broadcast{
					Client:  c,
					Message: []byte(c.ID + ":" + sendMsg.Content),
				},
			}
		} else if sendMsg.Type == 4 {
			r.Incr(ctx, c.ID)                      //发消息的次数
			r.Expire(ctx, c.ID, passTime).Result() // 防止过快“分手”，建立连接三个月过期
			log.Println(c.ID, "发送消息", sendMsg.Content)
			b := model.IsFriends(c.ID)
			err := model.FriendNews([]byte(sendMsg.Content), c.ToID, c.ID)
			if err != nil || b {
				replyMsg := ReplyMsg{}
				if b { //已是好友
					replyMsg = ReplyMsg{
						Code:    e.WebsocketFriendFail,
						Content: e.GetMsg(e.WebsocketFriendFail),
					}
				} else {
					replyMsg = ReplyMsg{
						Code:    e.WebsocketFriendFail,
						Content: fmt.Sprintf("%v", err),
					}
				}
				msg, _ := json.Marshal(replyMsg)
				_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			}
			Manager.FBroadcast <- &FBroadcast{
				FBroadcast: &Broadcast{
					Client:  c,
					Message: []byte(sendMsg.Content),
				},
			}
		}
	}
}

func (c *Client) Write() {
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
			log.Println(c.ID, "接受消息:", string(message))
			replyMsg := ReplyMsg{
				From:    c.ToID,
				Code:    e.WebsocketSuccessMessage,
				Content: string(message),
			}
			msg, _ := json.Marshal(replyMsg)
			_ = c.Socket.WriteMessage(websocket.TextMessage, msg)
		}
	}

}
