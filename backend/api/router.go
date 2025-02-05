package api

import (
	"chat/controllers"
	"chat/middlewares"
	"chat/repositories"
	"chat/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterRoutes(db *gorm.DB) (*gin.Engine, *services.WebSocketService) {
	r := gin.Default()

	// CORS ミドルウェアを適用
	r.Use(middlewares.CORSConfig())

	// DIの実装
	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	messageRepo := repositories.NewMessageRepository(db)
	messageService := services.NewMessageService(messageRepo)
	messageController := controllers.NewMessageController(messageService)

	spaceRepo := repositories.NewSpaceRepository(db)
	spaceService := services.NewSpaceService(spaceRepo)
	spaceController := controllers.NewSpaceController(spaceService)

	// WebSocket の DI 設定
	webSocketService := services.NewWebSocketService(messageRepo)
	webSocketController := controllers.NewWebSocketController(webSocketService)

	// API ルート
	// ルーティング設定
	r.POST("/api/register", userController.RegisterUser)
	r.POST("/api/login", userController.LoginUser)

	r.GET("/api/messages", messageController.GetMessages)
	r.POST("/api/messages/create", messageController.CreateMessage)
	r.DELETE("/api/messages", messageController.DeleteMessage)

	r.POST("/api/spaces", spaceController.CreateSpace)
	r.GET("/api/spaces/list", spaceController.GetSpaces)

	// WebSocket
	r.GET("/api/ws", webSocketController.HandleConnections)

	// ヘルスチェック
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "healthy"})
	})

	return r, webSocketService
}
