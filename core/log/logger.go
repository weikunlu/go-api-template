package log

import (
	"github.com/weikunlu/go-api-template/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"os/signal"
	"syscall"
)

func NewLogger(level zapcore.Level, messageType string, enableCaller bool) *zap.Logger {
	appConfig := config.GetAppConfig()

	zapCore := newZapCore(level, appConfig.AppLogPath)

	zapFields := zap.Fields(zap.String("msg_type", messageType), zap.Int("pid", os.Getpid()))

	if !enableCaller {
		return zap.New(zapCore, zapFields)
	}

	return zap.New(zapCore, zap.AddCaller(), zap.Development(), zapFields)
}

func newZapCore(level zapcore.Level, filePath string) zapcore.Core {
	hook := config.GetLogRollingConfig(filePath)

	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)

	cfg := config.GetAppConfig()
	zapcoreCfg := config.GetZapcoreEncoderConfig()

	// execute kill -SIGHUP <process_id> to rotate
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP)
	go func() {
		for {
			<-c
			hook.Rotate()
		}
	}()

	writers := []zapcore.WriteSyncer{zapcore.AddSync(&hook)}
	if cfg.AppEnv == "local" {
		writers = append(writers, zapcore.AddSync(os.Stdout))

	}

	return zapcore.NewCore(
		zapcore.NewJSONEncoder(zapcoreCfg),
		zapcore.NewMultiWriteSyncer(writers...),
		atomicLevel,
	)
}

