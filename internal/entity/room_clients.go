package entity

import "sync"

type RoomClients struct {
	Clients map[string]*Client
	Mutex   sync.RWMutex
}
