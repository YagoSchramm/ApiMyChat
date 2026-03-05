package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/YagoSchramm/ApiMyChat/internal/entity"
	"github.com/YagoSchramm/ApiMyChat/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type WSController struct {
	hub         *entity.Hub
	msgUsecase  *usecase.MessageUsecase
	authUsecase *usecase.AuthUsecase
	roomUsecase *usecase.RoomUsecase
}

func NewWSController(
	hub *entity.Hub,
	msg *usecase.MessageUsecase,
	auth *usecase.AuthUsecase,
	room *usecase.RoomUsecase,
) *WSController {
	return &WSController{
		hub:         hub,
		msgUsecase:  msg,
		authUsecase: auth,
		roomUsecase: room,
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type SendMessageRequest struct {
	RoomID  string `json:"room" binding:"required"`
	Content string `json:"content" binding:"required"`
}

type ConnectRequest struct {
	UserID string `json:"id" binding:"required"`
}

type WSInboundMessage struct {
	Type    string `json:"type"`
	RoomID  string `json:"room"`
	Content string `json:"content"`
}

type ConnectedUsersResponse struct {
	Users []string `json:"users"`
	Count int      `json:"count"`
}

func isDevBypass(ctx *gin.Context) bool {
	return ctx.Query("dev") == "1" || ctx.GetHeader("X-Dev-Bypass") == "1"
}

func (c *WSController) resolveUserID(ctx *gin.Context, rawID string) (string, error) {
	if rawID == "" {
		return "", nil
	}

	if isDevBypass(ctx) {
		return rawID, nil
	}

	if c.authUsecase == nil {
		return rawID, nil
	}

	return c.authUsecase.VerifyJWT(rawID)
}

func (c *WSController) Connect(ctx *gin.Context) {

	var req ConnectRequest
	req.UserID = ctx.Query("id")

	if req.UserID == "" {
		ctx.JSON(http.StatusBadRequest, "id required")
		return
	}

	userID, err := c.resolveUserID(ctx, req.UserID)
	if err != nil || userID == "" {
		ctx.JSON(http.StatusUnauthorized, "invalid token or user")
		return
	}

	roomIDs, err := c.loadUserRoomIDs(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "failed to load user rooms")
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	client := &entity.Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}
	c.hub.Connect(client)
	for _, roomID := range roomIDs {
		c.hub.JoinRoom(roomID, userID)
	}

	log.Println("user conectado:", userID)
	go client.WritePump()
	client.ReadPump(func(msg []byte) {
		c.handleInboundMessage(userID, msg)
	})
	c.hub.Disconnect(userID)
	conn.Close()

	log.Println("user saiu:", userID)
}

func (c *WSController) SendMessage(ctx *gin.Context) {
	rawID := ctx.Query("id")
	if rawID == "" {
		ctx.JSON(http.StatusUnauthorized, "id required")
		return
	}

	userID, err := c.resolveUserID(ctx, rawID)
	if err != nil || userID == "" {
		ctx.JSON(http.StatusUnauthorized, "invalid token or user")
		return
	}

	var req SendMessageRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, "invalid body")
		return
	}

	c.msgUsecase.SendMessage(userID, req.RoomID, req.Content)
	ctx.JSON(http.StatusAccepted, gin.H{"status": "sent"})
}

func (c *WSController) GetConnectedUsers(ctx *gin.Context) {
	users := c.hub.ConnectedUsers()
	ctx.JSON(http.StatusOK, ConnectedUsersResponse{
		Users: users,
		Count: len(users),
	})
}

func (c *WSController) loadUserRoomIDs(userID string) ([]string, error) {
	if c.roomUsecase == nil {
		return []string{}, nil
	}

	rooms, err := c.roomUsecase.GetRoomsByUid(userID)
	if err != nil {
		return nil, err
	}

	roomIDs := make([]string, 0, len(rooms))
	for _, room := range rooms {
		roomIDs = append(roomIDs, room.ID)
	}

	return roomIDs, nil
}

func (c *WSController) handleInboundMessage(userID string, msg []byte) {
	var payload WSInboundMessage
	if err := json.Unmarshal(msg, &payload); err != nil {
		return
	}

	switch payload.Type {
	case "", "message":
		if payload.RoomID == "" || payload.Content == "" {
			return
		}
		c.msgUsecase.SendMessage(userID, payload.RoomID, payload.Content)
	case "join":
		if payload.RoomID == "" {
			return
		}
		c.hub.JoinRoom(payload.RoomID, userID)
	case "leave":
		if payload.RoomID == "" {
			return
		}
		c.hub.LeaveRoom(payload.RoomID, userID)
	}
}
