package handler

import (
	"double-token-example/internal/logic"
	"double-token-example/internal/model"
	"double-token-example/pkg/response"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	orderlogic *logic.OrderLogic
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		orderlogic: logic.NewOrderLogic(),
	}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req model.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams)
		return
	}

	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Error(c, response.ErrUnauthorized)
		return
	}
	req.CreatorID = userID
	if err := h.orderlogic.CreateOrder(c.Request.Context(), &req); err != nil {
		response.Error(c, err)
		return
	}
	response.Success(c, nil)
}

func (h *OrderHandler) GetOrderList(c *gin.Context) {

}

func (h *OrderHandler) GetOrderDetail(c *gin.Context) {

}
