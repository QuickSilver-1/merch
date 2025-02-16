package realization

import (
	"database/sql"
	"errors"
	"fmt"
	"merch/internal/domain"
	e "merch/internal/presentation/customError"
	"merch/internal/presentation/postgres"
	"net/http"
)

// Inventory представляет собой структуру для работы с инвентарем
type Inventory struct{}

// NewInventory создает новый экземпляр Inventory
func NewInventory() *Inventory {
	return &Inventory{}
}

// Buy выполняет покупку предмета
func (s *Inventory) Buy(userId domain.UserId, subject domain.Item) error {
	LoggerService.Debug("Buying subject")
	tx, err := postgres.DbService.Db.Begin()

	if err != nil {
		return &e.TransactionError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Sprintf("Creating transaction error: %v", err),
		}
	}

	_, err = tx.Exec(`UPDATE Users SET "coins" = "coins" - $1 WHERE "id" = $2`, subject.Cost, userId)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			LoggerService.Error(fmt.Sprintf("Rollback error: %v", err))
		}
		return &e.DbQueryError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Sprintf("Db query error: %v", err),
		}
	}

	_, err = tx.Exec(`INSERT INTO Inventory ("subject_name", "user_id") VALUES ($1, $2)`, subject.Name, userId)
	if err != nil {
		err = tx.Rollback()
		if err != nil {
			LoggerService.Error(fmt.Sprintf("Rollback error: %v", err))
		}
		return &e.DbQueryError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Sprintf("Db query error: %v", err),
		}
	}

	err = tx.Commit()
	if err != nil {
		return &e.TransactionError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Sprintf("Commit error: %v", err),
		}
	}

	return nil
}

// GetSubjectByName получает предмет по его названию
func (s *Inventory) GetSubjectByName(name domain.UserEmail) (*domain.Item, error) {
	LoggerService.Debug("Getting subject")
	var subject domain.Item
	err := postgres.DbService.Db.QueryRow(`SELECT "id", "name", "cost" FROM Subject WHERE "name" = $1`, name).Scan(&subject.Id, &subject.Name, &subject.Cost)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, &e.RowsNotFoundError{
				Code: http.StatusBadRequest,
				Err:  "Subject not exists",
			}
		}

		return nil, &e.DbQueryError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Sprintf("Query db error: %v", err),
		}
	}

	return &subject, nil
}

// GetInventory получает инвентарь пользователя по его идентификатору
func (s *Inventory) GetInventory(userId domain.UserId) (*[]domain.Inventory, error) {
	rows, err := postgres.DbService.Db.Query(`SELECT "id", "subject_name", "user_id" FROM Inventory WHERE "user_id" = $1`, userId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, &e.DbQueryError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Sprintf("Query db error: %v", err),
		}
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			LoggerService.Error(fmt.Sprintf("Error closing rows: %v", err))
		}
	}()

	var inventory []domain.Inventory
	for rows.Next() {
		var subject domain.Inventory
		err = rows.Scan(&subject.Id, &subject.Subject, &subject.UserId)
		if err != nil {
			return nil, &e.DbQueryError{
				Code: http.StatusInternalServerError,
				Err:  fmt.Sprintf("Scanning subject error: %v", err),
			}
		}
		inventory = append(inventory, subject)
	}

	return &inventory, nil
}
