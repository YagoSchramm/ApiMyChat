package repository

import (
	"database/sql"
)

type MediaRepository struct {
	connection *sql.DB
}

func NewMediaRepository(db *sql.DB) MediaRepository {
	return MediaRepository{connection: db}
}
func (r *MediaRepository) Create(userId, messageId string, url string, mediaType string) error {
	query := `
	INSERT INTO medias (uid,message_id, url, type, created_at)
	VALUES ($1, $2, $3, $4, NOW())`
	_, err := r.connection.Exec(query, userId, messageId, url, mediaType)
	return err
}
func (r *MediaRepository) GetByMessageId(messageId string) ([]string, error) {
	query := `
	SELECT url
	FROM medias
	WHERE message_id = $1`
	rows, err := r.connection.Query(query, messageId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls := make([]string, 0)
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return urls, nil
}
