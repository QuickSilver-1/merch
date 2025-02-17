package realization

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"merch/internal/domain"
	e "merch/internal/presentation/customError"
	"merch/internal/presentation/postgres"
	"net/http"
	"time"
)

// User - структуру для работы с пользователями
type User struct{}

// NewUser создает новый экземпляр User
func NewUser() *User {
	return &User{}
}

// Create создает нового пользователя
func (s *User) Create(user domain.User) (*domain.UserId, error) {
	LoggerService.Debug("Creating user")
	var id uint64
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err := postgres.DbService.Db.QueryRowContext(ctx, `INSERT INTO Users ("email", "password", "coins") VALUES ($1, $2, $3) RETURNING "id"`, user.Email, user.Password, user.Coins).Scan(&id)

	if err != nil {
		return nil, &e.UserCreatingError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Sprintf("User creating error: %v", err),
		}
	}

	return &id, nil
}

// GetById получает пользователя по его идентификатору
func (s *User) GetById(id domain.UserId) (*domain.User, error) {
	LoggerService.Debug("Getting user by id")
	var user domain.User
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err := postgres.DbService.Db.QueryRowContext(ctx, `SELECT "id", "email", "password", "coins" FROM Users WHERE "id" = $1`, id).Scan(&user.Id, &user.Email, &user.Password, &user.Coins)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &e.RowsNotFoundError{
				Code: http.StatusBadRequest,
				Err:  "User not exists",
			}
		}

		return nil, &e.DbQueryError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Sprintf("Query db error: %v", err),
		}
	}

	return &user, nil
}

// GetByEmail получает пользователя по его email
func (s *User) GetByEmail(email domain.UserEmail) (*domain.User, error) {
	LoggerService.Info("Getting user by email")
	var user domain.User
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err := postgres.DbService.Db.QueryRowContext(ctx, `SELECT "id", "email", "password", "coins" FROM Users WHERE "email" = $1`, email).Scan(&user.Id, &user.Email, &user.Password, &user.Coins)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, &e.DbQueryError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Sprintf("Query db error: %v", err),
		}
	}

	return &user, nil
}
