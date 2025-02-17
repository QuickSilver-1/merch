package tests

import (
	"bytes"
	"encoding/json"
	"log"
	"merch/internal/domain"
	"merch/internal/presentation/postgres"
	"merch/internal/presentation/realization"
	"merch/internal/presentation/server"
	"merch/internal/services"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var (
	HOST   string
	PORT   string
	USER   string
	PASS   string
	NAME   string
	SECRET string
	Srv    *server.Server
)

func init() {
	err := godotenv.Load("../../.env")
	if err != nil {
		panic(err)
	}

	HOST = os.Getenv("DB_HOST")
	PORT = os.Getenv("DB_PORT")
	USER = os.Getenv("DB_USER")
	PASS = os.Getenv("DB_PASSWORD")
	NAME = os.Getenv("DB_NAME")
	SECRET = os.Getenv("SECRET_KEY")

	Srv = StartTestServer()
}

func StartTestServer() *server.Server {
	// Настройка логгера
	logger, err := realization.NewLogger()
	if err != nil {
		log.Fatalf("Could not create logger: %v", err)
	}

	// Настройка базы данных
	_, err = postgres.CreateDB(HOST, PORT, USER, PASS, NAME, logger)
	if err != nil {
		log.Fatalf("Could not create database connection: %v", err)
	}

	// Настройка сервисов
	userRepo := realization.NewUser()
	authRepo := realization.NewAuth(SECRET)
	transactionRepo := realization.NewTransaction()
	inventoryRepo := realization.NewInventory()

	userService := services.NewUserService(userRepo, authRepo, transactionRepo, inventoryRepo)
	moneyService := services.NewMoneyService(userRepo, transactionRepo, inventoryRepo)

	// Запуск сервера
	srv := server.NewServer()
	go func() {
		if err := srv.Start(moneyService, userService, SECRET, "8080"); err != nil {
			log.Fatalf("Could not start server: %v", err)
		}
	}()

	// Ждем запуска сервера
	time.Sleep(3 * time.Second)

	return srv
}

func TestCreateUserAndLogin(t *testing.T) {
	// Тест создания нового пользователя
	newUserData := domain.AuthorizationData{
		Username: "newuser@example.com",
		Password: "Password123",
	}

	payload, err := json.Marshal(newUserData)
	assert.NoError(t, err)

	resp, err := http.Post("http://localhost:8080/api/auth", "application/json", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			return
		}
	}()
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	assert.NotEmpty(t, result["token"])

	// Тест входа в существующего пользователя
	existingUserData := domain.AuthorizationData{
		Username: "newuser@example.com",
		Password: "Password123",
	}

	payload, err = json.Marshal(existingUserData)
	assert.NoError(t, err)

	resp, err = http.Post("http://localhost:8080/api/auth", "application/json", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			return
		}
	}()

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var loginResult map[string]string
	err = json.NewDecoder(resp.Body).Decode(&loginResult)
	assert.NoError(t, err)
	assert.NotEmpty(t, loginResult["token"])
}
