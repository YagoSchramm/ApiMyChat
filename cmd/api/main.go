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
		Usecase: uc,
	}
	r := gin.Default()
	r.POST("/send-code", EmailController.RequestOTP)
	r.POST("/verify-email", EmailController.VerifyCode)
	r.POST("/register", usercontroller.CreateUser)
	r.GET("/GetUserByID/:id", usercontroller.GetByID)

	r.Run(":8000")
}
