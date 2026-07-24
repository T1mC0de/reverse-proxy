package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoad(t *testing.T) {
	createTempFile := func(t *testing.T, content string) string {
		t.Helper()
		tmpDir := t.TempDir()
		tmpFile := filepath.Join(tmpDir, "config.yaml")
		err := os.WriteFile(tmpFile, []byte(content), 0644)
		if err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}
		return tmpFile
	}

	tests := []struct {
		name        string
		yamlContent string
		wantErr     bool
		validate    func(t *testing.T, c *Config)
	}{
		{
			name: "Valid config",
			yamlContent: `
server:
  port: ":8080"
transport:
  max_idle_conns: 100
  idle_conn_timeout: 90s
  tls_handshake_timeout: 10s
upstreams:
  - name: backend
    url: "http://localhost:9002"
routes:
  - name: api
    path: "/api"
    timeout: 3s
    upstreams:
      - "backend"
`,
			wantErr: false,
			validate: func(t *testing.T, c *Config) {
				if c.Server.Port != ":8080" {
					t.Errorf("expected port :8080, got %s", c.Server.Port)
				}
				if c.Transport.MaxIdleConns != 100 {
					t.Errorf("expected MaxIdleConns 100, got %d", c.Transport.MaxIdleConns)
				}
				if c.Transport.IdleConnTimeout != 90*time.Second {
					t.Errorf("expected IdleConnTimeout 90s, got %v", c.Transport.IdleConnTimeout)
				}
				if len(c.Routes[0].Upstreams) == 0 || c.Routes[0].Upstreams[0] != "http://localhost:9002" {
					t.Errorf("expected upstream URL 'http://localhost:9002', got %v", c.Routes[0].Upstreams)
				}
			},
		},
		{
			name:        "File not found",
			yamlContent: "",
			wantErr:     true,
		},
		{
			name: "Invalid YAML syntax",
			yamlContent: `
server:
  port: ":8080"
  invalid_yaml_structure: [
`,
			wantErr: true,
		},
		{
			name: "Missing upstream reference",
			yamlContent: `
upstreams:
  - name: backend
    url: "http://localhost:9002"
routes:
  - name: api
    path: "/api"
    upstreams:
      - "non_existent_backend"
`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var filePath string
			if tt.name != "File not found" {
				filePath = createTempFile(t, tt.yamlContent)
			} else {
				filePath = "/non/existent/path.yaml"
			}

			got, err := Load(filePath)

			if (err != nil) != tt.wantErr {
				t.Fatalf("Load() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && tt.validate != nil {
				tt.validate(t, got)
			}
		})
	}
}