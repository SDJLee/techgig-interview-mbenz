package logger

import (
	"fmt"
	"net/url"

	"github.com/SDJLee/mercedes-benz/util"
	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// As of now, the logs are written to log files.
// Future plan for this package is to push the logs into ELK stack.

var slogger *zap.SugaredLogger

type lumberjackSink struct {
	*lumberjack.Logger
}

func (lumberjackSink) Sync() error {
	return nil
}

func init() {
	logFile := viper.GetString(util.LogPath)
	// logFile := "./test.log"
	env := viper.GetString(util.AppEnv)
	devMode := false
	if env == "dev" || env == "" {
		devMode = true
	}
	logWriter := getLogWriter(logFile)
	encoderConfig := getEncoderConfig(devMode)
	zap.RegisterSink("lumberjack", func(*url.URL) (zap.Sink, error) {
		return lumberjackSink{
			Logger: logWriter,
		}, nil
	})
	cfg := zap.Config{
		Encoding:    "json",
		Level:       zap.NewAtomicLevelAt(zapcore.DebugLevel),
		OutputPaths: []string{fmt.Sprintf("lumberjack:%s", logFile)}, // add logger path
		// OutputPaths:   []string{"stdout", fmt.Sprintf("lumberjack:%s", logFile)}, // add logger path
		Development:   devMode,
		EncoderConfig: encoderConfig,
	}
	cfg.EncoderConfig = encoderConfig

	logger, err := cfg.Build(
		zap.Fields(zap.String("tag", "merc-benz-route-checker")),
	)
	if err != nil {
		fmt.Println("error creating shared logger", err)
		return
	}

	zap.RedirectStdLog(logger)
	slogger = logger.Sugar()
	slogger.Info("merc-benz-route-checker logger initialized")
}

func getEncoderConfig(devMode bool) zapcore.EncoderConfig {
	var encoderConfig zapcore.EncoderConfig
	if devMode {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderConfig = zap.NewProductionEncoderConfig()
	}
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.MessageKey = "message"
	encoderConfig.LevelKey = "level"
	encoderConfig.TimeKey = "time"
	encoderConfig.CallerKey = "caller"
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return encoderConfig
}

func getLogWriter(logFile string) *lumberjack.Logger {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    1,
		MaxBackups: 5,
		MaxAge:     30,
		Compress:   false,
	}
	return lumberJackLogger
}

func Logger() *zap.SugaredLogger {
	return slogger
}

func SubLogger(name string) *zap.SugaredLogger {
	return slogger.Named(name)
}
