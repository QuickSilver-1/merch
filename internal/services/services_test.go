package services

import (
	"errors"
	"merch/internal/domain"
	"merch/test/mocks"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	mockUserRepo        *mocks.MockUserRepo
	mockTransactionRepo *mocks.MockTransactionRepo
	mockInventoryRepo   *mocks.MockInventoryRepo
	mockAuthRepo        *mocks.MockAuthRepo
	moneyService        *MoneyService
	userService         *UserService
)

func setup() {
	mockUserRepo = new(mocks.MockUserRepo)
	mockTransactionRepo = new(mocks.MockTransactionRepo)
	mockInventoryRepo = new(mocks.MockInventoryRepo)
	mockAuthRepo = new(mocks.MockAuthRepo)

	moneyService = NewMoneyService(mockUserRepo, mockTransactionRepo, mockInventoryRepo)
	userService = NewUserService(mockUserRepo, mockAuthRepo, mockTransactionRepo, mockInventoryRepo)
}

type InvalidSubjectName = NoMoneyError

func TestMoneyService_BuyMerch(t *testing.T) {
	setup()

	tests := []struct {
		name                string
		user                *domain.User
		subject             *domain.Item
		inventory           domain.Inventory
		getSubjectByNameErr error
		getByIdErr          error
		buyErr              error
		expectErr           bool
		errType             interface{}
	}{
		{
			name:      "Successful Purchase",
			user:      &domain.User{Id: 1, Email: "test@example.com", Coins: 200},
			subject:   &domain.Item{Name: "t-shirt", Cost: 100},
			inventory: domain.Inventory{Subject: "t-shirt", UserId: 1},
			expectErr: false,
		},
		{
			name:                "Not Enough Coins",
			user:                &domain.User{Id: 1, Email: "test@example.com", Coins: 50},
			subject:             &domain.Item{Name: "t-shirt", Cost: 100},
			inventory:           domain.Inventory{Subject: "t-shirt", UserId: 1},
			getSubjectByNameErr: nil,
			getByIdErr:          nil,
			expectErr:           true,
			errType:             &NoMoneyError{},
		},
		{
			name:                "Invalid subject's name",
			user:                &domain.User{Id: 1, Email: "test@example.com", Coins: 50},
			subject:             nil,
			inventory:           domain.Inventory{Subject: "cup", UserId: 1},
			getSubjectByNameErr: &InvalidSubjectName{Code: http.StatusBadRequest},
			expectErr:           true,
			errType:             &InvalidSubjectName{},
		},
		{
			name:       "GetById Error",
			user:       nil,
			subject:    &domain.Item{Name: "t-shirt", Cost: 100},
			inventory:  domain.Inventory{Subject: "t-shirt", UserId: 1},
			getByIdErr: errors.New("user not found"),
			expectErr:  true,
			errType:    errors.New("user not found"),
		},
		{
			name:                "Buy Error",
			user:                &domain.User{Id: 1, Email: "test@example.com", Coins: 200},
			subject:             &domain.Item{Name: "t-shirt", Cost: 100},
			inventory:           domain.Inventory{Subject: "t-shirt", UserId: 1},
			getSubjectByNameErr: nil,
			getByIdErr:          nil,
			buyErr:              errors.New("could not buy item"),
			expectErr:           true,
			errType:             errors.New("could not buy item"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockInventoryRepo.On("GetSubjectByName", tt.inventory.Subject).Return(tt.subject, tt.getSubjectByNameErr)
			mockUserRepo.On("GetById", tt.inventory.UserId).Return(tt.user, tt.getByIdErr)
			if tt.subject != nil && tt.user != nil {
				mockInventoryRepo.On("Buy", tt.inventory.UserId, *tt.subject).Return(tt.buyErr)
			}

			err := moneyService.BuyMerch(tt.inventory)
			if tt.expectErr {
				assert.Error(t, err)
				assert.IsType(t, tt.errType, err)
				if errType, ok := err.(*NoMoneyError); ok {
					assert.Equal(t, http.StatusBadRequest, errType.GetCode())
				}
			} else {
				assert.NoError(t, err)
			}

			mockInventoryRepo.AssertExpectations(t)
			mockUserRepo.AssertExpectations(t)

			mockInventoryRepo.ExpectedCalls = nil
			mockUserRepo.ExpectedCalls = nil
		})
	}
}

