package repository

import (
	"database/sql"
	"time"
)

type FcmRepository struct {
	connection *sql.DB
}

func NewFcmRepository(db *sql.DB) FcmRepository {
	return FcmRepository{connection: db}
}
func (r *FcmRepository) Create(uid string, token string, createdAt time.Time) error {
	query := `
	INSERT INTO user_devices (user_id, fcm_token, created_at)
	SELECT $1, $2, $3
	WHERE NOT EXISTS (
		SELECT 1
		FROM user_devices
		WHERE user_id = $1 AND fcm_token = $2
	)`
	_, err := r.connection.Exec(query, uid, token, createdAt)
	return err
}

func (r *FcmRepository) GetTokensByUid(uid string) ([]string, error) {
	query := `
	SELECT fcm_token
	FROM user_devices
	WHERE user_id = $1`

	rows, err := r.connection.Query(query, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tokens := make([]string, 0, 2)
	for rows.Next() {
		var token string
		if err := rows.Scan(&token); err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tokens, nil
}

func (r *FcmRepository) Delete(token string) error {
	query := `
	DELETE FROM user_devices
	WHERE fcm_token = $1`
	_, err := r.connection.Exec(query, token)
	return err
}
