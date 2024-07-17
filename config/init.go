package cfg

import (
	"flag"
	"log"
	"os"

	"github.com/spf13/viper"
)

var (
	globalConf Conf
)

// 配置来决定项目运行模式

func Setup() *viper.Viper {
	log.Println("配置初始化")
	var config string
	// 定义命令行参数
	flag.StringVar(&config, "c", "", "choose config file.")
	flag.Parse()
	switch config {
	case "":
		config = ConfigDefaultFile
		AppModel = DevModel
	case ConfigDefaultFile, ConfigDevFile:
		AppModel = DevModel
	case ConfigReleaseFile:
		AppModel = ReleaseMode
	default:
		AppModel = DevModel
	}

	// 设置配置文件
	v := viper.New()
	v.SetConfigFile(config) // 配置文件名
	v.SetConfigType("yml")  // 配置文件类型
	v.AddConfigPath(".")    // 查找配置文件的路径，可以是绝对路径或相对路径

	// 检查配置文件是否存在
	if _, err := os.Stat(config); os.IsNotExist(err) {
		log.Fatalf("Config file '%s' does not exist", config)
	}
	err := v.ReadInConfig()
	if err != nil {
		log.Fatalf("Fatal error config file: %s \n", err)
	}
	var conf FileConfig
	if err = v.Unmarshal(&conf); err != nil {
		panic(err)
	}

	globalConf.L.Lock()
	globalConf.Config = conf
	globalConf.L.Unlock()

	//二次确定模式
	switch conf.App.Model {
	case DevModel:
		AppModel = DevModel
	case ReleaseMode:
		AppModel = ReleaseMode
	}

	return v
}

func GetGlobalConf() FileConfig {
	globalConf.L.RLock()
	defer globalConf.L.RUnlock()
	return globalConf.Config
}
