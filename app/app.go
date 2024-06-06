package app

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seedlings-calm/prst/common"
	"github.com/seedlings-calm/prst/db"
	"github.com/seedlings-calm/prst/middleware"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type BaseApp struct {
	//gin上下文本
	Context *gin.Context
	//日志
	Zap   *zap.Logger
	Mysql *gorm.DB
}

func (e *BaseApp) MakeMysql() *BaseApp {
	e.Mysql = db.Db
	return e
}

// MakeContext 设置http上下文
func (e *BaseApp) MakeContext(c *gin.Context) *BaseApp {
	e.Context = c
	return e
}

func (e *BaseApp) MakeLog() *BaseApp {
	e.Zap = common.Zap
	return e
}

// SuccessResponse 发送成功响应
func (e *BaseApp) SuccessResponse(msg string, data interface{}) {
	xq := middleware.GenerateXRequestIdFromContext(e.Context)
	if msg == "" {
		msg = "success"
	}
	res := &common.Response{
		RequestId: xq,
		Code:      http.StatusOK,
		Message:   msg,
		Data:      data,
	}
	// switch e.Context.ContentType() {
	// case "application/json":
	// case "application/xml":
	// case "application/x-protobuf":
	// default:
	// }
	res.JsonResponse(e.Context)
}

// ErrorResponse 发送错误响应
func (e *BaseApp) ErrorResponse(code int, message string) {
	xq := middleware.GenerateXRequestIdFromContext(e.Context)
	if code == 0 {
		code = http.StatusInternalServerError
	}
	res := &common.Response{
		RequestId: xq,
		Code:      code,
		Message:   message,
		Data:      nil,
	}
	res.JsonResponse(e.Context)
}
