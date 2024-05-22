package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	jwt "github.com/seedlings-calm/prst/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerApiExampleRouter)
}

// registerApiExampleRouter
func registerApiExampleRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {

	r := v1.Group("/example")
	{
		r.GET("/:name", GetExample)
		r.GET("", GetExample)
	}
}

// @Summary 展示例子
// @Description 展示例子
// @Tags Example
// @Param name path string false "名称"
// @Success 200 {object} Respo
// @Router /api/v1/example/{name} [get]
func GetExample(c *gin.Context) {
	c.Set("status", http.StatusOK)

	if c.Param("name") == "" {
		panic("么有参数")
	}
	c.AbortWithStatusJSON(http.StatusOK, Respo{
		Info:      fmt.Sprintf("hello to %s", c.Param("name")),
		RequestId: c.GetHeader(jwt.XRequestId),
	})
}

type Respo struct {
	Info      string
	RequestId string
}
