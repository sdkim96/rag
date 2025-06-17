package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sdkim96/rag-backend/router"
)

func main() {

	engine := gin.Default()
	v1 := engine.Group("/api/v1")

	router.InitDefaultRouter(v1)
	router.InitAuthRouter(v1)

	engine.Run(":8080") // listen and serve on
}
