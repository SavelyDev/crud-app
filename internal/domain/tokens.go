package domain

import "time"

type RefreshSession struct {
	Id        int       `json:"id"`
	UserId    int       `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
