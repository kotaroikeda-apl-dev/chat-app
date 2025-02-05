package controllers

import (
	"chat/models"
	"chat/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	Service *services.MessageService
}

func NewMessageController(service *services.MessageService) *MessageController {
	return &MessageController{Service: service}
}

func (c *MessageController) GetMessages(ctx *gin.Context) {
	spaceIdStr := ctx.Query("spaceId")
	spaceId, err := strconv.Atoi(spaceIdStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無効な spaceId"})
		return
	}

	messages, err := c.Service.GetMessages(spaceId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "メッセージ取得失敗"})
		return
	}

	ctx.JSON(http.StatusOK, messages)
}

// メッセージ作成API
func (c *MessageController) CreateMessage(ctx *gin.Context) {
	var msg models.Message // `models.Message` を使用

	if err := ctx.ShouldBindJSON(&msg); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "リクエストのパースに失敗しました"})
		return
	}
	id, err := c.Service.CreateMessage(msg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "メッセージの保存に失敗しました"})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "メッセージ作成成功", "id": id})
}

func (c *MessageController) DeleteMessage(ctx *gin.Context) {
	messageIDStr := ctx.Query("id")
	spaceIDStr := ctx.Query("spaceId")

	messageID, err := strconv.Atoi(messageIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無効なメッセージID"})
		return
	}

	spaceID, err := strconv.Atoi(spaceIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "無効なスペースID"})
		return
	}

	err = c.Service.DeleteMessage(messageID, spaceID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "メッセージ削除に失敗しました"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "メッセージ削除成功"})
}
