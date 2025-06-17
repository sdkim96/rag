package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sdkim96/rag-backend/core"
	"github.com/sdkim96/rag-backend/middleware"
)

func InitAuthRouter(rg *gin.RouterGroup) {
	rg.GET("login", middleware.BasicMiddleWare, loginEndpoint)
	rg.GET("me", middleware.AuthMiddleware, meEndpoint)

}
func loginEndpoint(c *gin.Context) {
	core.EncodeToken("sdkim96")
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
func meEndpoint(c *gin.Context) {
	// This endpoint requires authentication
	c.JSON(200, gin.H{
		"status": "authenticated",
	})
}
