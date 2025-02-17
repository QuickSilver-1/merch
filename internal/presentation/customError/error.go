package customError

import "merch/internal/domain"

// Алиасы типов для различных типов ошибок
type JWTGenerationError = domain.BaseError

type JWTDecodeError = domain.BaseError

type NeedAuthorization = domain.BaseError

type InvalidPasswordFormat = domain.BaseError

type MigratingError = domain.BaseError

type DbConnectionError = domain.BaseError

type DbQueryError = domain.BaseError

type UserCreatingError = domain.BaseError

type RowsNotFoundError = domain.BaseError

type TransactionError = domain.BaseError

type LoggerBuildError = domain.BaseError
