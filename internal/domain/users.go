package domain

import (
	_ "github.com/sirupsen/logrus"
)

type User struct {
	ID                 int
	Email              string
	Password           string
	Name               string
	DateOfBirth        string
	Token              string
	DaysToNotification int
	SubscribeUsers     []int
}

type RegisterRequest struct {
	Email       string `json:"email"`
	Name        string `json:"name"`
	Password    string `json:"password"`
	DateOfBirth string `json:"dateOfBirth"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	Email              string `json:"email"`
	Name               string `json:"name"`
	DaysToNotification int    `json:"daysToNotification"`
}

type UserInListResponse struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type ProfileResponse struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	Name         string `json:"name"`
	IsSubscribed bool   `json:"isSubscribed"`
}

type SubscribeRequest struct {
	UserId int `json:"userId"`
}

type SettingsRequest struct {
	DaysToNotification int    `json:"daysToNotification"`
	Email              string `json:"email,omitempty"`
}

type DefaultResponse struct {
	Success bool `json:"success"`
}
