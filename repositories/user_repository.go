package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"main.go/dto"
	"main.go/models"
)

type IMemoryRepository interface {
	CreateUser(user *models.User) (*models.User, error)
	DeleteUser(clientId string) error
	UserNum() int
	GetAllUser() (*map[string]*models.User, error)
	GetUserByClientId(Id string) (*models.User, error)
	GetUsersByRoomId(Id string) (*map[string]*models.User, error)
	MakeRoom(room *models.GameRoom, user *models.User) (*models.GameRoom, error)
	GetRoom(roomId string) (*models.GameRoom, error)
	JoinRoom(roomId string, user *models.User) (*map[string]*models.User, error)
	// SetRoomId(user *models.User, roomId string) error
	SetGame(room *models.GameRoom) error
	TempRoom(signal chan string, ch chan *models.GameRoom, userch chan *models.GameRoom, room *models.GameRoom) error
	RunGame(signal chan string, ch chan *models.GameRoom, room *models.GameRoom) error
	TestRoom(ch chan *models.GameRoom) (*models.GameRoom, error)
	UpdateRoom(room *models.GameRoom) error
	SetCh(room *models.GameRoom, signal chan string, ch chan *models.GameRoom, userch chan *models.GameRoom) error
	SaveLogRoom(room *models.GameRoom) error
	CheckCells(cellConnFrom, cellConnTo string, room *models.GameRoom) error
	CompareCellId(clientId, cellId string, room *models.GameRoom) error
	AddCell(cellConnFrom, cellConnTo string, room *models.GameRoom) error
}

// 現状、すべての変数にmuがつくため、効率が良くない
type MemoryRepository struct {
	memoryUser     map[string]*models.User     // ユーザーID → ユーザー
	memoryCell     map[string]*models.Cell     // セルID → セル
	memoryGameRoom map[string]*models.GameRoom // ルームID → ゲームルーム
	mu             sync.Mutex
}

func NewMemoryRepository(memoryUser map[string]*models.User, memoryCell map[string]*models.Cell, memoryGameRoom map[string]*models.GameRoom) IMemoryRepository {
	return &MemoryRepository{memoryUser: memoryUser, memoryCell: memoryCell, memoryGameRoom: memoryGameRoom}
}

func (s *MemoryRepository) CreateUser(user *models.User) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.memoryUser[(*user).ID] != nil {
		return nil, errors.New("user already exists")
	}
	s.memoryUser[(*user).ID] = user
	return user, nil
}

func (s *MemoryRepository) DeleteUser(clientId string) error {
	user, err := s.GetUserByClientId(clientId)
	s.mu.Lock()
	defer s.mu.Unlock()
	var room *models.GameRoom = nil
	if err != nil {
		//ユーザーはルームに所属していなかった
	} else {
		room = s.memoryGameRoom[user.RoomID]
	}
	if user != nil && room != nil {
		if user := room.Players[clientId]; user != nil {
			var hostTF bool
			if user.ID == room.HostPlayer.ID {
				hostTF = true
			} else {
				hostTF = false
			}
			delete(room.Players, clientId)
			go func(users *map[string]*models.User, roomId string, hostTF bool) {
				s.mu.Lock()
				defer s.mu.Unlock()
				//ここで、ホストの変更もしくはルームの削除を行う
				if len(*users) == 0 {
					delete(s.memoryGameRoom, user.RoomID)
				} else if hostTF {
					//host譲渡をここに書く また、hostが変更されたら、hostになった人に通知
					//gameroomないのhostを変更
					var firstUser *models.User
					for _, v := range *users {
						firstUser = v
						break // 最初の要素でループを抜ける
					}
					s.memoryGameRoom[roomId].HostPlayer.ID = firstUser.ID
					firstUser.Send(`{"type":"host"}`)
					firstUser.Send(`{"type":"message","message":"このルームのホストになりました。"}`)
				}
			}(&room.Players, user.RoomID, hostTF)
		}
	}
	if user := s.memoryUser[clientId]; user == nil {
		return errors.New("user already deleted")
	}
	delete(s.memoryUser, clientId)
	return nil
}

func (s *MemoryRepository) UserNum() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.memoryUser == nil {
		log.Fatal("userNum Error")
	}
	userNum := len(s.memoryUser)
	return userNum
}

