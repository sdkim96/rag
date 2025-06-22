package router

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/sdkim96/rag-backend/middleware"
	"github.com/sdkim96/rag-backend/models"
	"github.com/sdkim96/rag-backend/svc"
)

// rootString is a constant representing the root message ID for new conversations.
const rootString string = "00000000-0000-0000-0000-000000000000"

func InitConversationRouter(rg *gin.RouterGroup) {
	rg.POST("/new", newConversationEndpoint)
	rg.GET("/", middleware.AuthMiddleware, getConversationsEndpoint)
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

	username, exist := c.Get("UserName")
	if !exist {
		errorResp.Message = "인증되지 않은 사용자입니다."
		errorResp.Code = 401
		c.AbortWithStatusJSON(401, errorResp)
		return
	}
	usernameStr, ok := username.(string)
	if !ok {
		errorResp.Message = "유저 이름이 문자열이 아닙니다."
		errorResp.Code = 500
		c.AbortWithStatusJSON(500, errorResp)
		return
	}

	dto, err := svc.GetConversations(usernameStr)
	if err != nil {
		errorResp.Message = "대화 목록을 가져오는 중 오류가 발생했습니다."
		errorResp.Code = 500
		c.AbortWithStatusJSON(500, errorResp)
		return
	}

	c.JSON(200, &models.APIResponse{
		Status:  "ok",
		Message: "대화 목록을 가져왔습니다.",
		Code:    200,
		Data:    dto,
	})

}
