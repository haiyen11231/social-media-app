FROM golang:1.24.0-alpine3.21 AS builder

WORKDIR /app

RUN mkdir /var/log/entry/
COPY ../../ .

RUN --mount=type=cache,target=/go/pkg/mod/cache \
    go mod download
RUN --mount=type=cache,target=/go/pkg/mod/cache \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o web_app.linux -mod=mod cmd/web_app/main.go
CMD ["/app/web_app.linux"]

FROM builder AS test_env
RUN cp /app/configs/files/test.yml /app/config.yml
EXPOSE 8000
CMD ["/app/web_app.linux"]

FROM builder AS live_env
RUN cp /app/configs/files/live.yml /app/config.yml
CMD ["/app/newsfeed_service.linux"]