func (s *MemoryRepository) GetAllUser() (*map[string]*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.memoryUser == nil {
		return nil, errors.New("users not found")
	}
	return &s.memoryUser, nil
}

func (s *MemoryRepository) GetUserByClientId(Id string) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, ok := s.memoryUser[Id]
	if !ok {
		return nil, errors.New("invalid userId")
	}
	return user, nil
}

func (s *MemoryRepository) GetUsersByRoomId(Id string) (*map[string]*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	users, ok := s.memoryGameRoom[Id]
	if !ok {
		return nil, errors.New("nil pointer room")
	}
	return &users.Players, nil
}

func (s *MemoryRepository) MakeRoom(room *models.GameRoom, user *models.User) (*models.GameRoom, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.memoryGameRoom[(*room).ID] = room
	(*user).RoomID = (*room).ID
	return room, nil
}

func (s *MemoryRepository) JoinRoom(roomId string, user *models.User) (*map[string]*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	room := s.memoryGameRoom[roomId]
	if room == nil {
		return nil, errors.New("room not found")
	}
	tempUsers := make(map[string]*models.User)

	for k, v := range room.Players {
		tempUsers[k] = v
	}

	for k, v := range room.Observers {
		tempUsers[k] = v
	}
	(*room).Players[user.ID] = user
	(*user).RoomID = roomId
	return &tempUsers, nil
}

func (s *MemoryRepository) GetRoom(roomId string) (*models.GameRoom, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	room := s.memoryGameRoom[roomId]
	if room == nil {
		return nil, errors.New("empty room")
	}
	return room, nil
}

func (s *MemoryRepository) SetGame(room *models.GameRoom) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	room.Started = true
	//cellの初期化と、go funcでの試合開始、log保存開始
	// rand.Seed(time.Now().UnixNano())
	screenWidth := 360.0
	screenHeight := 640.0
	margin := 20.0
	for _, client := range room.Players {
		x := rand.Float64()*(screenWidth-2*margin) + margin
		y := rand.Float64()*(screenHeight-2*margin) + margin
		id := uuid.New().String()
		cell := models.Cell{
			ID:       id,
			PlayerID: client.ID,
			Hp:       10,
			X:        x,
			Y:        y,
			Rank:     1,
			Power:    1,
		}
		room.Cells[id] = &cell
	}
	// room.Cells
	// go func
	return nil
}

func (s *MemoryRepository) TempRoom(signal chan string, ch chan *models.GameRoom, userch chan *models.GameRoom, room *models.GameRoom) error {
	//userからの処理を一時的に保存
	realRoom := s.memoryGameRoom[room.ID]
	//一旦warning無視　プロトタイプ完成したら直す
	tempRoom := *realRoom
	//ゲーム中はループ
	go func() {
		for {
			tempRoom = *<-userch
		}
	}()
	for {
		msg := <-signal
		if msg == "update" {
			ch <- &tempRoom
		} else if msg == "end" {
			//roomIdを"-1にしたり終了処理行う"
			break
		}
	}
	return nil
}

func (s *MemoryRepository) RunGame(signal chan string, ch chan *models.GameRoom, room *models.GameRoom) error {
	updateTicker := time.NewTicker(30 * time.Millisecond)
	endTimer := time.After(time.Duration(room.TimeLeftSec) * time.Second)
	defer updateTicker.Stop()

	for {
		select {
		case <-updateTicker.C:
			signal <- "update"
			// 30msごとの処理（ゲームロジックなど）
			//ユーザーから送信された情報を基に、updateに一時的に構造体を作成し、30msごとに更新する
			//異常、チートな移動、変更がないか また、ここで変更をlogとして保存しておく
			//ほかの関数に一時的に保存し、この関数から信号を送信したらtestroomに保存したものを送信し、test検証後にアップデートする
			//ルームのプレイヤーが０になったらsavelog以外消す？再接続可能にするか
			room, err := s.TestRoom(ch)
			if err != nil {
				return err
			}
			//ここでアップデートする
			err = s.UpdateRoom(room)
			if err != nil {
				return err
			}
			s.Broadcast(room)
			s.SaveLogRoom(room)
		case <-endTimer:
			signal <- "end"
			// 120秒経過でルームを終了
			log.Println("ルームのタイムアウトにより終了します:", room.ID)
			// room.End()
			return nil
		}
	}
}

