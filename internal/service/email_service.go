package service

import (
	"gopkg.in/gomail.v2"
)

type GmailService struct {
	Email    string
	Password string
}

func (s *GmailService) Send(to, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.Email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", "Verificação de E-mail")
	m.SetBody("text/html", "<h3>Seu código de verificação para o MyChat é:</h3><h1>"+body+"</h1>")

	d := gomail.NewDialer("smtp.gmail.com", 587, s.Email, s.Password)

	return d.DialAndSend(m)
}
