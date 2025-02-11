package main

import (
	"chat/api"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	host     = os.Getenv("DATABASE_HOST")
	user     = os.Getenv("DATABASE_USER")
	password = os.Getenv("DATABASE_PASSWORD")
	dbName   = os.Getenv("DATABASE_NAME")
	port     = "5432"
	dsn      = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbName, port)
)

func main() {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("DB接続エラー: %v", err)
	}

	// ルートの登録
	r, webSocketService := api.RegisterRoutes(db)

	// **WebSocketのメッセージ処理をゴルーチンで実行**
	go webSocketService.HandleMessages()

	r.Run(":8080")

	log.Println("サーバーを起動中: ポート:8080")
}
