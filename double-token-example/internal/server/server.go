package server

import (
	"double-token-example/internal/config"
	"double-token-example/internal/db"
	"double-token-example/internal/kafka"
	"double-token-example/internal/middleware"
	"double-token-example/internal/router"

	"github.com/gin-gonic/gin"
)

type Server struct {
	Engine *gin.Engine
}

func init() {
	// 初始化数据库
	db.InitMySQL()
	db.InitMongoDB()
	db.InitRedisDB()

	//初始化kafka
	kafka.InitKafkaServer()
}

func NewServer() *Server {
	// 创建gin引擎
	engine := gin.Default()

	// 注册中间件
	engine.Use(middleware.CORS())
	// engine.Use(middleware.Logger())
	// engine.Use(middleware.Recovery())

	// 注册路由
	router.RegisterRoutes(engine)

	return &Server{
		Engine: engine,
	}
}

func (s *Server) Run() {
	// 启动服务器
	if err := s.Engine.Run(":" + config.Cfg.Server.Port); err != nil {
		panic(err)
	}
}
