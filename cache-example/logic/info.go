package logic

import (
	"cache-example/db"
	"cache-example/repository"
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandlerCache5(c *gin.Context) {

}

func HandlerCache4(c *gin.Context) {

}

func HandlerCache3(c *gin.Context) {

}

func HandlerCache2(c *gin.Context) {

}

func HandlerMysql1(c *gin.Context) {
	query := c.Query("id")
	if query == "" {
		_, _ = c.Writer.Write([]byte("Invalid id"))
	}
	id, err := strconv.ParseInt(query, 10, 64)
	if err != nil {
		_, _ = c.Writer.Write([]byte("Invalid id"))
		return
	}
	infoRepository := repository.NewInfoRepository()
	// 从Mysql中获取数据
	ret, err := infoRepository.GetFromMysql(id)
	if err != nil {
		log.Printf("Error getting from mysql: %v\n", err)
		_, _ = c.Writer.Write([]byte("Error getting from mysql"))
		return
	}
	log.Printf("%v\n", ret)
	_, _ = c.Writer.Write([]byte(fmt.Sprintf("%v", ret)))
}

func HandlerCache1(c *gin.Context) {
	ctx := context.Background()
	query := c.Query("id")
	if query == "" {
		_, _ = c.Writer.Write([]byte("Invalid id"))
	}
	id, err := strconv.ParseInt(query, 10, 64)
	if err != nil {
		_, _ = c.Writer.Write([]byte("Invalid id"))
		return
	}
	infoRepository := repository.NewInfoRepository()
	info := &db.Info{}

	// 从缓存中获取数据
	cache, err := infoRepository.GetFromCache(id, ctx)
	if err != nil {
		log.Printf("Error getting from cache: %v\n", err)
	}
	if cache != nil {
		info = cache
	} else {
		// 从Mysql中获取数据
		ret, err := infoRepository.GetFromMysql(id)
		if err != nil {
			log.Printf("Error getting from mysql: %v\n", err)
			_, _ = c.Writer.Write([]byte("Error getting from mysql"))
			return
		}
		info = ret
		// 保存到缓存中
		err = infoRepository.SaveToCache(info, ctx)
		if err != nil {
			log.Fatalf("Error saving to cache: %v\n", err)
		}
	}

	log.Printf("%v\n", info)
	_, _ = c.Writer.Write([]byte(fmt.Sprintf("%v", info)))

}
