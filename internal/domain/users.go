package domain

import (
	_ "github.com/sirupsen/logrus"
)

type User struct {
	ID             int
	Email          string
	Password       string
	Name           string
	DateOfBirth    string
	Token          string
	TelegramName   string
	SubscribeUsers []string
}

type RegisterRequest struct {
	Email        string `json:"email"`
	Name         string `json:"name"`
	Password     string `json:"password"`
	TelegramName string `json:"telegram_name"`
	DateOfBirth  string `json:"date_of_birth"`
}

type LoginRequest struct {
	Email        string `json:"email"`
	Password     string `json:"password"`
	TelegramName string `json:"telegram_name"`
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
}

type ProfileResponse struct {
	ID           int    `json:"id"`
	TelegramName string `json:"telegram_name"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	DateOfBirth  string `json:"date_of_birth"`
}

type SubscribeRequest struct {
	TelegramName  string `json:"telegram_name"`
	SubscribeUser string `json:"subscribe_user"`
}
