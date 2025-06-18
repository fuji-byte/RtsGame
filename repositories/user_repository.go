package repositories

import (
	"errors"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"main.go/models"
)

type IMemoryRepository interface {
	CreateUser(user *models.User) error
	DeleteUser(clientId string) error
	UserNum() int
	GetAllUser() (*map[string]*models.User, error)
	GetUserByClientId(Id string) (*models.User, error)
	GetUsersByRoomId(Id string) (*map[string]*models.User, error)
	MakeRoom(room *models.GameRoom, user *models.User) (*models.GameRoom, error)
	GetRoom(roomId string) (*models.GameRoom, error)
	JoinRoom(roomId string, user *models.User) (*models.GameRoom, error)
	// SetRoomId(user *models.User, roomId string) error
	StartGame(room *models.GameRoom) error
	RunGame(room *models.GameRoom)
}

type MemoryRepository struct {
	memoryUser     map[string]*models.User     // ユーザーID → ユーザー
	memoryCell     map[string]*models.Cell     // セルID → セル
	memoryGameRoom map[string]*models.GameRoom // ルームID → ゲームルーム
	mu             sync.Mutex
}

func NewMemoryRepository(memoryUser map[string]*models.User, memoryCell map[string]*models.Cell, memoryGameRoom map[string]*models.GameRoom) IMemoryRepository {
	return &MemoryRepository{memoryUser: memoryUser, memoryCell: memoryCell, memoryGameRoom: memoryGameRoom}
}

func (s *MemoryRepository) CreateUser(user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.memoryUser[(*user).ID] != nil {
		return errors.New("user already exists")
	}
	s.memoryUser[(*user).ID] = user
	return nil
}

func (s *MemoryRepository) DeleteUser(clientId string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
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

func (s *MemoryRepository) JoinRoom(roomId string, user *models.User) (*models.GameRoom, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	room := s.memoryGameRoom[roomId]
	if room == nil {
		return nil, errors.New("room not found")
	}
	(*room).Players[user.ID] = user
	(*user).RoomID = roomId
	return room, nil
}

func (s *MemoryRepository) GetRoom(roomId string) (*models.GameRoom, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	room := s.memoryGameRoom[roomId]
	if room == nil {
		return nil, errors.New("room not found")
	}
	return room, nil
}

func (s *MemoryRepository) StartGame(room *models.GameRoom) error {
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

func (s *MemoryRepository) RunGame(room *models.GameRoom) {
	updateTicker := time.NewTicker(30 * time.Millisecond)
	endTimer := time.After(120 * time.Second)
	defer updateTicker.Stop()

	for {
		select {
		case <-updateTicker.C:
			// 30msごとの処理（ゲームロジックなど）
			//ユーザーから送信された情報を基に、updateに一時的に構造体を作成し、30msごとに更新する
			// room.Update()

		case <-endTimer:
			// 120秒経過でルームを終了
			log.Println("ルームのタイムアウトにより終了します:", room.ID)
			// room.End()
			return
		}
	}
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
