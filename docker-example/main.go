package main

import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
	r.GET("/", func(context *gin.Context) {
		context.String(200, "看到消息就说明部署成功了")
	})
	r.Run(":5000")
}
