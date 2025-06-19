package dto

type MessageInput struct {
	Type    string `json:"type"`
	Option  string `json:"option"`
	RoomID  string `json:"roomId"`
	Message string `json:"message"`
	Error   error  `json:"err"`
}

type MessageOutput struct {
	Type    string `json:"type"`
	Option  string `json:"option"`
	RoomID  string `json:"roomId"`
	Message string `json:"message"`
	Error   error  `json:"err"`
}

type GameRoomInput struct {
	RoomName string //interfaceでもよい.プレイヤーが指定するルーム番号
	// Cells map[string]*models.Cell
	X      float64
	Y      float64
	CellId string
}
