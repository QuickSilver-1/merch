package services

import (
	"merch/internal/domain"
	"merch/internal/interfaces"
	"net/http"
)

const START_MONEY = 1000

// InvalidPassword - ошибка неверного пароля
type InvalidPassword struct {
	Err  string
	Code int
}

func (e *InvalidPassword) Error() string {
	return e.Err
}

func (e *InvalidPassword) GetCode() int {
	return e.Code
}

// UserService предоставляет методы для работы с пользователями
type UserService struct {
	user        interfaces.UserRepo
	auth        interfaces.AuthRepo
	transaction interfaces.TransactionRepo
	inventory   interfaces.InventoryRepo
}

// NewUserService создает новый экземпляр UserService
func NewUserService(user interfaces.UserRepo, auth interfaces.AuthRepo, transaction interfaces.TransactionRepo, inventory interfaces.InventoryRepo) *UserService {
	return &UserService{
		user: user,
		auth: auth,
		transaction: transaction,
		inventory: inventory,
	}
}

// Login выполняет авторизацию пользователя
func (s *UserService) Login(data domain.AuthorizationData) (*domain.Token, error) {
	user, err := s.user.GetByEmail(data.Username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		// Создание нового пользователя, если пользователь не найден
		user := domain.CreateUser(data.Username, data.Password, START_MONEY)
		id, err := s.user.Create(*user)
		if err != nil {
			return nil, err
		}
		data.Id = *id
	} else {
		// Проверка пароля
		if user.Password != data.Password {
			return nil, &InvalidPassword{
				Code: http.StatusUnauthorized,
				Err:  "Invalid password or email",
			}
		}
	}

	// Создание токена авторизации
	tokens, err := s.auth.CreateToken(data)
	if err != nil {
		return nil, err
	}

	return tokens, nil
}

// GetInfo получает информацию о пользователе по его идентификатору
func (s *UserService) GetInfo(userId uint64) (*domain.UserInfo, error) {
	user, err := s.user.GetById(userId)
	if err != nil {
		return nil, err
	}

	transactions, err := s.transaction.GetTransaction(user.Email)
	if err != nil {
		return nil, err
	}

	inventory, err := s.inventory.GetInventory(userId)
	if err != nil {
		return nil, err
	}

	return &domain.UserInfo{
		Coins:        user.Coins,
		Inventory:    *inventory,
		Transactions: *transactions,
	}, nil
}

// Token декодирует переданный токен и возвращает информацию о нем
func (s *UserService) Token(token domain.Token) (*domain.AuthorizationToken, error) {
	return s.auth.DecodeToken(token)
}

// Access проверяет не был ли токен доступа отозван
func (s *UserService) Access(token domain.Token, userId domain.UserId) (*domain.SuccessfulAuth, error) {
	return s.auth.Access(token, userId)
}
