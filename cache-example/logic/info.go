package logic

import (
	"cache-example/db"
	"cache-example/repository"
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func HandlerCache5(c *gin.Context) {

}

func HandlerCache4(c *gin.Context) {

}

func HandlerCache3(c *gin.Context) {

}

func HandlerDoubleWrite(c *gin.Context) {
	idStr := c.Query("id")
	name := c.Query("name")
	if idStr == "" || name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id or name"})
		return
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	infoRepository := repository.NewInfoRepository()
	info := &db.Info{
		ID:   id,
		Name: name,
	}
	//修改数据库
	//开启事务
	db.DB.Begin()
	err = infoRepository.UpdateToMysql(info)
	if err != nil {
		//回滚事务
		db.DB.Rollback()
		log.Printf("Error updating to mysql: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating to mysql"})
		return
	}
	//修改缓存
	err = infoRepository.SaveToCache(info, context.Background())
	if err != nil {
		//回滚事务
		db.DB.Rollback()
		log.Printf("Error saving to cache: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error saving to cache"})
		return
	}
	//提交事务
	db.DB.Commit()
	c.JSON(http.StatusOK, gin.H{"message": "Update success"})
}

func HandlerMysql1(c *gin.Context) {
	query := c.Query("id")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	id, err := strconv.ParseInt(query, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	infoRepository := repository.NewInfoRepository()
	// 从Mysql中获取数据
	ret, err := infoRepository.GetFromMysql(id)
	if err != nil {
		log.Printf("Error getting from mysql: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting from mysql"})
		return
	}
	log.Printf("%v\n", ret)
	c.JSON(http.StatusOK, ret)
}

func HandlerCache1(c *gin.Context) {
	ctx := context.Background()
	query := c.Query("id")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	id, err := strconv.ParseInt(query, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting from mysql"})
			return
		}
		info = ret
		// 保存到缓存中
		err = infoRepository.SaveToCache(info, ctx)
		if err != nil {
			log.Printf("Error saving to cache: %v\n", err)
		}
	}

	log.Printf("%v\n", info)
	c.JSON(http.StatusOK, info)
}
