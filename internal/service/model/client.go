package model

import "github.com/gorilla/websocket"

type Client struct {
	UserID string
	Conn   *websocket.Conn
	Send   chan []byte
}

func (c *Client) WritePump() {
	for msg := range c.Send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}
