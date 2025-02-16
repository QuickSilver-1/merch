package services

import (
	"merch/internal/domain"
	"merch/internal/interfaces"
	"net/http"
)

// NoMoneyError используется для обозначения ошибки недостатка средств
type NoMoneyError = InvalidPassword

// MoneyService предоставляет методы для работы с покупками и переводами средств
type MoneyService struct {
	user        interfaces.UserRepo
	transaction interfaces.TransactionRepo
	inventory   interfaces.InventoryRepo
}

// NewMoneyService создает новый экземпляр MoneyService
func NewMoneyService(user interfaces.UserRepo, transaction interfaces.TransactionRepo, inventory interfaces.InventoryRepo) *MoneyService {
	return &MoneyService{
		user:        user,
		transaction: transaction,
		inventory:   inventory,
	}
}

// BuyMerch выполняет покупку мерча
func (s *MoneyService) BuyMerch(inventory domain.Inventory) error {
	// Получение предмета по названию
	subject, err := s.inventory.GetSubjectByName(inventory.Subject)
	if err != nil {
		return err
	}

	// Получение данных пользователя по его идентификатору
	user, err := s.user.GetById(inventory.UserId)
	if err != nil {
		return err
	}

	// Проверка, хватает ли пользователю монет для покупки
	if user.Coins < subject.Cost {
		return &NoMoneyError{
			Code: http.StatusBadRequest,
			Err:  "Not enough coins to complete the transaction",
		}
	}

	// Покупка предмета
	err = s.inventory.Buy(inventory.UserId, *subject)
	if err != nil {
		return err
	}

	return nil
}

// MoneyTransfer выполняет перевод средств
func (s *MoneyService) MoneyTransfer(transaction domain.Transaction) error {
	// Получение данных пользователя по его email (имя отправителя)
	user, err := s.user.GetByEmail(transaction.SenderName)
	if err != nil {
		return err
	}

	// Проверка, хватает ли пользователю монет для перевода
	if user.Coins < transaction.Amount {
		return &NoMoneyError{
			Code: http.StatusBadRequest,
			Err:  "Not enough coins to complete the transaction",
		}
	}

	// Выполнение перевода средств
	return s.transaction.Transfer(transaction)
}
