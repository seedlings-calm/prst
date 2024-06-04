package router

import (
	"github.com/gin-gonic/gin"
	"github.com/seedlings-calm/prst/app/models"
	"github.com/seedlings-calm/prst/core"
	jwt "github.com/seedlings-calm/prst/middleware"
)

func init() {
	routerCheckRole = append(routerCheckRole, registerApiExampleRouter)
}

// registerApiExampleRouter
func registerApiExampleRouter(v1 *gin.RouterGroup, authMiddleware *jwt.GinJWTMiddleware) {

	r := v1.Group("/example")
	{
		r.GET("/:name/:phone", GetExample)
	}
}

// @Summary 展示例子
// @Description 展示例子
// @Tags Example
// @Param name path string false "名称"
// @Param phone path string false "手机号"
// @Success 200 {object} Respo
// @Router /api/v1/example/{name}/{phone} [get]
func GetExample(c *gin.Context) {

	data := models.Query{}
	ba := core.Ba.MakeContext(c).
		Bind(&data)

	err := ba.GetError()
	if err != nil {
		ba.ErrorResponse(500, err.Error())
		return
	}

	ba.SuccessResponse("", data)
}

type Respo struct {
	Info      string
	RequestId string
}
