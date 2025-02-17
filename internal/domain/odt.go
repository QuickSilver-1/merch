package domain

import "time"

// AuthorizationData представляет данные для авторизации пользователя
type AuthorizationData struct {
	Id       uint64 `json:"id"`       // Уникальный идентификатор
	Username string `json:"username"` // Имя пользователя
	Password string `json:"password"` // Пароль пользователя
}

// AuthorizationToken представляет токен авторизации пользователя
type AuthorizationToken struct {
	Id      uint64    `json:"id"`      // Уникальный идентификатор токена
	Email   string    `json:"email"`   // Email пользователя, связанный с токеном
	Expires time.Time `json:"expires"` // Дата окончания действия токена
}

// UserInfo содержит информацию о пользователе, включая монеты, инвентарь и транзакции
type UserInfo struct {
	Coins        int           `json:"coins"`       // Количество монет
	Inventory    []Inventory   `json:"inventory"`   // Инвентарь пользователя
	Transactions []Transaction `json:"coinHistory"` // История транзакций пользователя
}

// Алиасы типов для удобства использования
type Token = string

type UserId = uint64

type UserEmail = string

type Msg = string

type Amount = int

type SuccessfulAuth = bool

// BaseError представляет собой базовую ошибку с кодом и сообщением
type BaseError struct {
	Code int    // Код ошибки
	Err  string // Сообщение об ошибке
}

// Error возвращает сообщение об ошибке
func (e *BaseError) Error() string {
	return e.Err
}

// GetCode возвращает код ошибки
func (e *BaseError) GetCode() int {
	return e.Code
}
