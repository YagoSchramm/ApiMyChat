package entity

import "time"

type Message struct {
	ID        string    `json:"id"`
	SenderID  string    `json:"senderId"`
	RoomID    string    `json:"roomId"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	Status    string    `json:"status"`
}
