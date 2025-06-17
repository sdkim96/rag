package router

import "github.com/gin-gonic/gin"

func InitDefaultRouter(rg *gin.RouterGroup) {
	rg.GET("health", healthEndpoint)

}

func healthEndpoint(c *gin.Context) {
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
