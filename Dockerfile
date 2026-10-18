ARG GO_VERSION=1.25.14
ARG ALPINE_VERSION=3.23.4
ARG DELVE_VERSION=v1.27.2
ARG TARGETOS=linux
ARG TARGETARCH=amd64

FROM golang:${GO_VERSION}-alpine AS build-base
WORKDIR /app

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download
COPY . .

FROM build-base AS build-release
ARG TARGETOS
ARG TARGETARCH
ARG COMMIT_SHA=dirty
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -ldflags="-s -w -X main.commitSHA=${COMMIT_SHA}" -o main ./cmd/api

FROM build-base AS build-debug
ARG TARGETOS
ARG TARGETARCH
ARG COMMIT_SHA=dirty
ARG DELVE_VERSION
RUN GOBIN=/app go install github.com/go-delve/delve/cmd/dlv@${DELVE_VERSION}
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -gcflags="all=-N -l" -tags=debug \
    -ldflags="-X main.commitSHA=${COMMIT_SHA}" -o main ./cmd/api

FROM alpine:${ALPINE_VERSION} AS runtime-base
WORKDIR /app

RUN addgroup -S app \
    && adduser -S app -G app \
    && mkdir -p /app/etc \
    && chown -R app:app /app
COPY --chown=app:app ./cmd/api/etc/config.yaml ./etc/config.yaml

USER app
EXPOSE 8080

FROM runtime-base AS release
COPY --from=build-release --chown=app:app /app/main ./main
CMD ["./main"]

FROM runtime-base AS debug
COPY --from=build-debug --chown=app:app /app/main ./main
COPY --from=build-debug --chown=app:app /app/dlv ./dlv
EXPOSE 40000
CMD ["./dlv", "--listen=:40000", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./main"]

FROM release AS final

