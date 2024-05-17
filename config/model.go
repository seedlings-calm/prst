package cfg

import "go.uber.org/zap/zapcore"

type FileConfig struct {
	App       App       `json:"app" yaml:"app" mapstructure:"app"`
	ZapLogger ZapLogger `json:"logger" yaml:"logger" mapstructure:"logger"`
}

type App struct {
	Host string `json:"app" yaml:"host" mapstructure:"host"`
	Port string `json:"port" yaml:"port" mapstructure:"port"`
}

// ZapLogger 定义日志配置参数
type ZapLogger struct {
	// Stdout     string        ` json:"stdout" yaml:"stdout" mapstructure:"stdout"`             //日志输出方式 default:控制台
	FilePath   string        ` json:"filePath" yaml:"filePath" mapstructure:"filePath"`       //日志文件位置
	MaxSize    int           ` json:"maxSize" yaml:"maxSize" mapstructure:"maxSize" `         //  进行切割之前，日志文件最大值(单位：MB)
	MaxBackups int           ` json:"maxBackups" yaml:"maxBackups" mapstructure:"maxBackups"` //保留旧文件的最大个数
	MaxAge     int           ` json:"maxAge" yaml:"maxAge" mapstructure:"maxAge"`             //  保留旧文件的最大天数
	Level      zapcore.Level ` json:"level" yaml:"level" mapstructure:"level"`
	Compress   bool          ` json:"compress" yaml:"compress" mapstructure:"compress"` //是否压缩/归档旧文件
}
