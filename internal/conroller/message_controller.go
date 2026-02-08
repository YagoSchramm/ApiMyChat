package controller

import (
	"net/http"
	"strconv"

	"github.com/YagoSchramm/ApiMyChat/internal/usecase"
	"github.com/gin-gonic/gin"
)

type MessageController struct {
	Usecase *usecase.MessageUsecase
}

func (mctrl *MessageController) GetByRoom(ctx *gin.Context) {
	roomID := ctx.Query("room")
	if roomID == "" {
		ctx.JSON(http.StatusBadRequest, "room required")
		return
	}

	limit := 50
	if v := ctx.Query("limit"); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed > 0 {
			limit = parsed
		}
	}

	msgs, err := mctrl.Usecase.GetByRoom(roomID, limit)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, msgs)
}
