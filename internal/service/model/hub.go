package model

import "sync"

type Hub struct {
	Rooms       map[string]map[string]*Client
	OnlineUsers map[string]*Client
	Mutex       sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		Rooms:       make(map[string]map[string]*Client),
		OnlineUsers: make(map[string]*Client),
	}
}
func (h *Hub) JoinRoom(roomID string, c *Client) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	if _, ok := h.Rooms[roomID]; !ok {
		h.Rooms[roomID] = make(map[string]*Client)
	}

	h.Rooms[roomID][c.UserID] = c
	h.OnlineUsers[c.UserID] = c
}
func (h *Hub) Leave(roomID, userID string) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	if room, ok := h.Rooms[roomID]; ok {
		if client, ok := room[userID]; ok {
			close(client.Send)
			delete(room, userID)
		}
	}

	delete(h.OnlineUsers, userID)
}
func (h *Hub) Broadcast(roomID string, msg []byte) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	room, ok := h.Rooms[roomID]

	if !ok {
		return
	}

	for _, client := range room {
		select {
		case client.Send <- msg:
		default:
			close(client.Send)
			delete(room, client.UserID)
			delete(h.OnlineUsers, client.UserID)
		}
	}
}
