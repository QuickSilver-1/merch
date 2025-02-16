package postgres

import (
	"database/sql"
	"fmt"
	"merch/internal/interfaces"
	e "merch/internal/presentation/customError"
	"net/http"
)

var (
	DbService *DB
)

// DB - структура для работы с базой данных
type DB struct {
	Db     *sql.DB
	Logger interfaces.LoggerRepo
}

// CreateDB создает подключение к базе данных и возвращает экземпляр DB
func CreateDB(ip, port, user, pass, nameDB string, logger interfaces.LoggerRepo) (*DB, error) {
	logger.Debug("Database connection creating...")
	sqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", ip, port, user, pass, nameDB)
	conn, err := sql.Open("postgres", sqlInfo)

	if err != nil {
		return nil, &e.DbConnectionError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Sprintf("Database connection error: %v", err),
		}
	}

	DbService = &DB{
		Db:     conn,
		Logger: logger,
	}

	logger.Info("Database connection has been created")
	return DbService, nil
}

// CloseDB закрывает подключение к базе данных
func (db *DB) CloseDB() error {
	db.Logger.Debug("Closing database connection")
	return db.Db.Close()
}
