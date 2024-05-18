package cfg

import "github.com/gin-gonic/gin"

var (
	DevModel          = "dev"
	ReleaseMode       = "release"
	AppModel          string // 开发模式： dev,release   优先级: *.yml配置 > 文件分类配置（config.yml 为release, 其余配置文件为dev）
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
	case DevModel:
		return gin.DebugMode
	case ReleaseMode:
		return gin.ReleaseMode
	}
	return gin.DebugMode
}
