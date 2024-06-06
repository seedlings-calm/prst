package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/seedlings-calm/prst/common"
	"go.uber.org/zap"
)

func InitEngine() *gin.Engine {

	r := gin.New()

	if gin.Mode() == gin.DebugMode {
		r.Use(gin.Logger())
		r.Use(gin.Recovery())
	} else {
		r.Use(common.LoggerMiddleware())
		r.Use(Recovery())
	}

	r.Use(RequestId(XRequestId))
	r.Use(Options()).Use(Secure)
	return r
}

func Recovery() gin.HandlerFunc {
	DefaultErrorWriter := &PanicExceptionRecord{}
	return gin.RecoveryWithWriter(DefaultErrorWriter, func(c *gin.Context, err interface{}) {
		// 这里针对发生的panic等异常进行统一响应即可
		// 这里的 err 数据类型为 ：runtime.boundsError  ，需要转为普通数据类型才可以输出
		xq := GenerateXRequestIdFromContext(c)
		er := convertToError(err)
		common.NewResponse(http.StatusInternalServerError, fmt.Sprintf("%s", er), nil, xq).JsonResponse(c)
	})
}
func convertToError(err interface{}) error {
	switch e := err.(type) {
	case runtime.Error:
		return errors.New(e.Error())
	case error:
		return e
	default:
		return fmt.Errorf("caught an unknown error: %v", err)
	}
}

// PanicExceptionRecord  panic等异常记录
type PanicExceptionRecord struct{}

func (p *PanicExceptionRecord) Write(b []byte) (n int, err error) {
	errStr := string(b)
	err = errors.New(errStr)
	common.Zap.Error("gin-recovery", zap.Error(err))
	return len(errStr), err
}
