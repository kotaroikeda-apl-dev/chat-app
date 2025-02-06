package models

import "time"

type Message struct {
	ID        int       `json:"id"`
	SpaceID   int       `json:"space_id"`
	Username  string    `json:"username"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}
