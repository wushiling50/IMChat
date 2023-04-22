package service

import (
	"main/IMChat/model"
)

type Result struct {
	StartTime int64
	Content   interface{}
	From      string
}
type SendSortMsg struct {
	Msg      string `json:"msg"`
	Read     uint   `json:"read"`
	CreateAt int64  `json:"create_at"`
}

func FindHistory(toId, id string, time int64) (results []Result, err error) {
	var resultsEach []model.News //双方之间的消息记录

	//判断时间戳
	if time == 999999999 {
		err = model.DB.Table("news").Order("time desc").Where("send_id in ?", []string{id, toId}).
			Scan(&resultsEach).Error //获取消息
	} else {
		err = model.DB.Table("news").Order("time desc").Where("send_id in ?", []string{id, toId}).
			Where("time>=?", time).Scan(&resultsEach).Error //获取消息
	}

	for _, r := range resultsEach {
		sendSort := SendSortMsg{
			Msg:      r.Msg,
			Read:     r.ReadOrNot,
			CreateAt: r.Time,
		}
		result := Result{
			StartTime: r.Time,
			Content:   sendSort,
			From:      r.SendID,
		}
		results = append(results, result)
	}
	return results, err
}
