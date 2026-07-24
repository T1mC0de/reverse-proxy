# Reverse Proxy in Go

A reverse proxy I built to understand how things like Nginx work under the hood. It reads a YAML config, routes requests based on path prefixes, and forwards them to upstream servers.

## What it does

- Takes incoming HTTP requests and figures out where they should go based on the URL path
- If a route has multiple upstream servers, it balances requests between them
- Adds `X-Forwarded-For` and `X-Forwarded-Host` headers so backends know who the real client is
- Each route can have its own timeout
- Shuts down gracefully when you hit Ctrl+C — finishes processing current requests before exiting
- Logs everything in a readable format using `slog`

## How it's organized

```
cmd/
  proxy/      - starts the proxy
  servers/    - a dummy backend for testing (returns "Hello /path")
configs/      - YAML files that define routes and upstreams
internal/
  config/     - reads and validates the YAML config
  router/     - matches incoming request paths to routes
  forwarder/  - sends the request to an upstream and pipes the response back
  balancer/   - picks which upstream to use when there are several
  proxy/      - ties router + balancer + forwarder together
  runtime/    - builds everything at startup (dependency injection, basically)
  server/     - thin wrapper around http.Server with graceful shutdown
```

## Configuration

Here's what a config looks like:

```yaml
server:
  port: ":8080"

transport:
  max_idle_conns: 100
  idle_conn_timeout: 90s
  tls_handshake_timeout: 10s

upstreams:
  - name: backend
    url: "https://jsonplaceholder.typicode.com"
  - name: frontend
    url: "http://localhost:9003"

routes:
  - name: api
    path: "/users"
    timeout: 3s
    upstreams:
      - "backend"
  - name: static
    path: "/"
    timeout: 5s
    upstreams:
      - "frontend"
```

## Running it locally

First, start a test backend so you have something to proxy to:

```bash
go run ./cmd/servers/main.go
# Starts on :9003, responds with "Hello /your-path"
```

Then start the proxy:

```bash
go run ./cmd/proxy/main.go
# Starts on :8080
```

Now you can test it:

```bash
# This goes to jsonplaceholder
curl http://localhost:8080/users

# This goes to the local test server
curl http://localhost:8080/
```

## Things I'd like to add

The project works, but there's a bunch of stuff I want to improve:

- **Health checks for upstreams** — right now the balancer has no idea if a backend is down. Adding health checks and a circuit breaker would make it actually production-ready.
- **Hot config reload** — having to restart the proxy every time you change a route is annoying. Should be able to watch the config file and reload gracefully.
- **Tests** — I know, I know. The interfaces are designed to be mockable, I just need to actually write the tests.

## Why I built this

I've used Nginx without really understanding what happens inside. Building one from scratch — even a simple one — teaches you a lot about HTTP, connection pooling, and why reverse proxies are configured the way they are.