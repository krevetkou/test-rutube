package services

import (
	"errors"
	"github.com/krevetkou/test-rutube/internal/domain"
)

type UsersRepository interface {
	InsertUser(userReg domain.RegisterRequest) (domain.User, error)
	IsUserExists(tg string) bool
	GetAllUsers() ([]domain.ProfileResponse, error)
	GetUserByTg(tg string) (domain.User, error)
	CreateToken(tg string) (string, error)
	GetProfile(token domain.LoginResponse) (domain.ProfileResponse, error)
	IsTokenExists(token domain.LoginResponse) (bool, error)
	Subscribe(tgUser, subscribeUser string) error
	Unsubscribe(tgUser, subscribeUser string) error
}

type UsersService struct {
	Storage UsersRepository
}

func NewUserService(storage UsersRepository) UsersService {
	return UsersService{
		Storage: storage,
	}
}

func (s UsersService) Create(user domain.RegisterRequest) (domain.User, error) {
	// входящие параметры необходимо валидировать
	if user.Name == "" || user.Email == "" || user.Password == "" || user.DateOfBirth == "" || user.TelegramName == "" {
		return domain.User{}, domain.ErrFieldsRequired
	}

	isUserExists := s.Storage.IsUserExists(user.TelegramName)
	if isUserExists {
		return domain.User{}, domain.ErrExists
	}

	newUser, err := s.Storage.InsertUser(user)
	if err != nil {
		return domain.User{}, err
	}

	return newUser, nil
}

func (s UsersService) List() ([]domain.ProfileResponse, error) {
	users, err := s.Storage.GetAllUsers()
	if err != nil {
		return []domain.ProfileResponse{}, err
	}

	return users, nil
}

func (s UsersService) Login(user domain.LoginRequest) (domain.User, error) {
	isUserExists := s.Storage.IsUserExists(user.TelegramName)
	if !isUserExists {
		return domain.User{}, domain.ErrNotExists
	}

	userData, err := s.Storage.GetUserByTg(user.Email)
	if err != nil {
		return domain.User{}, domain.ErrNotExists
	}

	if user.Password != userData.Password {
		return domain.User{}, domain.ErrBadCredentials
	}

	return userData, nil
}

func (s UsersService) CreateToken(user domain.User) (string, error) {
	t, err := s.Storage.CreateToken(user.Email)
	if err != nil {
		return "", domain.ErrTokenNotCreated
	}

	return t, nil
}

func (s UsersService) GetProfile(token domain.LoginResponse) (domain.ProfileResponse, error) {
	user, err := s.Storage.GetProfile(token)
	if err != nil {
		return domain.ProfileResponse{}, domain.ErrNotFound
	}

	return user, nil
}

func (s UsersService) Subscribe(tgUser, subscribeUser string) error {
	err := s.Storage.Subscribe(tgUser, subscribeUser)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrExists):
			return domain.ErrExists
		case errors.Is(err, domain.ErrNotExists):
			return domain.ErrNotExists
		}
	}
	return nil
}

func (s UsersService) Unsubscribe(tgUser, subscribeUser string) error {
	err := s.Storage.Unsubscribe(tgUser, subscribeUser)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotExists):
			return domain.ErrNotExists
		}
	}
	return nil
}
