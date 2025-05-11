package router

import (
	"double-token-example/internal/handler"
	"double-token-example/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(engine *gin.Engine) {
	// 创建处理器实例
	userHandler := handler.NewUserHandler()
	goodsHandler := handler.NewGoodsHandler()
	orderHandler := handler.NewOrderHandler()

	// 公共路由组
	public := engine.Group("/api")
	{
		// 用户相关
		user := public.Group("/user")
		{
			user.POST("/login", userHandler.Login)
			user.POST("/register", userHandler.Register)
		}
	}

	// 需要认证的路由组
	private := engine.Group("/api")
	private.Use(middleware.AccessAuth())
	{
		// 用户相关
		user := private.Group("/user")
		{
			user.GET("/info", userHandler.GetUserInfo)
			user.POST("/logout", userHandler.Logout)
		}

		// 商品相关
		goods := private.Group("/goods")
		{
			goods.POST("", goodsHandler.CreateGoods)
			goods.GET("", goodsHandler.GetGoodsList)
			goods.GET("/:id", goodsHandler.GetGoodsDetail)
			goods.PUT("/:id", goodsHandler.UpdateGoods)
			goods.DELETE("/:id", goodsHandler.DeleteGoods)
		}

		// 订单相关
		order := private.Group("/order")
		{
			order.POST("/create", orderHandler.CreateOrder)
			order.GET("/list", orderHandler.GetOrderList)
			order.GET("/:id", orderHandler.GetOrderDetail)
		}
	}

	// 刷新token接口
	refresh := engine.Group("/api")
	refresh.Use(middleware.RefreshAuth())
	{
		refresh.POST("/refresh", userHandler.RefreshToken)
	}
}
