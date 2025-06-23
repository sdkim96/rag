package models

import "time"

type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Code    int    `json:"code"`
	Data    any    `json:"data,omitempty"`
}

type TokenDTO struct {
	Token string `json:"token"`
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
	Conversations []*ConversationMeta `json:"conversations"`
}

func MockGetConversation() *GetConversationsDTO {
	return &GetConversationsDTO{
		Conversations: []*ConversationMeta{
			{
				ID:        "mocked_conversation_id_1",
				Title:     "Mocked Conversation 1",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			{},
		},
	}
}
