package controller

import (
	"net/http"

	restful "github.com/emicklei/go-restful/v3"
	"github.com/xunlbz/go-restful-template/pkg/collector"
)

type metricsService struct {
	ec collector.EdgeCollector
}

func init() {
	service := new(metricsService)
	registerService("metricsService", service)
}

func (a *metricsService) init() {
	// 初始化收集器
	a.ec = collector.NewEdgeCollector()
	a.ec.Collect()
}

//GET http://localhost:8889/api/v1/metrics?collect=
func (a *metricsService) handlerMetrics(request *restful.Request, response *restful.Response) {
	module := request.QueryParameter("collect")
	mcs := a.ec.GetModuleMitrics(module)
	if mcs.Label == "" {
		writeErrorAsJSON(response, http.StatusBadRequest, errorNoModule)
		return
	}
	response.WriteAsJson(mcs)
}

// Build  create a new service that can handle REST requests for collector resources.
func (a *metricsService) Build() *restful.WebService {
	a.init()
	ws := new(restful.WebService)
	ws.Path("/api/v1/metrics")
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	tags := []string{"指标监控"}
	ws.Route(ws.GET("/").To(a.handlerMetrics).
		// docs
		Doc("get metrics").
		Param(SecurityHeader).
		Param(restful.QueryParameter("collect", "[os,ap,network,terminal,poe,serivce,docker]")).
		Metadata(KeyOpenAPITags, tags).
		Writes([]interface{}{}).
		Returns(200, "OK", []interface{}{}))

	return ws
}
