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
		return createTransactionError(err)
	}

	if err := s.updateUserCoins(tx, transaction.ReceiverName, transaction.Amount); err != nil {
		return s.rollbackTransaction(tx, err)
	}

	if err := s.updateUserCoins(tx, transaction.SenderName, -transaction.Amount); err != nil {
		return s.rollbackTransaction(tx, err)
	}

	if err := s.insertTransaction(tx, transaction); err != nil {
		return s.rollbackTransaction(tx, err)
	}

	if err := tx.Commit(); err != nil {
		return createCommitError(err)
	}

	return nil
}

func (s *Transaction) updateUserCoins(tx *sql.Tx, email string, amount int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err := tx.ExecContext(ctx, `UPDATE Users SET "coins" = "coins" + $1 WHERE "email" = $2`, amount, email)
	if err != nil {
		return createDbQueryError(err)
	}
	return nil
}

func (s *Transaction) insertTransaction(tx *sql.Tx, transaction domain.Transaction) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	_, err := tx.ExecContext(ctx, `INSERT INTO Transaction("sender_name", "receiver_name", "amount") VALUES ($1, $2, $3)`, transaction.SenderName, transaction.ReceiverName, transaction.Amount)
	if err != nil {
		return createDbQueryError(err)
	}
	return nil
}

func (s *Transaction) rollbackTransaction(tx *sql.Tx, originalErr error) error {
	if rbErr := tx.Rollback(); rbErr != nil {
		LoggerService.Error(fmt.Sprintf("Rollback error: %v", rbErr))
	}
	return originalErr
}

func createTransactionError(err error) error {
	return &e.TransactionError{
		Code: http.StatusInternalServerError,
		Err:  fmt.Sprintf("Creating transaction error: %v", err),
	}
}

func createDbQueryError(err error) error {
	return &e.DbQueryError{
		Code: http.StatusInternalServerError,
		Err:  fmt.Sprintf("Db query error: %v", err),
	}
}

func createCommitError(err error) error {
	return &e.TransactionError{
		Code: http.StatusInternalServerError,
		Err:  fmt.Sprintf("Commit error: %v", err),
	}
}

// GetTransaction получает транзакции для указанного пользователя по email
func (s *Transaction) GetTransaction(user domain.UserEmail) (*[]domain.Transaction, error) {
	LoggerService.Debug("Getting transaction")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	rows, err := postgres.DbService.Db.QueryContext(ctx, `SELECT "id", "sender_name", "receiver_name", "amount" FROM Transaction WHERE "sender_name" = $1 OR "receiver_name" = $1`, user)
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
