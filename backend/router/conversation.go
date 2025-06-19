package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/sdkim96/rag-backend/models"
)

// rootString is a constant representing the root message ID for new conversations.
const rootString string = "00000000-0000-0000-0000-000000000000"

func InitConversationRouter(rg *gin.RouterGroup) {
	rg.POST("/new", newConversationEndpoint)
	rg.GET("/", getConversationsEndpoint)
}

// ## New Conversation Endpoint
//
// This endpoint creates a new conversation and returns thre conversation ID and parent message ID.
func newConversationEndpoint(c *gin.Context) {

	resp := &models.APIResponse{
		Status:  "ok",
		Message: "New conversation created",
		Code:    200,
		Data: models.NewConversationDTO{
			ConversationID:  uuid.New().String(),
			ParentMessageID: rootString,
		},
	}
	c.JSON(200, resp)
}

func getConversationsEndpoint(c *gin.Context) {

	errorResp := &models.APIResponse{
		Status:  "error",
		Message: "Not implemented",
		Code:    501,
		Data:    nil,
	}
	c.JSON(501, errorResp)

}