func TestMoneyService_MoneyTransfer(t *testing.T) {
	setup()

	tests := []struct {
		name          string
		user          *domain.User
		transaction   domain.Transaction
		getByEmailErr error
		expectErr     bool
		errType       interface{}
	}{
		{
			name:        "Successful Transfer",
			user:        &domain.User{Email: "test@example.com", Coins: 200},
			transaction: domain.Transaction{SenderName: "test@example.com", Amount: 100},
			expectErr:   false,
		},
		{
			name:        "Not Enough Coins",
			user:        &domain.User{Email: "test@example.com", Coins: 50},
			transaction: domain.Transaction{SenderName: "test@example.com", Amount: 100},
			expectErr:   true,
			errType:     &NoMoneyError{},
		},
		{
			name:          "GetByEmail Error",
			user:          nil,
			transaction:   domain.Transaction{SenderName: "invalid@example.com", Amount: 100},
			getByEmailErr: errors.New("user not found"),
			expectErr:     true,
			errType:       errors.New("user not found"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo.On("GetByEmail", tt.transaction.SenderName).Return(tt.user, tt.getByEmailErr)
			if !tt.expectErr {
				mockTransactionRepo.On("Transfer", tt.transaction).Return(nil)
			}

			err := moneyService.MoneyTransfer(tt.transaction)
			if tt.expectErr {
				assert.Error(t, err)
				assert.IsType(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
			}

			mockUserRepo.AssertExpectations(t)
			mockTransactionRepo.AssertExpectations(t)

			mockUserRepo.ExpectedCalls = nil
			mockTransactionRepo.ExpectedCalls = nil
		})
	}
}

func TestUserService_Login(t *testing.T) {
	setup()

	tests := []struct {
		name           string
		authData       domain.AuthorizationData
		user           *domain.User
		getByEmailErr  error
		createErr      error
		createTokenErr error
		expectedId     domain.UserId
		expectedToken  domain.Token
		expectErr      bool
		errType        interface{}
	}{
		{
			name:          "Create New User",
			authData:      domain.AuthorizationData{Username: "newuser@example.com", Password: "Password123"},
			user:          nil,
			expectedId:    1,
			expectedToken: "testtoken",
			expectErr:     false,
		},
		{
			name:          "Invalid Password",
			authData:      domain.AuthorizationData{Username: "testuser@example.com", Password: "wrongpassword"},
			user:          &domain.User{Email: "testuser@example.com", Password: "correctpassword"},
			expectedId:    0,
			expectedToken: "",
			expectErr:     true,
			errType:       &InvalidPassword{},
		},
		{
			name:          "GetByEmail Error",
			authData:      domain.AuthorizationData{Username: "testuser@example.com", Password: "password"},
			user:          nil,
			getByEmailErr: errors.New("database error"),
			expectErr:     true,
			errType:       errors.New("database error"),
		},
		{
			name:      "Create Error",
			authData:  domain.AuthorizationData{Username: "newuser2@example.com", Password: "Password123"},
			user:      nil,
			createErr: errors.New("database error"),
			expectErr: true,
			errType:   errors.New("database error"),
		},
		{
			name:           "CreateToken Error",
			authData:       domain.AuthorizationData{Username: "newuser3@example.com", Password: "Password123"},
			user:           nil,
			createTokenErr: errors.New("token error"),
			expectErr:      true,
			errType:        errors.New("token error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo.On("GetByEmail", tt.authData.Username).Return(tt.user, tt.getByEmailErr)
			if tt.user == nil && tt.getByEmailErr == nil {
				mockUserRepo.On("Create", mock.Anything).Return(&tt.expectedId, tt.createErr)
			}
			if tt.createErr == nil {
				mockAuthRepo.On("CreateToken", mock.Anything).Return(&tt.expectedToken, tt.createTokenErr)
			}

			token, err := userService.Login(tt.authData)
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, token)
				assert.IsType(t, tt.errType, err)
				if errType, ok := err.(*InvalidPassword); ok {
					assert.Equal(t, http.StatusUnauthorized, errType.GetCode())
				}
			} else {
				assert.NoError(t, err)
				assert.Equal(t, &tt.expectedToken, token)
			}

			mockUserRepo.AssertExpectations(t)
			mockAuthRepo.AssertExpectations(t)

			mockUserRepo.ExpectedCalls = nil
			mockAuthRepo.ExpectedCalls = nil
		})
	}
}

