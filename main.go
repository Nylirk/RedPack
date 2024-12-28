package main

import (
	"RedPack/core"
	"RedPack/global"
	"RedPack/initialize"
)

func main() {
	global.DB = initialize.InitMysql()
	if global.DB != nil {
		initialize.CreateTables()
		//关闭数据库
		db, _ := global.DB.DB()
		defer db.Close()
	}
	global.REDIS = initialize.InitRedis()
	if global.REDIS != nil {
		defer global.REDIS.Close()
	}
	core.RunWebServer()
}
