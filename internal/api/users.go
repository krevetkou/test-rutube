package api

import (
	"encoding/json"
	"errors"
	"github.com/krevetkou/test-rutube/internal/domain"
	"log"
	"net/http"
)

type UsersService interface {
	Create(actor domain.RegisterRequest) (domain.User, error)
	List() ([]domain.ProfileResponse, error)
	Login(actor domain.LoginRequest) (domain.User, error)
	CreateToken(user domain.User) (string, error)
	GetProfile(token domain.LoginResponse) (domain.ProfileResponse, error)
	Subscribe(tgUser, subscribeUser string) error
	Unsubscribe(tgUser, subscribeUser string) error
}

type UsersHandler struct {
	Service UsersService
}

func NewUsersHandler(service UsersService) UsersHandler {
	return UsersHandler{
		Service: service,
	}
}

func (h UsersHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "content type not allowed", http.StatusUnsupportedMediaType)
		return
	}

	var newUser domain.RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		log.Println(err)
		return
	}

	createdUser, err := h.Service.Create(newUser)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrFieldsRequired):
			http.Error(w, "all required fields must have values", http.StatusUnprocessableEntity)
		case errors.Is(err, domain.ErrExists):
			http.Error(w, "user already exists", http.StatusConflict)
		default:
			http.Error(w, "unexpected error", http.StatusInternalServerError)
		}
		log.Println(err)

		return
	}

	t, err := h.Service.CreateToken(createdUser)
	if err != nil {
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(t))
	if err != nil {
		log.Println(err)
		return
	}
}

func (h UsersHandler) List(w http.ResponseWriter, r *http.Request) {
	users, err := h.Service.List()
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

func (h UsersHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "content type not allowed", http.StatusUnsupportedMediaType)
		return
	}

	var user domain.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		log.Println(err)
		return
	}

	createdUser, err := h.Service.Login(user)
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
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write(data)
	if err != nil {
		log.Println(err)
		return
	}
}

func (h UsersHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")

	regUser, err := h.Service.GetProfile(domain.LoginResponse{AccessToken: token})
	if err != nil {
		http.Error(w, "failed to get profile", http.StatusBadRequest)
		log.Println(err)
		return
	}

	data, err := json.Marshal(regUser)
	if err != nil {
		http.Error(w, "failed to create response data", http.StatusInternalServerError)
		log.Println(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")
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

	var subscribe domain.SubscribeRequest
	err := json.NewDecoder(r.Body).Decode(&subscribe)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		log.Println(err)
		return
	}

	err = h.Service.Subscribe(subscribe.TelegramName, subscribe.SubscribeUser)
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

	text := subscribe.TelegramName + " successfully subscribed to " + subscribe.SubscribeUser

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(text))
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

	var subscribe domain.SubscribeRequest
	err := json.NewDecoder(r.Body).Decode(&subscribe)
	if err != nil {
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		log.Println(err)
		return
	}

	err = h.Service.Unsubscribe(subscribe.TelegramName, subscribe.SubscribeUser)
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

	text := subscribe.TelegramName + " successfully unsubscribed from " + subscribe.SubscribeUser

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(text))
	if err != nil {
		log.Println(err)
		return
	}
}
