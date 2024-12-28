package core

import (
	"RedPack/initialize"
)

type server interface {
	ListenAndServe() error
}

func RunWebServer() {
	//初始化路由
	Router := initialize.Routers()
	//初始化服务
	s := initServer(":4001", Router)
	s.ListenAndServe().Error()
}
