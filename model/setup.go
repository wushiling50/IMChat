package model

import (
	"fmt"
	"strings"

	"gopkg.in/ini.v1"
)

var (
	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string
)

// mysql数据库的启动
func SetUp() {
	file, err := ini.Load("./IMChat/conf/conf.ini") //加载配置信息
	if err != nil {
		fmt.Println("配置文件读取有误,请检查配置文件路径")
	}
	LoadMysql(file) //读取配置信息

	path := strings.Join([]string{DbUser, ":", DbPassWord, "@tcp(", DbHost, ":", DbPort, ")/", DbName, "?charset=utf8mb4&parseTime=True"}, "")
	Database(path)
}

func LoadMysql(file *ini.File) {
	Db = file.Section("mysql").Key("DB").String()
	DbHost = file.Section("mysql").Key("DbHost").String()
	DbPort = file.Section("mysql").Key("DbPort").String()
	DbUser = file.Section("mysql").Key("DbUser").String()
	DbPassWord = file.Section("mysql").Key("DbPassWord").String()
	DbName = file.Section("mysql").Key("DbName").String()

}
