package core

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/seedlings-calm/prst/common"
	cfg "github.com/seedlings-calm/prst/config"
	"github.com/seedlings-calm/prst/db"
	"github.com/seedlings-calm/prst/middleware"
	"gorm.io/gorm"
)

var Ba = BaseApp{
	Log: common.NewGinLogger(cfg.Config.ZapLogger),
}

type BaseApp struct {
	//gin上下文本
	Context *gin.Context
	//日志
	Log    *common.GinLogger
	Mysql  *gorm.DB
	errors error
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

func (e *BaseApp) MakeLog(config cfg.ZapLogger) {
	e.Log = common.NewGinLogger(config)
}

func (e *BaseApp) AddError(err error) {
	if e.errors == nil {
		e.errors = err
	} else if err != nil {
		e.errors = fmt.Errorf("%v; %w", e.errors, err)
	}
}
func (e *BaseApp) GetError() error {
	defer func() {
		e.errors = nil
	}()
	return e.errors
}

// SuccessResponse 发送成功响应
func (e *BaseApp) SuccessResponse(msg string, data interface{}) {
	xq := middleware.GenerateXRequestIdFromContext(e.Context)
	if msg == "" {
		msg = "success"
	}
	res := Response{
		RequestId: xq,
		Code:      http.StatusOK,
		Message:   msg,
		Data:      data,
	}
	jsonResponse(e.Context, res)
}

// ErrorResponse 发送错误响应
func (e *BaseApp) ErrorResponse(code int, message string) {
	xq := middleware.GenerateXRequestIdFromContext(e.Context)
	res := Response{
		RequestId: xq,
		Code:      code,
		Message:   message,
		Data:      nil,
	}
	jsonResponse(e.Context, res)
}
