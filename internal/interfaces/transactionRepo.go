package interfaces

import "merch/internal/domain"

// TransactionRepo предоставляет методы для работы с транзакциями
type TransactionRepo interface {
	// Transfer выполняет перевод средств в рамках транзакции
	Transfer(transaction domain.Transaction) error

	// GetTransaction получает транзакции для указанного пользователя по email
	GetTransaction(email domain.UserEmail) (transactions *[]domain.Transaction, err error)
}
