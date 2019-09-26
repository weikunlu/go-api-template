package log

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strconv"
)

const loggerKey = iota

var logger *zap.Logger

func init() {
	logger = NewLogger(zapcore.InfoLevel, "service-log", true)
}

func NewContext(ctx *gin.Context, fields ...zapcore.Field) {
	ctx.Set(strconv.Itoa(loggerKey), WithContext(ctx).With(fields...))
}

func WithContext(ctx *gin.Context) *zap.Logger {
	if ctx == nil {
		return logger
	}
	l, _ := ctx.Get(strconv.Itoa(loggerKey))
	ctxLogger, ok := l.(*zap.Logger)
	if ok {
		return ctxLogger
	}
	return logger
}

func GetErrorMessageField(message string) zap.Field {
	return GetStringField("error_message", message)
}

func GetStringField(key string, value string) zap.Field {
	return zap.String(key, value)
}
