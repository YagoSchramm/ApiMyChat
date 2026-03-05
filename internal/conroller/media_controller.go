package controller

import (
	"net/http"
	"strings"

	"github.com/YagoSchramm/ApiMyChat/internal/usecase"
	"github.com/gin-gonic/gin"
)

type MediaController struct {
	Usecase usecase.MediaUsecase
}

type createMediaRequest struct {
	UserID    string `json:"userId"`
	UID       string `json:"uid"`
	RoomID    string `json:"roomId" binding:"required"`
	URL       string `json:"url" binding:"required"`
	Type      string `json:"type" binding:"required"`
	MediaType string `json:"mediaType"`
}

func (mctrl *MediaController) Create(ctx *gin.Context) {
	var req createMediaRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	userID := strings.TrimSpace(req.UserID)
	if userID == "" {
		userID = strings.TrimSpace(req.UID)
	}

	mediaType := strings.TrimSpace(req.Type)
	if mediaType == "" {
		mediaType = strings.TrimSpace(req.MediaType)
	}

	if userID == "" || req.RoomID == "" || req.URL == "" || mediaType == "" {
		ctx.JSON(http.StatusBadRequest, "userId, roomId, url and type required")
		return
	}

	messageID, err := mctrl.Usecase.Create(userID, req.RoomID, req.URL, mediaType)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "created", "messageId": messageID})
}

func (mctrl *MediaController) GetByMessageId(ctx *gin.Context) {
	messageID := strings.TrimSpace(ctx.Param("messageId"))
	if messageID == "" {
		messageID = strings.TrimSpace(ctx.Query("messageId"))
	}
	if messageID == "" {
		ctx.JSON(http.StatusBadRequest, "messageId required")
		return
	}

	mediaURLs, err := mctrl.Usecase.GetByMessageId(messageID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, mediaURLs)
}
