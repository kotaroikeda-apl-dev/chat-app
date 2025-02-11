package controllers

import (
	"chat/models"
	"chat/services"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MessageController struct {
	Service services.MessageService
}

func NewMessageController(service services.MessageService) *MessageController {
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
	fmt.Println("メッセージ作成エンドポイントにリクエストが来ました")

	var msg models.Message // `models.Message` を使用

	if err := ctx.ShouldBindJSON(&msg); err != nil {
		log.Printf("リクエストパースエラー: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "リクエストのパースに失敗しました"})
		return
	}
	log.Printf("受信したメッセージ: %+v", msg)
	id, err := c.Service.CreateMessage(msg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "メッセージの保存に失敗しました"})
		return
	}

	msg.ID = id
	ctx.JSON(http.StatusCreated, msg)
}

func (c *MessageController) DeleteMessage(ctx *gin.Context) {
	messageID, err1 := strconv.Atoi(ctx.Query("id"))
	spaceID, err2 := strconv.Atoi(ctx.Query("spaceId"))

	if err1 != nil || err2 != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "メッセージIDまたはスペースIDが無効です"})
		return
	}

	log.Printf("受信した削除リクエスト - メッセージID: %d, スペースID: %d", messageID, spaceID)

	if err := c.Service.DeleteMessage(messageID, spaceID); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "メッセージ削除成功"})
}
