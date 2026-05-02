package impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/YagoSchramm/ApiMyChat/domain/entities"
	"github.com/lib/pq"
)

type PostgresMessageRepository struct {
	db *sql.DB
}

func NewPostgresMessageRepository(db *sql.DB) *PostgresMessageRepository {
	return &PostgresMessageRepository{db: db}
}

func (r *PostgresMessageRepository) Create(ctx context.Context, message entities.Message) error {
	createdAt := message.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now().UTC()
	}

	_, err := r.db.ExecContext(
		ctx,
		`INSERT INTO messages (id, sender_id, room_id, content, status, created_at)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		message.ID.String(),
		message.SenderID,
		message.RoomID,
		message.Content,
		message.Status,
		createdAt,
	)
	if err != nil {
		return fmt.Errorf("insert message: %w", err)
	}

	return nil
}

func (r *PostgresMessageRepository) GetByRoomID(ctx context.Context, roomID string, limit int) (*[]entities.Message, error) {
	if limit <= 0 {
		limit = 50
	}

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT
			m.id,
			m.sender_id,
			m.room_id,
			COALESCE(u.name, '') AS username,
			m.content,
			COALESCE(array_agg(media.url ORDER BY media.created_at ASC) FILTER (WHERE media.url IS NOT NULL), '{}'::text[]) AS media_urls,
			m.created_at,
			m.status
		FROM messages m
		LEFT JOIN users u ON u.uid::text = m.sender_id
		LEFT JOIN medias media ON media.message_id::text = m.id
		WHERE m.room_id = $1
		GROUP BY m.id, m.sender_id, m.room_id, u.name, m.content, m.created_at, m.status
		ORDER BY m.created_at ASC, m.id ASC
		LIMIT $2`,
		roomID,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("query room messages: %w", err)
	}
	defer rows.Close()

	messages := make([]entities.Message, 0)
	for rows.Next() {
		var (
			message   entities.Message
			mediaURLs pq.StringArray
		)

		if err := rows.Scan(
			&message.ID,
			&message.SenderID,
			&message.RoomID,
			&message.UserName,
			&message.Content,
			pq.Array(&mediaURLs),
			&message.CreatedAt,
			&message.Status,
		); err != nil {
			return nil, fmt.Errorf("scan room message: %w", err)
		}

		message.MediaURLs = []string(mediaURLs)
		messages = append(messages, message)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate room messages: %w", err)
	}

	return &messages, nil
}

func (r *PostgresMessageRepository) GetLastByRoomID(ctx context.Context, roomID string) (*entities.Message, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT
			m.id,
			m.sender_id,
			m.room_id,
			COALESCE(u.name, '') AS username,
			m.content,
			COALESCE(array_agg(media.url ORDER BY media.created_at ASC) FILTER (WHERE media.url IS NOT NULL), '{}'::text[]) AS media_urls,
			m.created_at,
			m.status
		FROM messages m
		LEFT JOIN users u ON u.uid::text = m.sender_id
		LEFT JOIN medias media ON media.message_id::text = m.id
		WHERE m.room_id = $1
		GROUP BY m.id, m.sender_id, m.room_id, u.name, m.content, m.created_at, m.status
		ORDER BY m.created_at DESC, m.id DESC
		LIMIT 1`,
		roomID,
	)

	var (
		message   entities.Message
		mediaURLs pq.StringArray
	)

	if err := row.Scan(
		&message.ID,
		&message.SenderID,
		&message.RoomID,
		&message.UserName,
		&message.Content,
		pq.Array(&mediaURLs),
		&message.CreatedAt,
		&message.Status,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("scan last room message: %w", err)
	}

	message.MediaURLs = []string(mediaURLs)
	return &message, nil
}
