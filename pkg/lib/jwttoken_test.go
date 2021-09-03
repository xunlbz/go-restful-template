package lib

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"io"
	"log"
	"net/http"
	"strings"
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	restful "github.com/emicklei/go-restful/v3"
)

var (
	sharedSecret = []byte("shared-token")
)

// This example shows how to create a (Route) Filter that performs a JWT HS512 authentication.
//
// GET http://localhost:8080/secret
// and use shared-token as a shared secret.

func TestServer(t *testing.T) {
	ws := new(restful.WebService)
	ws.Route(ws.GET("/secret").Filter(authJWT).To(secretJWT))
	restful.Add(ws)
	genJWT("admin")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func keyFunc(token *jwt.Token) (i interface{}, err error) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		panic(err)
	}
	hash := sha256.Sum256(sharedSecret)
	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		panic(err)
	}
	log.Printf("signature: %x\n", sig)

	valid := ecdsa.VerifyASN1(&privateKey.PublicKey, hash[:], sig)
	log.Println("signature verified:", valid)
	return sharedSecret, nil
}

func secretJWT(req *restful.Request, resp *restful.Response) {
	io.WriteString(resp, "42")
}

// ValidJWT  valid JWT token
func validJWT(authHeader string) bool {
	log.Println("authHeader:", authHeader)
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return false
	}

	jwtToken := strings.Split(authHeader, " ")
	if len(jwtToken) < 2 {
		return false
	}
	log.Println("token:", jwtToken)
	parts := strings.Split(jwtToken[1], ".")
	err := jwt.SigningMethodHS512.Verify(strings.Join(parts[0:2], "."), parts[2], sharedSecret)
	if err != nil {
		return false
	}
	token, err := jwt.Parse(jwtToken[1], keyFunc)
	if err != nil {
		log.Println("parse error", err)
		return false
	}
	t := token.Claims.(jwt.MapClaims)
	log.Println(t)
	return true
}

func authJWT(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	authHeader := req.HeaderParameter("Authorization")

	if !validJWT(authHeader) {
		resp.WriteErrorString(401, "401: Not Authorized")
		return
	}

	chain.ProcessFilter(req, resp)
}

func genJWT(name string) {
	claims := &jwt.StandardClaims{
		Subject:   name,
		ExpiresAt: time.Now().Add(30 * time.Second).Unix(), // 过期时间，必须设置
		Issuer:    "wayclouds.com",                         // 可不必设置，也可以填充用户名，
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims) //生成token
	log.Println(token.SignedString(sharedSecret))
}
