package core

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	cst "github.com/sdkim96/rag-backend/constants"
)

type ClaimsDecoded struct {
	UserName  string    `json:"username"`
	ExpiredAt time.Time `json:"expired_at"`
	Issuer    string    `json:"issuer"`
}

func EncodeToken(username string, issuer string) (string, error) {

	var (
		bytesString []byte
		expAt       int64
		tkn         *jwt.Token
		s           string
		err         error
	)

	s = ""
	cfg := GetAppConfig()

	bytesString = []byte(cfg.TokenConfig.Secret)
	expAt = time.Now().Unix() + cfg.TokenConfig.Duration

	tkn = jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.RegisteredClaims{
			Subject:   username,
			Issuer:    issuer,
			ExpiresAt: jwt.NewNumericDate(time.Unix(expAt, 0)),
		},
	)
	s, err = tkn.SignedString(bytesString)
	log.Printf("%s 님의 JWT 토큰: %s", username, s)
	if err != nil {
		return s, err
	}
	return s, nil

}

func DecodeToken(token string) (*ClaimsDecoded, error) {
	tkn, err := jwt.Parse(
		token,
		func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				log.Printf("Unexpected signing method: %v", token.Header["alg"])
				return nil, jwt.ErrSignatureInvalid
			}
			cfg := GetAppConfig()
			return []byte(cfg.TokenConfig.Secret), nil
		},
	)
	if err != nil {
		log.Printf("Failed to parse JWT token: %v", err)
		return &ClaimsDecoded{}, err
	}
	// TODO: 여기서 토큰의 유효성을 검사하고, 필요한 정보를 추출하여 Token 구조체를 반환해야 합니다.
	sub, err := tkn.Claims.GetSubject()
	if err != nil {
		log.Printf("Failed to get subject from claims: %v", err)
		return &ClaimsDecoded{}, err
	}
	exp, err := tkn.Claims.GetExpirationTime()
	if err != nil {
		log.Printf("Failed to get expiration time from the claims: %v", err)
		return &ClaimsDecoded{}, err
	}
	iss, err := tkn.Claims.GetIssuer()
	if err != nil {
		log.Printf("Failed to get issuer from the claims: %v", err)
		return &ClaimsDecoded{}, err
	}

	log.Printf("claims: %s, %s, %s", sub, exp.Time, iss)
	return &ClaimsDecoded{
		UserName:  sub,
		ExpiredAt: exp.Time,
		Issuer:    iss,
	}, nil
}

func GetUserInfoFromGithub(githubToken string) (string, error) {

	responseModel := make(map[string]interface{})

	request, err := http.NewRequest("GET", cst.GithubUserInfoURL, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		return "", err
	}
	request.Header.Add("Authorization", "Bearer "+githubToken)
	request.Header.Add("Accept", "application/vnd.github+json")
	request.Header.Add("X-GitHub-Api-Version", "2022-11-28")

	httpClient := http.DefaultClient
	resp, err := httpClient.Do(request)
	if err != nil {
		log.Printf("Failed to get user info from GitHub: %v", err)
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to get user info from GitHub, status code: %d", resp.StatusCode)
		return "", http.ErrHandlerTimeout
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		return "", err
	}
	err = json.Unmarshal(bodyBytes, &responseModel)
	if err != nil {
		log.Printf("Failed to unmarshal response body: %v", err)
		return "", err
	}
	username, ok := responseModel["name"].(string)
	if !ok {
		log.Println("Failed to get username from response body")
		return "", fmt.Errorf("failed to get username from response body")
	}
	log.Printf("GitHub username: %s", username)
	return username, nil

}
