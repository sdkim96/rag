package router

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	cst "github.com/sdkim96/rag-backend/constants"
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
	core.EncodeToken("sdkim96", cst.InternalIssuer)
	c.JSON(200, gin.H{
		"status": "ok",
	})
}
func oauthGithubLoginEndpoint(c *gin.Context) {

	cfg := core.GetAppConfig()

	authurl := cst.GithubAuthorizationURL
	authurl += "?client_id=" + cfg.AuthConfig.GithubClientID
	authurl += "&redirect_uri=" + cfg.AuthConfig.GithubRedirectURL
	authurl += "&state=" + cfg.AuthConfig.GithubState
	for _, scp := range cfg.AuthConfig.GithubScopes {
		authurl += "&scope=" + scp
	}

	c.Redirect(302, authurl)
}
func oauthGithubCallbackEndpoint(c *gin.Context) {

	req := &models.OAuthCallbackRequest{
		Code:  c.Query("code"),
		State: c.Query("state"),
	}

	cfg := core.GetAppConfig()

	if req.State != cfg.AuthConfig.GithubState {
		log.Println("State mismatch in OAuth callback")
		c.JSON(400, &models.APIResponse{
			Status:  cst.Error,
			Message: "State mismatch in OAuth callback",
			Code:    400,
			Data:    nil,
		})
		return
	}

	tokenUrl := cst.GithubTokenURL
	tokenUrl += "?client_id=" + cfg.AuthConfig.GithubClientID
	tokenUrl += "&client_secret=" + cfg.AuthConfig.GithubClientSecret
	tokenUrl += "&code=" + req.Code
	tokenUrl += "&redirect_uri=" + cfg.AuthConfig.GithubRedirectURL

	request, err := http.NewRequest("POST", tokenUrl, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		c.JSON(500, &models.APIResponse{
			Status:  cst.Error,
			Message: "Failed to create request",
			Code:    500,
			Data:    nil,
		})
		return
	}
	request.Header.Add("Accept", "application/json")
	client := http.DefaultClient

	resp, err := client.Do(request)
	if err != nil {
		log.Printf("Failed to get access token: %v", err)
		c.JSON(500, &models.APIResponse{
			Status:  cst.Error,
			Message: "Failed to get access token",
			Code:    500,
			Data:    nil,
		})
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get access token, status code: %d", resp.StatusCode)
		c.JSON(500, &models.APIResponse{
			Status:  cst.Error,
			Message: "Failed to get access token",
			Code:    500,
			Data:    nil,
		})
		return
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		c.JSON(500, &models.APIResponse{
			Status:  cst.Error,
			Message: cst.InternalServerError,
			Code:    500,
			Data:    nil,
		})
		return
	}

	ghToken := &models.GithubToken{}
	err = json.Unmarshal(bodyBytes, ghToken)
	if err != nil {
		log.Printf("Failed to unmarshal GitHub token: %v", err)
		c.JSON(500, &models.APIResponse{
			Status:  cst.Error,
			Message: cst.InternalServerError,
			Code:    500,
			Data:    nil,
		})
		return
	}
	//TODO: Replace the gh token to the models.TokenDTO
	tkn, err := core.EncodeToken(ghToken.AccessToken, cst.GithubIssuer)
	if err != nil {
		log.Printf("Failed to encode token: %v", err)
		c.JSON(500, &models.APIResponse{
			Status:  cst.Error,
			Message: cst.InternalServerError,
			Code:    500,
			Data:    nil,
		})
		return
	}
	c.JSON(200, &models.APIResponse{
		Status:  cst.Ok,
		Message: cst.NewTokenCreated,
		Code:    200,
		Data: &models.TokenDTO{
			Token: tkn,
		},
	})

}
func meEndpoint(c *gin.Context) {
	// This endpoint requires authentication
	username, exist := c.Get("UserName")
	if !exist {
		log.Println("UserName not found in gin context")
		c.AbortWithStatusJSON(401, &models.APIResponse{
			Status:  cst.Error,
			Message: cst.UnAuthorizedUserError,
			Code:    401,
			Data:    nil,
		})
		return
	}
	usernameStr, ok := username.(string)
	if !ok {
		log.Printf("Type assertion failed for username: %v", username)
		// Return an error response if the type assertion fails
		c.AbortWithStatusJSON(500, &models.APIResponse{
			Status:  cst.Error,
			Message: cst.InternalServerError,
			Code:    500,
			Data:    nil,
		})
		return
	}
	log.Printf("Authenticated user: %s", usernameStr)

	c.JSON(200, &models.APIResponse{
		Status:  cst.Ok,
		Message: "User information retrieved successfully",
		Code:    200,
		Data: &models.UserDTO{
			UserName: usernameStr,
		},
	})
}
