package models

type Message struct {
	ID        int    `json:"id"`
	SpaceID   string `json:"space_id"`
	Username  string `json:"username"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
}
