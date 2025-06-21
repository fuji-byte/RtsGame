package models

import (
	"sync"
)

// redisで保管予定
// on memory
type GameRoom struct {
	ID string `gorm:"primaryKey"`
	// Key string
	RoomName    string //interfaceでもよい.プレイヤーが指定するルーム番号
	Players     map[string]*User
	HostPlayer  *User
	Cells       map[string]*Cell
	Started     bool //default false
	Mutex       sync.Mutex
	TimeLeftSec int
	Signal      chan string
	Ch          chan *GameRoom
	UserCh      chan *GameRoom
}
