package cfg

import "github.com/gin-gonic/gin"

var (
	DevModel          = "dev"
	ReleaseMode       = "release"
	AppModel          string // 开发模式： dev,release
	ConfigDefaultFile = "config.dev.yml"
	ConfigDevFile     = "config.dev.yml"
	ConfigReleaseFile = "config.yml"
)

var (
	//全局配置
	Config FileConfig
)

func ModelSwitchGinModel() string {
	switch AppModel {
	case "dev":
		return gin.DebugMode
	case "release":
		return gin.ReleaseMode
	}
	return gin.DebugMode
}
