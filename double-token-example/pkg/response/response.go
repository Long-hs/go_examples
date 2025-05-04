package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 响应结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// 错误码
const (
	CodeSuccess = 200
	CodeError   = 1
)

// 错误信息
var (
	ErrInvalidParams = "参数错误"
	ErrUnauthorized  = "未授权"
	ErrInternal      = "服务器内部错误"
)

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, err interface{}) {
	var code int
	var message string

	switch e := err.(type) {
	case string:
		code = CodeError
		message = e
	case error:
		code = CodeError
		message = e.Error()
	default:
		code = CodeError
		message = ErrInternal
	}

	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

// ErrorWithCode 带错误码的错误响应
func ErrorWithCode(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}
