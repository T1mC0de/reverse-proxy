package router

import (
	"build-your-own-reverse-proxy/src/internal/config"
	"strings"
)

type Router struct {
	routes []config.Route
}

func NewRouter(routes []config.Route) *Router {
	return &Router{
		routes: routes,
	}
}

func (r *Router) Match(path string) *config.Route {
	for _, route := range r.routes {
		if strings.HasPrefix(path, route.Path) {
			return &route
		}
	}
	return nil
}













