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
		Host       string
		Port       string
		Password   string
		DB         int
		UserExpiry int64
		Bloom      struct {
			Name          string
			ErrorRate     float64
			ExpectedItems int64
		}
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
	Kafka struct {
		Groups struct {
			GoodsGroup string
			OrderGroup string
		}
		Brokers []string
		Topics  struct {
			OrderTopic string
			GoodsTopic string
		}
	}
}

var Cfg Config

func init() {
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
	viper.SetDefault("redis.userExpiry", 86400)
	viper.SetDefault("redis.bloom.name", "bloom")
	viper.SetDefault("redis.bloom.errorRate", 0.001)
	viper.SetDefault("redis.bloom.expectedItems", 10000)

	viper.SetDefault("mysql.host", "127.0.0.1")
	viper.SetDefault("mysql.port", "8806")
	viper.SetDefault("mysql.user", "root")
	viper.SetDefault("mysql.password", "root")
	viper.SetDefault("mysql.dbname", "double_token_example")
	viper.SetDefault("mongodb.uri", "mongodb://localhost:27017/double_token_example")
	viper.SetDefault("mongodb.dbname", "double_token_example")
	viper.SetDefault("jwt.secret", "your-secret-key")
	viper.SetDefault("jwt.expireTime", 86400)

	viper.SetDefault("kafka.brokers", []string{"127.0.0.1:9092"})
	viper.SetDefault("kafka.groups.goodsGroup", "goods_group")
	viper.SetDefault("kafka.groups.orderGroup", "order_group")
	viper.SetDefault("kafka.topics.goodsTopic", "goods_topic")
	viper.SetDefault("kafka.topics.orderTopic", "order_topic")

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
