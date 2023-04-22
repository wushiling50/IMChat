package model

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var DB *gorm.DB
var mysqlLogger logger.Interface

func Database(connstring string) {
	fmt.Println(connstring)

	mysqlLogger = logger.Default.LogMode(logger.Info)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: connstring, // DSN data source name
	}), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: false,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: mysqlLogger,
	})
	if err != nil {
		panic("Mysql数据库连接错误")
	} else {
		fmt.Println("Mysql数据库连接成功")
	}

	mysqldb, err := db.DB()
	if err != nil {
		panic("连接db服务失败")
	}
	mysqldb.SetMaxIdleConns(20)                  // 设置连接池
	mysqldb.SetMaxOpenConns(100)                 //设置最大连接数
	mysqldb.SetConnMaxLifetime(time.Second * 30) //最大连接时间
	DB = db
	migration()
}
