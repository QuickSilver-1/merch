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

	"github.com/dgrijalva/jwt-go"
)

// Token - структура для JWT токена
type Token struct {
	Id    uint64 `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

// Auth - структуру для работы с авторизацией
type Auth struct {
	secretKey string
}

// NewAuth создает новый экземпляр Auth с заданным секретным ключом
func NewAuth(secret string) *Auth {
	LoggerService.Debug("Creating auth service")
	return &Auth{
		secretKey: secret,
	}
}

// CreateToken создает JWT токен для данного email
func (s *Auth) CreateToken(data domain.AuthorizationData) (*domain.Token, error) {
	LoggerService.Debug("Creating JWT")
	expires := time.Now().Add(time.Hour * 24 * 30) // Срок действия токена - 30 дней
	claims := &Token{
		Id:    data.Id,
		Email: data.Username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expires.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenS, err := token.SignedString([]byte(s.secretKey))

	if err != nil {
		return nil, &e.JWTGenerationError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Sprintf("Token generation error: %v", err),
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err = postgres.DbService.Db.ExecContext(ctx, ` DELETE FROM Token WHERE "user_id" = $1 `, data.Id)

	if err != nil {
		LoggerService.Error("Deleting token error")
	}

	ctx, cancel = context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err = postgres.DbService.Db.ExecContext(ctx, `INSERT INTO Token("user_id", "value") VALUES ($1, $2)`, data.Id, tokenS)

	if err != nil {
		return nil, &e.DbQueryError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Sprintf("Database query error: %v", err),
		}
	}

	return &tokenS, nil
}

// DecodeToken декодирует JWT токен
func (s *Auth) DecodeToken(tokenStr domain.Token) (*domain.AuthorizationToken, error) {
	LoggerService.Debug("Decoding JWT")
	claims := &Token{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, &e.NeedAuthorization{
			Code: http.StatusUnauthorized,
			Err:  "Bad token's signature",
		}
	}

	if !token.Valid {
		return nil, &e.NeedAuthorization{
			Code: http.StatusUnauthorized,
			Err:  "Invalid token",
		}
	}

	return &domain.AuthorizationToken{
		Id:    claims.Id,
		Email: claims.Email,
	}, nil
}

// Access проверяет доступ по токену для указанного пользователя
func (s *Auth) Access(token domain.Token, userId domain.UserId) (*domain.SuccessfulAuth, error) {
	LoggerService.Debug("Checking user's access")
	var notExists = false
	var exists = true

	var id uint64
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	err := postgres.DbService.Db.QueryRowContext(ctx, `SELECT "user_id" FROM Token WHERE "value" = $1`, token).Scan(&id)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &notExists, nil
		}

		return nil, &e.DbQueryError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Sprintf("Database query error: %v", err),
		}
	}

	if id != userId {
		return &notExists, nil
	}

	return &exists, nil
}
