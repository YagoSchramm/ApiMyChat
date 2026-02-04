package usecase

import (
	"fmt"

	"github.com/YagoSchramm/ApiMyChat/internal/entity"
	"github.com/YagoSchramm/ApiMyChat/internal/repository"
)

type UserUseCase struct {
	urepo repository.UserRepository
}

func NewUserUseCase(ur repository.UserRepository) UserUseCase {
	return UserUseCase{
		urepo: ur,
	}
}
func (uc *UserUseCase) CreateUser(user entity.User) (entity.User, error) {
	user1, err := uc.urepo.CreateUser(user)
	if err != nil {
		fmt.Println(err)
		return entity.User{}, nil
	}
	return user1, nil
}
func (uc *UserUseCase) GetByID(id string) (entity.User, error) {
	user, err := uc.urepo.GetByID(id)
	if err != nil {
		fmt.Println(err)
		return entity.User{}, nil
	}
	return user, nil
}
func (uc *UserUseCase) GetAll(id string) ([]entity.User, error) {
	return uc.urepo.GetAll(id)
}
