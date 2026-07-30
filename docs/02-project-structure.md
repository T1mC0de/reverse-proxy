# Project Structure

```text
reverse-proxy/
├── cmd/
│   ├── proxy/
│   │   └── main.go
│   └── backend/
│       └── main.go          # a toy backend used for the test bench
├── configs/
│   ├── config.yaml
│   └── config_advanced.yaml
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── balancer/
│   │   └── balancer.go
│   ├── forwarder/
│   │   └── forwarder.go
│   ├── router/
│   │   └── router.go
│   ├── runtime/
│   │   └── runtime.go
│   ├── server/
│   │   └── server.go
│   └── proxy/
│       └── proxy.go
├── go.mod
└── go.sum
```

## Project Structure Explained

Let's go through the project structure and explain what each folder/file is for.

- `cmd/` - contains the entry points for the application. Each subfolder represents a separate command. In our case, we have two commands: `proxy` and `backend`. The `proxy` command is the main reverse proxy application, while the `backend` command is a toy backend used for testing purposes.
- `configs/` - contains the configuration files for the application. The `config.yaml` file is the main configuration file, while `config_advanced.yaml` is an advanced configuration file that demonstrates more complex routing and upstream configurations.
- `internal/` - contains the internal packages of the application. These packages are not meant

More information about internal packages:

- `config/` - contains the code for reading and validating the configuration files.
- `balancer/` - contains the code for load balancing between multiple upstreams.
- `forwarder/` - contains the code for forwarding requests to the appropriate upstream and returning the response back to the client.
- `router/` - contains the code for routing incoming requests to the appropriate upstream based on the request path and configuration.
- `runtime/` - contains the code for initializing the application at startup, including dependency injection and setting up the necessary components.
- `server/` - contains the code for running the HTTP server and handling graceful shutdowns.
- `proxy/` - contains the code that ties together the router, balancer, and forwarder to create the complete reverse proxy functionality.

If you find yourself struggling to understand the project structure, don't worry! The following chapters will guide you through the code step by step, explaining how each component works and how they fit together to create a fully functional reverse proxy. So far, some words like "upstream" and "route" may be unfamiliar, but they will become clear as we progress through the guide.

Also, if you want to use your own implementation of the reverse proxy with your own structure, you can find a chapter about testing your own implementation at the end of the guide.
