package e

var MsgFlags = map[int]string{
	SUCCESS:       "操作成功",
	ERROR:         "操作失败",
	InvalidParams: "请求参数错误",

	ErrorExistUserName: "用户名已存在",
	ErrorNotExistUser:  "用户不存在",
	ErrorNotCompare:    "不匹配",

	ErrorAuthCheckTokenFail:    "Token鉴权失败",
	ErrorAuthCheckTokenTimeout: "Token已超时",
	ErrorAuthToken:             "Token生成失败",
	ErrorAuth:                  "Token错误",

	WebsocketSuccessMessage: "解析content信息",
	WebsocketSuccess:        "发送信息，请求历史纪录操作成功",
	WebsocketEnd:            "请求历史纪录，但没有更多记录了",
	WebsocketOnlineReply:    "针对回复信息在线应答成功",
	WebsocketOfflineReply:   "针对回复信息离线回答成功",
	WebsocketLimit:          "请求收到限制",
	WebsocketGroupReply:     "发送群聊消息成功",
	WebsocketFriendReply:    "好友申请发送成功",
	WebsocketFriendFail:     "好友申请操作错误",
}

// GetMsg 获取状态码对应信息
func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}
	return MsgFlags[ERROR]
}
