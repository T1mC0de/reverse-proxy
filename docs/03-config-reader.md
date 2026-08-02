# Config reader

## Why do we need a config?

If you have used, for example, Nginx before, you may have noticed that it uses configuration files `nginx.conf` to define how the reverse proxy should behave. For example, you can define routes, upstreams, and other settings in a configuration file. This allows you to change the behavior of the reverse proxy without having to recompile the code.

## Config file in our project

So, in our case, to keep things simple, we will use YAML files to define the configuration of our reverse proxy. This way, we can easily change the behavior of the reverse proxy without having to modify the code. Let's take a look at the configuration file that we will use in our project.

```yaml
# config.yaml
server:
  port: ":8080"

transport:
  max_idle_conns: 100
  idle_conn_timeout: 90s
  tls_handshake_timeout: 10s

upstreams:
  - name: backend1
    url: "https://dummyjson.com"
  - name: backend2
    url: "https://jsonplaceholder.typicode.com"
  - name: frontend
    url: "http://localhost:9003"

routes:
  - name: api
    path: "/posts"
    timeout: 3s
    upstreams:
      - "backend1"
      - "backend2"

  - name: frontend
    path: "/"
    timeout: 3s
    upstreams:
      - "frontend"
```

Now, let's break down the configuration file and explain what each section does.

#### Server section

First of all, we define the `server` section, which contains the port on which the reverse proxy will listen for incoming requests. In this case, we set it to `:8080`, which means that the reverse proxy will listen on port 8080 on all available network interfaces.

#### Transport section

Secondly, we define the `transport` section, which contains settings for the HTTP transport. In this case, we set the maximum number of idle connections to 100, the idle connection timeout to 90 seconds, and the TLS handshake timeout to 10 seconds.

#### Upstreams section

The `upstreams` section defines the upstream servers that the reverse proxy will forward requests to. In this case, we define three upstreams: `backend1`, `backend2`, and `frontend`. Each upstream has a name and a URL.

#### Routes section

Finally, we define the `routes` section, which contains the routes that the reverse proxy will handle. In this case, we define two routes: `api` and `frontend`. Each route has a name, a path, a timeout, and a list of upstreams. The reverse proxy will forward requests to the upstreams based on the path of the request.

## Parsing the config file

To parse the configuration file, we will use the `gopkg.in/yaml.v3` package, which provides a simple way to parse YAML files in Go (if you are using another language - most likely your language has a similar package).
We will define a structures that matches the structures of the configuration file, and then we will use the `yaml.Unmarshal` function to parse the configuration file into our structures.

```go
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
 MaxIdleConns        int        `yaml:"max_idle_conns"`
 IdleConnTimeout     time.Duration `yaml:"idle_conn_timeout"`
 TLSHandshakeTimeout time.Duration `yaml:"tls_handshake_timeout"`
}

type Upstream struct {
 Name string `yaml:"name"`
 URL  string `yaml:"url"`
}

type Route struct {
 Name      string     `yaml:"name"`
 Path      string     `yaml:"path"`
 Timeout   time.Duration `yaml:"timeout"`
 Upstreams []string   `yaml:"upstreams"`
}

type Config struct {
 Server     Server      `yaml:"server"`
 Transport  Transport   `yaml:"transport"`
 Upstreams  []Upstream  `yaml:"upstreams"`
 Routes     []Route     `yaml:"routes"`
}
```

Now, we need to implement a function that will read the configuration file and parse it into our `Config` structure.
Here is my implementation of the `Load` function that does just that:

```go
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
```

Notice, in the function I also used map to store URLs of upstreams in `Routes`, so that we can easily replace the names of upstreams in `Routes` with their corresponding URLs. This way, when we handle a request, we can easily forward it to the correct upstream server.
