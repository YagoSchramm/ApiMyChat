package usecase

import (
	"encoding/json"
	"log"
	"time"

	"github.com/YagoSchramm/ApiMyChat/internal/entity"
	"github.com/YagoSchramm/ApiMyChat/internal/repository"
	"github.com/google/uuid"
)

type MessageUsecase struct {
	repo        repository.MessageRepository
	mediaRepo   repository.MediaRepository
	hub         *entity.Hub
	roomUsecase *RoomUsecase
	fcmUsecase  *FCMUsecase
}

func NewMessageUsecase(
	r repository.MessageRepository,
	mediaRepo repository.MediaRepository,
	hub *entity.Hub,
	roomUsecase *RoomUsecase,
	fcmUsecase *FCMUsecase,
) *MessageUsecase {
	return &MessageUsecase{
		repo:        r,
		mediaRepo:   mediaRepo,
		hub:         hub,
		roomUsecase: roomUsecase,
		fcmUsecase:  fcmUsecase,
	}
}

func (u *MessageUsecase) SendMessage(sender, room, content string) {

	msg := entity.Message{
		ID:        uuid.New().String(),
		SenderID:  sender,
		RoomID:    room,
		Content:   content,
		CreatedAt: time.Now(),
	}

	if err := u.repo.Create(&msg); err != nil {
		log.Printf("failed to persist message %s: %v", msg.ID, err)
		return
	}

	jsonMsg, err := json.Marshal(msg)
	if err != nil {
		log.Printf("failed to encode message %s: %v", msg.ID, err)
		return
	}

	u.hub.Broadcast(room, jsonMsg)
	u.notifyOfflineParticipants(msg)
}

func (u *MessageUsecase) GetByRoom(roomID string, limit int) ([]entity.Message, error) {
	messages, err := u.repo.GetByRoom(roomID, limit)
	if err != nil {
		return nil, err
	}

	for i := range messages {
		urls, err := u.mediaRepo.GetByMessageId(messages[i].ID)
		if err != nil {
			return nil, err
		}
		messages[i].MediaURLs = urls
	}

	return messages, nil
}
func (u *MessageUsecase) GetLastByRoom(roomID string) (*entity.Message, error) {
	msg, err := u.repo.LastByRoom(roomID)
	if err != nil {
		return nil, err
	}

	urls, err := u.mediaRepo.GetByMessageId(msg.ID)
	if err != nil {
		return nil, err
	}
	msg.MediaURLs = urls

	return msg, nil
}
func (u *MessageUsecase) UpdateStatus(messageID, status string) error {
	return u.repo.UpdateStatus(messageID, status)
}
func (u *MessageUsecase) notifyOfflineParticipants(msg entity.Message) {
	if u.roomUsecase == nil || u.fcmUsecase == nil {
		return
	}

	room, err := u.roomUsecase.GetRoomById(msg.RoomID)
	if err != nil {
		log.Printf("failed to load room %s users for fcm: %v", msg.RoomID, err)
		return
	}

	offlineUsers := make([]string, 0, len(room.Users))
	for _, user := range room.Users {
		if user.UID == msg.SenderID {
			continue
		}
		if u.hub.IsOnline(user.UID) {
			continue
		}
		offlineUsers = append(offlineUsers, user.UID)
	}

	if len(offlineUsers) == 0 {
		log.Printf("fcm skipped for message %s: no offline participants in room %s", msg.ID, msg.RoomID)
		return
	}

	log.Printf("fcm dispatch for message %s: room=%s offline_users=%d", msg.ID, msg.RoomID, len(offlineUsers))
	if err := u.fcmUsecase.NotifyUsers(
		offlineUsers,
		"Nova mensagem",
		msg.Content,
		map[string]string{
			"room_id":    msg.RoomID,
			"sender_id":  msg.SenderID,
			"message_id": msg.ID,
		},
	); err != nil {
		log.Printf("failed to notify offline users via fcm: %v", err)
		return
	}

	log.Printf("fcm sent successfully for message %s", msg.ID)
}
