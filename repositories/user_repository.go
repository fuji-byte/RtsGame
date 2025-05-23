package repositories

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"main.go/models"
)

type IMemoryRepository interface {
	CreateUser(user *models.User) error
	DeleteUser(clientId string) error
	UserNum() int
	GetAllUser() (*map[string]*models.User, error)
	GetUser(Id string) (*models.User, error)
	GetUserFromRoom(Id string) (*map[string]*models.User, error)
	MakeRoom(roomId string, clientId string) *models.GameRoom
	GetRoom(roomId string) (*models.GameRoom, error)
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
	if s.memoryUser[user.ID] != nil {
		return errors.New("user already exists")
	}
	s.memoryUser[user.ID] = user
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

func (s *MemoryRepository) GetUser(Id string) (*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	user, ok := s.memoryUser[Id]
	if !ok {
		return nil, errors.New("invalid userId")
	}
	return user, nil
}

func (s *MemoryRepository) GetUserFromRoom(Id string) (*map[string]*models.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	users, ok := s.memoryGameRoom[Id]
	if !ok {
		return nil, errors.New("nil pointer room")
	}
	return &users.Players, nil
}

func (s *MemoryRepository) MakeRoom(roomId string, clientId string) *models.GameRoom {
	newRoom := &models.GameRoom{ID: roomId, Players: make(map[string]*models.User), Cells: make(map[string]*models.Cell), Started: false, TimeLeftSec: 120}
	user, err := s.GetUser(clientId)
	s.mu.Lock()
	defer s.mu.Unlock()
	if err != nil {
		return nil
	}
	newRoom.Players[clientId] = user
	s.memoryGameRoom[roomId] = newRoom
	return newRoom
}

func (s *MemoryRepository) GetRoom(roomId string) (*models.GameRoom, error) {
	room := s.memoryGameRoom[roomId]
	if room == nil {
		return nil, errors.New("room not found")
	}
	return room, nil
}

type SessionRepository interface {
	SetSession(userID string, data string) error
	GetSession(userID string) (string, error)
	DeleteSession(userID string) error
}

type redisSessionRepository struct {
	client *redis.Client
}

func NewRedisSessionRepository(client *redis.Client) SessionRepository {
	return &redisSessionRepository{client: client}
}

func (r *redisSessionRepository) SetSession(userID string, data string) error {
	return r.client.Set(context.Background(), "session:"+userID, data, 24*time.Hour).Err()
}

func (r *redisSessionRepository) GetSession(userID string) (string, error) {
	return r.client.Get(context.Background(), "session:"+userID).Result()
}

func (r *redisSessionRepository) DeleteSession(userID string) error {
	return r.client.Del(context.Background(), "session:"+userID).Err()
}