func TestUserService_GetInfo(t *testing.T) {
	setup()

	tests := []struct {
		name              string
		user              *domain.User
		transactions      *[]domain.Transaction
		inventory         *[]domain.Inventory
		getByIdErr        error
		getTransactionErr error
		getInventoryErr   error
		expectErr         bool
		errType           interface{}
	}{
		{
			name:         "Successful GetInfo",
			user:         &domain.User{Id: 1, Email: "testuser@example.com", Coins: 100},
			transactions: &[]domain.Transaction{},
			inventory:    &[]domain.Inventory{},
			expectErr:    false,
		},
		{
			name:         "GetById Error",
			user:         nil,
			transactions: &[]domain.Transaction{},
			inventory:    &[]domain.Inventory{},
			getByIdErr:   errors.New("user not found"),
			expectErr:    true,
			errType:      errors.New("user not found"),
		},
		{
			name:              "GetTransaction Error",
			user:              &domain.User{Id: 1, Email: "testuser@example.com", Coins: 100},
			transactions:      nil,
			inventory:         &[]domain.Inventory{},
			getTransactionErr: errors.New("transaction error"),
			expectErr:         true,
			errType:           errors.New("transaction error"),
		},
		{
			name:            "GetInventory Error",
			user:            &domain.User{Id: 1, Email: "testuser@example.com", Coins: 100},
			transactions:    &[]domain.Transaction{},
			inventory:       nil,
			getInventoryErr: errors.New("inventory error"),
			expectErr:       true,
			errType:         errors.New("inventory error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUserRepo.On("GetById", mock.Anything).Return(tt.user, tt.getByIdErr)
			mockTransactionRepo.On("GetTransaction", mock.Anything).Return(tt.transactions, tt.getTransactionErr)
			mockInventoryRepo.On("GetInventory", mock.Anything).Return(tt.inventory, tt.getInventoryErr)

			userInfo, err := userService.GetInfo(uint64(1))
			if tt.expectErr {
				assert.Error(t, err)
				assert.Nil(t, userInfo)
				assert.IsType(t, tt.errType, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, userInfo)
				assert.Equal(t, tt.user.Coins, userInfo.Coins)
				assert.Equal(t, tt.transactions, &userInfo.Transactions)
				assert.Equal(t, tt.inventory, &userInfo.Inventory)
			}

			mockUserRepo.AssertExpectations(t)
			mockTransactionRepo.AssertExpectations(t)
			mockInventoryRepo.AssertExpectations(t)

			mockUserRepo.ExpectedCalls = nil
			mockTransactionRepo.ExpectedCalls = nil
			mockInventoryRepo.ExpectedCalls = nil
		})
	}
}

func TestUserService_Token(t *testing.T) {
	setup()

	var token domain.Token = "testtoken"
	expectedData := &domain.AuthorizationToken{Id: 1, Email: "testuser", Expires: time.Now().Add(time.Hour)}

	mockAuthRepo.On("DecodeToken", token).Return(expectedData, nil)

	data, err := userService.Token(token)
	assert.NoError(t, err)
	assert.Equal(t, expectedData, data)

	mockAuthRepo.AssertExpectations(t)
}

func TestUserService_Access(t *testing.T) {
	setup()

	var token domain.Token = "testtoken"
	userId := domain.UserId(1)
	var expectedAccess domain.SuccessfulAuth = true

	mockAuthRepo.On("Access", token, userId).Return(&expectedAccess, nil)

	access, err := userService.Access(token, userId)
	assert.NoError(t, err)
	assert.Equal(t, &expectedAccess, access)

	mockAuthRepo.AssertExpectations(t)
}
