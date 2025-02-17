package mocks

import (
	"merch/internal/domain"

	"github.com/stretchr/testify/mock"
)

// MockTransactionRepo - мок-объект для интерфейса TransactionRepo
type MockTransactionRepo struct {
	mock.Mock
}

func (m *MockTransactionRepo) Transfer(transaction domain.Transaction) error {
	args := m.Called(transaction)
	return args.Error(0)
}

func (m *MockTransactionRepo) GetTransaction(email domain.UserEmail) (*[]domain.Transaction, error) {
	args := m.Called(email)
	return args.Get(0).(*[]domain.Transaction), args.Error(1)
}
