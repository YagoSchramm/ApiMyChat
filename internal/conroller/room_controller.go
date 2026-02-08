package controller

import (
	"net/http"

	"github.com/YagoSchramm/ApiMyChat/internal/usecase"
	"github.com/gin-gonic/gin"
)

type RoomController struct {
	Usecase usecase.RoomUsecase
}

type CreateRoomRequest struct {
	UserAID string `json:"userAId" binding:"required"`
	UserBID string `json:"userBId" binding:"required"`
}

func (rctrl *RoomController) CreateRoom(ctx *gin.Context) {
	var req CreateRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	room, err := rctrl.Usecase.CreateRoom(req.UserAID, req.UserBID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, room)
}
