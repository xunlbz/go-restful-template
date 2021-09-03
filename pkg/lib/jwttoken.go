package lib

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/xunlbz/go-restful-template/pkg/log"
)

type TokenType int

type Token struct {
	AccessToken  string `json:"token" description:"Expired after 30 minutes "`
	RefreshToken string `json:"refreshToken" description:"Expired after 2 hours "`
}

var (
	AccessToken         TokenType = 0
	RefreshToken        TokenType = 1
	secretKey           []byte    = []byte("www.github.com")
	AccessTokenExpired            = time.Minute * 30
	RefreshTokenExpired           = time.Hour * 2
)

type Claims struct {
	Type TokenType
	jwt.StandardClaims
}

func GenSecretKey() {
	secretKey = []byte("edge_admin_template")
}

func CreateToken(name string) (t Token, err error) {
	claims := &Claims{
		Type: AccessToken,
		StandardClaims: jwt.StandardClaims{
			Subject:   name,
			ExpiresAt: time.Now().Add(AccessTokenExpired).Unix(), // 过期时间，必须设置
			Issuer:    "wayclouds.com",                           // 可不必设置，也可以填充用户名，
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims) //生成token
	accessToken, err := token.SignedString(secretKey)
	if err != nil {
		log.Errorf("sign token error: %v", err)
		return t, err
	}
	claims.Type = RefreshToken
	claims.StandardClaims.ExpiresAt = time.Now().Add(RefreshTokenExpired).Unix()
	token = jwt.NewWithClaims(jwt.SigningMethodHS512, claims) //生成token
	refreshToken, err := token.SignedString(secretKey)
	if err != nil {
		log.Errorf("sign token error %v", err)
		return t, err
	}
	t = Token{AccessToken: accessToken, RefreshToken: refreshToken}
	return t, nil
}

func secretKeyFunc(t *jwt.Token) (interface{}, error) {
	return secretKey, nil
}

func ParseToken(token string) (c Claims, err error) {
	claims := new(Claims)

	t, err := jwt.ParseWithClaims(token, claims, secretKeyFunc)
	if err != nil {
		return *claims, err
	}
	if !t.Valid {
		return *claims, err
	}
	return *claims, nil
}
