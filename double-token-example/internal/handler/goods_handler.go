package handler

import (
	"double-token-example/internal/logic"
	"double-token-example/internal/model"
	"double-token-example/pkg/response"

	"github.com/gin-gonic/gin"
)

type GoodsHandler struct {
	goodsLogic *logic.GoodsLogic
}

func NewGoodsHandler() *GoodsHandler {
	return &GoodsHandler{
		goodsLogic: logic.NewGoodsLogic(),
	}
}

// CreateGoods 创建商品
func (h *GoodsHandler) CreateGoods(c *gin.Context) {
	var req model.CreateGoodsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams)
		return
	}

	// 获取当前用户ID
	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Error(c, response.ErrUnauthorized)
		return
	}

	// 设置创建者ID
	req.CreatorID = userID
	if err := h.goodsLogic.CreateGoods(c.Request.Context(), &req); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// GetGoodsList 获取商品列表
func (h *GoodsHandler) GetGoodsList(c *gin.Context) {
	var req model.GetGoodsListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Error(c, response.ErrInvalidParams)
		return
	}

	list, total, err := h.goodsLogic.GetGoodsList(c.Request.Context(), &req)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"list":  list,
		"total": total,
	})
}

// GetGoodsDetail 获取商品详情
func (h *GoodsHandler) GetGoodsDetail(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, response.ErrInvalidParams)
		return
	}

	detail, err := h.goodsLogic.GetGoodsDetail(c.Request.Context(), id)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, detail)
}

// UpdateGoods 更新商品
func (h *GoodsHandler) UpdateGoods(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, response.ErrInvalidParams)
		return
	}

	var req model.UpdateGoodsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams)
		return
	}

	// 获取当前用户ID
	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Error(c, response.ErrUnauthorized)
		return
	}

	// 设置更新者ID
	req.UpdaterID = userID

	if err := h.goodsLogic.UpdateGoods(c.Request.Context(), id, &req); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// DeleteGoods 删除商品
func (h *GoodsHandler) DeleteGoods(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.Error(c, response.ErrInvalidParams)
		return
	}

	// 获取当前用户ID
	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Error(c, response.ErrUnauthorized)
		return
	}

	if err := h.goodsLogic.DeleteGoods(c.Request.Context(), id, userID); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}
