package controller

import (
	"sync"

	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	restful "github.com/emicklei/go-restful/v3"
	spec "github.com/go-openapi/spec"
	log "github.com/xunlbz/go-restful-template/pkg/log"
)

var services = make(map[string]Service)

// Service
type Service interface {
	Build() *restful.WebService
}

func registerService(name string, service Service) {
	if _, ok := services[name]; ok {
		log.Errorf("service %s exists, please check\n", name)
		return
	}
	services[name] = service
}

//RegisterWebService  add webservice
func RegisterWebService(wsContainer *restful.Container) {
	log.Info("RegisterWebService starting")
	wsContainer.Add(pingService())
	var sw sync.WaitGroup
	for name, service := range services {
		sw.Add(1)
		go func(name string, service Service) {
			defer sw.Done()
			ws := service.Build()
			wsContainer.Add(ws)
			log.Infof("register service %s", name)
		}(name, service)
	}
	sw.Wait()
	wsContainer.Add(swaggerService())
	log.Info("RegisterWebService end")
	addFilter(wsContainer)
}

func addFilter(wsContainer *restful.Container) {
	log.Info("add Filter start")
	handleJWTAuthFilter(wsContainer)
	log.Info("add Filter end")
}

// 服务状态检测接口
func pingService() *restful.WebService {
	ws := new(restful.WebService)
	path := "/api/v1/ping"
	addAuthExcludePath(path)
	ws.Path(path)
	ws.Consumes(restful.MIME_JSON)
	ws.Produces(restful.MIME_JSON)
	tags := []string{"系统状态"}
	f := func(request *restful.Request, response *restful.Response) {
		log.Debug("get ping")
		response.Write([]byte("pong"))
	}
	ws.Route(ws.GET("/").To(f).
		// docs
		Doc("get ping").
		Metadata(restfulspec.KeyOpenAPITags, tags))
	return ws
}

func swaggerService() *restful.WebService {
	//4.swagger
	apiPath := "/swagger.json"
	config := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(),
		APIPath:                       apiPath,
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}
	addAuthExcludePath(apiPath)
	return restfulspec.NewOpenAPIService(config)
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "Edge Admin API",
			Description: "Resource for EdgeAdminApi",
			Version:     "1.0.0",
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  "www.wayclouds.com",
					Email: "info@wayclouds.com",
				},
			},
		},
	}
	swo.Tags = []spec.Tag{{TagProps: spec.TagProps{
		Name:        "Edge Admin Api",
		Description: "Resource for Edge Admin Api"}}}

}
