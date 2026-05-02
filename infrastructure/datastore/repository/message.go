package repository

import (
	"context"

	"github.com/YagoSchramm/ApiMyChat/domain/entities"
)

type MessageRepository interface {
	Create(
		ctx context.Context,
		message entities.Message,
	) error

	GetByRoomID(
		ctx context.Context,
		roomID string,
		limit int,
	) (*[]entities.Message, error)

	GetLastByRoomID(
		ctx context.Context,
		roomID string,
	) (*entities.Message, error)
}
