package domain

import "time"

type User struct {
	Id           int       `json:"-"`
	Name         string    `json:"name" binding:"required"`
	Email        string    `json:"email" binding:"required"`
	PasswordHash string    `json:"password_hash" binding:"required"`
	Registered   time.Time `json:"-"`
}

type SignInInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
