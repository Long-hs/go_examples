package db

import (
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoDB *mongo.Client

// InitMongoDB 初始化MongoDB连接
func InitMongoDB() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/double_token_example")

	var err error
	mongoDB, err = mongo.Connect(nil, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// 测试连接
	err = mongoDB.Ping(nil, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("MongoDB connected successfully")
}

// GetMongoDB 获取MongoDB连接
func GetMongoDB() *mongo.Client {
	return mongoDB
}

// GetMongoDBCollection 获取MongoDB集合
func GetMongoDBCollection(collectionName string) *mongo.Collection {
	return mongoDB.Database("double_token_example").Collection(collectionName)
}
