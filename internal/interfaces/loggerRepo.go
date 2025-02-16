package interfaces

import "merch/internal/domain"

// Интерфейс LoggerRepo определяет методы для логирования сообщений
type LoggerRepo interface {
	Fatal(msg domain.Msg)
	Error(msg domain.Msg)
	Warn(msg domain.Msg)
	Info(msg domain.Msg)
	Debug(msg domain.Msg)
}
