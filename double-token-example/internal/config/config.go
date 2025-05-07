package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string
	}
	Redis struct {
		Host     string
		Port     string
		Password string
		DB       int
	}
	MySQL struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
	MongoDB struct {
		URI    string
		DBName string
	}
	JWT struct {
		Secret     string
		ExpireTime int64
	}
	Kakfa struct {
		Brokers []string
		Topic   string
		Group   string
	}
}

var Cfg Config

func Init() {
	// 设置配置文件路径
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")

	// 设置默认值
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("redis.host", "127.0.0.1")
	viper.SetDefault("redis.port", "6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("mysql.host", "127.0.0.1")
	viper.SetDefault("mysql.port", "8806")
	viper.SetDefault("mysql.user", "root")
	viper.SetDefault("mysql.password", "root")
	viper.SetDefault("mysql.dbname", "double_token_example")
	viper.SetDefault("mongodb.uri", "mongodb://localhost:27017/double_token_example")
	viper.SetDefault("mongodb.dbname", "double_token_example")
	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.expire_time", 86400)

	viper.SetDefault("kafka.brokers", []string{"127.0.0.1:9092"})
	viper.SetDefault("kafka.group", "double_token_example_group")

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: %v\n", err)
	}

	// 读取环境变量
	viper.AutomaticEnv()

	// 绑定配置到结构体
	if err := viper.Unmarshal(&Cfg); err != nil {
		log.Fatalf("Error unmarshaling config: %v", err)
	}

	// 打印配置信息
	log.Printf("Config loaded: %+v\n", Cfg)
}
