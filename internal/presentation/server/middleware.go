package server

import (
	"fmt"
	"merch/internal/presentation/realization"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// LoggerMiddleware возвращает middleware, который логирует информацию о запросах
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next()

		realization.LoggerService.Info(fmt.Sprintf("Completed %s %s with %d in %v",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(start)))
	}
}

// AuthMiddleware проверяет JWT токен в заголовке авторизации
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Извлекаем токен из заголовка авторизации
		tokenStr := ctx.Request.Header.Get("Authorization")

		parts := strings.Split(tokenStr, " ")
		if len(parts) != 2 {
			ctx.JSON(http.StatusUnauthorized, STATUS_UNAUTHORIZED)
			return
		}

		tokenStr = parts[1]
		token, err := UserService.Token(tokenStr)

		if err != nil {
			answerError(ctx, err)
			return
		}

		if token.Expires.Before(time.Now()) {
			ctx.JSON(http.StatusUnauthorized, fmt.Sprintf("%s - %s", STATUS_UNAUTHORIZED, "token expired"))
			return
		}

		exists, err := UserService.Access(tokenStr, token.Id)
		if err != nil {
			answerError(ctx, err)
			return
		}

		if !*exists {
			ctx.JSON(http.StatusUnauthorized, fmt.Sprintf("%s - %s", STATUS_UNAUTHORIZED, "you are trying to log into someone else's account or the token for this account has been revoked"))
			return
		}

		ctx.Set("token", token)
		ctx.Next()
	}
}
