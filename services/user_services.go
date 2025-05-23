package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"main.go/models"
	"main.go/repositories"
)

type IMemoryService interface {
	CreateUser(lientId string, conn *websocket.Conn) error
	DeleteUser(lientId string) error
	UserNum() int
	GetUser(Id string) (*map[string]*models.User, error)
	GetUserFromRoom(Id string) (*map[string]*models.User, error)
	MakeRoom(clientId string) (string, error)
	JoinRoom(roomId string, user *models.User) error
}

type MemoryService struct {
	memoryRepository repositories.IMemoryRepository
}

func NewMemoryService(memoryRepository repositories.IMemoryRepository) IMemoryService {
	return &MemoryService{memoryRepository: memoryRepository}
}

func (s *MemoryService) CreateUser(clientId string, conn *websocket.Conn) error {
	newUser := models.User{ID: clientId, Name: "none", Color: "blue", IsOnline: true, Conn: conn, RoomID: "-1"}
	return s.memoryRepository.CreateUser(&newUser)
}

func (s *MemoryService) DeleteUser(clientId string) error {
	return s.memoryRepository.DeleteUser(clientId)
}

func (s *MemoryService) UserNum() int {
	return s.memoryRepository.UserNum()
}

func (s *MemoryService) GetUser(Id string) (*map[string]*models.User, error) {
	if Id == "all" {
		return s.memoryRepository.GetAllUser()
	} else {
		user, err := s.memoryRepository.GetUser(Id)
		if err != nil {
			return nil, err
		}
		tempMap := map[string]*models.User{
			user.ID: user,
		}
		return &tempMap, nil
	}
}

func (s *MemoryService) GetUserFromRoom(Id string) (*map[string]*models.User, error) {
	return s.memoryRepository.GetUserFromRoom(Id)
}

func (s *MemoryService) MakeRoom(clientId string) (string, error) {
	//clietIdのroomが存在していないか
	user, err := s.memoryRepository.GetUser(clientId)
	if err != nil {
		return "", err
	}
	if (*user).RoomID != "-1" {
		return "", errors.New("user already in a Room")
	}
	roomId := uuid.New().String()
	(*user).RoomID = roomId
	room := s.memoryRepository.MakeRoom(roomId, clientId)
	room.HostPlayer = (user)
	return roomId, nil
}

func (s *MemoryService) JoinRoom(roomId string, user *models.User) error {
	if (*user).RoomID != "-1" {
		return errors.New("user already in a Room")
	}
	room, err := s.memoryRepository.GetRoom(roomId)
	if err != nil {
		return err
	}
	room.Players[user.ID] = user
	(*user).RoomID = roomId
	return nil
}

type SessionService struct {
	sessionRepo repositories.SessionRepository
}

func NewSessionService(repo repositories.SessionRepository) *SessionService {
	return &SessionService{sessionRepo: repo}
}

func (s *SessionService) Login(userID string, sessionData string) error {
	return s.sessionRepo.SetSession(userID, sessionData)
}

func (s *SessionService) Logout(userID string) error {
	return s.sessionRepo.DeleteSession(userID)
}

func (s *SessionService) GetSessionData(userID string) (string, error) {
	return s.sessionRepo.GetSession(userID)
}
