package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	logging "zap-example/config"
	"zap-example/middleware"
)

func main() {
	// 读取配置文件路径参数
	configPath := flag.String("config", "config.yaml", "配置文件路径")
	flag.Parse()

	// 初始化配置
	if err := initConfig(*configPath); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
	}

	// 初始化日志系统
	logger, err := logging.NewLogger(getLogConfig())
	if err != nil {
		log.Fatalf("初始化日志失败: %v", err)
	}
	defer logger.Sync() // 确保所有日志都被写入

	// 创建Gin引擎 (注意: 这里不再使用gin.Default()，避免重复的日志中间件)
	r := gin.New()

	// 添加中间件
	// 1. 日志中间件：记录HTTP请求信息
	r.Use(middleware.Logger(logger))
	// 2. 异常恢复中间件：捕获并记录panic
	r.Use(middleware.Recovery(logger))

	// 设置路由
	setupRoutes(r)
	viper.GetString("server.address")
	err = r.Run(viper.GetString("server.address"))
	if err != nil {
		logger.Fatal("服务启动失败", zap.Error(err))
	}
	logger.Info("服务启动成功", zap.String("address", viper.GetString("server.address")))
}

// initConfig 初始化配置系统
func initConfig(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")
	viper.AutomaticEnv() // 自动读取环境变量

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("配置文件未找到，使用默认配置: %s\n", configPath)
			return nil
		}
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	fmt.Printf("使用配置文件: %s\n", viper.ConfigFileUsed())
	return nil
}

// getLogConfig 获取日志配置
func getLogConfig() logging.Config {
	return logging.Config{
		Level:         viper.GetString("log.level"),
		Dir:           viper.GetString("log.dir"),
		MaxSize:       viper.GetInt("log.max_size"),
		MaxBackups:    viper.GetInt("log.max_backups"),
		MaxAge:        viper.GetInt("log.max_age"),
		Compress:      viper.GetBool("log.compress"),
		Development:   viper.GetBool("log.development"),
		DisableCaller: viper.GetBool("log.disable_caller"),
		UseLocalTime:  viper.GetBool("log.use_local_time"),
		UseUTCTime:    viper.GetBool("log.use_utc_time"),
		SplitByLevel:  viper.GetBool("log.split_by_level"),
		LevelFiles:    viper.GetStringMapString("log.level_files"),
	}
}

// setupRoutes 设置应用路由
func setupRoutes(r *gin.Engine) {
	// 健康检查路由
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "服务正常运行"})
	})

	// 模拟错误路由，用于测试异常捕获
	r.GET("/error", func(c *gin.Context) {
		// 记录不同级别的日志
		c.Set("user_id", 123) // 示例：设置上下文信息
		c.Next()              // 调用后续中间件

		// 触发panic
		panic("这是一个测试错误")
	})

	// API组
	api := r.Group("/api")
	{
		api.GET("/health", healthCheck)
		api.GET("/users", listUsers)
		api.GET("/users/:id", getUser)
	}
}

// healthCheck 健康检查处理器
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": "1.0.0",
		"time":    time.Now().Format(time.RFC3339),
	})
}

// listUsers 获取用户列表
func listUsers(c *gin.Context) {
	// 模拟数据库查询
	time.Sleep(100 * time.Millisecond)

	users := []gin.H{
		{"id": 1, "name": "张三", "email": "zhangsan@example.com"},
		{"id": 2, "name": "李四", "email": "lisi@example.com"},
	}

	c.JSON(http.StatusOK, users)
}

// getUser 获取单个用户
func getUser(c *gin.Context) {
	id := c.Param("id")

	// 模拟数据库查询
	time.Sleep(50 * time.Millisecond)

	if id == "1" {
		c.JSON(http.StatusOK, gin.H{
			"id":    1,
			"name":  "张三",
			"email": "zhangsan@example.com",
		})
		return
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
}
