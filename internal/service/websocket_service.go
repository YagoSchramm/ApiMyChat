package service

import (
	"fmt"
	"sync"

	"github.com/YagoSchramm/ApiMyChat/internal/service/model"
)

type WebSocketService struct {
	Rooms map[string]*model.RoomClients
	Mutex sync.RWMutex
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		Rooms: make(map[string]*model.RoomClients),
	}
}

func (s *WebSocketService) JoinRoom(roomID string, client *model.Client) {

	s.Mutex.Lock()
	room, exists := s.Rooms[roomID]
	if !exists {
		room = &model.RoomClients{
			Clients: make(map[string]*model.Client),
		}
		s.Rooms[roomID] = room
	}
	s.Mutex.Unlock()

	room.Mutex.Lock()
	room.Clients[client.UserID] = client
	room.Mutex.Unlock()

	fmt.Println("User entrou na room:", roomID)
}
func (s *WebSocketService) Broadcast(roomID string, message []byte) {

	s.Mutex.RLock()
	room, exists := s.Rooms[roomID]
	s.Mutex.RUnlock()

	if !exists {
		return
	}

	room.Mutex.RLock()
	for _, client := range room.Clients {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(room.Clients, client.UserID)
		}
	}
	room.Mutex.RUnlock()
}
func (s *WebSocketService) LeaveRoom(roomID string, userID string) {
	s.Mutex.RLock()
	room, exists := s.Rooms[roomID]
	s.Mutex.RUnlock()

	if !exists {
		return
	}

	room.Mutex.Lock()
	delete(room.Clients, userID)
	room.Mutex.Unlock()

	fmt.Println("User saiu da room:", roomID)
}
