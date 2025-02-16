package domain

// Transaction - объект для переводов днег между пользователями
type Transaction struct {
	Id           uint64 `json:"id"`            // Уникальный идентификатор транзакции
	SenderName   string `json:"sender_name"`   // Имя отправителя
	ReceiverName string `json:"receiver_name"` // Имя получателя
	Amount       int    `json:"amount"`        // Сумма транзакции
}

// CreateTransaction создает новую транзакцию с указанными параметрами
func CreateTransaction(senderName, receiverName string, amount int) *Transaction {
	return &Transaction{
		SenderName:   senderName,
		ReceiverName: receiverName,
		Amount:       amount,
	}
}
