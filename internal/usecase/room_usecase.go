package usecase

import (
	"fmt"
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

func (u *RoomUsecase) CreateRoom(name string, userIDs []string) (entity.Room, error) {
	participants := uniqueNonEmptyUsers(userIDs...)
	if len(participants) < 2 {
		return entity.Room{}, fmt.Errorf("room requires two valid users")
	}

	room := entity.Room{
		Name:      name,
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
	}

	if err := u.repo.Create(&room, participants); err != nil {
		return entity.Room{}, err
	}

	return room, nil
}
func (u *RoomUsecase) GetRoomById(uid string) (entity.RoomWithUsers, error) {
	room, err := u.repo.GetRoomById(uid)
	if err != nil {
		return entity.RoomWithUsers{}, err
	}
	return room, nil
}
func (u *RoomUsecase) GetRoomsByUid(uid string) ([]entity.RoomWithUsers, error) {
	rooms, err := u.repo.GetRoomsByUid(uid)
	if err != nil {
		return []entity.RoomWithUsers{}, err
	}
	return rooms, nil
}
func uniqueNonEmptyUsers(userIDs ...string) []string {
	seen := make(map[string]struct{}, len(userIDs))
	unique := make([]string, 0, len(userIDs))

	for _, userID := range userIDs {
		if userID == "" {
			continue
		}
		if _, ok := seen[userID]; ok {
			continue
		}
		seen[userID] = struct{}{}
		unique = append(unique, userID)
	}

	return unique
}
