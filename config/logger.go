package config

import (
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var zapcoreEncoderConfig zapcore.EncoderConfig

func init() {
	zapcoreEncoderConfig = zapcore.EncoderConfig{
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
		EncodeCaller:   zapcore.FullCallerEncoder,
		EncodeName:     zapcore.FullNameEncoder,
	}
}

func GetZapcoreEncoderConfig() zapcore.EncoderConfig {
	return zapcoreEncoderConfig
}

func GetLogRollingConfig(filePath string) lumberjack.Logger {
	return lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    512, // unit size ：M
		MaxBackups: 12,
		MaxAge:     7,
		Compress:   true,
	}
}
