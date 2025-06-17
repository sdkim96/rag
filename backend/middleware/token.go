package middleware

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/sdkim96/rag-backend/core"
)

func BasicMiddleWare(c *gin.Context) {
	log.Println("Middleware called")
	c.Next()
}

func AuthMiddleware(c *gin.Context) {

	log.Println("인증 필요한 요청입니다.")

	baseHeaderPrefix := core.GetAppConfig().AuthConfig.HeaderPrefix
	tokenString := c.GetHeader("Authorization")
	whetherHas := strings.HasPrefix(tokenString, baseHeaderPrefix)
	if !whetherHas {
		log.Println("인증 헤더가 없습니다.")
		c.AbortWithStatusJSON(401, gin.H{
			"error": "인증 헤더가 없습니다.",
		})
		return
	}
	tokenString = strings.TrimPrefix(tokenString, baseHeaderPrefix)

	tkn, err := core.DecodeToken(tokenString)
	if err != nil {
		log.Printf("토큰 파싱 실패: %v", err)
		c.AbortWithStatusJSON(401, gin.H{
			"error": "토큰 파싱 실패",
		})
		return
	}

	// tkn.Username을 context에 저장
	c.Set("username", tkn.Username)
	c.Next()
}
