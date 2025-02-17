package mocks

import (
	"merch/internal/domain"

	"github.com/stretchr/testify/mock"
)

// MockInventoryRepo - мок-объект для интерфейса InventoryRepo
type MockInventoryRepo struct {
	mock.Mock
}

func (m *MockInventoryRepo) GetSubjectByName(name string) (*domain.Item, error) {
	args := m.Called(name)
	return args.Get(0).(*domain.Item), args.Error(1)
}

func (m *MockInventoryRepo) Buy(userId domain.UserId, subject domain.Item) error {
	args := m.Called(userId, subject)
	return args.Error(0)
}

func (m *MockInventoryRepo) GetInventory(userId domain.UserId) (*[]domain.Inventory, error) {
	args := m.Called(userId)
	return args.Get(0).(*[]domain.Inventory), args.Error(1)
}
