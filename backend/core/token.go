package core

import (
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	UserName  string    `json:"username"`
	ExpiredAt time.Time `json:"expired_at"`
}

func EncodeToken(username string) (string, error) {

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

func DecodeToken(token string) (*Token, error) {
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
		log.Printf("토큰 파싱 실패: %v", err)
		return &Token{}, err
	}
	// TODO: 여기서 토큰의 유효성을 검사하고, 필요한 정보를 추출하여 Token 구조체를 반환해야 합니다.
	sub, err := tkn.Claims.GetSubject()
	if err != nil {
		log.Printf("토큰에서 subject 추출 실패: %v", err)
		return &Token{}, err
	}
	exp, err := tkn.Claims.GetExpirationTime()
	if err != nil {
		log.Printf("토큰에서 만료 시간 추출 실패: %v", err)
		return &Token{}, err
	}

	log.Printf("토큰에서 추출한 subject: %s, %s", sub, exp.Time)
	return &Token{
		UserName:  sub,
		ExpiredAt: exp.Time,
	}, nil
}
