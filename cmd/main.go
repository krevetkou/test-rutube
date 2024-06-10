package main

import (
	"errors"
	"github.com/bxcodec/faker/v3"
	"github.com/go-chi/chi/v5"
	"github.com/krevetkou/test-rutube/internal/api"
	"github.com/krevetkou/test-rutube/internal/domain"
	"github.com/krevetkou/test-rutube/internal/services"
	"github.com/krevetkou/test-rutube/internal/storage"
	"github.com/rs/cors"
	"log"
	"math/rand/v2"
	"net/http"
	"time"
)

func main() {
	usersStorage := storage.NewStorage()
	userService := services.NewUserService(usersStorage)
	userHandler := api.NewUsersHandler(userService)

	insertUsers(usersStorage)

	r := chi.NewRouter()
	r.Route("/user", func(r chi.Router) {
		r.Get("/list-today", userHandler.ListToday)
		r.Get("/list", userHandler.List)
		r.Get("/info", userHandler.GetUserInfo)
		r.Post("/login", userHandler.Login)
		r.Post("/subscribe", userHandler.Subscribe)
		r.Post("/unsubscribe", userHandler.Unsubscribe)
		r.Post("/settings", userHandler.Settings)
	})

	handler := cors.Default().Handler(r)
	err := http.ListenAndServe(":8080", handler)
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

	for i := 0; i < rand.IntN(10)+5; i++ {
		users = append(users, domain.RegisterRequest{
			Email:       faker.Email(),
			Name:        faker.Name(),
			Password:    faker.Password(),
			DateOfBirth: faker.Date(),
		})
	}
	users = append(users, domain.RegisterRequest{
		Email:       faker.Email(),
		Name:        faker.Name(),
		Password:    faker.Password(),
		DateOfBirth: time.Now().Format("2006-01-02"),
	})
	users = append(users, domain.RegisterRequest{
		Email:       "test@test.ru",
		Name:        faker.Name(),
		Password:    "testtest",
		DateOfBirth: time.Now().Format("2006-01-02"),
	})

	for _, user := range users {
		_, err := storage.InsertUser(user)
		if err != nil {
			log.Printf("insert user error: %s", err)
		}
	}
}
