package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"main.go/controllers"
	"main.go/infra"
	"main.go/models"
	"main.go/repositories"
	"main.go/services"
)

func main() {

	users := make(map[string]*models.User, 0)
	cells := make(map[string]*models.Cell, 0)
	gameRooms := make(map[string]*models.GameRoom, 0)
	redisClient := infra.NewRedisClient()
	RedisSessionRepository := repositories.NewRedisSessionRepository(redisClient)
	RedisSessionService := services.NewSessionService(RedisSessionRepository)
	RedisSessionController := controllers.NewSessionHandler(RedisSessionService)
	userMemoryRepository := repositories.NewMemoryRepository(users, cells, gameRooms)
	userMemoryService := services.NewMemoryService(userMemoryRepository)
	userController := controllers.NewMemoryController(userMemoryService)

	r := gin.Default()
	// r.Use(cors.Default())
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"},
		// AllowOrigins:     []string{"http://www.touhobby.com"},
		AllowMethods:     []string{"GET", "POST", "PUT"},            // 単純メソッドのみ許可
		AllowHeaders:     []string{"Content-Type", "Authorization"}, // 必要なヘッダーのみ許可
		AllowCredentials: true,
		MaxAge:           86400, // プリフライトリクエストを24時間キャッシュ
	}))

	r.GET("/ws", userController.HandleWebSocket)
	r.GET("/login", RedisSessionController.Login)
	r.Run(":8080") // localhost:8080 でサーバーを立てます。
}

//cookieにより一度ログインしたら自動認証できるようにする
