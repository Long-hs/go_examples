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
	r.GET("/cache2", logic.HandlerCache2)
	r.GET("/cache3", logic.HandlerCache3)
	r.GET("/cache4", logic.HandlerCache4)
	r.GET("/cache5", logic.HandlerCache5)
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
