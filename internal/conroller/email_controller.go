package controller

import (
	"fmt"
	"net/http"

	"github.com/YagoSchramm/ApiMyChat/internal/entity"
	"github.com/YagoSchramm/ApiMyChat/internal/usecase"
	"github.com/gin-gonic/gin"
)

type EmailController struct {
	UseCase *usecase.OTPUseCase
}

func (ctrl *EmailController) RequestOTP(c *gin.Context) {
	var req entity.EmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "E-mail inválido"})
		return
	}

	if err := ctrl.UseCase.ExecuteSend(req.Email); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Falha ao enviar e-mail"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Código enviado!"})
}
func (ctrl *EmailController) VerifyCode(c *gin.Context) {
	var req entity.VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Dados inválidos"})
		return
	}
	savedCode, err := ctrl.UseCase.Repo.GetOTP(req.Email)
	if err != nil || savedCode != req.Code {
		c.JSON(401, gin.H{"error": "Código inválido ou expirado"})
		return
	}

	c.JSON(200, gin.H{"message": "E-mail verificado com sucesso!"})
}
