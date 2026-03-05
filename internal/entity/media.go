package entity

import "time"

type Media struct {
	ID        string    `json:"id"`
	URL       string    `json:"url"`
	Type      string    `json:"type"`
	MessageID string    `json:"message_id"`
	UserID    string    `json:"user_id"`
	createdAt time.Time `json:"created_at"`
}
