package repository

import "database/sql"

type MessageRepository struct {
	connection *sql.DB
}

func NewMessageRepository(connection *sql.DB) MessageRepository {
	return MessageRepository{connection: connection}
}
