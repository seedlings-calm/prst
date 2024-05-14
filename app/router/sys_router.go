package router

import (
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"

	swaggerfiles "github.com/swaggo/files"

	_ "github.com/seedlings-calm/prst/docs" // swagger读取文档配置路径
	jwt "github.com/seedlings-calm/prst/middleware"
)

func InitSysRouter(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) *gin.RouterGroup {
	g := r.Group("")
	// 静态文件
	sysStaticFileRouter(g)
	// swagger；注意：生产环境可以注释掉
	sysSwaggerRouter(g)
	return g
}

func sysStaticFileRouter(r *gin.RouterGroup) {
	r.Static("/static", "./static")
}

func sysSwaggerRouter(r *gin.RouterGroup) {
	r.GET("/swagger/prst/*any", ginSwagger.WrapHandler(swaggerfiles.NewHandler(), ginSwagger.InstanceName("prst")))
}
