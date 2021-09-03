package controller

import (
	"net/http"

	restful "github.com/emicklei/go-restful/v3"
	"github.com/xunlbz/go-restful-template/pkg/lib"
	"github.com/xunlbz/go-restful-template/pkg/log"
	"github.com/xunlbz/go-restful-template/pkg/model"
)

type userService struct {
	userDao model.UserDao
}

func init() {
	service := new(userService)
	registerService("userService", service)
}

func (a *userService) init() {
	a.userDao = model.NewUserDao()
}

//PUT http://localhost:8889/api/v1/user
func (a *userService) updatePassword(request *restful.Request, response *restful.Response) {
	user := new(User)
	err := request.ReadEntity(user)
	if err != nil {
		log.Error(err)
		writeErrorAsJSON(response, http.StatusBadRequest, errorParseParams)
		return
	}
	if !a.userDao.Exists(user.Username) {
		writeErrorAsJSON(response, http.StatusBadRequest, errorUserNotExists)
		return
	}
	if user.Password != "" {
		a.userDao.Update(user.Username, user.Password)
		token, _ := lib.CreateToken(user.Username)
		response.WriteAsJson(token)
		return
	}
	writeErrorAsJSON(response, http.StatusBadRequest, errorUsernamePassword)
}

// Build  create a new service that can handle REST requests for collector resources.
func (a *userService) Build() *restful.WebService {
	a.init()
	ws := new(restful.WebService)
	path := "/api/v1/user"
	ws.Path(path)
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	tags := []string{"用户管理"}

	ws.Route(ws.PUT("/").To(a.updatePassword).
		// docs
		Doc("update user password").
		Metadata(KeyOpenAPITags, tags).
		Param(SecurityHeader).
		Reads(User{}).
		Writes(lib.Token{}).
		Returns(200, "OK", lib.Token{}))
	return ws
}
