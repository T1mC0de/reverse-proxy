package forwarder

import (
	"context"
	"io"
	"net/http"
	"log/slog"
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

func (f *Forwarder) Forward(ctx context.Context, w http.ResponseWriter, r *http.Request, upstreamURL string) error {
	targetPath := upstreamURL + r.URL.Path
	requestBody := r.Body

	req, err := http.NewRequestWithContext(ctx, r.Method, targetPath, requestBody)
	if err != nil {
		f.logger.Error("Failed to create request", "error", err)
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return err
	}
	
	req.Header = r.Header.Clone()
	req.Header.Set("X-Forwarded-For", r.RemoteAddr)
	req.Header.Set("X-Forwarded-Host", r.Host)

	resp, err := f.transport.RoundTrip(req)
	if err != nil {
		f.logger.Error("Failed to forward request", "error", err)
		f.logger.Debug("Forwarding request", "headers", req.Header, "URL", req.URL.String())
		return err
	}
	defer resp.Body.Close()

	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	w.WriteHeader(resp.StatusCode)
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		return err
	}

	return nil
}