package handler

import (
	"double-token-example/internal/logic"
	"double-token-example/internal/model"
	"double-token-example/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userLogic *logic.UserLogic
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userLogic: logic.NewUserLogic(),
	}
}

// RegisterRequest 注册请求
type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req model.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams)
		return
	}

	if err := h.userLogic.Register(c.Request.Context(), req.Username, req.Password, req.Phone, req.Email); err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, nil)
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req model.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, response.ErrInvalidParams)
		return
	}

	token, err := h.userLogic.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		response.Error(c, err)
		return
	}

	response.Success(c, gin.H{
		"token": token,
	})
}

// GetUserInfo 获取用户信息
func (h *UserHandler) GetUserInfo(c *gin.Context) {
	userID := c.GetInt64("user_id")
	if userID == 0 {
		response.Error(c, response.ErrUnauthorized)
		return
	}

	user, err := h.userLogic.GetUserInfo(c.Request.Context(), userID)
	if err != nil {
		response.Error(c, err)
		return
	}

	// 过滤敏感信息
	response.Success(c, gin.H{
		"id":       user.ID,
		"username": user.Username,
		"phone":    user.Phone,
		"email":    user.Email,
		"status":   user.Status,
	})
}
