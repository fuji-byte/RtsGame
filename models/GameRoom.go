package models

// redisでキャッシュで保管予定（db実装後）
// on memory

type GameRoom struct {
	ID string `gorm:"primaryKey"`
	// Key string
	RoomName    string //interfaceでもよい.プレイヤーが指定するルーム番号
	Players     map[string]*User
	Observers   map[string]*User
	HostPlayer  *User
	Cells       map[string]*Cell
	CellConn    map[string]map[string]int //cell connectionとその進捗
	Started     bool                      //default false
	TimeLeftSec float32
	Signal      chan string
	Ch          chan *GameRoom
	UserCh      chan *GameRoom
}
