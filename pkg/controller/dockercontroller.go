package controller

import (
	"net/http"
	"time"

	restful "github.com/emicklei/go-restful/v3"
	"github.com/xunlbz/go-restful-template/pkg/lib"
)

type dockerService struct {
	cli lib.DockerClient
}

func init() {
	service := new(dockerService)
	registerService("dockerService", service)
}

func (a *dockerService) init() {
	a.cli = lib.NewDockerClient()
}

//PUT http://localhost:8889/api/v1/docker/restart/{container}
func (a *dockerService) restart(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("container")
	err := a.cli.ContainerRestart(id, time.Second*10)
	if err != nil {
		writeErrorAsJSON(response, http.StatusBadRequest, errorDockerRestart)
		return
	}
	res := BoolRes{Result: true}
	response.WriteAsJson(res)
}

// Build  create a new service that can handle REST requests for collector resources.
func (a *dockerService) Build() *restful.WebService {
	a.init()
	ws := new(restful.WebService)
	ws.Path("/api/v1/docker")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	tags := []string{"容器管理"}
	ws.Route(ws.PUT("/restart/{container}").To(a.restart).
		// docs
		Doc("restart container by id").
		Param(SecurityHeader).
		Metadata(KeyOpenAPITags, tags).
		Writes(BoolRes{}).
		Returns(200, "OK", BoolRes{}))

	return ws
}
