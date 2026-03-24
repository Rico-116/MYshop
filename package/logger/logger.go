package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

var Log *zap.Logger
var Sugar *zap.SugaredLogger

func Init(env string) error {
	level := zap.NewAtomicLevel()
	//日志级别配置
	switch strings.ToLower(env) {
	case "dev", "debug":
		level.SetLevel(zap.DebugLevel)
	case "test":
		level.SetLevel(zap.InfoLevel)
	case "prod", "release":
		level.SetLevel(zap.InfoLevel)
	default:
		level.SetLevel(zap.InfoLevel)
	}
	//编码器配置
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
	//文件输出（日志轮转）
	fileWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "./log/app.log",
		MaxSize:    20,
		MaxBackups: 10,
		MaxAge:     30,
		Compress:   true,
	})
	// 4. 控制台输出
	consoleWriter := zapcore.AddSync(os.Stdout)

	var fileEncoder zapcore.Encoder
	var consoleEncoder zapcore.Encoder

	if strings.ToLower(env) == "dev" || strings.ToLower(env) == "debug" {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleEncoder = zapcore.NewConsoleEncoder(encoderConfig)
		fileEncoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		consoleEncoder = zapcore.NewJSONEncoder(encoderConfig)
		fileEncoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 5. 多路输出
	core := zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, consoleWriter, level),
		zapcore.NewCore(fileEncoder, fileWriter, level),
	)

	// 6. 创建 logger
	Log = zap.New(
		core,
		zap.AddCaller(),                   // 打印调用文件和行号
		zap.AddCallerSkip(1),              // 封装后避免 caller 偏移
		zap.Development(),                 // 开发模式更友好
		zap.AddStacktrace(zap.ErrorLevel), // error 及以上自动带堆栈
	)

	Sugar = Log.Sugar()
	return nil
}
func Sync() {
	_ = Log.Sync()
}
