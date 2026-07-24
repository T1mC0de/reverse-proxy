# Build Your Own Reverse Proxy

A lightweight reverse proxy built in Go as a learning project. The proxy reads a YAML config, routes requests by path prefix, proxies them to upstream services, and shuts down gracefully on `SIGINT`.

## Project Structure

- `cmd/proxy` — reverse proxy entry point
- `cmd/servers` — simple example backend server
- `configs` — YAML config files
- `internal/config` — configuration loading and validation
- `internal/router` — route lookup by request path
- `internal/forwarder` — sending requests an upstream and returning responses to the client
- `internal/proxy` — router + forwarder glue
- `internal/runtime` — application runtime layer assembly
- `internal/server` — `http.Server` wrapper
- `internal/balancer` — balance request between several upstreams
  

## How It Works

1. The application reads the config and builds the runtime.
2. The router looks up a route by path.
3. The balancer selects upstream from the route.
4. The forwarder builds a new HTTP request and sends it to the upstream.
5. The upstream response is returned to the client without additional business logic.

## Configuration

The project includes two example configs:

- `configs/config.yaml` — main example with an external backend and a local frontend
- `configs/config_advanced.yaml` — example with multiple upstreams in one route


-------
------
-----
----
---

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