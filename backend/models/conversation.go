package models

type DTO interface {
	Mock() DTO
}

type NewConversationDTO struct {
	ConversationID  string `json:"conversation_id"`
	ParentMessageID string `json:"parent_message_id"`
}

func (n NewConversationDTO) Mock() DTO {
	return NewConversationDTO{
		ConversationID:  "mocked_conversation_id",
		ParentMessageID: "mocked_parent_message_id",
	}
}

type GetConversationsDTO struct {
}
