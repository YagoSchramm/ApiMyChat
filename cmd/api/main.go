package main

import (
	"log"
	"net/http"
	"os"

	"github.com/YagoSchramm/ApiMyChat/internal/config"
	controller "github.com/YagoSchramm/ApiMyChat/internal/conroller"
	"github.com/YagoSchramm/ApiMyChat/internal/db"
	"github.com/YagoSchramm/ApiMyChat/internal/entity"
	"github.com/YagoSchramm/ApiMyChat/internal/repository"
	"github.com/YagoSchramm/ApiMyChat/internal/service"
	"github.com/YagoSchramm/ApiMyChat/internal/usecase"
	"github.com/gin-gonic/gin"
)

func main() {
	if err := config.LoadDotEnv(".env"); err != nil && !os.IsNotExist(err) {
		log.Fatalf("failed to load .env: %v", err)
	}

	dbConn, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}
	cache := repository.NewMemoryCache()
	emailSrv := &service.GmailService{
		Email:    mustGetEnv("GMAIL_EMAIL"),
		Password: mustGetEnv("GMAIL_APP_PASSWORD"),
	}
	supabaseService := service.NewSupabaseAuthService(
		mustGetEnv("SUPABASE_URL"),
		mustGetEnv("SUPABASE_KEY"),
	)
	authUsecase := usecase.NewAuthUsecase(supabaseService)
	otpUseCase := &usecase.OTPUseCase{
		Repo:  cache,
		Email: emailSrv,
	}
	EmailController := &controller.EmailController{
		UseCase: otpUseCase,
	}
	urepo := repository.NewUserRepository(dbConn)
	uc := usecase.NewUserUseCase(urepo)
	usercontroller := &controller.UserController{
		Usecase:     uc,
		AuthUsecase: *authUsecase,
	}
	hub := entity.NewHub()
	msgRepo := repository.NewMessageRepository(dbConn)
	msgUse := usecase.NewMessageUsecase(msgRepo, hub)
	msgController := &controller.MessageController{
		Usecase: msgUse,
	}

	wsController := controller.NewWSController(hub, msgUse, authUsecase)
	roomRepo := repository.NewRoomRepository(dbConn)
	roomUsecase := usecase.NewRoomUsecase(roomRepo)
	roomController := &controller.RoomController{
		Usecase: roomUsecase,
	}

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization, X-Dev-Bypass")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.POST("/send-code", EmailController.RequestOTP)
	r.POST("/verify-email", EmailController.VerifyCode)
	r.POST("/CreateUser", usercontroller.CreateUser)
	r.POST("/login", usercontroller.Login)
	r.POST("/CreateRoom", roomController.CreateRoom)
	r.GET("/GetRoomsByUid/:uid", roomController.GetRoomsByUid)
	r.GET("/GetRoomByUid/:id", roomController.GetRoomById)
	r.GET("/GetUserByID/:id", usercontroller.GetByID)
	r.GET("/GetAll/:id", usercontroller.GetAll)
	r.GET("/messages", msgController.GetByRoom)
	r.GET("/ws/connect", wsController.Connect)
	r.POST("/ws/message", wsController.SendMessage)
	port := getEnvOrDefault("API_PORT", "8000")
	r.Run(":" + port)
}

func mustGetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("missing required env var: %s", key)
	}
	return value
}

func getEnvOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
