package models

type OAuthCallbackRequest struct {
	Code  string `json:"code"`
	State string `json:"state"`
}

type GetConversationByIDReq struct {
	ConversationID string `json:"conversation_id"`
}
