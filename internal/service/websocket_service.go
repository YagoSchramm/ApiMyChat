package service

import (
	"fmt"
	"sync"

	"github.com/YagoSchramm/ApiMyChat/internal/entity"
)

type WebSocketService struct {
	Rooms map[string]*entity.RoomClients
	Mutex sync.RWMutex
}

func NewWebSocketService() *WebSocketService {
	return &WebSocketService{
		Rooms: make(map[string]*entity.RoomClients),
	}
}

func (s *WebSocketService) JoinRoom(roomID string, client *entity.Client) {

	s.Mutex.Lock()
	room, exists := s.Rooms[roomID]
	if !exists {
		room = &entity.RoomClients{
			Clients: make(map[string]*entity.Client),
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
