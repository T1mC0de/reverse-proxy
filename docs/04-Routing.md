# Routing: Match Rules

This is where the proxy stops being just an HTTP client and becomes a router. The router holds the list of routes exactly as they came out of the config, and matching follows that order — the first route whose `Path` prefix matches wins:

```go
// internal/router/router.go
package router

import (
 "build-your-own-reverse-proxy/src/internal/config"
 "strings"
)

type Router struct {
 routes []config.Route
}

func NewRouter(routes []config.Route) *Router {
 return &Router{
  routes: routes,
 }
}

func (r *Router) Match(path string) *config.Route {
 for _, route := range r.routes {
  if strings.HasPrefix(path, route.Path) {
   return &route
  }
 }
 return nil
}
```

Now, we can use the `Router` in our `ReverseProxy` to match incoming requests to the appropriate route and forward them to the correct upstream.
