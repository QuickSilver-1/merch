package mocks

import (
	"merch/internal/domain"

	"github.com/stretchr/testify/mock"
)

// MockAuthRepo - мок-объект для интерфейса AuthRepo
type MockAuthRepo struct {
	mock.Mock
}

func (m *MockAuthRepo) CreateToken(data domain.AuthorizationData) (*domain.Token, error) {
	args := m.Called(data)
	return args.Get(0).(*domain.Token), args.Error(1)
}

func (m *MockAuthRepo) DecodeToken(token domain.Token) (*domain.AuthorizationToken, error) {
	args := m.Called(token)
	return args.Get(0).(*domain.AuthorizationToken), args.Error(1)
}

func (m *MockAuthRepo) Access(token domain.Token, userId domain.UserId) (*domain.SuccessfulAuth, error) {
	args := m.Called(token, userId)
	return args.Get(0).(*domain.SuccessfulAuth), args.Error(1)
}
