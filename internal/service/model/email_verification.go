package model

type EmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type VerifyRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}
