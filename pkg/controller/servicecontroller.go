package controller

import (
	"net/http"

	restful "github.com/emicklei/go-restful/v3"
	"github.com/xunlbz/go-restful-template/pkg/lib"
	"github.com/xunlbz/go-restful-template/pkg/log"
)

type unitService struct {
}

func init() {
	service := new(unitService)
	registerService("unitService", service)
}

func (a *unitService) init() {
}

//PUT http://localhost:8889/api/v1/service/restart/{name}
func (a *unitService) restart(request *restful.Request, response *restful.Response) {
	target := request.PathParameter("name")
	log.Infof("restart service %s", target)
	err := lib.RestartUnit(target)
	if err != nil {
		log.Error(err)
		writeErrorAsJSON(response, http.StatusBadRequest, errorServiceRestart)
		return
	}
	res := BoolRes{Result: true}
	response.WriteAsJson(res)
}

// Build  create a new service that can handle REST requests for collector resources.
func (a *unitService) Build() *restful.WebService {
	a.init()
	ws := new(restful.WebService)
	ws.Path("/api/v1/service")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	tags := []string{"服务管理"}

	ws.Route(ws.PUT("/restart/{name}").To(a.restart).
		// docs
		Doc("restart service by name").
		Param(SecurityHeader).
		Metadata(KeyOpenAPITags, tags).
		Writes([]interface{}{}).
		Returns(200, "OK", []interface{}{}))

	return ws
}
