package entity

import (
	"time"
)

type UserChat struct {
	ID        int64      `json:"id"`
	From      int64      `json:"from"`
	To        int64      `json:"to"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
