package cfg

import "go.uber.org/zap/zapcore"

type FileConfig struct {
	App       App         `json:"app" yaml:"app" mapstructure:"app"`
	Ssl       Ssl         `json:"ssl" yaml:"ssl" mapstructure:"ssl"`
	ZapLogger ZapLogger   `json:"logger" yaml:"logger" mapstructure:"logger"`
	Mysql     MysqlConfig `json:"mysql" yaml:"mysql" mapstructure:"mysql"`
	Redis     []Redis     `mapstructure:"redis" json:"redis" yaml:"redis"`
}

type Ssl struct {
	Key    string `json:"key" yaml:"key" mapstructure:"key"`
	Pem    string `json:"pem" yaml:"pem" mapstructure:"pem"`
	Domain string `json:"domain" yaml:"domain" mapstructure:"domain"`
}
type App struct {
	Model  string `json:"model" yaml:"model" mapstructure:"model"`
	Host   string `json:"host" yaml:"host" mapstructure:"host"`
	Port   int    `json:"port" yaml:"port" mapstructure:"port"`
	Enable bool   `json:"enable" yaml:"enable" mapstructure:"enable"`
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

type MysqlConfig struct {
	Host     string `json:"host" yaml:"host" mapstructure:"host"`
	Port     int    `json:"port" yaml:"port" mapstructure:"port"`
	User     string `json:"user" yaml:"user" mapstructure:"user"`
	Password string `json:"password" yaml:"password" mapstructure:"password"`
	DBName   string `json:"dbName" yaml:"dbName" mapstructure:"dbName"`
	OpenConn int    `json:"openConn" yaml:"openConn" mapstructure:"openConn"`
	IdleConn int    `json:"idleConn" yaml:"idleConn" mapstructure:"idleConn"`
}

type Redis struct {
	Addr     string `mapstructure:"addr" json:"addr" yaml:"addr"`             // 服务器地址
	Port     int    `json:"port" yaml:"port" mapstructure:"port"`             //端口
	Password string `mapstructure:"password" json:"password" yaml:"password"` // 密码
	DB       int    `mapstructure:"db" json:"db" yaml:"db"`                   // 单实例模式下redis的哪个数据库
	Name     string `mapstructure:"name" json:"name" yaml:"name"`             //
}
