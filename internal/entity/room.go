package entity

import "time"

type Room struct {
	Name           string
	ID             string
	RoomClient     RoomClients
	CreatedAt      time.Time
	LastConnection map[string]time.Time
}
