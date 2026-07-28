package config

import (
	"os"
	"time"
	"gopkg.in/yaml.v3"
	"fmt"
)

type Server struct {
	Port string `yaml:"port"`
}

type Transport struct {
	MaxIdleConns        int    		  `yaml:"max_idle_conns"`
	IdleConnTimeout     time.Duration `yaml:"idle_conn_timeout"`
	TLSHandshakeTimeout time.Duration `yaml:"tls_handshake_timeout"`
}

type Upstream struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Route struct {
	Name      string   		`yaml:"name"`
	Path      string   		`yaml:"path"`
	Timeout   time.Duration `yaml:"timeout"`
	Upstreams []string 		`yaml:"upstreams"`
}

type Config struct {
	Server     Server      `yaml:"server"`
	Transport  Transport   `yaml:"transport"`
	Upstreams  []Upstream  `yaml:"upstreams"`
	Routes     []Route     `yaml:"routes"`
}

func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var c Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&c); err != nil {
		return nil, err
	}

	upstreamMap := make(map[string]string)
	for i := range c.Upstreams {
		upstreamMap[c.Upstreams[i].Name] = c.Upstreams[i].URL
	}

	for i := range c.Routes {
		for j := range c.Routes[i].Upstreams {
			var exists bool
			c.Routes[i].Upstreams[j], exists = upstreamMap[c.Routes[i].Upstreams[j]]
			if !exists {
				return nil, fmt.Errorf("upstream '%s' referenced in route '%s' not found", c.Routes[i].Upstreams[j], c.Routes[i].Name)
			}
		}
	}

	return &c, nil
}