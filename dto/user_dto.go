package dto

import "main.go/models"

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

// type GameRoomInput struct {
// 	RoomName string `json:"roomName"` //interfaceでもよい.プレイヤーが指定するルーム番号
// 	// Cells map[string]*models.Cell`json:"cells"`
// 	X      float64 `json:"x"`
// 	Y      float64 `json:"y"`
// 	CellId string  `json:"cellid"`
// }
