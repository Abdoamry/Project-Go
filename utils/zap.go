package logger

import (
	"go.uber.org/zap"
)

var Log *zap.Logger

func InitLogger() {
	logger, _ := zap.NewProduction()
	Log = logger
	defer logger.Sync()
	Log.Info("The application is Running")
}