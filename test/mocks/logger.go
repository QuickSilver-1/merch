package mocks

import (
	"github.com/stretchr/testify/mock"
)

// LoggerRepo - мок-объект для интерфейса LoggerRepo
type LoggerRepo struct {
	mock.Mock
}

func (m *LoggerRepo) Debug(message string) {
	m.Called(message)
}

func (m *LoggerRepo) Info(message string) {
	m.Called(message)
}

func (m *LoggerRepo) Warn(message string) {
	m.Called(message)
}

func (m *LoggerRepo) Error(message string) {
	m.Called(message)
}

func (m *LoggerRepo) Fatal(message string) {
	m.Called(message)
}
