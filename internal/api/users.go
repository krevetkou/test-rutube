package api

import (
	"encoding/json"
	"errors"
	"github.com/krevetkou/test-rutube/internal/domain"
	"log"
	"net/http"
	"time"
)

const CookieExpired = time.Hour * 24 * 365 * 10
const AuthCookieName = "token"

type UsersService interface {
	GetUsersByToken(token string) ([]domain.ProfileResponse, error)
	GetAllUsers() ([]domain.UserInListResponse, error)
	Login(actor domain.LoginRequest) (domain.UserResponse, error)
	CreateToken(email string) (string, error)
	GetUserInfo(token string) (domain.UserResponse, error)
	Subscribe(token string, userId int) error
	Unsubscribe(token string, userId int) error
	Settings(token string, daysToBirthday int, email string) error
}

type UsersHandler struct {
	Service UsersService
}

func NewUsersHandler(service UsersService) UsersHandler {
	return UsersHandler{
		Service: service,
	}
}

func (h UsersHandler) ListToday(w http.ResponseWriter, r *http.Request) {
	users, err := h.Service.GetAllUsers()
	if err != nil {
		http.Error(w, "failed to get users", http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(users)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to create response data", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}

func (h UsersHandler) List(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(AuthCookieName)

	if err != nil {
		http.Error(w, "need to login", http.StatusUnauthorized)
		return
	}

	users, err := h.Service.GetUsersByToken(cookie.Value)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to get users", http.StatusInternalServerError)
		return
	}
	data, err := json.Marshal(users)
	if err != nil {
		log.Println(err)
		http.Error(w, "failed to create response data", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}

func (h UsersHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "content type not allowed", http.StatusUnsupportedMediaType)
		return
	}

	var request domain.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		log.Println(err)
		return
	}

	createdUser, err := h.Service.Login(request)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotExists):
			http.Error(w, "user doesn't exist", http.StatusConflict)
		case errors.Is(err, domain.ErrBadCredentials):
			http.Error(w, "incorrect password", http.StatusBadRequest)
		default:
			http.Error(w, "unexpected error", http.StatusInternalServerError)
		}
		log.Println(err)

		return
	}

	data, err := json.Marshal(createdUser)
	if err != nil {
		http.Error(w, "failed to create response data", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	token, err := h.Service.CreateToken(request.Email)
	if err != nil {
		log.Println(err)
		return
	}

	authCookie := http.Cookie{
		Name:     AuthCookieName,
		Value:    token,
		Path:     ".app.localhost",
		Secure:   false,
		Expires:  time.Now().Local().Add(CookieExpired),
		SameSite: http.SameSiteNoneMode,
	}
	http.SetCookie(w, &authCookie)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}

func (h UsersHandler) GetUserInfo(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(AuthCookieName)

	if err != nil {
		http.Error(w, "need to login", http.StatusUnauthorized)
		return
	}

	user, err := h.Service.GetUserInfo(cookie.Value)
	if err != nil {
		http.Error(w, "failed to get profile", http.StatusBadRequest)
		log.Println(err)
		return
	}

	data, err := json.Marshal(user)
	if err != nil {
		http.Error(w, "failed to create response data", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}

func (h UsersHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "content type not allowed", http.StatusUnsupportedMediaType)
		return
	}

	var request domain.SubscribeRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		log.Println(err)
		return
	}

	cookie, err := r.Cookie(AuthCookieName)

	if err != nil {
		http.Error(w, "need to login", http.StatusUnauthorized)
		return
	}

	err = h.Service.Subscribe(cookie.Value, request.UserId)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrExists):
			http.Error(w, "subscribe exist", http.StatusConflict)
		case errors.Is(err, domain.ErrNotExists):
			http.Error(w, "user doesn't exist", http.StatusConflict)
		default:
			http.Error(w, "unexpected error", http.StatusInternalServerError)
		}
		log.Println(err)

		return
	}

	data, err := json.Marshal(domain.DefaultResponse{Success: true})
	if err != nil {
		http.Error(w, "failed to create response data", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}

func (h UsersHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "content type not allowed", http.StatusUnsupportedMediaType)
		return
	}

	var request domain.SubscribeRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		log.Println(err)
		return
	}

	cookie, err := r.Cookie(AuthCookieName)

	if err != nil {
		http.Error(w, "need to login", http.StatusUnauthorized)
		return
	}

	err = h.Service.Unsubscribe(cookie.Value, request.UserId)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotExists):
			http.Error(w, "user or subscribe doesn't exist", http.StatusConflict)
		default:
			http.Error(w, "unexpected error", http.StatusInternalServerError)
		}
		log.Println(err)

		return
	}

	data, err := json.Marshal(domain.DefaultResponse{Success: true})
	if err != nil {
		http.Error(w, "failed to create response data", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}

func (h UsersHandler) Settings(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "content type not allowed", http.StatusUnsupportedMediaType)
		return
	}

	var request domain.SettingsRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		log.Println(err)
		return
	}

	cookie, err := r.Cookie(AuthCookieName)

	if err != nil {
		http.Error(w, "need to login", http.StatusUnauthorized)
		return
	}

	err = h.Service.Settings(cookie.Value, request.DaysToNotification, request.Email)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrNotExists):
			http.Error(w, "user doesn't exist", http.StatusConflict)
		default:
			http.Error(w, "unexpected error", http.StatusInternalServerError)
		}
		log.Println(err)

		return
	}

	data, err := json.Marshal(domain.DefaultResponse{Success: true})
	if err != nil {
		http.Error(w, "failed to create response data", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}
