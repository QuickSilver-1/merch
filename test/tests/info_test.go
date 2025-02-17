package tests

import (
	"bytes"
	"encoding/json"
	"merch/internal/domain"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserInfoRetrieval(t *testing.T) {
	// Создание пользователя и получение его токена
	user := domain.AuthorizationData{Username: "info@example.com", Password: "Password123"}
	token := createUserAndGetToken(t, user)

	// Тест 1: Вывод информации о пользователе
	retrieveUserInfo(t, token, http.StatusOK)
}

func createUserAndGetToken(t *testing.T, user domain.AuthorizationData) string {
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

func retrieveUserInfo(t *testing.T, token string, expectedStatus int) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "http://localhost:8080/api/info", nil)
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
