package main

import (
	"github.com/gin-gonic/gin"
	"github.com/seedlings-calm/prst/app/router"
	"github.com/seedlings-calm/prst/common"
	cfg "github.com/seedlings-calm/prst/config"
	"github.com/seedlings-calm/prst/middleware"
)

var AppRouters = make([]func(r *gin.Engine, mw *middleware.GinJWTMiddleware), 0)

func init() {
	//  注册路由 fixme 其他应用的路由，在本目录新建文件放在init方法
	AppRouters = append(AppRouters, router.InitRouter)
}

// @title prst API
// @version 0.0.1
// @description gin框架API
func main() {
	//初始化配置信息
	_ = cfg.Setup()
	//初始化jwt
	jwtMW, err := common.JWTInit()
	if err != nil {
		panic("初始化jwt失败")
	}
	//初始化zaplogger
	_ = common.LoggerInit()

	gin.SetMode(cfg.ModelSwitchGinModel())

	r := gin.New()
	//适配gin的运行模式

	middleware.InitMiddleWare(r)

	// r.Use(zapLogger.Middleware())

	//加载所有路由
	for _, f := range AppRouters {
		f(r, jwtMW)
	}
	r.Run(cfg.Config.App.Host + ":" + cfg.Config.App.Port)
}
