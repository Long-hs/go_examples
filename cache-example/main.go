package main

import (
	"cache-example/db"
	"cache-example/logic"

	"github.com/gin-gonic/gin"
)

func main() {
	db.NewMysqlDB()
	db.NewRedisDB()

	db.KafkaServer = db.NewKafkaServer([]string{"127.0.0.1:9092"}, []string{"cache_example"}, "cache_example_group")

	r := gin.Default()
	// 缓存回溯
	r.GET("/cache1", logic.HandlerCache1)
	r.GET("/mysql1", logic.HandlerMysql1)
	// 双写
	r.GET("/doubleWrite", logic.HandlerDoubleWrite)
	//读更新写删除
	r.GET("/readUpdate", logic.HandlerRU)
	r.GET("/writeDelete", logic.HandlerWD)

	// 延时双删
	r.GET("/delayedDoubleDel", logic.HandlerDelayedDoubleDel)

	//异步更新
	r.GET("/asyncUpdate", logic.HandlerAsyncUpdate)
	err := r.Run(":8080")
	if err != nil {
		panic(err)
	}
}
