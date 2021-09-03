package main

import (
	"flag"
	"fmt"
	"net/http"

	assetfs "github.com/elazarl/go-bindata-assetfs"
	restful "github.com/emicklei/go-restful/v3"
	bindata "github.com/xunlbz/go-restful-template/bindata"
	"github.com/xunlbz/go-restful-template/cmd"
	"github.com/xunlbz/go-restful-template/pkg/collector"
	"github.com/xunlbz/go-restful-template/pkg/controller"
	"github.com/xunlbz/go-restful-template/pkg/database"
	"github.com/xunlbz/go-restful-template/pkg/log"
	"github.com/xunlbz/go-restful-template/pkg/websocket"
)

var (
	logFile  string
	logLevel string
	port     int
	help     bool
	swagger  bool
	version  bool
)

func init() {
	flag.StringVar(&logFile, "log_file", "", "set log file")
	flag.StringVar(&logLevel, "log_level", "info", "set log level debug, info, error")
	flag.IntVar(&port, "port", 8889, "set web port,default is 8889")
	flag.BoolVar(&help, "help", false, "help")
	flag.BoolVar(&swagger, "swagger", false, "swagger ui enable")
	flag.BoolVar(&version, "version", false, "get version")
}

func main() {
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	if version {
		fmt.Println(cmd.Version())
		return
	}
	startServer()

}

func registerService(wsContainer *restful.Container) {
	//1.apiService
	controller.RegisterWebService(wsContainer)
	//2.staticFileService
	log.Infof("register staticFileService")
	if swagger {
		// fs := assetfs.AssetFS{
		// 	Asset:     swagbindata.Asset,
		// 	AssetDir:  swagbindata.AssetDir,
		// 	AssetInfo: swagbindata.AssetInfo,
		// 	Prefix:    "swaggerui",
		// }
		// http.Handle("/", http.FileServer(&fs))
		http.Handle("/", http.FileServer(http.Dir("./swaggerui")))
	} else {
		fs := assetfs.AssetFS{
			Asset:     bindata.Asset,
			AssetDir:  bindata.AssetDir,
			AssetInfo: bindata.AssetInfo,
			Prefix:    "dist",
		}
		http.Handle("/", http.FileServer(&fs))
	}
	//3.websocket
	log.Info("register websocket")
	http.HandleFunc("/ws", websocket.HandleServer)
}

// Global Filter
func globalFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	log.Debugf("[Request starting] %s,%s", req.Request.Method, req.Request.URL)
	chain.ProcessFilter(req, resp)
	log.Debugf("[Request end] %s,%s, [Response status] %d", req.Request.Method, req.Request.URL, resp.StatusCode())

}

func startServer() {
	wsContainer := restful.DefaultContainer
	//1.setLoggerLevel
	log.SetLoggerConfig(logLevel, logFile)
	log.Info("edge admin starting!")
	//2.Create or Connect sqlite database
	database.Open()
	//3.Rigster Colloector
	collector.Register()
	//4.add global filter before any webservice
	restful.Filter(globalFilter)
	//5.Add container filter to enable CORS
	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{"X-My-Header"},
		AllowedHeaders: []string{"Content-Type", "Accept", "Authorization"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		CookiesAllowed: false,
		Container:      wsContainer}
	wsContainer.Filter(cors.Filter)
	//6.Add container filter to respond to OPTIONS
	wsContainer.Filter(wsContainer.OPTIONSFilter)
	//7.registerService and add webServiceFileter
	registerService(wsContainer)
	log.Infof("edge admin started listening on localhost:%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
