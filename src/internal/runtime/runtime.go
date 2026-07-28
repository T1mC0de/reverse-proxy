package runtime

import (
	"build-your-own-reverse-proxy/src/internal/balancer"
	"build-your-own-reverse-proxy/src/internal/config"
	"build-your-own-reverse-proxy/src/internal/forwarder"
	"build-your-own-reverse-proxy/src/internal/proxy"
	"build-your-own-reverse-proxy/src/internal/router"
	"log/slog"
	"net/http"
)

type Runtime struct {
	Proxy *proxy.Proxy
	logger *slog.Logger
}

func Build(cfg *config.Config, logger *slog.Logger) (*Runtime, error) {
	router := router.NewRouter(cfg.Routes)
	forwarder := forwarder.NewForwarder(&http.Transport{
		MaxIdleConns:        cfg.Transport.MaxIdleConns,
		IdleConnTimeout:     cfg.Transport.IdleConnTimeout,
		TLSHandshakeTimeout: cfg.Transport.TLSHandshakeTimeout,
	}, logger)

	balancer := balancer.NewRoundRobinBalancer()

	proxy := proxy.NewProxy(router, forwarder, balancer, logger)

	return &Runtime{
		Proxy: proxy,
		logger: logger,
	}, nil
}




