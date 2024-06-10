package storage

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/krevetkou/test-rutube/internal/domain"
	"time"
)

const DefaultDays = 2

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

	ifExists := s.IsUserExists(userReg.Email)
	if ifExists {
		return domain.User{}, domain.ErrExists
	}

	if len(s.users) > 0 {
		lastID = s.users[len(s.users)-1:][0].ID
	}

	user := domain.User{
		ID:                 lastID + 1,
		Email:              userReg.Email,
		Password:           userReg.Password,
		Name:               userReg.Name,
		DateOfBirth:        userReg.DateOfBirth,
		DaysToNotification: DefaultDays,
	}

	s.users = append(s.users, user)

	return user, nil
}

func (s *Storage) IsUserExists(email string) bool {
	for i := range s.users {
		if s.users[i].Email == email {
			return true
		}
	}

	return false
}

func (s *Storage) GetAllUsers() ([]domain.UserInListResponse, error) {
	users := make([]domain.UserInListResponse, 0)

	for i, val := range s.users {
		isToday := val.DateOfBirth == time.Now().Format("2006-01-02")

		if isToday {
			users = append(users, domain.UserInListResponse{
				Email: s.users[i].Email,
				Name:  s.users[i].Name,
			})
		}
	}

	return users, nil
}

func (s *Storage) GetUsersByToken(token string) ([]domain.ProfileResponse, error) {
	users := make([]domain.ProfileResponse, 0)
	currentUser, err := getUserByToken(s.users, token)
	if err != nil {
		return []domain.ProfileResponse{}, err
	}

	for _, user := range s.users {
		users = append(users, domain.ProfileResponse{
			ID:           user.ID,
			Email:        user.Email,
			Name:         user.Name,
			IsSubscribed: contains(currentUser.SubscribeUsers, user.ID),
		})
	}

	return users, nil
}

func (s *Storage) GetUserByEmail(email string) (domain.User, error) {
	var user *domain.User
	for i := range s.users {
		if s.users[i].Email == email {
			user = &s.users[i]
		}
	}

	if user == nil {
		return domain.User{}, domain.ErrNotFound
	}

	return *user, nil
}

func (s *Storage) CreateToken(email string) (string, error) {
	payload := jwt.MapClaims{
		"sub": email,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)

	t, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", domain.ErrTokenNotCreated
	}

	for i := range s.users {
		if s.users[i].Email == email {
			s.users[i].Token = t
		}
	}

	return t, nil
}

func (s *Storage) GetUserInfo(token string) (domain.UserResponse, error) {
	for _, user := range s.users {
		if user.Token == token {
			return domain.UserResponse{
				Email:              user.Email,
				Name:               user.Name,
				DaysToNotification: user.DaysToNotification,
			}, nil
		}
	}

	return domain.UserResponse{}, domain.ErrNotFound
}

func (s *Storage) IsTokenExists(token string) (bool, error) {
	for i := range s.users {
		if s.users[i].Token == token {
			return true, nil
		}
	}

	return false, domain.ErrNotExists
}

func (s *Storage) Subscribe(token string, userId int) error {
	ifExistsUser, err := s.IsTokenExists(token)
	if !ifExistsUser {
		return domain.ErrNotExists
	}
	if err != nil {
		return err
	}

	for i := range s.users {
		if s.users[i].Token == token {
			for _, id := range s.users[i].SubscribeUsers {
				if id == userId {
					return domain.ErrExists
				}
			}
			s.users[i].SubscribeUsers = append(s.users[i].SubscribeUsers, userId)
		}
	}
	return nil
}

func (s *Storage) Unsubscribe(token string, userId int) error {
	ifExistsUser, err := s.IsTokenExists(token)
	if !ifExistsUser {
		return domain.ErrNotExists
	}
	if err != nil {
		return err
	}

	for i := range s.users {
		if s.users[i].Token == token {
			for ind, id := range s.users[i].SubscribeUsers {
				if id == userId {
					s.users[i].SubscribeUsers = append(s.users[i].SubscribeUsers[:ind], s.users[i].SubscribeUsers[ind+1:]...)
					return nil
				}
			}
			return domain.ErrNotExists
		}
	}
	return nil
}

func (s *Storage) Settings(token string, daysToNotification int, email string) error {
	ifExistsUser, err := s.IsTokenExists(token)
	if !ifExistsUser {
		return domain.ErrNotExists
	}
	if err != nil {
		return err
	}

	for _, user := range s.users {
		if user.Token == token {
			user.DaysToNotification = daysToNotification
			user.Email = email
		}
	}

	return nil
}

func contains(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func getUserByToken(users []domain.User, token string) (domain.User, error) {
	for _, user := range users {
		if user.Token == token {
			return user, nil
		}
	}

	return domain.User{}, domain.ErrNotExists
}
