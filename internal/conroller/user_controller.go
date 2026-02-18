package controller

import (
	"fmt"
	"net/http"

	"github.com/YagoSchramm/ApiMyChat/internal/entity"
	"github.com/YagoSchramm/ApiMyChat/internal/usecase"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	Usecase     usecase.UserUseCase
	AuthUsecase usecase.AuthUsecase
}

func (uctrl *UserController) CreateUser(c *gin.Context) {
	var req entity.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := uctrl.AuthUsecase.Supabase.CreateUser(req.Email, req.Password)
	if err != nil {
		fmt.Println("erro no supabase")
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	newUser := entity.User{UID: resp.ID, Email: req.Email, Name: req.Name, Description: req.Description, CreatedAt: req.CreatedAt}
	if _, err := uctrl.Usecase.CreateUser(newUser); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusCreated, "Usuário criado com sucesso!")

}
func (uctrl *UserController) GetByID(c *gin.Context) {
	id := c.Param("id")
	var user entity.User
	user, err := uctrl.Usecase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	user1 := user
	c.JSON(http.StatusAccepted, user1)
}
func (uctrl *UserController) GetByEmail(c *gin.Context) {
	email := c.Param("email")
	var user entity.User
	user, err := uctrl.Usecase.GetByID(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	user1 := user
	c.JSON(http.StatusAccepted, user1)
}
func (uctrl *UserController) GetAll(c *gin.Context) {
	id := c.Param("id")
	var userList []entity.User
	userList, err := uctrl.Usecase.GetAll(id)
	if err != nil {

		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusAccepted, userList)
}

func (uctrl *UserController) Login(ctx *gin.Context) {

	var req entity.LoginRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
		return
	}

	token, err := uctrl.AuthUsecase.Login(req.Email, req.Password)
	if err != nil {
		fmt.Println("erro no supabase")
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	var user entity.User
	if token != "" {
		user, err := uctrl.Usecase.GetByEmail(req.Email)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"uid":         user.UID,
			"email":       user.Email,
			"name":        user.Name,
			"description": user.Description,
			"createdAt":   user.CreatedAt,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"uid":         user.UID,
		"email":       user.Email,
		"name":        user.Name,
		"description": user.Description,
		"createdAt":   user.CreatedAt,
	})

}
