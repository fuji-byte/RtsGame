package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"main.go/dto"
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
	// HandshakeTimeout	ハンドシェイクのタイムアウト時間	通常 5秒〜10秒程度が一般的
	// ReadBufferSize / WriteBufferSize	I/Oバッファのサイズ（バイト単位）	多数接続時や大規模アプリで最適化可能
	// WriteBufferPool	書き込みバッファのプール	接続数が多く、GC削減したいときに有効
	// Subprotocols	WebSocketのサブプロトコル（例: chat, json, etc.）	クライアントとプロトコルネゴシエート可
	// Error	エラーハンドリングのカスタム処理	403 や 500 に独自のHTML出力など可能
	// CheckOrigin	オリジンチェック（CORS対応）	本番ではここでOriginを検証すべき
	// EnableCompression	per-message圧縮を有効化（RFC 7692）	クライアントがサポートしていれば圧縮される
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

	// 新しいクライアントを登録
	clientId := uuid.New().String()
	user, err := c.service.CreateUser(clientId, conn)
	go user.StartWriter()
	if err != nil {
		log.Fatal("CreateUser Error")
	}
	userNum := c.userNumControll()
	fmt.Println("クライアントが接続しました.現在", userNum, "人")

	//websocket接続切断時の処理
	defer func() {
		//user情報削除、または一定時間保持
		//room情報 models.gameroomの情報変更,models.user[clientId]の削除
		//データベース使うならuser.IsOnline falseにする

		// user, err := c.service.GetUserByClientId(clientId)
		// if err != nil {
		// 	fmt.Println("user取得エラー")
		// }
		users, err := c.service.GetUsersByRoomId(user.RoomID)
		if err != nil {
			//roomが存在していない
		}
		err = c.service.DeleteUser(clientId)
		if err != nil {
			log.Fatal("DeleteUser Error")
		}
		c.userNumControll()
		if users != nil {
			broadcast(*users, "roomNum", "message", len(*users))
			fmt.Println(users)
		}
		fmt.Println("切断後処理完了:", clientId)
	}()

	//接続状態時の処理
	for {
		_, msg, err := user.Conn.ReadMessage()
		if err != nil {
			fmt.Println("接続が切断されました:", err)
			return
		}
		// fmt.Printf("受信メッセージ: %s\n", msg)
		fmt.Printf("メッセージを受信しました\n")
		var raw dto.TypeInput
		// var receivedMsg dto.MessageInput
		//shouldbindingするのか
		if err := json.Unmarshal(msg, &raw); err != nil {
			fmt.Println("JSON のパースに失敗1:", err)
			continue
		}
		switch raw.Type {
		case "message":
			var receivedMsg dto.MessageInput
			if err := json.Unmarshal(raw.Data, &receivedMsg); err != nil {
				fmt.Println("MessageInput の解析に失敗:", err)
				continue
			}
			switch receivedMsg.Type {
			case "makeRoom":
				roomId, err := c.service.MakeRoom(clientId)
				if err != nil {
					fmt.Println("MakeRoom Error:", err)
					msg := `{"type":"errorMessage","message":"ルームを作成できませんでした。","error": "make a room Error"}`
					user.Send(msg)
					continue
				}
				user.Send(fmt.Sprintf(`{"type":"roomId","message":"%v"}`, roomId))
				user.Send(`{"type":"roomNum", "message":"1"}`)
				user.Send(`{"type":"host"}`)
			case "joinRoom":
				// needed roomId
				user, err := c.service.GetUserByClientId(clientId)
				if err != nil {
					fmt.Println("no users:", err)
					user.Send(`{"type":"errorMessage","message":"ルームを作成できませんでした。","error": "get a user Error"}`)
					continue
				}
				tempUsers, err := c.service.JoinRoom(receivedMsg.RoomID, user)
				if err != nil {
					fmt.Println("join room Error:", err)
					user.Send(`{"type":"errorMessage","message":"ルームに参加できませんでした。","error": "join in the room Error"}`)
					continue
				}
				room, err := c.service.GetUsersByRoomId(receivedMsg.RoomID)
				if err != nil {
					fmt.Println("no users:", err)
					user.Send(`{"type":"errorMessage","message":"ルームメンバーを取得できませんでした。","error": "couldn't get the room member(s) Error"}`)
					continue
				}
				msg := fmt.Sprintf(`%vが参加しました`, user.Name)
				broadcast(*room, "roomNum", "message", len(*room))
				broadcast(*tempUsers, "message", "message", msg)
				user.Send(fmt.Sprintf(`{"type":"roomId","message":"%v"}`, user.RoomID))
			case "match":
				//host playerがmatchを送信したら
				user, err := c.service.GetUserByClientId(clientId)
				if err != nil {
					fmt.Println("no users:", err)
					user.Send(`{"type":"errorMessage","message":"ユーザーを取得できませんでした。","error": "couldn't get the users Error"}`)
					continue
				}
				err = c.service.StartGame(user)
				if err != nil {
					user.Send(`{"type":"errorMessage","message":"ゲームのスタートに失敗しました","error": "failed to start the game or invalid host Error"}`)
					continue
				}
				room, err := c.service.GetUsersByRoomId(user.RoomID)
				if err != nil {
					fmt.Println("no users:", err)
					user.Send(`{"type":"errorMessage","message":"ルームメンバーを取得できませんでした。","error": "couldn't get the room member(s) Error"}`)
					continue
				}
				broadcast(*room, "gameStart", "gameStart", "game start")
				//この辺でgo funcで継続的にbroadcastするか
			case "observe":
				//roomのobserverに追加
				//終わるまでか、観戦キャンセルされるまでずっとブロードキャスト
				continue
			default:
				user.Send(`{"type":"errorMessage","message":"タイプが適切ではありません","error": "type Error"}`)
				continue
			}

		case "game":
			var receivedMsg dto.GameRoomInput
			if err := json.Unmarshal(raw.Data, &receivedMsg); err != nil {
				fmt.Println("GameRoomInput の解析に失敗:", err)
				continue
			}
			switch raw.Type {
			case "playing":
				user, err := c.service.GetUserByClientId(clientId)
				if err != nil {
					fmt.Println("no users:", err)
					user.Send(`{"type":"errorMessage","message":"ユーザーを取得できませんでした。","error": "couldn't get the users Error"}`)
					continue
				}
				//cellの線を切るときとつなぐときの二つ作る。オプションで指定してもらう
				err = c.service.UpdateRoom(user.ID, user.RoomID, receivedMsg.CellConnFrom, receivedMsg.CellConnTo, receivedMsg.CellId) // c.service.GetRoomInfo()
				if err != nil {
					user.Send(`{"type":"errorMessage","message":"ルームデータをアップデートできませんでした。","error": "couldn't update the room Error"}`)
					continue
				}
				continue
			default:
				user.Send(`{"type":"errorMessage","message":"タイプが適切ではありません","error": "type Error"}`)
				continue
			}
		default:
			fmt.Println("不明なタイプ:", raw.Type)
			continue
		}
	}
}

