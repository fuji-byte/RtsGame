package models

type Message struct {
	Type    string `json:"type"`
	RoomID  string `json:"roomId"`
	Message string `json:"message"`
	Error   error  `json:err`
}
