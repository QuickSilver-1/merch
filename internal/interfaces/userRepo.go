package interfaces

import "merch/internal/domain"

// UserRepo предоставляет методы для работы с пользователями
type UserRepo interface {
	// Create создает нового пользователя и возвращает его идентификатор
	Create(user domain.User) (id *domain.UserId, err error)

	// GetById получает пользователя по его идентификатору
	GetById(id domain.UserId) (user *domain.User, err error)

	// GetByEmail получает пользователя по его email
	GetByEmail(email domain.UserEmail) (user *domain.User, err error)
}
