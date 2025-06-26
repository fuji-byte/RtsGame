package dto

import (
	"encoding/json"

	"main.go/models"
)

type TypeInput struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type MessageInput struct {
	Type     string          `json:"type"`
	Option   string          `json:"option"`
	RoomID   string          `json:"roomId"`
	Message  string          `json:"message"`
	Error    error           `json:"err"`
	GameRoom models.GameRoom `json:"gameRoom"`
}

type MessageOutput struct {
	Type    string `json:"type"`
	Option  string `json:"option"`
	RoomID  string `json:"roomId"`
	Message string `json:"message"`
	Error   error  `json:"err"`
}

type GameRoomInput struct {
	CellConnFrom string `json:"cellConnFrom"`
	CellConnTo   string `json:"cellConnTo"`
	CellId       string `json:"cellid"`
}
