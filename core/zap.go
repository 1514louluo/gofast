package core

import (
	"GF_PROJECT_NAME/global"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Zap() *zap.SugaredLogger {
//	logPath := global.
	hook := lumberjack.Logger {
		Filename: global.CONFIG.Zap.File,
		MaxSize: global.CONFIG.Zap.MaxSize, // megabytes
		MaxBackups: global.CONFIG.Zap.MaxBackups, // 最多保留300个备份
		MaxAge: global.CONFIG.Zap.MaxAge,  // days
		Compress: global.CONFIG.Zap.Compress,
	}
	syncWriter := zapcore.AddSync(&hook)
	var level zapcore.Level
	switch global.CONFIG.Zap.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}
	encoderConfig := zap.NewProductionEncoderConfig()
	//时间格式
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		syncWriter,
		level,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	return logger.Sugar()
}
