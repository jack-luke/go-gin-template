# Go Gin App

A starter setup for a production Gin application, including structured logging, Prometheus HTTP metrics, Kubernetes health probes.

First, clone the repository and 'cd' into it. 

## Build
All builds use [Chainguard's Static Base Image](https://images.chainguard.dev/directory/image/static/overview) for security.

### Ko
To build with Ko, pushing to the local Docker image store tagged as `gin:latest`.
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