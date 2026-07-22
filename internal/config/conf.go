package config

import (
	"os"
	"time"
	"gopkg.in/yaml.v3"
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
	Upstreams []Upstream 	`yaml:"upstreams"`
}

type Config struct {
	Server     Server      `yaml:"server"`
	Transport  Transport   `yaml:"transport"`
	Upstreams  []Upstream  `yaml:"upstreams"`
	Routes     []Route     `yaml:"routes"`
}


// loads without validation
func (c *Config) Load(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&c); err != nil {
		return err
	}

	return nil
}