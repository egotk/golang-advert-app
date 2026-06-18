package corehttp

import (
	"net/http"
)

type Route struct {
	Method     string
	Path       string
	Handler    http.HandlerFunc
	Middleware []Middleware
}

func NewRoute(
	method string,
	path string,
	handler http.HandlerFunc,
	middleware ...Middleware,
) Route {
	return Route{
		Method:     method,
		Path:       path,
		Handler:    handler,
		Middleware: middleware,
	}
}
