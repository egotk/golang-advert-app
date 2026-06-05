package corehttp

import (
	"fmt"
	"net/http"
)

type ApiVersion string

var (
	ApiV1 = ApiVersion("v1")
	ApiV2 = ApiVersion("v2")
)

type APIVersionRouter struct {
	*http.ServeMux
	apiVersion ApiVersion
}

func NewAPIVersionRouter(apiVersion ApiVersion) *APIVersionRouter {
	return &APIVersionRouter{
		ServeMux:   http.NewServeMux(),
		apiVersion: apiVersion,
	}
}

func (r *APIVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)
		handler := ChainMiddleware(route.Handler, route.Middleware...)
		r.Handle(pattern, handler)
	}
}
