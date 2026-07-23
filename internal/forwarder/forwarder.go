package forwarder

import (
	"context"
	"io"
	"net/http"
	"log/slog"
	"net/url"
)

type Forwarder struct {
	transport *http.Transport
	logger    *slog.Logger
}

func NewForwarder(transport *http.Transport, logger *slog.Logger) *Forwarder {
	return &Forwarder{
		transport: transport,
		logger:    logger,
	}
}

func (f *Forwarder) Forward(ctx context.Context, w http.ResponseWriter, r *http.Request, upstreamURL string) {
	targetPath, err := url.JoinPath(upstreamURL, r.URL.Path)
	if err != nil {
		f.logger.Error("Failed to join URL path", "error", err)
		http.Error(w, "Failed to join URL path", http.StatusInternalServerError)
		return
	}
	requestBody := r.Body

	req, err := http.NewRequestWithContext(ctx, r.Method, targetPath, requestBody)
	if err != nil {
		f.logger.Error("Failed to create request", "error", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}
	
	req.Header = r.Header.Clone()
	req.Header.Set("X-Forwarded-For", r.RemoteAddr)
	req.Header.Set("X-Forwarded-Host", r.Host)

	resp, err := f.transport.RoundTrip(req)
	if err != nil {
		f.logger.Error("Failed to forward request", "error", err)
		f.logger.Debug("Forwarding request", "headers", req.Header, "URL", req.URL.String())
		http.Error(w, "Failed to forward request", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	f.logger.Debug("Received response", "status", resp.StatusCode, "headers", resp.Header, "URL", req.URL.String())

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		http.Error(w, "Failed to copy response body", http.StatusBadGateway)
		return
	}
}