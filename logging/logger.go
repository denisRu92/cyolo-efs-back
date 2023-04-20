package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerContextKey string

const (
	LogBaseFieldsKey = "Log-base-fields"
	logLevelEnvVar   = "LOG_LEVEL"
	appIdEnvVar      = "APP_ID"
	serviceKey       = "service"
)

var (
	Log           *zap.SugaredLogger
	FastLog       *zap.Logger
	atom          zap.AtomicLevel
	baseFieldsKey = loggerContextKey(LogBaseFieldsKey)
)

// init the loggers and defaults. LOG_LEVEL (default: info) and APP_ID (default: "") environment variables should be set.
func init() {
	atom = zap.NewAtomicLevel()
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	logger := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		atom,
	))

	atom.SetLevel(zap.InfoLevel)
	logLevel := os.Getenv(logLevelEnvVar)
	if logLevel != "" {
		err := atom.UnmarshalText([]byte(strings.ToLower(logLevel)))
		if err != nil {
			logger.Fatal("Invalid Log level")
		}
	}

	appId := os.Getenv(appIdEnvVar)
	Log = logger.Sugar().With(serviceKey, appId)
	FastLog = logger.With(zap.String(serviceKey, appId))
}
