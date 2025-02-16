package interfaces

import "merch/internal/domain"

// AuthRepo предоставляет методы для работы с авторизационными токенами и доступом
type AuthRepo interface {
	// CreateToken создает авторизационный токен
	CreateToken(data domain.AuthorizationData) (token *domain.Token, err error)

	// DecodeToken декодирует авторизационный токен
	DecodeToken(token domain.Token) (data *domain.AuthorizationToken, err error)

	// Access проверяет не был ли токен отозван
	Access(token domain.Token, userId domain.UserId) (exists *domain.SuccessfulAuth, err error)
}
