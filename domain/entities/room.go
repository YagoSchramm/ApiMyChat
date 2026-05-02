package entities

import "time"

type RoomClients map[string]any

type Room struct {
	Name           string
	ID             string
	RoomClient     RoomClients
	CreatedAt      time.Time
	LastConnection map[string]time.Time
}
