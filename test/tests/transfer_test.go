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

func TestMoneyTransferScenarios(t *testing.T) {
	// Создание пользователей и получение их токенов
	user1 := domain.AuthorizationData{Username: "transfer1@example.com", Password: "Password123"}
	user2 := domain.AuthorizationData{Username: "transfer2@example.com", Password: "Password123"}
	token1 := CreateUserAndGetToken(t, user1)
	_ = CreateUserAndGetToken(t, user2)

	// Тест 1: Перевод денег другому пользователю
	transferData := server.SenderTransaction{ToUser: "transfer2@example.com", Amount: 100}
	transferMoney(t, token1, transferData, http.StatusOK)

	// Тест 2: Перевод денег другому пользователю при нехватке денег
	transferData.Amount = 10000 // Увеличиваем сумму перевода
	transferMoney(t, token1, transferData, http.StatusBadRequest)

	// Тест 3: Перевод денег самому себе
	transferData = server.SenderTransaction{ToUser: "transfer1@example.com", Amount: 100}
	transferMoney(t, token1, transferData, http.StatusBadRequest)

	// Тест 4: Перевод денег несуществующему пользователю
	transferData = server.SenderTransaction{ToUser: "nonexistent@example.com", Amount: 100}
	transferMoney(t, token1, transferData, http.StatusBadRequest)
}

func CreateUserAndGetToken(t *testing.T, user domain.AuthorizationData) string {
	payload, err := json.Marshal(user)
	assert.NoError(t, err)
	resp, err := http.Post("http://localhost:8080/api/auth", "application/json", bytes.NewBuffer(payload))
	assert.NoError(t, err)
	defer func() {
		err := resp.Body.Close()
		if err != nil {
			return
		}
	}()
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	var result map[string]string
	err = json.NewDecoder(resp.Body).Decode(&result)
	assert.NoError(t, err)
	return result["token"]
}

func transferMoney(t *testing.T, token string, transaction server.SenderTransaction, expectedStatus int) {
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
	assert.Equal(t, expectedStatus, resp.StatusCode)
}