func (c *MemoryController) userNumControll() int {
	//s.memoryService.userNum()を取得して、User全員にブロードキャスト
	users, err := c.service.GetAllUser()
	if err != nil {
		fmt.Println("no users", err)
		return -1
	}
	broadcast(*users, "userNum", "message", len(*users))
	return len(*users)
}

// broadcast wants users map[string]*models.User, messageType string, content(int float32 string)
func broadcast[T int | float32 | string](
	users map[string]*models.User,
	messageType1 string,
	messageType2 string,
	content T,
) {
	message := fmt.Sprintf(`{"type":"%v", "%v": "%v"}`, messageType1, messageType2, content)
	// jsonBytes, err := json.Marshal(message)
	// if err != nil {
	// 	fmt.Println("Marshalエラー:", err)
	// 	return
	// }
	for _, user := range users {
		if user == nil {
			fmt.Printf("send Message Error to user\n")
		}
		user.Send(message)
	}
}

// type SessionHandler struct {
// 	service *services.SessionService
// }

// func NewSessionHandler(service *services.SessionService) *SessionHandler {
// 	//redisは接続の状態を保存するもので、接続を保存するものではない
// 	return &SessionHandler{service: service}
// }

// func (h *SessionHandler) Login(c *gin.Context) {
// 	userID := c.Query("user_id")
// 	data := "logged_in"

// 	err := h.service.Login(userID, data)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to login"})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"status": "ok"})
// }
