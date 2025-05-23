package models

import "github.com/gorilla/websocket"

type User struct {
	ID       string
	Name     string
	Color    string
	IsOnline bool
	Conn     *websocket.Conn
	RoomID   string
}
