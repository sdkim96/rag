package middleware

import (
	"log"
	"strings"

	"github.com/gin-gonic/gin"

	cst "github.com/sdkim96/rag-backend/constants"
	"github.com/sdkim96/rag-backend/core"
	"github.com/sdkim96/rag-backend/models"
)

func BasicMiddleWare(c *gin.Context) {
	log.Println("Middleware called")
	c.Next()
}

func MockAuthMiddleware(c *gin.Context) {
	user := core.GetAppConfig().Mock.MockUser
	c.Set("UserName", user)
	c.Next()
}

func AuthMiddleware(c *gin.Context) {
	var UserName string = ""
	log.Println("This request needs authentication")

	baseHeaderPrefix := core.GetAppConfig().AuthConfig.HeaderPrefix
	tokenString := c.GetHeader("Authorization")

	ok := strings.HasPrefix(tokenString, baseHeaderPrefix)
	if !ok {
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

	// TODO: Cache the token information on redis or local session not to fetch user info from GitHub every time
	switch tkn.Issuer {
	case cst.GithubIssuer:
		ghUser, err := core.GetUserInfoFromGithub(tkn.UserName)
		if err != nil {
			log.Printf("Failed to fetch github: %v", err)
			c.AbortWithStatusJSON(401, &models.APIResponse{
				Status:  cst.Error,
				Message: "Failed to fetch user info from GitHub",
				Code:    401,
				Data:    nil,
			})
			return
		}
		UserName = ghUser
	case cst.InternalIssuer:
		log.Printf("Internal Issuer: %s", tkn.UserName)
		UserName = tkn.UserName
	default:
		log.Printf("Internal Issuer")
		UserName = tkn.UserName
	}

	c.Set("UserName", UserName)
	c.Next()
}
