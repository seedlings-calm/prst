package main

import (
	"github.com/gin-gonic/gin"
	"github.com/seedlings-calm/prst/app/router"
)

var AppRouters = make([]func(r *gin.Engine), 0)

func setup() {
	//  注册路由 fixme 其他应用的路由，在本目录新建文件放在init方法
	AppRouters = append(AppRouters, router.InitRouter)
	//初始化gin配置
}

// @title prst API
// @version 0.0.1
// @description gin框架API
func main() {
	r := gin.Default()
	setup()
	for _, f := range AppRouters {
		f(r)
	}
	r.Run()
}
