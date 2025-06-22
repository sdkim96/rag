package svc

import (
	"log"

	"github.com/sdkim96/rag-backend/db"
	"github.com/sdkim96/rag-backend/models"
)

// `GetConversations` retrieves all conversations binded to a specific user.
func GetConversations(username string) *models.GetConversationsDTO {

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
		cvsDTO.ID = cvs.ID
		cvsDTO.Title = cvs.Title
		cvsDTO.CreatedAt = cvs.CreatedAt
		cvsDTO.UpdatedAt = cvs.UpdatedAt

		// This `append` operation returns pointer to slice.
		conversationDTOs = append(conversationDTOs, cvsDTO)

	}
	log.Printf("Retrived %d conversations for user: %s", len(conversationDTOs), username)

	return &models.GetConversationsDTO{
		Conversations: conversationDTOs,
	}
}

func GetConversationByID(username string, conversationID string) any {

	return nil
}
