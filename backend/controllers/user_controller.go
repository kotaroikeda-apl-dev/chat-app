package controllers

import (
	"chat/models"
	"chat/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	Service *services.UserService
}

func NewUserController(service *services.UserService) *UserController {
	return &UserController{Service: service}
}

// ユーザー登録API
func (c *UserController) RegisterUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "リクエストのパースに失敗しました"})
		return
	}

	if err := c.Service.RegisterUser(user); err != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "ユーザー登録成功"})
}

func (c *UserController) LoginUser(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "リクエストのパースに失敗しました"})
		return
	}

	token, err := c.Service.AuthenticateUser(user)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
