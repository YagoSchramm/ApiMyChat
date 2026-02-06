package entity

import (
	"time"
)

type User struct {
	UID         string    `json:"id"`
	Email       string    `json:"email" binding:"required,email"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	CreatedAt   time.Time `json:"createdAt" binding:"required"`
	Password    string    `json:"password" binding:"required"`
}
