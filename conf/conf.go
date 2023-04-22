package conf

import (
	"fmt"

	"gopkg.in/ini.v1"
)

var (
	AppMode  string
	HttpPort string
)

// 服务的启动
func SetUp() {
	file, err := ini.Load("./IMChat/conf/conf.ini") //读取配置信息
	if err != nil {
		fmt.Println("配置文件读取有误,请检查配置文件路径")
	}
	LoadServer(file) //读取配置信息
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("service").Key("AppMode").String()
	HttpPort = file.Section("service").Key("HttpPort").String()
}
