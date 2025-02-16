package server

import (
	"merch/internal/presentation/postgres"
	"merch/internal/presentation/realization"
	"merch/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Переменные для доступа к сервисам
var (
	BuyService  *services.MoneyService
	UserService *services.UserService
	SecretKey   string
)

// Константы http ответов
const (
	STATUS_UNAUTHORIZED    = "Authorization required"
	STATUS_INTERNAL_SERVER = "Sorry, something went wrong, we are already solving the problem"
	STATUS_BAD_REQUEST     = "Invalid data"
)

// Server определяет сервер с сервисами
type Server struct {
	srv *gin.Engine
}

// NewServer создает новый экземпляр Server
func NewServer() *Server {
	gin.SetMode(gin.ReleaseMode)
	srv := gin.New()

	h := NewHandlers()

	srv.Use(LoggerMiddleware())
	srv.GET("/api/auth", h.Auth)

	srv.Use(AuthMiddleware())
	srv.POST("/api/info", h.GetInfo)
	srv.POST("/api/sendCoin", h.SendCoin)
	srv.GET("/api/buy/:item", h.BuyMerch)

	// Обработка несуществующих маршрутов
	srv.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusBadRequest, "Invalid route")
	})

	realization.LoggerService.Info("Server has been created")
	return &Server{
		srv: srv,
	}
}

// Start запускает сервер
func (s *Server) Start(money *services.MoneyService, user *services.UserService, secret, port string) error {
	BuyService = money
	UserService = user
	SecretKey = secret

	realization.LoggerService.Debug("Starting server")
	err := s.srv.Run(":" + port)

	if err != nil {
		return err
	}

	realization.LoggerService.Info("Stopping server")
	return nil
}

// Shutdown завершает работу сервера и закрывает подключение к базе данных
func (s *Server) Shutdown() error {
	return postgres.DbService.CloseDB()
}
