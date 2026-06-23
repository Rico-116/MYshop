package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

var Log *zap.Logger
var Sugar *zap.SugaredLogger
var accessLogMode = "slow"
var slowRequestThreshold = time.Second

func Init(env string) error {
	if env == "" {
		env = os.Getenv("APP_ENV")
	}
	if env == "" {
		env = "dev"
	}
	env = strings.ToLower(env)
	accessLogMode = strings.ToLower(os.Getenv("ACCESS_LOG"))
	if accessLogMode == "" {
		accessLogMode = "slow"
	}
	if slowMs, err := strconv.Atoi(os.Getenv("SLOW_REQUEST_MS")); err == nil && slowMs > 0 {
		slowRequestThreshold = time.Duration(slowMs) * time.Millisecond
	}

	level := zap.NewAtomicLevel()
	//日志级别配置
	switch env {
	case "dev", "debug":
		level.SetLevel(zap.DebugLevel)
	case "test":
		level.SetLevel(zap.InfoLevel)
	case "prod", "release":
		level.SetLevel(zap.InfoLevel)
	default:
		level.SetLevel(zap.InfoLevel)
	}
	if err := os.MkdirAll("./log", 0755); err != nil {
		return err
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

	if env == "dev" || env == "debug" {
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
	options := []zap.Option{
		zap.AddCaller(),                   // 打印调用文件和行号
		zap.AddStacktrace(zap.ErrorLevel), // error 及以上自动带堆栈
		zap.Fields(
			zap.String("app", "myshop"),
			zap.String("env", env),
		),
	}
	if env == "dev" || env == "debug" {
		options = append(options, zap.Development())
	}
	Log = zap.New(core, options...)

	Sugar = Log.Sugar()
	return nil
}
func Sync() {
	if Log == nil {
		return
	}
	_ = Log.Sync()
}

func InitLogger() error {
	var err error

	// 开发阶段用 NewDevelopment，日志更清楚
	Log, err = zap.NewDevelopment()
	if err != nil {
		return err
	}

	return nil
}

func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = fmt.Sprintf("%d", start.UnixNano())
		}
		c.Header("X-Request-ID", requestID)
		c.Set("request_id", requestID)

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		fields := []zap.Field{
			zap.String("request_id", requestID),
			zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path),
			zap.String("query", c.Request.URL.RawQuery),
			zap.Int("status", status),
			zap.Duration("latency", latency),
			zap.Float64("latency_ms", float64(latency.Microseconds())/1000),
			zap.String("client_ip", c.ClientIP()),
			zap.String("user_agent", c.Request.UserAgent()),
		}
		if len(c.Errors) > 0 {
			fields = append(fields, zap.String("errors", c.Errors.String()))
		}

		switch {
		case status >= http.StatusInternalServerError:
			Log.Error("http request completed", fields...)
		case status >= http.StatusBadRequest:
			Log.Warn("http request completed", fields...)
		case accessLogMode == "all":
			Log.Info("http request completed", fields...)
		case accessLogMode == "slow" && latency >= slowRequestThreshold:
			Log.Warn("slow http request completed", fields...)
		default:
			return
		}
	}
}

func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				Log.Error("http panic recovered",
					zap.String("panic", fmt.Sprint(err)),
					zap.String("request_id", c.GetString("request_id")),
					zap.String("method", c.Request.Method),
					zap.String("path", c.Request.URL.Path),
					zap.String("client_ip", c.ClientIP()),
					zap.ByteString("stack", debug.Stack()),
				)
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"code": http.StatusInternalServerError,
					"msg":  "服务器内部错误",
				})
			}
		}()
		c.Next()
	}
}
