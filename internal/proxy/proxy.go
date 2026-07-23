package proxy

import (
	"build-your-own-reverse-proxy/internal/router"
	"build-your-own-reverse-proxy/internal/forwarder"
	"log/slog"
	"net/http"
	"context"
)

type Proxy struct {
	router    *router.Router
	forwarder *forwarder.Forwarder
	logger    *slog.Logger
}

func NewProxy(router *router.Router, forwarder *forwarder.Forwarder, logger *slog.Logger) *Proxy {
	return &Proxy{
		router:    router,
		forwarder: forwarder,
		logger:    logger,
	}
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	route := p.router.Match(r.URL.Path)
	if route == nil {
		p.logger.Debug("Route not found", "path", r.URL.Path)
		http.Error(w, "Route not found", http.StatusNotFound)
		return
	}

	p.logger.Debug("Route detected", "route_path", route.Path, "upstream", route.Upstreams[0], "request_path", r.URL.Path)

	ctx, cancel := context.WithTimeout(r.Context(), route.Timeout)
	defer cancel()

	p.forwarder.Forward(ctx, w, r, route.Upstreams[0])
}


