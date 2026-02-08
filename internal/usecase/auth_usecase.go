package usecase

import (
	"github.com/YagoSchramm/ApiMyChat/internal/service"
)

type AuthUsecase struct {
	Supabase *service.SupabaseService
}

func NewAuthUsecase(s *service.SupabaseService) *AuthUsecase {
	return &AuthUsecase{Supabase: s}
}

func (u *AuthUsecase) Login(email, password string) (string, error) {

	res, err := u.Supabase.Login(email, password)
	if err != nil {
		return "", err
	}

	return res.AccessToken, nil
}
func (u *AuthUsecase) VerifyJWT(id string) (string, error) {
	res, err := u.Supabase.UserExists(id)

	if err != nil {
		return "", err
	}
	if res == true {
		return id, nil
	}
	return "", err
}
