package server

import (
	"crypto/sha256"
	"encoding/hex"
	e "merch/internal/presentation/customError"
	"net/http"
	"regexp"
	"strings"
	"unicode"
)

// ValidPass проверяет пароль и возвращает его хэш, если он валиден
func ValidPass(pass string) (string, error) {
	if !isValidLength(pass) {
		return "", &e.InvalidPasswordFormat{
			Code: http.StatusBadRequest,
			Err:  "The password length must be from 3 to 30 characters inclusive",
		}
	}

	if !containsUpperAndDigit(pass) {
		return "", &e.InvalidPasswordFormat{
			Code: http.StatusBadRequest,
			Err:  "Password must contain at least 1 capital letter and 1 number",
		}
	}

	if containsInvalidChars(pass) {
		return "", &e.InvalidPasswordFormat{
			Code: http.StatusBadRequest,
			Err:  "The password must consist of letters of the Latin alphabet, numbers and symbols _!@#&*-",
		}
	}

	hash := GenHash(pass)
	return hash, nil
}

// GenHash создает хэш пароля
func GenHash(str string) string {
	hasher := sha256.New()
	hasher.Write([]byte(str)) // Преобразуем строку в хэш
	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash) // Возвращаем хэш в виде шестнадцатеричной строки
}

// IsValidEmail проверяет, является ли строка допустимым адресом электронной почты
func IsValidEmail(email string) bool {
	if len(email) > 256 {
		return false
	}

	// Регулярное выражение для проверки адреса электронной почты
	var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

	return emailRegex.MatchString(email)
}

// isValidLength проверяет длину пароля
func isValidLength(pass string) bool {
	return len(pass) >= 3 && len(pass) <= 30
}

// containsUpperAndDigit проверяет наличие заглавной буквы и цифры
func containsUpperAndDigit(pass string) bool {
	hasUpper, hasDigit := false, false
	for _, char := range pass {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasDigit = true
		}
	}
	return hasUpper && hasDigit
}

// containsInvalidChars проверяет наличие недопустимых символов в пароле
func containsInvalidChars(pass string) bool {
	for _, char := range pass {
		if !unicode.IsLetter(char) && !unicode.IsDigit(char) && !strings.Contains("_!@#&*-", string(char)) {
			return true
		}
	}
	return false
}
