package server

import (
	"encoding/json"
	"fmt"
	"merch/internal/domain"
	e "merch/internal/presentation/customError"
	"merch/internal/presentation/realization"
	"net/http"

	"github.com/gin-gonic/gin"
)

type inventoryForm struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
}

type recieverTransaction struct {
	FromUser string `json:"fromUser"`
	Amount   int    `json:"amount"`
}

type senderTransaction struct {
	ToUser string `json:"toUser"`
	Amount int    `json:"amount"`
}

type coinHistory struct {
	Received []recieverTransaction `json:"received"`
	Sent     []senderTransaction   `json:"sent"`
}

type userForm struct {
	Coins       int             `json:"coins"`
	Inventory   []inventoryForm `json:"inventory"`
	CoinHistory coinHistory     `json:"coinHistory"`
}

// Handlers определяет хендлеры для обработки HTTP-запросов
type Handlers struct{}

// NewHandlers создает новый экземпляр Handlers
func NewHandlers() *Handlers {
	return &Handlers{}
}

// GetInfo возвращает информацию о пользователе
func (*Handlers) GetInfo(ctx *gin.Context) {
	token := getJWT(ctx)

	if token == nil {
		answerError(ctx, &e.NeedAuthorization{
			Code: http.StatusUnauthorized,
			Err:  "No cookie",
		})
		return
	}

	info, err := UserService.GetInfo(token.Id)
	if err != nil {
		answerError(ctx, err)
		return
	}

	inventory := make(map[string]int)
	for _, subject := range info.Inventory {
		quantity, exists := inventory[subject.Subject]

		if !exists {
			quantity = 0
		}

		inventory[subject.Subject] = quantity + 1
	}

	var userInfo userForm
	for name, quantity := range inventory {
		userInfo.Inventory = append(userInfo.Inventory, inventoryForm{
			Type:     name,
			Quantity: quantity,
		})
	}

	var coinTransactions coinHistory
	for _, transaction := range info.Transactions {
		if transaction.SenderName == token.Email {
			coinTransactions.Sent = append(coinTransactions.Sent, senderTransaction{
				ToUser: transaction.ReceiverName,
				Amount: transaction.Amount,
			})
		} else {
			coinTransactions.Received = append(coinTransactions.Received, recieverTransaction{
				FromUser: transaction.SenderName,
				Amount:   transaction.Amount,
			})
		}
	}

	userInfo.CoinHistory = coinTransactions
	ctx.JSON(http.StatusOK, userInfo)
}

// SendCoin обрабатывает запрос на отправку монет
func (*Handlers) SendCoin(ctx *gin.Context) {
	token := getJWT(ctx)

	if token == nil {
		answerError(ctx, &e.NeedAuthorization{
			Code: http.StatusUnauthorized,
			Err:  "No cookie",
		})
		return
	}

	var transaction domain.Transaction
	err := json.NewDecoder(ctx.Request.Body).Decode(&transaction)
	defer func() {
		err := ctx.Request.Body.Close()
		if err != nil {
			realization.LoggerService.Error(fmt.Sprintf("Error closing body: %v", err))
		}
	}()

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, map[string]string{"errors": STATUS_UNAUTHORIZED})
		return
	}

	transaction.SenderName = token.Email
	err = BuyService.MoneyTransfer(transaction)

	if err != nil {
		answerError(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

// BuyMerch обрабатывает запрос на покупку товара
func (*Handlers) BuyMerch(ctx *gin.Context) {
	item := ctx.Param("item")
	token := getJWT(ctx)

	if token == nil {
		answerError(ctx, &e.NeedAuthorization{
			Code: http.StatusUnauthorized,
			Err:  "No cookie",
		})
		return
	}

	inventory := domain.Inventory{
		Subject: item,
		UserId:  token.Id,
	}

	err := BuyService.BuyMerch(inventory)
	if err != nil {
		answerError(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

// Auth обрабатывает запрос на авторизацию пользователя
func (*Handlers) Auth(ctx *gin.Context) {
	var data domain.AuthorizationData
	err := json.NewDecoder(ctx.Request.Body).Decode(&data)
	defer func() {
		err := ctx.Request.Body.Close()
		if err != nil {
			realization.LoggerService.Error(fmt.Sprintf("Error closing body: %v", err))
		}
	}()

	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"errors": fmt.Sprintf("%s: %s", STATUS_BAD_REQUEST, "check input data")})
		return
	}

	if !IsValidEmail(data.Username) {
		ctx.JSON(http.StatusBadRequest, map[string]string{"errors": fmt.Sprintf("%s: %s invalid email", STATUS_BAD_REQUEST, data.Username)})
		return
	}

	pass, err := ValidPass(data.Password)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, map[string]string{"errors": fmt.Sprintf("%s: %v", STATUS_BAD_REQUEST, err)})
		return
	}

	data.Password = pass
	token, err := UserService.Login(data)
	if err != nil {
		answerError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, map[string]string{"token": *token})
}

// answerError обрабатывает ошибки и возвращает соответствующий HTTP-статус
func answerError(ctx *gin.Context, err error) {
	baseErr := err.(*e.BaseError)

	switch baseErr.GetCode() {
	case http.StatusUnauthorized:
		ctx.JSON(http.StatusUnauthorized, STATUS_UNAUTHORIZED)
	case http.StatusInternalServerError:
		realization.LoggerService.Error(baseErr.Error())
		ctx.JSON(http.StatusInternalServerError, STATUS_INTERNAL_SERVER)
	case http.StatusBadRequest:
		ctx.JSON(http.StatusBadRequest, map[string]string{"errors": fmt.Sprintf("%s: %s", STATUS_BAD_REQUEST, baseErr.Error())})
	}
}

// getJWT извлекает JWT токен из контекста
func getJWT(ctx *gin.Context) *realization.Token {
	tokenStr, exists := ctx.Get("token")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, map[string]string{"errors": STATUS_UNAUTHORIZED})
		return nil
	}

	token := tokenStr.(realization.Token)
	return &token
}
