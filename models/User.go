package models

import (
	"fmt"

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
	// Mu       sync.Mutex
	SendCh chan []byte
}

func (u *User) StartWriter() {
	defer func() {
		// WebSocket接続を閉じる
		if u.Conn != nil {
			u.Conn.Close()
		}
		// SendCh を閉じる
		close(u.SendCh)
	}()

	for msg := range u.SendCh {
		if err := u.Conn.WriteMessage(websocket.TextMessage, msg); err != nil {
			fmt.Println("Write error:", err)
			break // 書き込みエラーが出たらループ終了
		}
	}
}

// Send は SendCh にメッセージを送る
func (u *User) Send(msg []byte) {
	select {
	case u.SendCh <- msg:
	default:
		fmt.Println("SendCh is full, dropping message")
	}
}
