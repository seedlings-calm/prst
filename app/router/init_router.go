package router

import (
	"github.com/gin-gonic/gin"
)

// InitRouter 路由初始化，不要怀疑，这里用到了
func InitRouter(r *gin.Engine) {
	// the jwt middleware
	var authMiddleware *interface{}

	// 注册系统路由
	InitSysRouter(r, authMiddleware)

	// 注册业务路由
	// TODO: 这里可存放业务路由，里边并无实际路由只有演示代码
	InitExamplesRouter(r, authMiddleware)
}
