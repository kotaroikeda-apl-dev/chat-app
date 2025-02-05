package controllers

import (
	"chat/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type SpaceController struct {
	Service *services.SpaceService
}

func NewSpaceController(service *services.SpaceService) *SpaceController {
	return &SpaceController{Service: service}
}

// スペース作成エンドポイント
func (c *SpaceController) CreateSpace(ctx *gin.Context) {
	var data struct {
		Name string `json:"name"`
	}
	if err := ctx.ShouldBindJSON(&data); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "リクエストのパースに失敗しました"})
		return
	}

	err := c.Service.CreateSpace(data.Name)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "スペースの作成に失敗しました"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "スペースが作成されました"})
}

// スペース一覧取得エンドポイント
func (c *SpaceController) GetSpaces(ctx *gin.Context) {
	spaces, err := c.Service.GetSpaces()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "スペース一覧の取得に失敗しました"})
		return
	}

	ctx.JSON(http.StatusOK, spaces)
}
