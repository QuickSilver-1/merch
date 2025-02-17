package realization

import (
	"fmt"
	"merch/internal/interfaces"
	e "merch/internal/presentation/customError"
	"net/http"

	"go.uber.org/zap"
)

var (
	LoggerService interfaces.LoggerRepo
)

// Logger - структура для логгирования с использованием zap
type Logger struct {
	logger *zap.Logger
}

// NewLogger создает новый экземпляр логгера с заданным путем для вывода логов
func NewLogger() (interfaces.LoggerRepo, error) {
	config := zap.NewDevelopmentConfig()
	config.OutputPaths = []string{"stdout"}

	logger, err := config.Build()
	if err != nil {
		return nil, &e.LoggerBuildError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Sprintf("Failed to configure logger: %v", err),
		}
	}

	LoggerService = &Logger{
		logger: logger,
	}

	return LoggerService, nil
}

func (l *Logger) Fatal(msg string) {
	l.logger.Fatal(msg)
}

func (l *Logger) Error(msg string) {
	l.logger.Error(msg)
}

func (l *Logger) Warn(msg string) {
	l.logger.Warn(msg)
}

func (l *Logger) Info(msg string) {
	l.logger.Info(msg)
}

func (l *Logger) Debug(msg string) {
	l.logger.Debug(msg)
}
