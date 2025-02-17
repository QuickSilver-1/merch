package tests

import (
	"bytes"
	"encoding/json"
	"merch/internal/domain"
	"merch/internal/presentation/server"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMerchPurchaseScenarios(t *testing.T) {
	// Создание пользователей и получение их токенов
	user1 := domain.AuthorizationData{Username: "buy1@example.com", Password: "Password123"}
	token1 := CreateUserAndGetToken(t, user1)
	user2 := domain.AuthorizationData{Username: "buy2@example.com", Password: "Password123"}
	_ = CreateUserAndGetToken(t, user2)

	// Тест 1: Покупка мерча
	item := "cup"
	buyMerch(t, token1, item, http.StatusOK)

	// // Тест 2: Покупка мерча при нехватке денег
	transaction := server.SenderTransaction{ToUser: user2.Username, Amount: 950}
	transfer(t, token1, transaction)
	expensiveItem := "t-shirt"
	buyMerch(t, token1, expensiveItem, http.StatusBadRequest)
}

func buyMerch(t *testing.T, token string, item string, expectedStatus int) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8080/api/buy/"+item, nil)
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			return
		}
	}()
	assert.Equal(t, expectedStatus, resp.StatusCode)
}

func transfer(t *testing.T, token string, transaction server.SenderTransaction) {
	payload, err := json.Marshal(transaction)
	assert.NoError(t, err)
	client := &http.Client{}
	req, err := http.NewRequest("POST", "http://localhost:8080/api/sendCoin", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	assert.NoError(t, err)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			return
		}
	}()
}
