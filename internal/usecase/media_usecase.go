package usecase

import (
	"time"

	"github.com/YagoSchramm/ApiMyChat/internal/entity"
	"github.com/YagoSchramm/ApiMyChat/internal/repository"
	"github.com/google/uuid"
)

type MediaUsecase struct {
	mediaRepo   repository.MediaRepository
	messageRepo repository.MessageRepository
}

func NewMediaUsecase(mediaRepo repository.MediaRepository, messageRepo repository.MessageRepository) MediaUsecase {
	return MediaUsecase{
		mediaRepo:   mediaRepo,
		messageRepo: messageRepo,
	}
}

func (u *MediaUsecase) Create(userID, roomID, url, mediaType string) (string, error) {
	msg := entity.Message{
		ID:        uuid.New().String(),
		SenderID:  userID,
		RoomID:    roomID,
		Content:   "",
		CreatedAt: time.Now(),
	}

	if err := u.messageRepo.Create(&msg); err != nil {
		return "", err
	}

	if err := u.mediaRepo.Create(userID, msg.ID, url, mediaType); err != nil {
		return "", err
	}

	return msg.ID, nil
}

func (u *MediaUsecase) GetByMessageId(messageID string) ([]string, error) {
	return u.mediaRepo.GetByMessageId(messageID)
}
