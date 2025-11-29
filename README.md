# Go Gin App

A starter setup for a production Gin application, including structured logging, Prometheus HTTP metrics, Kubernetes health probes.

## Getting Started
First, clone the project.
New HTTP route handlers and middleware can be added to the `controllers/` and `middleware/` directories respectively.
See [Gin Docs - Examples](https://gin-gonic.com/en/docs/examples/) for a guide on how to implement these and some of the features Gin offers. 

`main.go` is the main entrypoint, so new functionality is attached to the server here.

### Project Structure
```bash
.
├── .ko.yaml        # Ko build config
├── controllers/    # HTTP route handlers
│   └── health.go 
├── Dockerfile      # Docker build config
├── go.mod
├── go.sum
├── main.go         # Project entrypoint; sets up and runs Gin server
├── middleware/     # Gin middleware for logging, metrics, etc.
│   ├── errors.go
│   ├── logging.go
│   ├── metrics.go
│   └── security.go
└── README.md
```

### Default Behaviour

* The Gin server is started in release mode.
* The Gin server is set to trust no proxies.
* Kubernetes liveliness and readiness probes are mounted on `/healthz` and `/readyz` respectively. Checks should be added to the readiness probe as the app grows.
* Security headers are applied to all requests.
* Go structured logging used throughout, with key-value log format by default.
* Basic Prometheus HTTP metrics on request count, latency and in-flight are collected, and made available on `/metrics`.

## Build
All builds use [Chainguard's Static Base Image](https://images.chainguard.dev/directory/image/static/overview) for security.

### Ko
To build with [Ko](https://ko.build), pushing to the local Docker image store tagged as `gin:latest`.
```bash
export KO_DOCKER_REPO=gin
ko build --bare --local .
```
Configure Ko build in the `.ko.yaml` file.

### Docker
```bash
docker build . -t gin:latest
```

## Authors

* Jack Luke