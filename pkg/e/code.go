package e

import (
	"encoding/json"
	"fmt"
	"main/IMChat/serializer"
)

const (
	SUCCESS       = 200
	ERROR         = 500
	InvalidParams = 400

	//成员错误
	ErrorExistUserName  = 10001
	ErrorNotExistUser   = 10002
	ErrorPasswordSame   = 10003
	ErrorFailEncryption = 10004
	ErrorNotCompare     = 10005

	//Token错误
	ErrorAuthCheckTokenFail    = 30001 //token 错误
	ErrorAuthCheckTokenTimeout = 30002 //token 过期
	ErrorAuthToken             = 30003
	ErrorAuth                  = 30004

	//数据库错误
	ErrorDatabase = 40001
	ErrorRedis    = 40201

	WebsocketSuccessMessage = 50001
	WebsocketSuccess        = 50002
	WebsocketEnd            = 50003
	WebsocketOnlineReply    = 50004
	WebsocketOfflineReply   = 50005
	WebsocketLimit          = 50006
	WebsocketGroupReply     = 50007
	WebsocketFriendReply    = 50008
	WebsocketFriendFail     = 50009
)

// 返回错误信息
func ErrorResponse(err error) serializer.Response {
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return serializer.Response{
			Status: 40001,
			Msg:    "JSON类型不匹配",
			Error:  fmt.Sprint(err),
		}
	}
	return serializer.Response{
		Status: 40001,
		Msg:    "参数错误",
		Error:  fmt.Sprint(err),
	}
}
