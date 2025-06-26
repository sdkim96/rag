package svc

import (
	"log"

	"github.com/sdkim96/rag-backend/db"
	"github.com/sdkim96/rag-backend/models"
)

// `GetConversations` retrieves all conversations binded to a specific user.
func GetConversations(username string) *models.GetConversationsData {

	h := db.NewHandler()
	conversations := make([]*db.Conversation, 0)
	conversationDTOs := make([]*models.ConversationMeta, 0)

	// This query retrives all conversations except for the soft-deleted ones.
	h.Where(
		&db.Conversation{
			UserName: username,
		},
	).Find(&conversations)

	for _, cvs := range conversations {
		cvsDTO := &models.ConversationMeta{}
		cvsDTO.ConversationID = cvs.ID
		cvsDTO.Title = cvs.Title
		cvsDTO.CreatedAt = cvs.CreatedAt
		cvsDTO.UpdatedAt = cvs.UpdatedAt

		// This `append` operation returns pointer to slice.
		conversationDTOs = append(conversationDTOs, cvsDTO)

	}
	log.Printf("Retrived %d conversations for user: %s", len(conversationDTOs), username)

	return &models.GetConversationsData{
		Conversations: conversationDTOs,
	}
}

// SELECT
// t1.id,
// t1.title,
// t1.root_message_id,
// t1.created_at,
// t1.updated_at,
// t2.id as message_id,
// t2.parent_message_id,
// t2.content,
// t2.role,
// t2.created_at as message_created_at,
// t2.updated_at as message_updated_at
// FROM conversations t1
// LEFT JOIN conversation_messages t2
// ON t1.id=t2.conversation_id
// WHERE t1.user_name="" AND t1.id ="";
func GetConversationByID(username string, conversationID string) (*models.GetConversationByIDData, error) {

	h := db.NewHandler()
	results := make(map[string]interface{})

	err := h.Table("conversations").
		Select("conversations.*, conversation_messages.id as message_id, conversation_messages.content, ...").
		Joins("LEFT JOIN conversation_messages ON conversations.id = conversation_messages.conversation_id").
		Where("conversations.user_name = ? AND conversations.id = ?", username, conversationID).
		Scan(&results).Error
	if err != nil {
		log.Printf("Error retrieving conversation by ID: %v", err)
		return nil, err
	}

	return nil, nil
}
