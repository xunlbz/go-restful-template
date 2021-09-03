package controller

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"

	restful "github.com/emicklei/go-restful/v3"
	"github.com/xunlbz/go-restful-template/pkg/lib"
	"github.com/xunlbz/go-restful-template/pkg/log"
	"github.com/xunlbz/go-restful-template/pkg/model"
)

type authService struct {
	userDao model.UserDao
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RefreshToken struct {
	Token string `json:"token"`
}

func init() {
	service := new(authService)
	registerService("authService", service)
}

func (a *authService) init() {
	a.userDao = model.NewUserDao()
}

//POST http://localhost:8889/api/v1/auth/admin
func (a *authService) register(request *restful.Request, response *restful.Response) {
	user := new(User)
	err := request.ReadEntity(user)
	if err != nil {
		log.Error(err)
		writeErrorAsJSON(response, http.StatusBadRequest, errorParseParams)
		return
	}
	admin := "admin"
	// only support register admin
	if a.userDao.Exists(admin) {
		writeErrorAsJSON(response, http.StatusBadRequest, errorUserExists)
		return
	}

	a.userDao.Insert(admin, user.Password)
	token, err := lib.CreateToken(admin)
	if err != nil {
		log.Error(err)
		writeErrorAsJSON(response, http.StatusInternalServerError, NewSystemError(err))
		return
	}
	response.WriteAsJson(token)

}

//POST http://localhost:8889/api/v1/auth/login
func (a *authService) login(request *restful.Request, response *restful.Response) {
	user := new(User)
	err := request.ReadEntity(user)
	if err != nil {
		log.Error(err)
		writeErrorAsJSON(response, http.StatusBadRequest, errorParseParams)
		return
	}
	dbuser := a.userDao.GetOne(user.Username)
	if dbuser.ID > 0 && fmt.Sprintf("%x", md5.Sum([]byte(user.Password))) == dbuser.Password {
		token, err := lib.CreateToken(user.Username)
		if err != nil {
			log.Error(err)
			writeErrorAsJSON(response, http.StatusInternalServerError, NewSystemError(err))
			return
		}
		response.WriteAsJson(token)
	} else {
		writeErrorAsJSON(response, http.StatusBadRequest, errorUsernamePassword)
	}
}

//PUT http://localhost:8889/api/v1/auth/token
func (a *authService) refresh(request *restful.Request, response *restful.Response) {
	rftoken := new(RefreshToken)
	err := request.ReadEntity(rftoken)
	log.Debugf("refresh Token is %s", rftoken.Token)
	if err != nil || rftoken.Token == "" {
		writeErrorAsJSON(response, http.StatusBadRequest, errorSetParams)
		return
	}
	c, err := lib.ParseToken(rftoken.Token)
	if err != nil {
		writeErrorAsJSON(response, http.StatusBadRequest, errorExpiredAuth)
		return
	}
	token, _ := lib.CreateToken(c.Subject)
	response.WriteAsJson(token)
}

//GET http://localhost:8889/api/v1/auth/auth
func (a *authService) auth(request *restful.Request, response *restful.Response) {
	authHeader := request.HeaderParameter("Authorization")
	log.Debugf("auth Header: %s", authHeader)
	token := request.HeaderParameter("Authorization")
	log.Debugf("auth Header: %s", token)
	if strings.HasPrefix(token, "Bearer ") {
		jwtToken := strings.Split(token, " ")
		if len(jwtToken) < 2 {
			writeErrorAsJSON(response, http.StatusUnauthorized, errorNoAuth)
			return
		}
		token = jwtToken[1]
	}

	if token == "" {
		writeErrorAsJSON(response, http.StatusUnauthorized, errorNoAuth)
		return
	}

	c, err := lib.ParseToken(token)
	if err != nil {
		writeErrorAsJSON(response, http.StatusUnauthorized, errorExpiredAuth)
		return
	}
	if c.Type != lib.AccessToken {
		writeErrorAsJSON(response, http.StatusUnauthorized, errorNoAuth)
		return
	}
	response.WriteHeaderAndEntity(http.StatusNoContent, nil)
}

// Build  create a new service that can handle REST requests for collector resources.
func (a *authService) Build() *restful.WebService {
	a.init()
	ws := new(restful.WebService)
	path := "/api/v1/auth"
	ws.Path(path)
	addAuthExcludePath(path)
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	tags := []string{"服务认证"}

	ws.Route(ws.POST("/login").To(a.login).
		// docs
		Doc("login by username/password").
		Metadata(KeyOpenAPITags, tags).
		Reads(User{}).
		Writes(lib.Token{}).
		Returns(200, "OK", lib.Token{}))
	ws.Route(ws.PUT("/token").To(a.refresh).
		// docs
		Doc("refresh token").
		Metadata(KeyOpenAPITags, tags).
		Reads(RefreshToken{}).
		Writes(lib.Token{}).
		Returns(200, "OK", lib.Token{}))
	ws.Route(ws.POST("/user").To(a.register).
		// docs
		Doc("register user").
		Metadata(KeyOpenAPITags, tags).
		Reads(User{}).
		Writes(lib.Token{}).
		Returns(200, "OK", lib.Token{}))
	ws.Route(ws.GET("/auth").To(a.auth).
		// docs
		Doc("validate token for thirdparty").
		Param(SecurityHeader).
		Metadata(KeyOpenAPITags, tags).
		Writes(nil).
		Returns(204, "", nil))

	return ws
}
