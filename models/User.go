package models

import (
	"github.com/gorilla/websocket"
)

// on memory
type User struct {
	ID       string //clientId　一応primaryKey
	Name     string
	Color    string
	IsOnline bool
	Conn     *websocket.Conn
	RoomID   string
}
