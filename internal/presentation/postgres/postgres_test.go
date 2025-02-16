package postgres

import (
	"merch/test/mocks"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateDB(t *testing.T) {
    mockLogger := new(mocks.LoggerRepo)
    mockLogger.On("Debug", "Database connection creating...").Once()
    mockLogger.On("Info", "Database connection has been created").Once()

    db, err := CreateDB("localhost", "5432", "user", "password", "dbname", mockLogger)
    assert.NoError(t, err)
    assert.NotNil(t, db)

    mockLogger.AssertExpectations(t)
}

func TestCloseDB_Success(t *testing.T) {
    mockLogger := new(mocks.LoggerRepo)
    mockLogger.On("Debug", "Closing database connection").Once()

    mockDB, sqlMock, err := sqlmock.New()
    assert.NoError(t, err)
    defer mockDB.Close()

    sqlMock.ExpectClose().WillReturnError(nil)

    db := &DB{
        Db:     mockDB,
        Logger: mockLogger,
    }

    err = db.CloseDB()
    assert.NoError(t, err)

    mockLogger.AssertExpectations(t)
    sqlMock.ExpectationsWereMet()
}
