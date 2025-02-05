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

func (c *UserController) RegisterUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "リクエストのパースに失敗しました"})
		return
	}

	err := c.Service.CreateUser(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "ユーザー登録に失敗しました"})
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

	authenticated, err := c.Service.AuthenticateUser(user.Username, user.Password)
	if err != nil || !authenticated {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "認証失敗"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "ログイン成功"})
}
