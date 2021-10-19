package logger

import (
	"encoding/json"
	"fmt"
	"net"
	"net/url"
	"os"

	"github.com/SDJLee/mercedes-benz/util"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// As of now, the logs are written to log files.
// If ELK stack is available, configurations are available in this package to push the logs into ELK stack.

var slogger *zap.SugaredLogger
var queue = make(chan []byte, 10000)
var tag string = "merc-benz-route-checker"

type lumberjackSink struct {
	*lumberjack.Logger
}

func (lumberjackSink) Sync() error {
	return nil
}

func init() {
	env := util.GetEnv()
	devMode := false
	if env == util.EnvDev || env == "" {
		devMode = true
	}
	logFile := fmt.Sprintf("/var/log/%v.log", tag)
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
	var logger *zap.Logger
	var err error
	fmt.Printf("ship logs? '%s' \n", os.Getenv(util.ShipLogs))
	if os.Getenv(util.ShipLogs) == "true" {
		logger, err = cfg.Build(
			zap.Fields(zap.String("tag", tag)),
			// Enable this in an environment where ELK stack is available and respective configurations are added
			zap.Hooks(logstashHook),
		)
		// Enable this in an environment where ELK stack is available and respective configurations are added
		go logstashEmitter()
	} else {
		logger, err = cfg.Build(
			zap.Fields(zap.String("tag", tag)),
		)
	}
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

// lumberjack configuration for log rotation and writing
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

// logstash hook to push logs into logstash
func logstashHook(e zapcore.Entry) error {
	serialized, err := format(&e)
	if err != nil {
		return err
	}
	queue <- serialized
	return nil
}

// emitter transports logs to logstash
func logstashEmitter() {
	conn, err := net.Dial("tcp", "logstash:8089")
	defer func() {
		if conn != nil {
			conn.Close()
		}
	}()
	for msg := range queue {
		if err != nil {
			continue
		}
		_, _ = conn.Write(msg)
	}

}

// formats the logs
func format(e *zapcore.Entry) ([]byte, error) {
	fields := make(map[string]string)

	fields["level"] = e.Level.CapitalString()
	fields["message"] = e.Message
	fields["ex"] = e.Stack
	fields["timestamp"] = e.Time.String()
	fields["caller"] = e.Caller.String()
	fields["tag"] = tag

	serialized, err := json.Marshal(fields)
	if err != nil {
		return nil, fmt.Errorf("couldn't convert log message to json: %s", err)
	}

	// append newline to message so logstash doesn't choke on it
	serialized = append(serialized, "\n"...)
	return serialized, nil
}
