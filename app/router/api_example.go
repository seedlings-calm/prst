package router

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerApiExampleRouter)
}

// registerApiExampleRouter
func registerApiExampleRouter(v1 *gin.RouterGroup, authMiddleware *interface{}) {
	r := v1.Group("/example")
	{
		r.GET("/:name", GetExample)
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
	c.AbortWithStatusJSON(http.StatusOK, Respo{
		Info: fmt.Sprintf("hello to %s", c.Param("name")),
	})
}

type Respo struct {
	Info string
}
