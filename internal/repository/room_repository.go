package repository

import (
	"database/sql"

	"github.com/YagoSchramm/ApiMyChat/internal/entity"
)

type RoomRepository struct {
	connection *sql.DB
}

func NewRoomRepository(db *sql.DB) RoomRepository {
	return RoomRepository{connection: db}
}

func (r *RoomRepository) Create(room *entity.Room) error {
	query := `
	INSERT INTO rooms (id, user_a_id, user_b_id, created_at)
	VALUES ($1,$2,$3,$4)`

	_, err := r.connection.Exec(
		query,
		room.ID,
		room.UserAID,
		room.UserBID,
		room.CreatedAt,
	)
	return err
}
