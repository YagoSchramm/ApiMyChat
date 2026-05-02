package repository

import (
	"context"

	"github.com/YagoSchramm/ApiMyChat/domain/entities"
)

type UserRepository interface {
	// GetById attempts to get the user on the database from the given ID
	GetById(
		ctx context.Context,
		id string,
	) (*entities.User, error)

	// GetByEmail attempts to get the user on the database from the given email
	GetByEmail(
		ctx context.Context,
		email string,
	) (*entities.User, error)

	// GetAll returns all users from the database excluding the user with the given ID.
	GetAll(
		ctx context.Context,
		id string,
	) (*[]entities.User, error)

	// UpdateUser attempts to update the username and description
	UpdateUser(
		ctx context.Context,
		userInformation entities.UpdateUserDTO,
	) error
}
