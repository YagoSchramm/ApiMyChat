package entity

import (
	"sort"
	"sync"
)

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

func (h *Hub) Connect(c *Client) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	h.disconnectLocked(c.UserID)

	h.OnlineUsers[c.UserID] = c
}

func (h *Hub) JoinRoom(roomID string, userID string) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	if _, ok := h.Rooms[roomID]; !ok {
		h.Rooms[roomID] = make(map[string]*Client)
	}

	client, ok := h.OnlineUsers[userID]
	if !ok {
		return
	}

	h.Rooms[roomID][userID] = client
}

func (h *Hub) LeaveRoom(roomID, userID string) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	if room, ok := h.Rooms[roomID]; ok {
		delete(room, userID)
		if len(room) == 0 {
			delete(h.Rooms, roomID)
		}
	}
}

func (h *Hub) Disconnect(userID string) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()

	h.disconnectLocked(userID)
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
			h.disconnectLocked(client.UserID)
		}
	}
}

func (h *Hub) ConnectedUsers() []string {
	h.Mutex.RLock()
	defer h.Mutex.RUnlock()

	users := make([]string, 0, len(h.OnlineUsers))
	for userID := range h.OnlineUsers {
		users = append(users, userID)
	}

	sort.Strings(users)
	return users
}

func (h *Hub) IsOnline(userID string) bool {
	h.Mutex.RLock()
	defer h.Mutex.RUnlock()

	_, ok := h.OnlineUsers[userID]
	return ok
}

func (h *Hub) disconnectLocked(userID string) {
	if client, ok := h.OnlineUsers[userID]; ok {
		close(client.Send)
		_ = client.Conn.Close()
		delete(h.OnlineUsers, userID)
	}

	for roomID, room := range h.Rooms {
		delete(room, userID)
		if len(room) == 0 {
			delete(h.Rooms, roomID)
		}
	}
}
