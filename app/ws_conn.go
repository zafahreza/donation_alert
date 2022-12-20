package app

import "github.com/gorilla/websocket"

func NewWebSocket() *websocket.Upgrader {
	return &websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
}
