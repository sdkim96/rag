package models

type APIResponse struct {
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Code    int       `json:"code"`
	Data    DataField `json:"data,omitempty"`
}

type DataField interface {
	Mock()
}

// GET /token
type TokenData struct {
	Token string `json:"token"`
}

func (t *TokenData) Mock() {}

// GET /me
type UserData struct {
	UserName string `json:"user_name"`
}

func (u *UserData) Mock() {}

// POST /conversations/new
type NewConversationData struct {
	ConversationID  string `json:"conversation_id"`
	ParentMessageID string `json:"parent_message_id"`
}

func (n *NewConversationData) Mock() {}

// GET /conversations
type GetConversationsData struct {
	Conversations []*ConversationMeta `json:"conversations"`
}

func (g *GetConversationsData) Mock() {}

// Get /conversations/:ConversationID
type GetConversationByIDData struct {
	Conversation
}

func (g *GetConversationByIDData) Mock() {}
