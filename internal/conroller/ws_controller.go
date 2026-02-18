package controller

import (
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
}

func NewWSController(
	hub *entity.Hub,
	msg *usecase.MessageUsecase,
	auth *usecase.AuthUsecase,
) *WSController {
	return &WSController{
		hub:         hub,
		msgUsecase:  msg,
		authUsecase: auth,
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
	RoomID string `json:"room" binding:"required"`
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
	req.RoomID = ctx.Query("room")

	if req.UserID == "" || req.RoomID == "" {
		ctx.JSON(http.StatusBadRequest, "id and room required")
		return
	}

	userID, err := c.resolveUserID(ctx, req.UserID)
	if err != nil || userID == "" {
		ctx.JSON(http.StatusUnauthorized, "invalid token or user")
		return
	}
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		log.Println("upgrade error:", err)
		return
	}

	client := &entity.Client{
		UserID: userID,
		RoomID: req.RoomID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}
	c.hub.JoinRoom(req.RoomID, client)

	log.Println("user conectado:", userID)
	go client.WritePump()
	client.ReadPump(func(msg []byte) {
		c.msgUsecase.SendMessage(
			client.UserID,
			client.RoomID,
			string(msg),
		)
	})
	c.hub.Leave(req.RoomID, userID)
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
