package controller

import (
	"net/http"

	"github.com/YagoSchramm/ApiMyChat/internal/entity"
	"github.com/YagoSchramm/ApiMyChat/internal/usecase"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	Usecase usecase.UserUseCase
}

func (uctrl *UserController) CreateUser(c *gin.Context) {
	var req entity.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, "Erro de registro")
		return
	}
	if _, err := uctrl.Usecase.CreateUser(req); err != nil {
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
func (uctrl *UserController) GetAll(id string, c *gin.Context) {
	var req string
	var userList []entity.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	userList, err := uctrl.Usecase.GetAll(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusAccepted, userList)
}
