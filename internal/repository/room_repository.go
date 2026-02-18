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

func (r *RoomRepository) Create(room *entity.Room, userIDs []string) error {
	tx, err := r.connection.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	roomQuery := `
	INSERT INTO rooms (name,id, created_at)
	VALUES ($1, $2,$3)`
	_, err = tx.Exec(roomQuery, room.Name, room.ID, room.CreatedAt)
	if err != nil {
		return err
	}

	userQuery := `
	INSERT INTO room_users (room_id, user_id, joined_at, left_at)
	VALUES ($1, $2, $3, $4)`

	for _, userID := range userIDs {
		_, err = tx.Exec(userQuery, room.ID, userID, room.CreatedAt, nil)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
func (r *RoomRepository) GetRoomsByUid(uid string) ([]entity.RoomWithUsers, error) {
	roomQuery := `
	SELECT
		r.id,
		r.name,
		r.created_at,
		u.uid,
		u.email,
		u.name,
		u.description,
		u.created_at
	FROM rooms r
	INNER JOIN room_users ru_filter ON ru_filter.room_id = r.id AND ru_filter.user_id = $1
	INNER JOIN room_users ru ON ru.room_id = r.id
	INNER JOIN users u ON u.uid = ru.user_id
	ORDER BY r.created_at DESC`

	rows, err := r.connection.Query(roomQuery, uid)
	if err != nil {
		return []entity.RoomWithUsers{}, err
	}
	defer rows.Close()

	rooms := make([]entity.RoomWithUsers, 0)
	roomIndex := make(map[string]int)

	for rows.Next() {
		var roomID string
		var user entity.User
		var roomObj entity.RoomWithUsers
		err = rows.Scan(
			&roomID,
			&roomObj.Name,
			&roomObj.CreatedAt,
			&user.UID,
			&user.Email,
			&user.Name,
			&user.Description,
			&user.CreatedAt,
		)
		if err != nil {
			return []entity.RoomWithUsers{}, err
		}

		idx, ok := roomIndex[roomID]
		if !ok {
			roomObj.ID = roomID
			roomObj.Users = make([]entity.User, 0, 2)
			rooms = append(rooms, roomObj)
			idx = len(rooms) - 1
			roomIndex[roomID] = idx
		}

		rooms[idx].Users = append(rooms[idx].Users, user)
	}

	if err = rows.Err(); err != nil {
		return []entity.RoomWithUsers{}, err
	}

	return rooms, nil
}
func (r *RoomRepository) GetRoomById(uid string) (entity.RoomWithUsers, error) {
	roomQuery := `
	SELECT
		r.id,
		r.name,
		r.created_at,
		u.uid,
		u.email,
		u.name,
		u.description,
		u.created_at
	FROM rooms r
	INNER JOIN room_users ru ON ru.room_id = r.id
	INNER JOIN users u ON u.uid = ru.user_id
	WHERE r.id = $1`

	rows, err := r.connection.Query(roomQuery, uid)
	if err != nil {
		return entity.RoomWithUsers{}, err
	}
	defer rows.Close()

	room := entity.RoomWithUsers{
		Users: make([]entity.User, 0, 2),
	}
	found := false
	for rows.Next() {
		var user entity.User
		err = rows.Scan(
			&room.ID,
			&room.Name,
			&room.CreatedAt,
			&user.UID,
			&user.Email,
			&user.Name,
			&user.Description,
			&user.CreatedAt,
		)
		if err != nil {
			return entity.RoomWithUsers{}, err
		}
		found = true
		room.Users = append(room.Users, user)
	}

	if err = rows.Err(); err != nil {
		return entity.RoomWithUsers{}, err
	}
	if !found {
		return entity.RoomWithUsers{}, sql.ErrNoRows
	}

	return room, nil
}
