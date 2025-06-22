package models

import "time"

type ConversationMeta struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ConverSation struct {
	ConversationMeta
	Messages []*Message `json:"messages"`
}
