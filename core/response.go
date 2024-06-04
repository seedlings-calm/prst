package core

import (
	"github.com/gin-gonic/gin"
)

// Response 定义通用的 HTTP 响应结构
type Response struct {
	RequestId string      `json:"requestId,omitempty"`
	Code      int         `json:"code"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
}

// NewResponse 创建一个新的响应
func NewResponse(code int, message string, data interface{}, reqId string) *Response {
	return &Response{
		Code:      code,
		Message:   message,
		Data:      data,
		RequestId: reqId,
	}
}

// 发送标准化的 JSON 响应
func jsonResponse(c *gin.Context, res Response) {
	c.AbortWithStatusJSON(res.Code, res)
}
