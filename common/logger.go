package common

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// ZapLogger 定义日志配置参数
type ZapLogger struct {
	FilePath   string //日志文件位置
	MaxSize    int    //  进行切割之前，日志文件最大值(单位：MB)
	MaxBackups int    //保留旧文件的最大个数
	MaxAge     int    //  保留旧文件的最大天数
	Level      zapcore.Level
	Compress   bool //是否压缩/归档旧文件
}

type GinLogger struct {
	Logger *zap.Logger
}

func NewGinLogger(config ZapLogger) *GinLogger {
	core := newCore(config)

	logger := zap.New(core)

	return &GinLogger{Logger: logger}
}
func newCore(config ZapLogger) zapcore.Core {
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename:   config.FilePath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleDebugging := zapcore.Lock(os.Stdout)

	fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
	fileDebugging := zapcore.AddSync(lumberjackLogger)

	//写入文件和输出到控制台
	core := zapcore.NewTee(
		//控制台
		zapcore.NewCore(consoleEncoder, consoleDebugging, zapcore.DebugLevel),
		//文件
		zapcore.NewCore(fileEncoder, fileDebugging, zapcore.DebugLevel),
	)
	return core
}
func (g *GinLogger) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		// 记录请求信息
		g.Logger.Info(
			c.Request.Method,
			zap.String("path", c.Request.URL.Path),
			zap.String("ip", c.ClientIP()),
			zap.Duration("latency", time.Since(start)),
			zap.Int("status", c.Writer.Status()),
		)
		// 处理请求
		c.Next()

	}
}

// TODO:  根据配置文件，更改初始化
func ZapLoggerInit() *GinLogger {

	config := ZapLogger{
		FilePath:   "./log/prst_zap.log",
		Level:      zapcore.InfoLevel,
		MaxSize:    1, // MB
		MaxBackups: 3,
		MaxAge:     28,    // Days
		Compress:   false, // disabled by default
	}

	return NewGinLogger(config)

}
