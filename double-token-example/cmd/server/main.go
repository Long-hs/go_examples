package main

import (
	"double-token-example/internal/server"
)

func main() {

	// 初始化服务
	srv := server.NewServer()

	// 启动服务
	srv.Run()
}
