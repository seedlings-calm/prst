package main

import (
	"github.com/gin-gonic/gin"
	"github.com/seedlings-calm/prst/app/handle"
	"github.com/seedlings-calm/prst/app/router"
	"github.com/seedlings-calm/prst/common"
	"github.com/seedlings-calm/prst/middleware"
)

var AppRouters = make([]func(r *gin.Engine, mw *middleware.GinJWTMiddleware), 0)

func setup() {
	//  注册路由 fixme 其他应用的路由，在本目录新建文件放在init方法
	AppRouters = append(AppRouters, router.InitRouter)
}

// @title prst API
// @version 0.0.1
// @description gin框架API
func main() {
	//初始化jwt
	jwtMW, err := handle.JWTInit()
	if err != nil {
		panic("初始化jwt失败")
	}
	//初始化zaplogger
	zapLogger := common.ZapLoggerInit()

	r := gin.New()
	//初始化gin配置
	r.Use(gin.Recovery())
	r.Use(middleware.CheckOptions())
	//设置日志记录
	r.Use(zapLogger.Middleware())

	setup()
	//加载所有路由
	for _, f := range AppRouters {
		f(r, jwtMW)
	}
	r.Run(":8081")
}
