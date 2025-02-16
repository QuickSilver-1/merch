package interfaces

import "merch/internal/domain"

// InventoryRepo предоставляет методы для работы с инвентарем
type InventoryRepo interface {
	// Buy выполняет покупку предмета пользователем
	Buy(userId domain.UserId, subject domain.Item) error

	// GetSubjectByName получает предмет по его названию
	GetSubjectByName(email domain.UserEmail) (subject *domain.Item, err error)

	// GetInventory получает инвентарь пользователя по его идентификатору
	GetInventory(userId domain.UserId) (inventory *[]domain.Inventory, err error)
}
