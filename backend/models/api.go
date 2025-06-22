package models

import "time"

type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
}

type NewConversationDTO struct {
	ConversationID  string `json:"conversation_id"`
	ParentMessageID string `json:"parent_message_id"`
}

func MockNewConversation() *NewConversationDTO {
	return &NewConversationDTO{
		ConversationID:  "mocked_conversation_id",
		ParentMessageID: "mocked_parent_message_id",
	}
}

type GetConversationsDTO struct {
	Conversations []Conversation `json:"conversations"`
}

func MockGetConversation() *GetConversationsDTO {
	return &GetConversationsDTO{
		Conversations: []Conversation{
			{
				ID:        "mocked_conversation_id_1",
				Title:     "Mocked Conversation 1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{
				ID:        "mocked_conversation_id_2",
				Title:     "Mocked Conversation 2",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
		},
	}
}
