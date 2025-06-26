package models

import "time"

type ConversationMeta struct {
	ConversationID string    `json:"conversation_id"`
	Title          string    `json:"title"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type FeedbackApplied struct {
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Message struct {
	ID          string           `json:"id"`
	ParentID    string           `json:"parent_id"`
	ChildrenIDs []string         `json:"children_ids"`
	Role        string           `json:"role"`
	Content     string           `json:"content"`
	Feedback    *FeedbackApplied `json:"feedback,omitempty"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}

type ConversationMappings struct {
	Messages      []Message `json:"messages"`
	RootMessageID string    `json:"root_message_id"`
}

type Conversation struct {
	ConversationMeta
	Mapping ConversationMappings `json:"mapping"`
}
