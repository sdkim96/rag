package router

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	cst "github.com/sdkim96/rag-backend/constants"
	"github.com/sdkim96/rag-backend/middleware"
	"github.com/sdkim96/rag-backend/models"
	"github.com/sdkim96/rag-backend/svc"
)

func InitConversationsRouter(rg *gin.RouterGroup) {
	rg.POST("/new", newConversationEndpoint)
	rg.GET("/", middleware.AuthMiddleware, getConversationsEndpoint)
	rg.GET("/:ConversationID", middleware.AuthMiddleware, getConversationByIDEndpoint)
}

// ## New Conversation Endpoint
//
// This endpoint creates a new conversation and returns thre conversation ID and parent message ID.
func newConversationEndpoint(c *gin.Context) {

	resp := &models.APIResponse{
		Status:  "ok",
		Message: "New conversation created",
		Code:    200,
		Data: &models.NewConversationData{
			ConversationID:  uuid.New().String(),
			ParentMessageID: cst.RootMessageID,
		},
	}
	c.JSON(200, resp)
}

func getConversationsEndpoint(c *gin.Context) {

	errorResp := &models.APIResponse{
		Status:  cst.Error,
		Message: cst.InternalServerError,
		Code:    500,
		Data:    nil,
	}

	username, exist := c.Get("UserName")
	if !exist {
		log.Println("UserName not found in gin context")
		errorResp.Message = cst.UnAuthorizedUserError
		errorResp.Code = 401
		c.AbortWithStatusJSON(401, errorResp)
		return
	}
	usernameStr, ok := username.(string)
	if !ok {
		log.Printf("Type assertion failed for username: %v", username)
		c.AbortWithStatusJSON(500, errorResp)
		return
	}

	c.JSON(200, &models.APIResponse{
		Status:  cst.Ok,
		Message: cst.ConversationListRetrieved,
		Code:    200,
		Data:    svc.GetConversations(usernameStr),
	})

}

func getConversationByIDEndpoint(c *gin.Context) {
	errorResp := &models.APIResponse{
		Status:  cst.Error,
		Message: cst.InternalServerError,
		Code:    500,
		Data:    nil,
	}

	req := &models.GetConversationByIDReq{}
	err := c.ShouldBindUri(req)
	log.Printf("GetConversationByIDReq: %+v", req)

	if err != nil {
		errorResp.Message = cst.EntityError
		errorResp.Code = 422
		c.AbortWithStatusJSON(422, errorResp)
		return
	}
	username, exist := c.Get("UserName")
	if !exist {
		log.Println("UserName not found in gin context")
		errorResp.Message = cst.UnAuthorizedUserError
		errorResp.Code = 401
		c.AbortWithStatusJSON(401, errorResp)
		return
	}
	usernameStr, ok := username.(string)
	if !ok {
		log.Printf("Type assertion failed for username: %v", username)
		c.AbortWithStatusJSON(500, errorResp)
		return
	}
	data, err := svc.GetConversationByID(usernameStr, req.ConversationID)
	if err != nil {
		log.Printf("Error retrieving conversation by ID: %v", err)
		errorResp.Message = cst.InternalServerError
		c.AbortWithStatusJSON(500, errorResp)
		return
	}

	c.JSON(200, &models.APIResponse{
		Status:  cst.Ok,
		Message: cst.ConversationListRetrieved,
		Code:    200,
		Data:    data,
	})

}
