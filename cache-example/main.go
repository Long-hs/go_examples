package main

import (
	"cache-example/db"
	"cache-example/logic"
	"github.com/gin-gonic/gin"
)

func main() {
	db.NewMysqlDB()
	db.NewRedisDB()

	r := gin.Default()
	// 缓存回溯
	r.GET("/cache1", logic.HandlerCache1)
	r.GET("/mysql1", logic.HandlerMysql1)
	// 双写
	r.GET("/doubleWrite", logic.HandlerDoubleWrite)
	//读更新写删除
	r.GET("/readUpdate", logic.HandlerRU)
	r.GET("/writeDelete", logic.HandlerWD)

	r.GET("/cache4", logic.HandlerCache4)
	r.GET("/cache5", logic.HandlerCache5)
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
