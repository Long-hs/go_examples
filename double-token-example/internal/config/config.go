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
		Secret                 string
		AccessTokenExpireTime  int64
		RefreshTokenExpireTime int64
		AccessTokenType        string
		RefreshTokenType       string
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
