package usecase

import (
	"time"

	"github.com/YagoSchramm/ApiMyChat/internal/entity"
	"github.com/YagoSchramm/ApiMyChat/internal/repository"
	"github.com/google/uuid"
)

type RoomUsecase struct {
	repo repository.RoomRepository
}

func NewRoomUsecase(r repository.RoomRepository) RoomUsecase {
	return RoomUsecase{repo: r}
}

func (u *RoomUsecase) CreateRoom(userAID, userBID string) (entity.Room, error) {
	room := entity.Room{
		ID:        uuid.New().String(),
		UserAID:   userAID,
		UserBID:   userBID,
		CreatedAt: time.Now(),
	}

	if err := u.repo.Create(&room); err != nil {
		return entity.Room{}, err
	}

	return room, nil
}
