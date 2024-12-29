package core

import (
	"RedPack/initialize"
	"fmt"
	"net"
)

type server interface {
	ListenAndServe() error
}

func RunWebServer() {
	//初始化路由
	Router := initialize.Routers()
	//初始化服务
	port := 4001
	for i := 0; i < 3; i++ {
		port = port + i
		if checkPortInUse(port) {
			addr := fmt.Sprintf(":%d", port)
			s := initServer(addr, Router)
			s.ListenAndServe().Error()
		}
	}
}

func checkPortInUse(port int) bool {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false // 端口不可用
	}
	// 端口可用，关闭监听器
	listener.Close()
	return true
}
