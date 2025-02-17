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
			ctx.JSON(http.StatusUnauthorized, map[string]string{"errors": STATUS_UNAUTHORIZED})
			ctx.Abort()
			return
		}

		token, err := UserService.Token(parts[1])

		if err != nil {
			answerError(ctx, err)
			return
		}

		if time.Now().Before(token.Expires) {
			ctx.JSON(http.StatusUnauthorized, map[string]string{"errors": fmt.Sprintf("%s - %s", STATUS_UNAUTHORIZED, "token expired")})
			ctx.Abort()
			return
		}

		exists, err := UserService.Access(parts[1], token.Id)
		if err != nil {
			answerError(ctx, err)
			return
		}

		if !*exists {
			ctx.JSON(http.StatusUnauthorized, map[string]string{"errors": fmt.Sprintf("%s - %s", STATUS_UNAUTHORIZED, "you are trying to log into someone else's account or the token for this account has been revoked")})
			ctx.Abort()
			return
		}

		ctx.Set("token", parts[1])
		ctx.Next()
	}
}
