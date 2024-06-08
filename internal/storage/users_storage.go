package storage

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/krevetkou/test-rutube/internal/domain"
	"strings"
	"time"
)

var jwtSecretKey = []byte("very-secret-key")

type Storage struct {
	users []domain.User
}

func NewStorage() *Storage {
	return &Storage{
		users: make([]domain.User, 0),
	}
}

func (s *Storage) InsertUser(userReg domain.RegisterRequest) (domain.User, error) {
	var lastID int

	ifExists := s.IsUserExists(userReg.TelegramName)
	if ifExists {
		return domain.User{}, domain.ErrExists
	}

	if len(s.users) > 0 {
		lastID = s.users[len(s.users)-1:][0].ID
	}

	user := domain.User{
		ID:           lastID + 1,
		Email:        userReg.Email,
		Password:     userReg.Password,
		Name:         userReg.Name,
		DateOfBirth:  userReg.DateOfBirth,
		TelegramName: userReg.TelegramName,
	}

	s.users = append(s.users, user)

	return user, nil
}

func (s *Storage) IsUserExists(tg string) bool {
	for i := range s.users {
		if strings.Contains(s.users[i].TelegramName, tg) {
			return true
		}
	}

	return false
}

func (s *Storage) GetAllUsers() ([]domain.ProfileResponse, error) {
	users := make([]domain.ProfileResponse, 0)
	for i := range s.users {
		user := domain.ProfileResponse{
			Email:        s.users[i].Email,
			Name:         s.users[i].Name,
			DateOfBirth:  s.users[i].DateOfBirth,
			TelegramName: s.users[i].TelegramName,
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *Storage) GetUserByTg(tg string) (domain.User, error) {
	var user *domain.User
	for i := range s.users {
		if s.users[i].TelegramName == tg {
			user = &s.users[i]
		}
	}

	if user == nil {
		return domain.User{}, domain.ErrNotFound
	}

	return *user, nil
}

func (s *Storage) CreateToken(tg string) (string, error) {
	payload := jwt.MapClaims{
		"sub": tg,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", domain.ErrTokenNotCreated
	}

	for i := range s.users {
		if s.users[i].TelegramName == tg {
			s.users[i].Token = t
		}
	}

	return t, nil
}

func (s *Storage) GetProfile(token domain.LoginResponse) (domain.ProfileResponse, error) {
	var user *domain.ProfileResponse
	for i := range s.users {
		if s.users[i].Token == token.AccessToken {
			user = &domain.ProfileResponse{
				Email:        s.users[i].Email,
				Name:         s.users[i].Name,
				DateOfBirth:  s.users[i].DateOfBirth,
				TelegramName: s.users[i].TelegramName,
			}
		}
	}

	if user == nil {
		return domain.ProfileResponse{}, domain.ErrNotFound
	}

	return *user, nil
}

func (s *Storage) IsTokenExists(token domain.LoginResponse) (bool, error) {
	for i := range s.users {
		if strings.Contains(s.users[i].Token, token.AccessToken) {
			return true, nil
		}
	}

	return false, domain.ErrNotExists
}

func (s *Storage) Subscribe(tgUser, subscribeUser string) error {
	ifExistsUser := s.IsUserExists(tgUser)
	if !ifExistsUser {
		return domain.ErrNotExists
	}

	ifExistsSubscribeUser := s.IsUserExists(subscribeUser)
	if !ifExistsSubscribeUser {
		return domain.ErrNotExists
	}

	for i := range s.users {
		if s.users[i].TelegramName == tgUser {
			for _, user := range s.users[i].SubscribeUsers {
				if user == subscribeUser {
					return domain.ErrExists
				}
			}
			s.users[i].SubscribeUsers = append(s.users[i].SubscribeUsers, subscribeUser)
		}
	}

	return nil
}

func (s *Storage) Unsubscribe(tgUser, subscribeUser string) error {
	ifExistsUser := s.IsUserExists(tgUser)
	if !ifExistsUser {
		return domain.ErrNotExists
	}

	ifExistsSubscribeUser := s.IsUserExists(subscribeUser)
	if !ifExistsSubscribeUser {
		return domain.ErrNotExists
	}

	for i := range s.users {
		if s.users[i].TelegramName == tgUser {
			for ind, user := range s.users[i].SubscribeUsers {
				if user == subscribeUser {
					s.users[i].SubscribeUsers = append(s.users[i].SubscribeUsers[:ind], s.users[i].SubscribeUsers[ind+1:]...)
					return nil
				}
			}
			return domain.ErrNotExists
		}
	}

	return nil
}
