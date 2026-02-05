package repository

import (
	"database/sql"

	"github.com/YagoSchramm/ApiMyChat/internal/entity"
)

type MessageRepository struct {
	connection *sql.DB
}

func NewMessageRepository(db *sql.DB) MessageRepository {
	return MessageRepository{connection: db}
}

func (r *MessageRepository) Create(msg *entity.Message) error {
	query := `
		INSERT INTO messages (id, room_id, sender_id, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.connection.Exec(
		query,
		msg.ID,
		msg.RoomID,
		msg.SenderID,
		msg.Content,
		msg.CreatedAt,
	)

	return err
}

func (r *MessageRepository) GetByRoom(roomID string, limit int) ([]entity.Message, error) {
	query := `
		SELECT id, room_id, sender_id, content, created_at
		FROM messages
		WHERE room_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.connection.Query(query, roomID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var messages []entity.Message

	for rows.Next() {
		var msg entity.Message
		if err := rows.Scan(
			&msg.ID,
			&msg.RoomID,
			&msg.SenderID,
			&msg.Content,
			&msg.CreatedAt,
		); err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}

	return messages, nil
}
