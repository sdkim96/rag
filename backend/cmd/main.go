package main

import (
	"github.com/gin-gonic/gin"

	"github.com/sdkim96/rag-backend/core"
	"github.com/sdkim96/rag-backend/db"
	"github.com/sdkim96/rag-backend/router"
)

func main() {

	engine := gin.Default()
	core.InitEnv(".backend.env")

	v1 := engine.Group("/api/v1")
	conversations := v1.Group("/conversations")

	router.InitDefaultRouter(v1)
	router.InitAuthRouter(v1)
	router.InitConversationsRouter(conversations)

	db.NewHandler().MigrateTables()

	engine.Run(":8080") // listen and serve on
}
