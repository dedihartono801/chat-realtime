package entity

import (
	"time"
)

type Message struct {
	ID          int64      `json:"id"`
	UserChatID  int64      `json:"user_chat_id"`
	MessageText string     `json:"message_text"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
}
