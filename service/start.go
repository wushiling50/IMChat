package service

import (
	"encoding/json"
	"log"
	"main/IMChat/model"
	"main/IMChat/pkg/e"
	"main/IMChat/pkg/util"
	"regexp"
	"time"

	"github.com/gorilla/websocket"
)

func (manager *ClientManager) Start() {
	for {
		log.Println("xxxxx监听管道通信xxxxx")
		select {
		case conn := <-Manager.Register: // 建立连接
			log.Printf("有新连接: %v", conn.ID)
			Manager.Clients[conn.ID] = conn //将连接放到用户管理上
			replyMsg := &ReplyMsg{
				Code:    e.WebsocketSuccess,
				Content: "已连接至服务器",
			}
			msg, _ := json.Marshal(replyMsg)
			_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)

		case conn := <-Manager.Unregister: // 断开连接
			log.Printf("断开连接:%v", conn.ID)
			if _, ok := Manager.Clients[conn.ID]; ok {
				replyMsg := &ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "断开连接",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
				close(conn.Send)
				delete(Manager.Clients, conn.ID)
			}

		case broadcast := <-Manager.Broadcast: //1->2
			message := broadcast.Message
			toId := broadcast.Client.ToID           //2->1
			flag := false                           // 判断在线状态，默认不在线
			for id, conn := range Manager.Clients { //将内容发送到目标用户的管道中进行读取
				if id != toId {
					continue
				}
				select {
				case conn.Send <- message:
					flag = true
				default:
					close(conn.Send)
					delete(Manager.Clients, conn.ID)
				}
			}

			id := broadcast.Client.ID //1->2
			if flag {
				log.Println("对方在线")
				replyMsg := &ReplyMsg{
					From:    id,
					Code:    e.WebsocketOnlineReply,
					Content: "对方在线应答",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)

				err := model.CreateNews(id, string(message), 1, time.Now().Unix())
				if err != nil {
					log.Println("CreateMsg Err", err)
				}
			} else {
				log.Println("对方不在线")
				replyMsg := ReplyMsg{
					From:    id,
					Code:    e.WebsocketOfflineReply,
					Content: "对方不在线应答",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				err := model.CreateNews(id, string(message), 0, time.Now().Unix())
				if err != nil {
					log.Println("CreateMsg Err", err)
				}
			}
		case gbroadcast := <-Manager.GBroadcast:
			message := gbroadcast.GBroadcast.Message
			Id := gbroadcast.GBroadcast.Client.ID
			toId := gbroadcast.GBroadcast.Client.ToID //群号

			gidRe, err := regexp.Compile(toId)
			if err != nil {
				log.Println("正则表达式解析失败", err)
			}

			for id, conn := range Manager.Clients { //将内容发送到目标用户的管道中进行读取
				s1 := gidRe.FindString(id)
				if id == Id || s1 != toId {
					continue
				}
				select {
				case conn.Send <- message:
					time.Sleep(1 * time.Millisecond)
				default:
					close(conn.Send)
					delete(Manager.Clients, conn.ID)
				}
			}

			replyMsg := &ReplyMsg{
				From:    Id,
				Code:    e.WebsocketGroupReply,
				Content: "发送群聊消息",
			}
			msg, _ := json.Marshal(replyMsg)
			_ = gbroadcast.GBroadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)

			_, rmsg, _ := model.MsgSplit(string(message))
			err = model.CreateNews(Id, rmsg, 1, time.Now().Unix())
			if err != nil {
				log.Println("CreateMsg Err", err)
			}

		case fbroadcast := <-Manager.FBroadcast:
			message := fbroadcast.FBroadcast.Message
			Id := fbroadcast.FBroadcast.Client.ID     //1->2
			toId := fbroadcast.FBroadcast.Client.ToID //2->1
			umsg := util.FriendHash(message)
			if string(message) == "好友请求" {
				for id, conn := range Manager.Clients { //将内容发送到目标用户的管道中进行读取
					if id != toId {
						continue
					}
					select {
					case conn.Send <- message:
						time.Sleep(1 * time.Millisecond)
					default:
						close(conn.Send)
						delete(Manager.Clients, conn.ID)
					}
				}
				replyMsg := &ReplyMsg{
					From:    Id,
					Code:    e.WebsocketFriendReply,
					Content: "发送好友申请",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = fbroadcast.FBroadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)

				err := model.CreateNews(Id, umsg, 1, time.Now().Unix())
				if err != nil {
					log.Println("CreateMsg Err", err)
				}
			} else if string(message) == "Y" {
				for id, conn := range Manager.Clients { //将内容发送到目标用户的管道中进行读取
					if id != toId {
						continue
					}
					select {
					case conn.Send <- message:
						time.Sleep(1 * time.Millisecond)
					default:
						close(conn.Send)
						delete(Manager.Clients, conn.ID)
					}
				}
				replyMsg := &ReplyMsg{
					From:    Id,
					Code:    e.WebsocketFriendReply,
					Content: "好友申请通过",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = fbroadcast.FBroadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)

				err := model.AddFriends(Id)
				if err != nil {
					log.Println("CreateMsg Err", err)
				}
				err = model.CreateNews(Id, umsg, 1, time.Now().Unix())
				if err != nil {
					log.Println("CreateMsg Err", err)
				}
			} else if string(message) == "N" {
				for id, conn := range Manager.Clients { //将内容发送到目标用户的管道中进行读取
					if id != toId {
						continue
					}
					select {
					case conn.Send <- message:
						time.Sleep(1 * time.Millisecond)
					default:
						close(conn.Send)
						delete(Manager.Clients, conn.ID)
					}
				}
				replyMsg := &ReplyMsg{
					From:    Id,
					Code:    e.WebsocketFriendReply,
					Content: "好友申请拒绝",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = fbroadcast.FBroadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)

				var news model.News
				model.DB.Model(&model.News{}).Preload("News").Where("send_id=?", toId).
					Where("msg=?", util.FriendHash([]byte("好友请求"))).Delete(&news)

				err := model.CreateNews(Id, umsg, 1, time.Now().Unix())
				if err != nil {
					log.Println("CreateMsg Err", err)
				}
				model.DB.Model(&model.News{}).Preload("News").Where("send_id=?", Id).
					Where("msg=?", util.FriendHash(message)).Delete(&news)
			}
		}
	}
}
