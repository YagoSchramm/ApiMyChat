package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email" binding:"required,email"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	CreatedAt   time.Time `json:"createdAt" binding:"required"`
	Password    string    `json:"password" binding:"required"`
}

// UserCredentials are the credentials for authentication
type UserCredentials struct {
	// Name in the credential entity.
	Name string `json:"name"`

	// Email in the credential entity.
	Email string `json:"email" binding:"required,email"`

	// Password in the credential entity.
	Password string `json:"password" binding:"required"`
}

// UpdateUserDTO is the data transfer object for the update user endpoint
type UpdateUserDTO struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
}
