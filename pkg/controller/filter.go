package controller

import (
	"net/http"
	"strings"

	restful "github.com/emicklei/go-restful/v3"
	lib "github.com/xunlbz/go-restful-template/pkg/lib"
	log "github.com/xunlbz/go-restful-template/pkg/log"
)

var excludePath = []string{}

// WebService Basic Auth Filter
func webserviceBasicAuthFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	log.Debugf("[webservice-basic-auth-filter (logger)] %s,%s", req.Request.Method, req.Request.URL)
	// usr/pwd = admin/admin
	u, p, ok := req.Request.BasicAuth()
	if !ok || u != "admin" || p != "admin" {
		resp.AddHeader("WWW-Authenticate", "Basic realm=Protected Area")
		resp.WriteHeaderAndJson(http.StatusUnauthorized, errorNoAuth, restful.MIME_JSON)
		log.Infof("[webservice-basic-auth-filter (logger)] %s,%s Not Authorized", req.Request.Method, req.Request.URL)
		return
	}
	chain.ProcessFilter(req, resp)
}

func webserviceJWTAuthFilter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	if checkExcludePath(req.Request.URL.String()) || req.Request.Method == http.MethodOptions {
		chain.ProcessFilter(req, resp)
		return
	}
	log.Debugf("[webservice-jwt-auth-filter (logger)] %s,%s", req.Request.Method, req.Request.URL)
	token := req.HeaderParameter("Authorization")
	log.Debugf("auth Header: %s", token)
	if strings.HasPrefix(token, "Bearer ") {
		jwtToken := strings.Split(token, " ")
		if len(jwtToken) < 2 {
			writeErrorAsJSON(resp, http.StatusUnauthorized, errorNoAuth)
			return
		}
		token = jwtToken[1]
	}

	if token == "" {
		writeErrorAsJSON(resp, http.StatusUnauthorized, errorNoAuth)
		return
	}

	c, err := lib.ParseToken(token)
	if err != nil {
		writeErrorAsJSON(resp, http.StatusUnauthorized, errorExpiredAuth)
		return
	}
	if c.Type != lib.AccessToken {
		writeErrorAsJSON(resp, http.StatusUnauthorized, errorNoAuth)
		return
	}
	log.Debugf("[webservice-jwt-auth-filter (logger)] auth success %s,%s", req.Request.Method, req.Request.URL)
	chain.ProcessFilter(req, resp)
}

func checkExcludePath(uri string) bool {
	for _, path := range excludePath {
		if strings.HasPrefix(uri, path) {
			return true
		}
	}
	return false
}

func addAuthExcludePath(uri string) {
	excludePath = append(excludePath, uri)
}

func handleJWTAuthFilter(wsContainer *restful.Container) {
	lib.GenSecretKey() //初始化jwt key
	wsContainer.Filter(webserviceJWTAuthFilter)
}
