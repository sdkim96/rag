package router

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/sdkim96/rag-backend/core"
	"github.com/sdkim96/rag-backend/middleware"
	"github.com/sdkim96/rag-backend/models"
)

func InitAuthRouter(rg *gin.RouterGroup) {
	rg.GET("login", middleware.BasicMiddleWare, loginEndpoint)
	rg.GET("oauth/github/login", middleware.BasicMiddleWare, oauthGithubLoginEndpoint)
	rg.GET("oauth/github/callback", middleware.BasicMiddleWare, oauthGithubCallbackEndpoint)

	rg.GET("me", middleware.AuthMiddleware, meEndpoint)

}
func loginEndpoint(c *gin.Context) {
	core.EncodeToken("sdkim96")
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
func oauthGithubLoginEndpoint(c *gin.Context) {

	cfg := core.GetAppConfig()

	url := "https://github.com/login/oauth/authorize"
	url += "?client_id=" + cfg.AuthConfig.GithubClientID
	url += "&redirect_uri=" + cfg.AuthConfig.GithubRedirectURL
	url += "&state=" + cfg.AuthConfig.GithubState
	for _, scp := range cfg.AuthConfig.GithubScopes {
		url += "&scope=" + scp
	}

	c.Redirect(302, url)
}
func oauthGithubCallbackEndpoint(c *gin.Context) {

	req := &models.OAuthCallbackRequest{
		Code:  c.Query("code"),
		State: c.Query("state"),
	}

	cfg := core.GetAppConfig()

	if req.State != cfg.AuthConfig.GithubState {
		log.Println("State mismatch in OAuth callback")
		c.JSON(400, gin.H{
			"error": "Invalid state parameter",
		})
		return
	}

	tokenUrl := "https://github.com/login/oauth/access_token"

	clientID := cfg.AuthConfig.GithubClientID
	clientSecret := cfg.AuthConfig.GithubClientSecret
	code := req.Code
	redirectURI := cfg.AuthConfig.GithubRedirectURL

	tokenUrl += "?client_id=" + clientID
	tokenUrl += "&client_secret=" + clientSecret
	tokenUrl += "&code=" + code
	tokenUrl += "&redirect_uri=" + redirectURI

	request, err := http.NewRequest("POST", tokenUrl, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		c.JSON(500, gin.H{
			"error": "Failed to create request",
		})
		return
	}
	request.Header.Add("Accept", "application/json")

	client := http.DefaultClient

	resp, err := client.Do(request) // This should be replaced with actual HTTP request logic
	if err != nil {
		log.Printf("Failed to get access token: %v", err)
		c.JSON(500, gin.H{
			"error": "Failed to get access token",
		})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get access token, status code: %d", resp.StatusCode)
		c.JSON(500, gin.H{
			"error": "Failed to get access token",
		})
		return
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		c.JSON(500, gin.H{
			"error": "Failed to read response body",
		})
		return
	}

	ghToken := &models.GithubToken{}
	err = json.Unmarshal(bodyBytes, ghToken)
	if err != nil {
		log.Printf("Failed to unmarshal GitHub token: %v", err)
		c.JSON(500, gin.H{
			"error": "Failed to parse access token",
		})
		return
	}
	//TODO: Replace the gh token to the models.TokenDTO
	c.JSON(200, ghToken)

}
func meEndpoint(c *gin.Context) {
	// This endpoint requires authentication
	c.JSON(200, gin.H{
		"status": "authenticated",
	})
}
