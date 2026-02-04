package entity

import "time"

type Message struct {
	ID        string
	SenderID  string
	RoomID    string
	Content   string
	CreatedAt time.Time
}
