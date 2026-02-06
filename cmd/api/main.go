package main

import (
	controller "github.com/YagoSchramm/ApiMyChat/internal/conroller"
	"github.com/YagoSchramm/ApiMyChat/internal/db"
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

	r := gin.Default()
	r.POST("/send-code", EmailController.RequestOTP)
	r.POST("/verify-email", EmailController.VerifyCode)
	r.POST("/register", usercontroller.CreateUser)
	r.GET("/GetUserByID/:id", usercontroller.GetByID)
	r.POST("/login", usercontroller.Login)
	r.Run(":8000")
}
