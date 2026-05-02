package impl

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/YagoSchramm/ApiMyChat/domain/entities"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) GetById(ctx context.Context, id string) (*entities.User, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT uid, email, name, COALESCE(description, ''), COALESCE(password, ''), created_at
		 FROM users
		 WHERE uid = $1`,
		id,
	)

	var user entities.User
	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Description,
		&user.Password,
		&user.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("scan user by id: %w", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	row := r.db.QueryRowContext(
		ctx,
		`SELECT uid, email, name, COALESCE(description, ''), COALESCE(password, ''), created_at
		 FROM users
		 WHERE email = $1`,
		email,
	)

	var user entities.User
	if err := row.Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Description,
		&user.Password,
		&user.CreatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, fmt.Errorf("scan user by email: %w", err)
	}

	return &user, nil
}

func (r *PostgresUserRepository) GetAll(ctx context.Context, id string) (*[]entities.User, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT uid, email, name, COALESCE(description, ''), COALESCE(password, ''), created_at
		 FROM users
		 WHERE uid::text <> $1
		 ORDER BY name ASC, created_at DESC`,
		id,
	)
	if err != nil {
		return nil, fmt.Errorf("query users: %w", err)
	}
	defer rows.Close()

	users := make([]entities.User, 0)
	for rows.Next() {
		var user entities.User
		if err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.Name,
			&user.Description,
			&user.Password,
			&user.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan user list item: %w", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate users: %w", err)
	}

	return &users, nil
}

func (r *PostgresUserRepository) UpdateUser(ctx context.Context, userInformation entities.UpdateUserDTO) error {
	result, err := r.db.ExecContext(
		ctx,
		`UPDATE users
		 SET name = $1, description = $2
		 WHERE uid = $3`,
		userInformation.Name,
		userInformation.Description,
		userInformation.ID,
	)
	if err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check updated user rows: %w", err)
	}

	if affected == 0 {
		return fmt.Errorf("update user: %w", sql.ErrNoRows)
	}

	return nil
}
