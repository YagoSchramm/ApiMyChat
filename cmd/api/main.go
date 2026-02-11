package main

import (
	controller "github.com/YagoSchramm/ApiMyChat/internal/conroller"
	"github.com/YagoSchramm/ApiMyChat/internal/db"
	"github.com/YagoSchramm/ApiMyChat/internal/entity"
	"github.com/YagoSchramm/ApiMyChat/internal/repository"
	"github.com/YagoSchramm/ApiMyChat/internal/service"
	"github.com/YagoSchramm/ApiMyChat/internal/usecase"
	"github.com/gin-gonic/gin"
)

func main() {
	db, err := db.ConnectDB()
	if err != nil {
		panic(err)
	}
	cache := repository.NewMemoryCache()
	emailSrv := &service.GmailService{
		Email:    "mychatnoreplyapi@gmail.com",
		Password: "hkrolzunebsgtgfj",
	}
	supabaseService := service.NewSupabaseAuthService(
		"https://tefldqfpeckuzzfrooch.supabase.co",
		"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZSIsInJlZiI6InRlZmxkcWZwZWNrdXp6ZnJvb2NoIiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImlhdCI6MTc2OTY4MjQxNCwiZXhwIjoyMDg1MjU4NDE0fQ.O-450kZNaEQyAw7JOEAY_w4OuSfe2-NT32BLSC6J2xY",
	)
	authUsecase := usecase.NewAuthUsecase(supabaseService)
	otpUseCase := &usecase.OTPUseCase{
		Repo:  cache,
		Email: emailSrv,
	}
	EmailController := &controller.EmailController{
		UseCase: otpUseCase,
	}
	urepo := repository.NewUserRepository(db)
	uc := usecase.NewUserUseCase(urepo)
	usercontroller := &controller.UserController{
		Usecase:     uc,
		AuthUsecase: *authUsecase,
	}
	hub := entity.NewHub()
	msgRepo := repository.NewMessageRepository(db)
	msgUse := usecase.NewMessageUsecase(msgRepo, hub)
	msgController := &controller.MessageController{
		Usecase: msgUse,
	}

	wsController := controller.NewWSController(hub, msgUse, authUsecase)
	roomRepo := repository.NewRoomRepository(db)
	roomUsecase := usecase.NewRoomUsecase(roomRepo)
	roomController := &controller.RoomController{
		Usecase: roomUsecase,
	}

	r := gin.Default()
	r.POST("/send-code", EmailController.RequestOTP)
	r.POST("/verify-email", EmailController.VerifyCode)
	r.POST("/CreateUser", usercontroller.CreateUser)
	r.GET("/GetUserByID/:id", usercontroller.GetByID)
	r.POST("/login", usercontroller.Login)
	r.POST("/CreateRoom", roomController.CreateRoom)
	r.GET("/messages", msgController.GetByRoom)
	r.GET("/ws/connect", wsController.Connect)
	r.POST("/ws/message", wsController.SendMessage)
	r.Run(":8000")
}
