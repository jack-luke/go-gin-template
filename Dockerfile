FROM --platform=$BUILDPLATFORM golang:alpine AS build

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY main.go .
COPY ./middleware ./middleware/
COPY ./controllers ./controllers/

RUN go build -ldflags="-s -w -extldflags '-static'" -o /bin/gin

FROM cgr.dev/chainguard/static:latest

WORKDIR /gin

COPY --from=build /bin/gin ./gin

ENTRYPOINT [ "./gin" ]