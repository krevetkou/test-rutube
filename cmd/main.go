package main

import (
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/krevetkou/test-rutube/internal/api"
	"github.com/krevetkou/test-rutube/internal/domain"
	"github.com/krevetkou/test-rutube/internal/services"
	"github.com/krevetkou/test-rutube/internal/storage"
	"log"
	"net/http"
)

func main() {
	usersStorage := storage.NewStorage()
	userService := services.NewUserService(usersStorage)
	userHandler := api.NewUsersHandler(userService)

	insertUsers(usersStorage)

	r := chi.NewRouter()
	r.Route("/user", func(r chi.Router) {
		r.Post("/register", userHandler.Create)
		r.Get("/list", userHandler.List)
		r.Post("/login", userHandler.Login)
		r.Get("/profile", userHandler.GetProfile)
		r.Post("/subscribe", userHandler.Subscribe)
		r.Post("/unsubscribe", userHandler.Unsubscribe)
	})

	err := http.ListenAndServe(":8080", r)
	if errors.Is(err, http.ErrServerClosed) {
		log.Println("server closed")
		return
	}
	if err != nil {
		log.Printf("server error: %s", err)
	}
}

func insertUsers(storage *storage.Storage) {
	users := make([]domain.RegisterRequest, 0)
	users = append(users, domain.RegisterRequest{
		Email:        "l@mail.ru",
		Name:         "l",
		Password:     "lll",
		DateOfBirth:  "01.01.99",
		TelegramName: "l",
	}, domain.RegisterRequest{
		Email:        "m@mail.ru",
		Name:         "m",
		Password:     "mmm",
		DateOfBirth:  "03.12.66",
		TelegramName: "m",
	},
		domain.RegisterRequest{
			Email:        "n@mail.ru",
			Name:         "n",
			Password:     "nnn",
			DateOfBirth:  "30.07.02",
			TelegramName: "n",
		},
		domain.RegisterRequest{
			Email:        "o@mail.ru",
			Name:         "o",
			Password:     "ooo",
			DateOfBirth:  "09.06.10",
			TelegramName: "o",
		},
		domain.RegisterRequest{
			Email:        "p@mail.ru",
			Name:         "p",
			Password:     "ppp",
			DateOfBirth:  "02.10.75",
			TelegramName: "p",
		})

	for _, user := range users {
		_, err := storage.InsertUser(user)
		if err != nil {
			log.Printf("insert user error: %s", err)
		}

		_, err = storage.CreateToken(user.TelegramName)
		if err != nil {
			log.Printf("create token error: %s", err)
		}
	}
}
