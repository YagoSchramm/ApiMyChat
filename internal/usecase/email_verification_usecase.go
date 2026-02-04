package usecase

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/YagoSchramm/ApiMyChat/internal/repository"
)

type EmailService interface {
	Send(to, body string) error
}

type OTPUseCase struct {
	Repo  *repository.MemoryCache
	Email EmailService
}

func (u *OTPUseCase) ExecuteSend(email string) error {
	// Gerar código
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// Salvar (Ex: no Redis ou Map) com expiração de 10 min
	err := u.Repo.SaveOTP(email, code, 10*time.Minute)
	if err != nil {
		return err
	}

	// Enviar e-mail
	return u.Email.Send(email, "Seu código é: "+code)
}
