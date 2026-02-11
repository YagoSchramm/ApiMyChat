package usecase

import (
	"encoding/json"
	"time"

	"github.com/YagoSchramm/ApiMyChat/internal/entity"
	"github.com/YagoSchramm/ApiMyChat/internal/repository"
	"github.com/google/uuid"
)

type MessageUsecase struct {
	repo repository.MessageRepository
	hub  *entity.Hub
}

func NewMessageUsecase(r repository.MessageRepository, hub *entity.Hub) *MessageUsecase {
	return &MessageUsecase{repo: r, hub: hub}
}
func (u *MessageUsecase) SendMessage(sender, room, content string) {

	msg := entity.Message{
		ID:        uuid.New().String(),
		SenderID:  sender,
		RoomID:    room,
		Content:   content,
		CreatedAt: time.Now(),
	}

	u.repo.Create(&msg)

	jsonMsg, _ := json.Marshal(msg)
	u.hub.Broadcast(room, jsonMsg)
}

func (u *MessageUsecase) GetByRoom(roomID string, limit int) ([]entity.Message, error) {
	return u.repo.GetByRoom(roomID, limit)
}
