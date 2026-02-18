package controller

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/YagoSchramm/ApiMyChat/internal/entity"
	"github.com/YagoSchramm/ApiMyChat/internal/usecase"
	"github.com/gin-gonic/gin"
)

type RoomController struct {
	Usecase usecase.RoomUsecase
}

type CreateRoomRequest struct {
	Name    string   `json:"name" binding:"required"`
	UserIDs []string `json:"userIds" binding:"required,min=2"`
}

type GetRoomsByUidResponse struct {
	Contacts []entity.RoomWithUsers `json:"contacts"`
	Groups   []entity.RoomWithUsers `json:"groups"`
}

func (rctrl *RoomController) CreateRoom(ctx *gin.Context) {
	var req CreateRoomRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	room, err := rctrl.Usecase.CreateRoom(req.Name, req.UserIDs)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusCreated, room)
}
func (rctrl *RoomController) GetRoomById(ctx *gin.Context) {
	uid := ctx.Param("id")
	room, err := rctrl.Usecase.GetRoomById(uid)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, "room not found")
			return
		}
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, room)
}
func (rctrl *RoomController) GetRoomsByUid(ctx *gin.Context) {
	uid := ctx.Param("uid")
	roomList, err := rctrl.Usecase.GetRoomsByUid(uid)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err)
		return
	}

	contacts := make([]entity.RoomWithUsers, 0)
	groups := make([]entity.RoomWithUsers, 0)

	for _, room := range roomList {
		if len(room.Users) == 2 {
			contacts = append(contacts, room)
			continue
		}
		groups = append(groups, room)
	}

	ctx.JSON(http.StatusOK, GetRoomsByUidResponse{
		Contacts: contacts,
		Groups:   groups,
	})
}
