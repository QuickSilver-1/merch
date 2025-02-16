package domain

// Inventory - структура для работы с записями об инвенторе пользователей
type Inventory struct {
	Id      uint64 `json:"id"`      // Уникальный идентификатор инвентаря
	Subject string `json:"subject"` // Название купленного мерча
	UserId  uint64 `json:"user"`    // Идентификатор пользователя, купившего мерч
}

// CreateInventory создает новый Inventory
func CreateInventory(subject string, userId uint64) *Inventory {
	return &Inventory{
		Subject: subject,
		UserId:  userId,
	}
}
