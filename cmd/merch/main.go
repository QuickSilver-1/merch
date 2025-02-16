package main

import (
	"fmt"
	"merch/internal/presentation/postgres"
	"merch/internal/presentation/realization"
	"merch/internal/presentation/server"
	"merch/internal/services"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	logger, err := realization.NewLogger("../../../log/log.log")
	if err != nil {
		fmt.Println(err)
		return
	}

	err = godotenv.Load("../../../.env")
	if err != nil {
		logger.Error("Failed to load env fail")
		return
	}
	
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")

	db, err := postgres.CreateDB(host, port, user, password, name, logger)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	err = db.CreateSchema()
	if err != nil {
		logger.Error(err.Error())
		return
	}

	userRepo := realization.NewUser()
	transactionRepo := realization.NewTransaction()
	inventoryRepo := realization.NewInventory()
	moneyService := services.NewMoneyService(userRepo, transactionRepo, inventoryRepo)

	secretKey := os.Getenv("SECRET_KEY")
	authRepo := realization.NewAuth(secretKey)
	userService := services.NewUserService(userRepo, authRepo, transactionRepo, inventoryRepo)

	serverPort := os.Getenv("SERVER_PORT")
	srv := server.NewServer()
	err = srv.Start(moneyService, userService, secretKey, serverPort)
	if err != nil {
		logger.Error(fmt.Sprintf("Critical server error: %v", err))
	}

	err = srv.Shutdown()
	if err != nil {
		logger.Error(fmt.Sprintf("Shutdown server error: %v", err))
	}

	logger.Info("Graceful shutdown completed successfully")
}
