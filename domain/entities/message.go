package entities

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID        uuid.UUID `json:"id"`
	SenderID  string    `json:"senderId"`
	RoomID    string    `json:"roomId"`
	UserName  string    `json:"username"`
	Content   string    `json:"content"`
	MediaURLs []string  `json:"mediaUrls,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	Status    string    `json:"status"`
}
