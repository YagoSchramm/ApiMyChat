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
		return entity.User{}, err
	}
	return user1, nil
}
func (uc *UserUseCase) GetByID(id string) (entity.User, error) {
	user, err := uc.urepo.GetByID(id)
	if err != nil {
		fmt.Println(err)
		return entity.User{}, err
	}
	return user, nil
}
func (uc *UserUseCase) UpdateUser(user entity.UpdateUserModel) (entity.UpdateUserModel, error) {
	user1, err := uc.urepo.UpdateUser(user)
	if err != nil {
		fmt.Println(err)
		return entity.UpdateUserModel{}, err
	}
	return user1, nil
}

func (uc *UserUseCase) GetByEmail(email string) (entity.User, error) {
	user, err := uc.urepo.GetByEmail(email)
	if err != nil {
		fmt.Println(err)
		return entity.User{}, err
	}
	return user, nil
}
func (uc *UserUseCase) GetAll(id string) ([]entity.User, error) {
	return uc.urepo.GetAll(id)
}
