package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/YagoSchramm/ApiMyChat/internal/usecase"
	"github.com/gin-gonic/gin"
)

type FCMController struct {
	Usecase     *usecase.FCMUsecase
	AuthUsecase *usecase.AuthUsecase
}

type saveFCMTokenRequest struct {
	Uid      string `json:"uid"`
	ID       string `json:"id"`
	Token    string `json:"token"`
	FCMToken string `json:"fcmToken"`
}

type deleteFCMTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

func NewFCMController(fcmUsecase *usecase.FCMUsecase, authUsecase *usecase.AuthUsecase) *FCMController {
	return &FCMController{
		Usecase:     fcmUsecase,
		AuthUsecase: authUsecase,
	}
}

func (c *FCMController) SaveToken(ctx *gin.Context) {
	var req saveFCMTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	uid := strings.TrimSpace(req.Uid)
	if uid == "" {
		uid = strings.TrimSpace(req.ID)
	}
	if uid == "" {
		uid = strings.TrimSpace(ctx.Query("id"))
	}

	token := strings.TrimSpace(req.Token)
	if token == "" {
		token = strings.TrimSpace(req.FCMToken)
	}

	if uid == "" || token == "" {
		ctx.JSON(http.StatusBadRequest, "uid and token required")
		return
	}

	if err := c.Usecase.SaveToken(uid, token); err != nil {
		fmt.Print(err.Error())
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "saved"})
}

func (c *FCMController) DeleteToken(ctx *gin.Context) {
	var req deleteFCMTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	if err := c.Usecase.DeleteToken(req.Token); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "deleted"})
}
