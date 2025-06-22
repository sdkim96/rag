package svc

import (
	"log"

	"github.com/sdkim96/rag-backend/db"
	"github.com/sdkim96/rag-backend/models"
)

func GetConversations(username string) (*models.GetConversationsDTO, error) {

	h := db.NewHandler()

	cvs := &db.Conversation{}
	cvsDTO := &models.Conversation{}
	h.First(cvs, "username = ?", username)

	cvsDTO.ID = cvs.ID
	cvsDTO.Title = cvs.Title
	cvsDTO.CreatedAt = cvs.CreatedAt
	cvsDTO.UpdatedAt = cvs.UpdatedAt

	log.Printf(
		"%s 유저에 대해 다음 대화 확인: %s, %s, %s, %s",
		username,
		cvsDTO.ID,
		cvsDTO.Title,
		cvsDTO.CreatedAt,
		cvsDTO.UpdatedAt,
	)

	return models.MockGetConversation(), nil
}
