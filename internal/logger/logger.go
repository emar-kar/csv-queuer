package logger

import (
	"os"
	"path"
	"sync"

	"github.com/kyokomi/emoji"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var once sync.Once

type Logger struct {
	Info  *zap.Logger
	Error *zap.Logger
}

var logger *Logger

// InitLogger creates *Logger singleton object within the given folder path.
// *Logger is a wrapper around two *zap.Loggers. It provides separated objects
// to log info and errors separatly.
func InitLogger(logFolder string) {
	defer func() {
		if r := recover(); r != nil {
			emoji.Printf("cannot initialize logger :see-no-evil_monkey:: %s\n", r)
		}
	}()

	if logger == nil {
		createLogger(logFolder)
	}
}

func createLogger(logFolder string) {
	once.Do(func() {
		if _, err := os.Stat(logFolder); err == nil || os.IsNotExist(err) {
			os.RemoveAll(logFolder)
		}
		if err := os.MkdirAll(logFolder, 0777); err != nil {
			panic(err)
		}

		infoLog, err := buildLogger(logFolder, "access.log")
		if err != nil {
			panic(err)
		}
		errorsLog, err := buildLogger(logFolder, "errors.log")
		if err != nil {
			panic(err)
		}

		logger = &Logger{Info: infoLog, Error: errorsLog}
	})
}

// GetLogger returns logger object.
func GetLogger() *Logger {
	return logger
}

func buildLogger(logFolder, logFile string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logFilePath := path.Join(logFolder, logFile)

	if _, err := os.Create(logFilePath); err != nil {
		return nil, err
	}

	cfg.OutputPaths = []string{
		logFilePath,
	}

	return cfg.Build()
}
