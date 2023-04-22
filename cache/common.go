package cache

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	logging "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

// RedisClient Redis缓存客户端单例
var (
	RedisClient *redis.Client
	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string
)

// Redis 在中间件中初始化redis链接  防止循环导包，所以放在这里
func init() {
	file, err := ini.Load("./IMChat/conf/conf.ini") //加载配置信息
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径:", err)
	}
	LoadRedisData(file) //读取配置信息
	Redis()
}

func LoadRedisData(file *ini.File) {
	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPw = file.Section("redis").Key("RedisPw").String()
	RedisDbName = file.Section("redis").Key("RedisDbName").String()
}

// Redis 在中间件中初始化redis链接
func Redis() {
	ctx := context.Background()
	db, _ := strconv.ParseUint(RedisDbName, 10, 64)
	client := redis.NewClient(&redis.Options{
		Addr: RedisAddr,
		// Password: RedisPw,
		DB: int(db),
	})
	_, err := client.Ping(ctx).Result()
	if err != nil {
		logging.Info(err)
		panic(err)
	}
	RedisClient = client
}
