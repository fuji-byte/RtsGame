package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"main.go/models"
	"main.go/repositories"
)

type IMemoryService interface {
	CreateUser(clientId string, conn *websocket.Conn) error
	DeleteUser(clientId string) error
	UserNum() int
	GetAllUser() (*map[string]*models.User, error)
	GetUserByClientId(clientId string) (*models.User, error)
	MakeRoom(clientId string) (string, error)
	GetUsersByRoomId(Id string) (*map[string]*models.User, error)
	JoinRoom(roomId string, user *models.User) error
	StartGame(user *models.User) error
	GetRoomInfo(roomId string) (*models.GameRoom, error)
}

type MemoryService struct {
	memoryRepository repositories.IMemoryRepository
}

func NewMemoryService(memoryRepository repositories.IMemoryRepository) IMemoryService {
	return &MemoryService{memoryRepository: memoryRepository}
}

func (s *MemoryService) CreateUser(clientId string, conn *websocket.Conn) error {
	newUser := models.User{ID: clientId, Name: "guest", Color: "blue", IsOnline: true, Conn: conn, RoomID: "-1"}
	return s.memoryRepository.CreateUser(&newUser)
}

func (s *MemoryService) DeleteUser(clientId string) error {
	return s.memoryRepository.DeleteUser(clientId)
}

func (s *MemoryService) UserNum() int {
	return s.memoryRepository.UserNum()
}

func (s *MemoryService) GetAllUser() (*map[string]*models.User, error) {
	return s.memoryRepository.GetAllUser()
}

func (s *MemoryService) GetUserByClientId(clientId string) (*models.User, error) {
	user, err := s.memoryRepository.GetUserByClientId(clientId)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *MemoryService) MakeRoom(clientId string) (string, error) {
	//clietIdのroomが存在していないか
	user, err := s.memoryRepository.GetUserByClientId(clientId)
	if err != nil {
		return "", err
	}
	if (*user).RoomID != "-1" {
		return "", errors.New("user already in a Room")
	}
	roomId := uuid.New().String()
	//memoryの編集はrepositoryのほうがいいかも
	// err = s.memoryRepository.SetRoomId(user, roomId)
	// if err != nil {
	// 	return "", err
	// }

	players := map[string]*models.User{
		clientId: user,
	}
	newRoom := &models.GameRoom{ID: roomId, RoomName: "", Players: players, HostPlayer: user, Cells: make(map[string]*models.Cell), Started: false, TimeLeftSec: 120}
	// user, err := s.GetUserByClientId(clientId)
	_, err = s.memoryRepository.MakeRoom(newRoom, user)
	if err != nil {
		return "", err
	}
	return roomId, nil
}

func (s *MemoryService) GetUsersByRoomId(Id string) (*map[string]*models.User, error) {
	return s.memoryRepository.GetUsersByRoomId(Id)
}

func (s *MemoryService) JoinRoom(roomId string, user *models.User) error {
	if (*user).RoomID != "-1" {
		return errors.New("user already in a Room")
	}
	// room, err := s.memoryRepository.GetRoom(roomId)
	_, err := s.memoryRepository.JoinRoom(roomId, user)
	if err != nil {
		return err
	}
	// room.Players[user.ID] = user
	// (*user).RoomID = roomId
	return nil
}

func (s *MemoryService) StartGame(user *models.User) error {
	//hostかどうか、ほかにプレイヤーが一人以上いるか
	roomId := user.RoomID
	room, err := s.memoryRepository.GetRoom(roomId)
	if err != nil {
		return err
	}
	if user.ID != room.HostPlayer.ID {
		return errors.New("the user is not host")
	}
	if len(room.Players) <= 1 {
		return errors.New("the room doesn't exist member")
	}
	//ルーム処理
	err = s.memoryRepository.StartGame(room)
	if err != nil {
		return err
	}
	ch := make(chan *models.GameRoom)
	go s.memoryRepository.TempRoom(ch, room)
	go s.memoryRepository.RunGame(ch, room)
	return nil
}

func (s *MemoryService) GetRoomInfo(roomId string) (*models.GameRoom, error) {
	roomInfo, err := s.memoryRepository.GetRoom(roomId)
	if err != nil {
		return nil, err
	}
	return roomInfo, err
}

// func (s *MemoryService) RunGame(room *models.GameRoom) error {

// 	go s.memoryRepository.RunGmae()
// }

// type SessionService struct {
// 	sessionRepo repositories.SessionRepository
// }

// func NewSessionService(repo repositories.SessionRepository) *SessionService {
// 	return &SessionService{sessionRepo: repo}
// }

// func (s *SessionService) Login(userID string, sessionData string) error {
// 	return s.sessionRepo.SetSession(userID, sessionData)
// }

// func (s *SessionService) Logout(userID string) error {
// 	return s.sessionRepo.DeleteSession(userID)
// }

// func (s *SessionService) GetSessionData(userID string) (string, error) {
// 	return s.sessionRepo.GetSession(userID)
// }
