package entity

import "time"

type RoomWithUsers struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	Users     []User    `json:"users"`
}
