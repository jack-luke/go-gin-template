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

## Build
All container builds use [Chainguard's Static Base Image](https://images.chainguard.dev/directory/image/static/overview) for security.

```
go build . -o gin
```

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

## Features

* The Gin server is started in release mode.
* The Gin server is set to trust no proxies.
* Security headers are applied to all requests.
* Error responses are returned to the client as JSON.

### Endpoints
| Path | Description |
| --- | --- |
| `/healthz` | Kubernetes liveliness probe; simply returns a 200 OK response |
| `/readyz` | Kubernetes readiness probe; by default, returns a 200 OK response, and should be configured to include app readiness checks. |
| `/metrics` | Prometheus metrics endpoint. Returns application metrics in Prometheus format. |

### Logging
All logging is done to STDOUT in a structured manner using the [Go slog library](https://pkg.go.dev/log/slog).

HTTP requests are logged with the following fields:
| Field | Type | Description | Example
| --- | --- | --- | --- |
| `time` | string | Timestamp in RFC 3339 format.| 2025-11-29T18:01:38.300Z |
| `level` | int | Log severity level. | WARN |
| `msg` | string | HTTP response code and its meaning. | HTTP 503 (Service Unavailable) |
| `status` | int | HTTP response code. | 200 |
| `method` | string | HTTP request method. | GET |
| `path` | string | Request route. | /healthz |
| `client_ip` | string | Client IP address. | 127.0.0.1 |
| `duration` | string | Time taken for the request to complete. | 107.747µs | 
| `error` | string | The error message, if there was one. | "Error #01: MQTT broker not connected\n" |

### Metrics
Prometheus metrics are recorded using the default registry.

| Name | Type | Labels | Description |
| --- | --- | --- | --- |
| `http_requests_total` | Counter | method, route, status | Total number of HTTP requests. |
| `http_request_duration_seconds` | Histogram | method, route, status | HTTP request duration in seconds. |
| `http_in_flight_requests` | Gauge | | Number of requests currently being handled by the service. |

#### Using a Non-default Registry
```go
// create new registry
reg := prometheus.NewRegistry()

// pass the registerer to the middleware
r.Use(middleware.PrometheusMetrics(reg))

// attach the registry to the metrics endpoint
r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(reg, promhttp.HandlerOpts{})))
```

## Authors

* Jack Luke