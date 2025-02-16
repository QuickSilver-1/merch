package postgres

import (
	"fmt"
	e "merch/internal/presentation/customError"
	"net/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// CreateSchema выполняет миграции базы данных для создания схемы
func (db *DB) CreateSchema() error {
	db.Logger.Debug("Migrating...")

	// Создаем экземпляр драйвера для PostgreSQL
	driver, err := postgres.WithInstance(db.Db, &postgres.Config{})
	if err != nil {
		db.Logger.Error(fmt.Sprintf("Creating driver PostgreSQL fatal error: %v", err))
		return &e.MigratingError{
			Code: http.StatusInternalServerError,
			Err:  "Creating driver PostgreSQL error",
		}
	}

	// Создаем мигратор с указанным источником миграций и базой данных
	m, err := migrate.NewWithDatabaseInstance("file:/../../../internal/presentation/migrations", "postgres", driver)
	if err != nil {
		db.Logger.Error(fmt.Sprintf("Creating migrator error: %v", err))
		return &e.MigratingError{
			Code: http.StatusInternalServerError,
			Err:  "Creating migrator error",
		}
	}

	// Применяем миграции к базе данных
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		db.Logger.Error(fmt.Sprintf("Error applying migrations: %v", err))
		return &e.MigratingError{
			Code: http.StatusInternalServerError,
			Err:  "Error applying migrations",
		}
	}

	db.Logger.Debug("Migrations successfully applied!")
	return nil
}
