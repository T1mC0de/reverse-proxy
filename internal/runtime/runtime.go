package runtime

import (
	"build-your-own-reverse-proxy/internal/config"
	"build-your-own-reverse-proxy/internal/forwarder"
	"build-your-own-reverse-proxy/internal/proxy"
	"build-your-own-reverse-proxy/internal/router"
	"log/slog"
	"net/http"
)

type Runtime struct {
	proxy *proxy.Proxy
}

func NewRuntime(proxy *proxy.Proxy) *Runtime {
	return &Runtime{
		proxy: proxy,
	}
}

func (r *Runtime) Build(cfg *config.Config, logger *slog.Logger) (*Runtime, error) {
	router := router.NewRouter(cfg.Routes)
	forwarder := forwarder.NewForwarder(&http.Transport{
		MaxIdleConns:        cfg.Transport.MaxIdleConns,
		IdleConnTimeout:     cfg.Transport.IdleConnTimeout,
		TLSHandshakeTimeout: cfg.Transport.TLSHandshakeTimeout,
	})

	proxy := proxy.NewProxy(router, forwarder, logger)

	

	return &Runtime{
		proxy: proxy,
	}, nil
}




