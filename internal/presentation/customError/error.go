package customError

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

// Алиасы типов для различных типов ошибок
type JWTGenerationError = BaseError

type JWTDecodeError = BaseError

type NeedAuthorization = BaseError

type InvalidPasswordFormat = BaseError

type MigratingError = BaseError

type DbConnectionError = BaseError

type DbQueryError = BaseError

type UserCreatingError = BaseError

type RowsNotFoundError = BaseError

type TransactionError = BaseError

type LoggerBuildError = BaseError
