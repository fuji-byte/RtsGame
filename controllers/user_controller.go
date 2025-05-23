package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"main.go/models"
	"main.go/services"
)

type IMemoryController interface {
	HandleWebSocket(ctx *gin.Context)
}

type MemoryController struct {
	service services.IMemoryService
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewMemoryController(service services.IMemoryService) IMemoryController {
	return &MemoryController{service: service}
}

func (c *MemoryController) HandleWebSocket(ctx *gin.Context) {

	// HTTP 接続を WebSocket にアップグレード
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println("WebSocket 接続エラー:", err)
		return
	}
	defer conn.Close()

	fmt.Println("クライアントが接続しました")
	// 新しいクライアントを登録
	clientId := uuid.New().String()
	err = c.service.CreateUser(clientId, conn)
	if err != nil {
		log.Fatal("CreateUser Error")
	}
	c.userNumControll()
	//websocket接続切断時の処理
	defer func() {
		//user情報削除、または一定時間保持
		//room情報 models.gameroomの情報変更,models.user[clientId]の削除
		err := c.service.DeleteUser(clientId)
		if err != nil {
			log.Fatal("DeleteUser Error")
		}
		c.userNumControll()
		fmt.Println("切断後処理完了:", clientId)
	}()

	//接続状態時の処理
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("接続が切断されました:", err)
			return
		}
		fmt.Printf("受信メッセージ: %s\n", msg)
		var receivedMsg models.Message
		//shouldbindingするのか
		if err := json.Unmarshal(msg, &receivedMsg); err != nil {
			fmt.Println("JSON のパースに失敗:", err)
			continue
		}
		switch receivedMsg.Type {
		case "makeRoom":
			roomId, err := c.service.MakeRoom(clientId)
			if err != nil {
				fmt.Println("MakeRoom Error:", err)
				conn.WriteMessage(websocket.TextMessage, []byte(`"message":"ルームを作成できませんでした。","error": "makeRoom Error"`))
				continue
			}
			users, err := c.service.GetUserFromRoom(roomId)
			if err != nil {
				fmt.Println("no users:", err)
				return
			}
			broadcast(*users, "roomId", roomId)
		case "joinRoom":
			users, err := c.service.GetUser(clientId)
			if err != nil {
				fmt.Println("no users:", err)
				return
			}
			var user *models.User
			for _, tempUser := range *users {
				user = tempUser
			}
			err = c.service.JoinRoom(receivedMsg.RoomID, user)
			if err != nil {
				fmt.Println("join room Error:", err)
				conn.WriteMessage(websocket.TextMessage, []byte(`"message":"ルームに参加できませんでした。","error": "joinRoom Error"`))
				continue
			}
			room, err := c.service.GetUserFromRoom(receivedMsg.RoomID)
			if err != nil {
				fmt.Println("no users:", err)
				return
			}
			broadcast(*users, "roomNum", len(*room))
		case "match":
			return
		}
	}
}

func (c *MemoryController) userNumControll() {
	//s.memoryService.userNum()を取得して、User全員にブロードキャスト
	users, err := c.service.GetUser("all")
	if err != nil {
		fmt.Println("no users", err)
		return
	}
	broadcast(*users, "userNum", len(*users))
}

// broadcast wants users map[string]*models.User, messageType string, content(int float32 string)
func broadcast[T int | float32 | string](
	users map[string]*models.User,
	messageType string,
	content T,
) {
	message := fmt.Sprintf(`{"type":"%v", "%v": "%v"}`, messageType, messageType, content)

	for _, user := range users {
		err := user.Conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			fmt.Printf("send Message Error to user %s: %v\n", user.ID, err)
		}
	}
}

type SessionHandler struct {
	service *services.SessionService
}

func NewSessionHandler(service *services.SessionService) *SessionHandler {
	//redisは接続の状態を保存するもので、接続を保存するものではない
	return &SessionHandler{service: service}
}

func (h *SessionHandler) Login(c *gin.Context) {
	userID := c.Query("user_id")
	data := "logged_in"

	err := h.service.Login(userID, data)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
