package domain

// User - объект пользователя системы
type User struct {
	Id       uint64 `json:"id"`       // Уникальный идентификатор пользователя
	Email    string `json:"email"`    // Email пользователя
	Password string `json:"password"` // Пароль пользователя (хранится только в зашифрованном виде)
	Coins    int    `json:"coins"`    // Количество монет у пользователя
}

// CreateUser создает нового пользователя
func CreateUser(email, pass string, coins int) *User {
	return &User{
		Email:    email,
		Password: pass,
		Coins:    coins,
	}
}
