package router

import (
	"github.com/gin-gonic/gin"
	"github.com/seedlings-calm/prst/app/api"
	jwt "github.com/seedlings-calm/prst/common"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerApiExampleRouter)
}

// registerApiExampleRouter
func registerApiExampleRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {
	exa := api.Example{}
	r := v1.Group("/example")
	{
		r.GET("/:name/:phone", exa.GetExample)
		r.GET("/redis", exa.GetRedis)
	}
}