func (s *MemoryRepository) TestRoom(ch chan *models.GameRoom) (*models.GameRoom, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	newRoom := <-ch
	room := s.memoryGameRoom[newRoom.ID]
	//ここでnewRoomとroomを比較し、異常がないか検知する。現時点でのチート対策はない
	return room, nil
}

func (s *MemoryRepository) UpdateRoom(newRoom *models.GameRoom) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.memoryGameRoom[newRoom.ID] == nil {
		return errors.New("nil pointer in UpdateRoom")
	}
	room := s.memoryGameRoom[newRoom.ID]
	if room == nil {
		return errors.New("nil pointer in UpdateRoom")
	}
	*room = *newRoom
	room.TimeLeftSec -= 0.03
	return nil
}

func (s *MemoryRepository) SaveLogRoom(room *models.GameRoom) error {

	//プロトタイプ完成後に実装。
	//データベースなどにjsonで予定
	return nil
}

func (s *MemoryRepository) Broadcast(room *models.GameRoom) error {
	sendMessage := &dto.GameRoomOutput{
		Cells:       room.Cells,
		CellConn:    room.CellConn,
		TimeLeftSec: room.TimeLeftSec,
	}
	jsonData, err := json.Marshal(*sendMessage)
	if err != nil {
		fmt.Println(err)
	}
	msg := fmt.Sprintf(`{"type":"gameUpdate","message":%s}`, string(jsonData))
	for _, v := range room.Players {
		v.Send(msg)
	}
	for _, v := range room.Observers {
		v.Send(msg)
	}
	return nil
}

func (s *MemoryRepository) SetCh(room *models.GameRoom, signal chan string, ch chan *models.GameRoom, userch chan *models.GameRoom) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if room == nil {
		return errors.New("empty room")
	}
	(*room).Signal = signal
	(*room).Ch = ch
	(*room).UserCh = userch
	return nil
}

// roomにcellConnFrom, cellConnToが存在しているか調べる
func (s *MemoryRepository) CheckCells(cellConnFrom, cellConnTo string, room *models.GameRoom) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if room == nil {
		return errors.New("empty room")
	}
	if room.Cells[cellConnFrom] == nil || room.Cells[cellConnTo] == nil {
		return errors.New("invalid cellId")
	}
	return nil
}

// cellのプレイヤーIDとクライアントIDを比較
func (s *MemoryRepository) CompareCellId(clientId, cellId string, room *models.GameRoom) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if room == nil {
		return errors.New("empty room")
	}
	cell := (*room).Cells[cellId]
	if (*cell).PlayerID != clientId {
		return errors.New("you can't operate this cell")
	}
	return nil
}

func (s *MemoryRepository) AddCell(cellConnFrom, cellConnTo string, room *models.GameRoom) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if room == nil {
		return errors.New("empty room")
	}
	//ここは参照しているだけなので安全
	for _, v := range room.CellConn[cellConnFrom] {
		if v == cellConnTo {
			return errors.New("cell conn already exists")
		}
	}
	newRoom := *room
	// CellConn のディープコピー
	newCellConn := make(map[string][]string, len(room.CellConn))
	for key, slice := range room.CellConn {
		newCellConn[key] = append([]string(nil), slice...)
	}
	// 新しい接続を追加
	newCellConn[cellConnFrom] = append(newCellConn[cellConnFrom], cellConnTo)
	newRoom.CellConn = newCellConn

	//検証後にroomに保存する
	userch := (*room).UserCh
	userch <- &newRoom
	return nil
}

// type SessionRepository interface {
// 	SetSession(userID string, data string) error
// 	GetSession(userID string) (string, error)
// 	DeleteSession(userID string) error
// }

// type redisSessionRepository struct {
// 	client *redis.Client
// }

// func NewRedisSessionRepository(client *redis.Client) SessionRepository {
// 	return &redisSessionRepository{client: client}
// }

// func (r *redisSessionRepository) SetSession(userID string, data string) error {
// 	return r.client.Set(context.Background(), "session:"+userID, data, 24*time.Hour).Err()
// }

// func (r *redisSessionRepository) GetSession(userID string) (string, error) {
// 	return r.client.Get(context.Background(), "session:"+userID).Result()
// }

// func (r *redisSessionRepository) DeleteSession(userID string) error {
// 	return r.client.Del(context.Background(), "session:"+userID).Err()
// }
