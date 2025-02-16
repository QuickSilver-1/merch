package interfaces

// ServerRepo предоставляет методы для управления сервером
type ServerRepo interface {
	// Start запускает сервер
	Start() error

	// Shutdown останавливает сервер по паттерну graceful shutdown
	Shutdown() error
}
