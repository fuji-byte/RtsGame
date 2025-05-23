package models

import "sync"

type GameRoom struct {
	ID          string
	Players     map[string]*User
	HostPlayer  *User
	Cells       map[string]*Cell
	Started     bool //default false
	Mutex       sync.Mutex
	TimeLeftSec int
}
