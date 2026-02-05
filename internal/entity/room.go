package entity

import "time"

type Room struct {
	ID        string
	UserAID   string
	UserBID   string
	CreatedAt time.Time
}
