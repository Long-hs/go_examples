package main

import (
	"double-token-example/internal/config"
	"double-token-example/internal/server"
)

func main() {
	// 1. 加载配置
	config.Init()

	// 2. 初始化服务
	srv := server.NewServer()

	// 3. 启动服务
	srv.Run()
}
