package mocks

import (
	"merch/internal/domain"

	"github.com/stretchr/testify/mock"
)

// MockUserRepo - мок-объект для интерфейса UserRepo
type MockUserRepo struct {
    mock.Mock
}

func (m *MockUserRepo) Create(user domain.User) (*domain.UserId, error) {
    args := m.Called(user)
    return args.Get(0).(*domain.UserId), args.Error(1)
}

func (m *MockUserRepo) GetById(id domain.UserId) (*domain.User, error) {
    args := m.Called(id)
    return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepo) GetByEmail(email domain.UserEmail) (*domain.User, error) {
    args := m.Called(email)
    return args.Get(0).(*domain.User), args.Error(1)
}