package common

import (
	"os"
	"time"

	"github.com/gin-gonic/gin"
	cfg "github.com/seedlings-calm/prst/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type GinLogger struct {
	Logger *zap.Logger
}

func NewGinLogger(config cfg.ZapLogger) *GinLogger {
	core := newCore(config)

	logger := zap.New(core)

	return &GinLogger{Logger: logger}
}

func LumberjackLogger(config cfg.ZapLogger) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   config.FilePath,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
	}
}
func newCore(config cfg.ZapLogger) zapcore.Core {

	if cfg.AppModel == cfg.ReleaseMode {
		// 设置 Gin 的 panic handler
		gin.DefaultWriter = LumberjackLogger(config)
		gin.DefaultErrorWriter = LumberjackLogger(config)

		fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig())
		fileDebugging := zapcore.AddSync(LumberjackLogger(config))
		return zapcore.NewTee(
			zapcore.NewCore(fileEncoder, fileDebugging, config.Level),
		)
	}
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig())

	//输出到控制台
	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleDebugging, config.Level),
	)
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

func LoggerInit() *GinLogger {
	config := cfg.Config.ZapLogger
	if IsEmptyStruct(config) {
		config = LoggerDefault()
	}

	return NewGinLogger(config)

}

// 文件存储日志  配置
func fileEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006/01/02 - 15:04:05.000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

// 控制台输出  配置
func consoleEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseColorLevelEncoder, //使用带颜色编码日志级别
		EncodeTime:     zapcore.TimeEncoderOfLayout("2006/01/02 - 15:04:05.000"),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func ConfigLevelSwitchZapLevel(i int8) zapcore.Level {
	switch i {
	case -1:
		return zapcore.DebugLevel
	case 0:
		return zapcore.InfoLevel
	case 1:
		return zapcore.WarnLevel
	case 2:
		return zapcore.ErrorLevel
	case 3:
		return zapcore.DPanicLevel
	case 4:
		return zapcore.PanicLevel
	case 5:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

func LoggerDefault() cfg.ZapLogger {
	return cfg.ZapLogger{
		FilePath:   "logs/prst_logs.log",
		Level:      zapcore.InfoLevel,
		MaxSize:    1, // MB
		MaxBackups: 3,
		MaxAge:     7,     // Days
		Compress:   false, // disabled by default
	}
}
