package services

import (
	"errors"
	"github.com/krevetkou/test-rutube/internal/domain"
)

type UsersRepository interface {
	InsertUser(userReg domain.RegisterRequest) (domain.User, error)
	IsUserExists(email string) bool
	GetUsersByToken(token string) ([]domain.ProfileResponse, error)
	GetAllUsers() ([]domain.UserInListResponse, error)
	GetUserByEmail(email string) (domain.User, error)
	CreateToken(email string) (string, error)
	GetUserInfo(token string) (domain.UserResponse, error)
	IsTokenExists(token string) (bool, error)
	Subscribe(token string, id int) error
	Unsubscribe(token string, id int) error
	Settings(token string, daysToBirthday int, email string) error
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
	if user.Name == "" || user.Email == "" || user.Password == "" || user.DateOfBirth == "" {
		return domain.User{}, domain.ErrFieldsRequired
	}

	isUserExists := s.Storage.IsUserExists(user.Email)
	if isUserExists {
		return domain.User{}, domain.ErrExists
	}

	newUser, err := s.Storage.InsertUser(user)
	if err != nil {
		return domain.User{}, err
	}

	return newUser, nil
}

func (s UsersService) GetAllUsers() ([]domain.UserInListResponse, error) {
	users, err := s.Storage.GetAllUsers()
	if err != nil {
		return []domain.UserInListResponse{}, err
	}

	return users, nil
}

func (s UsersService) GetUsersByToken(token string) ([]domain.ProfileResponse, error) {
	users, err := s.Storage.GetUsersByToken(token)
	if err != nil {
		return []domain.ProfileResponse{}, err
	}

	return users, nil
}

func (s UsersService) Login(user domain.LoginRequest) (domain.UserResponse, error) {
	isUserExists := s.Storage.IsUserExists(user.Email)
	if !isUserExists {
		return domain.UserResponse{}, domain.ErrNotExists
	}

	userData, err := s.Storage.GetUserByEmail(user.Email)
	if err != nil {
		return domain.UserResponse{}, domain.ErrNotExists
	}

	if user.Password != userData.Password {
		return domain.UserResponse{}, domain.ErrBadCredentials
	}

	return domain.UserResponse{
		Email:              userData.Email,
		DaysToNotification: userData.DaysToNotification,
		Name:               userData.Name,
	}, nil
}

func (s UsersService) CreateToken(email string) (string, error) {
	t, err := s.Storage.CreateToken(email)
	if err != nil {
		return "", domain.ErrTokenNotCreated
	}

	return t, nil
}

func (s UsersService) GetUserInfo(token string) (domain.UserResponse, error) {
	user, err := s.Storage.GetUserInfo(token)
	if err != nil {
		return domain.UserResponse{}, domain.ErrNotFound
	}

	return user, nil
}

func (s UsersService) Subscribe(token string, userId int) error {
	err := s.Storage.Subscribe(token, userId)
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

func (s UsersService) Unsubscribe(token string, userId int) error {
	err := s.Storage.Unsubscribe(token, userId)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotExists):
			return domain.ErrNotExists
		}
	}
	return nil
}

func (s UsersService) Settings(token string, daysToBirthday int, email string) error {
	err := s.Storage.Settings(token, daysToBirthday, email)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotExists):
			return domain.ErrNotExists
		}
	}
	return nil
}
