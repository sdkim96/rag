package db

import (
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	ID           string `gorm:"primaryKey"`
	UserName     string `gorm:"uniqueIndex"`
	PasswordHash string
	Email        string `gorm:"uniqueIndex"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type Conversation struct {
	ID            string `gorm:"primaryKey"`
	Title         string
	UserID        string
	RootMessageID string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type ConversationMessage struct {
	ID              string `gorm:"primaryKey"`
	ConversationID  string
	ParentMessageID string
	Content         string
	Role            string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

type DBHandler struct {
	*gorm.DB
}

func NewHandler() *DBHandler {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}

	return &DBHandler{DB: db}
}

func (h *DBHandler) MigrateTables() {

	err := h.AutoMigrate(
		&User{},
		&Conversation{},
		&ConversationMessage{},
	)
	if err != nil {
		panic("failed to migrate database tables: " + err.Error())
	}
}
