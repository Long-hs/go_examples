package logic

import (
	"cache-example/db"
	"cache-example/repository"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/IBM/sarama"

	"github.com/gin-gonic/gin"
)

func HandlerAsyncUpdate(c *gin.Context) {
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
	marshaledInfo, err := json.Marshal(info)
	if err != nil {
		log.Printf("Error marshaling info: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error marshaling info"})
		return
	}
	msg := &sarama.ProducerMessage{
		Topic: db.KafkaServer.Topics[0],
		Value: sarama.StringEncoder(marshaledInfo),
	}
	partition, offset, err := db.KafkaServer.SyncProducer.SendMessage(msg)
	if err != nil {
		log.Printf("Failed to send message to Kafka: %v, info: %+v", err, info)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message to Kafka"})
		return
	}
	log.Printf("Successfully sent message to Kafka: topic=%s, partition=%d, offset=%d, info: %+v",
		db.KafkaServer.Topics[0], partition, offset, info)

	// 修改缓存
	err = infoRepository.SaveToCache(info, context.Background())
	if err != nil {
		log.Printf("Error saving to cache: %v\n", err)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Update success"})
}

func HandlerDelayedDoubleDel(c *gin.Context) {
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
	// 删除缓存
	err = infoRepository.DeleteFromCache(id, context.Background())
	if err != nil {
		log.Printf("Error deleting from cache: %v\n", err)
	}
	// 修改数据库
	err = infoRepository.UpdateToMysql(info)
	if err != nil {
		log.Printf("Error updating to mysql: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating to mysql"})
		return
	}
	// 删除缓存
	go func() {
		time.Sleep(1 * time.Millisecond)
		err = infoRepository.DeleteFromCache(id, context.Background())
		if err != nil {
			log.Printf("Error deleting from cache: %v\n", err)
		}
	}()

	c.JSON(http.StatusOK, gin.H{"message": "Update success"})
}

func HandlerRU(c *gin.Context) {
	// 复用HandlerCache1
	HandlerCache1(c)
}

func HandlerWD(c *gin.Context) {
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
	err = infoRepository.UpdateToMysql(info)
	if err != nil {
		log.Printf("Error updating to mysql: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error updating to mysql"})
	}
	//删除缓存
	err = infoRepository.DeleteFromCache(id, context.Background())
	if err != nil {
		log.Printf("Error deleting from cache: %v\n", err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Update success"})
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
