package repository

import (
	"database/sql"
	"fmt"

	"github.com/YagoSchramm/ApiMyChat/internal/entity"
)

type UserRepository struct {
	connection *sql.DB
}

func NewUserRepository(connection *sql.DB) UserRepository {
	return UserRepository{connection: connection}
}
func (ur *UserRepository) GetByID(id string) (entity.User, error) {
	user := entity.User{}
	err := ur.connection.QueryRow("SELECT * FROM users WHERE uid = $1", id).Scan(&user.UID, &user.Email, &user.Name)
	return user, err
}

func (ur *UserRepository) CreateUser(user entity.User) (entity.User, error) {
	var id int
	query, err := ur.connection.Prepare("INSERT INTO users (uid,name,email,createdAt,description) VALUES ($1, $2,$3,$4) ON CONFLICT DO NOTHING RETURNING id")
	if err != nil {
		fmt.Println(err)
		return entity.User{}, err
	}
	err = query.QueryRow(user.UID, user.Name, user.Email, user.CreatedAt, user.Description).Scan(&id)
	if err != nil {
		fmt.Printf("Usuário %d cadastrado com sucesso!", id)
		return user, nil
	}
	return entity.User{}, nil
}
func (ur *UserRepository) GetAll(id string) ([]entity.User, error) {
	query := "select * from users"
	rows, err := ur.connection.Query(query)
	if err != nil {
		fmt.Println(err)
		return []entity.User{}, nil
	}
	var userList []entity.User
	var userObj entity.User
	for rows.Next() {
		err = rows.Scan(
			&userObj.UID,
			&userObj.Name,
			&userObj.Email,
			&userObj.CreatedAt,
			&userObj.Description,
		)
		if err != nil {
			fmt.Println(err)
			return []entity.User{}, nil
		}
		if userObj.UID != id {
			userList = append(userList, userObj)
		}
	}
	return userList, nil
}
