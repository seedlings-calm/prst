package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger"

	swaggerfiles "github.com/swaggo/files"

	jwt "github.com/seedlings-calm/prst/common"
	cfg "github.com/seedlings-calm/prst/config"
	_ "github.com/seedlings-calm/prst/docs" // swagger读取文档配置路径
)

func InitSysRouter(r *gin.Engine, authMiddleware *jwt.GinJWTMiddleware) *gin.RouterGroup {
	g := r.Group("")
	// 静态文件
	sysStaticFileRouter(g)
	// swagger；注意：生产环境可以注释掉
	if cfg.AppModel != cfg.ReleaseMode {
		sysSwaggerRouter(g)
	}
	sysPrometheusRouter(g)
	g.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"msg": "pong"})
	})
	return g
}

func sysPrometheusRouter(r *gin.RouterGroup) {
	// 设置 Prometheus 监听路径
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))
}

func sysStaticFileRouter(r *gin.RouterGroup) {
	r.Static("/static", "./static")
}

func sysSwaggerRouter(r *gin.RouterGroup) {
	r.GET("/swagger/prst/*any", ginSwagger.WrapHandler(swaggerfiles.NewHandler(), ginSwagger.InstanceName("prst")))
}
