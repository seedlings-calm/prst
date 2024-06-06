package api

import (
	"github.com/gin-gonic/gin"
	"github.com/seedlings-calm/prst/app"
	"github.com/seedlings-calm/prst/app/models"
	"go.uber.org/zap"
)

type Example struct {
	app.BaseApp
}

// @Summary 展示例子
// @Description 展示例子
// @Tags Example
// @Param name path string false "名称"
// @Param phone path string false "手机号"
// @Success 200 {object} core.Response{}
// @Router /api/v1/example/{name}/{phone} [get]
func (e Example) GetExample(c *gin.Context) {

	data := models.Query{}
	ba := e.MakeContext(c).MakeLog()

	err := ba.Bind(&data)
	panic("sdfsdfsdf")
	if err != nil {
		ba.Zap.Error("get_example", zap.Error(err))
		ba.ErrorResponse(500, err.Error())
		return
	}

	ba.SuccessResponse("", data)
}
