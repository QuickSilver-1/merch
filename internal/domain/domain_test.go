package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserCreate(t *testing.T) {
    user := CreateUser("test@example.com", "Password123", 1000)
    assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Password123", user.Password)
	assert.Equal(t, 1000, user.Coins)
}

func TestTransactionCreate(t *testing.T) {
    transaction := CreateTransaction("sender@example.com", "reciever@example.com", 100)
    assert.Equal(t, "sender@example.com", transaction.SenderName)
	assert.Equal(t, "reciever@example.com", transaction.ReceiverName)
	assert.Equal(t, 100, transaction.Amount)
}

func TestItemCreate(t *testing.T) {
    item := CreateItem("cup", 20)
    assert.Equal(t, "cup", item.Name)
	assert.Equal(t, 20, item.Cost)
}

func TestInventoryCreate(t *testing.T) {
    inventory := CreateInventory("cup", 1)
    assert.Equal(t, "cup", inventory.Subject)
	assert.Equal(t, uint64(1), inventory.UserId)
}