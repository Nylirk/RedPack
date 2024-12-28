package initialize

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Addr     string
	Password string
	DB       int
}

func InitRedis() *redis.Client {
	config := Redis{
		Addr:     "121.40.25.175:6379",
		Password: "",
		DB:       0,
	}
	client := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     10,
		MinIdleConns: 5,
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("redis服务连接失败，错误：", err)
		panic(err)
	} else {
		fmt.Println("redis服务连接成功")
		return client
	}
}
