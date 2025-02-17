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

// Transaction - структура для работы с транзакциями
type Transaction struct{}

// NewTransaction создает новый экземпляр Transaction
func NewTransaction() *Transaction {
	return &Transaction{}
}

// Transfer выполняет перевод средств
func (s *Transaction) Transfer(transaction domain.Transaction) error {
	LoggerService.Debug("Transferring")
	tx, err := postgres.DbService.Db.Begin()

	if err != nil {
		return &e.TransactionError{
			Code: http.StatusInternalServerError,
			Err:  fmt.Sprintf("Creating transaction error: %v", err),
		}
	}

	res, err := tx.Exec(`UPDATE Users SET "coins" = "coins" + $1 WHERE "email" = $2`, transaction.Amount, transaction.ReceiverName)
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

	_, err = tx.Exec(`UPDATE Users SET "coins" = "coins" - $1 WHERE "email" = $2`, transaction.Amount, transaction.SenderName)
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

	amount, err := res.RowsAffected()

	if err != nil {
		err = tx.Rollback()
		if err != nil {
			LoggerService.Error(fmt.Sprintf("Rollback error: %v", err))
		}
		return &e.DbQueryError{
			Code: http.StatusInternalServerError,
			Err: fmt.Sprintf("Db query error: %v", err),
		}
	}

	if amount == 0 {
		err = tx.Rollback()
		if err != nil {
			LoggerService.Error(fmt.Sprintf("Rollback error: %v", err))
		}
		return &e.TransactionError{
			Code: http.StatusBadRequest,
			Err: "Invalid user's email",
		}
	}

	_, err = tx.Exec(`INSERT INTO Transaction("sender_name", "receiver_name", "amount") VALUES ($1, $2, $3)`, transaction.SenderName, transaction.ReceiverName, transaction.Amount)
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

// GetTransaction получает транзакции для указанного пользователя по email
func (s *Transaction) GetTransaction(user domain.UserEmail) (*[]domain.Transaction, error) {
	LoggerService.Debug("Getting transaction")
	rows, err := postgres.DbService.Db.Query(`SELECT "id", "sender_name", "receiver_name", "amount" FROM Transaction WHERE "sender_name" = $1 OR "receiver_name" = $1`, user)
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

	var transactions []domain.Transaction
	for rows.Next() {
		var transaction domain.Transaction
		err = rows.Scan(&transaction.Id, &transaction.SenderName, &transaction.ReceiverName, &transaction.Amount)
		if err != nil {
			return nil, &e.DbQueryError{
				Code: http.StatusInternalServerError,
				Err:  fmt.Sprintf("Scanning subject error: %v", err),
			}
		}
		transactions = append(transactions, transaction)
	}

	return &transactions, nil
}
