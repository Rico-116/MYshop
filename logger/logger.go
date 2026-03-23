package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger
var Sugar *zap.SugaredLogger

func Init() error {
	var err error

	// 生产环境风格日志
	Log, err = zap.NewProduction()
	if err != nil {
		return err
	}

	Sugar = Log.Sugar()
	return nil
}

func Sync() {
	if Log != nil {
		_ = Log.Sync()
	}
}
