FROM golang:1.22.2 AS builder

WORKDIR /app

RUN mkdir /var/log/entry/
COPY ../../ .

RUN --mount=type=cache,target=/go/pkg/mod/cache \
    go mod download
RUN --mount=type=cache,target=/go/pkg/mod/cache \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go buil -o web_app.linux -mod=mod cmd/web_app/main.go
CMD ["/app/web_app.linux"]

FROM builder AS test_env
RUN cp /app/configs/files/test.yml /app/config.yml
CMD ["/app/web_app.linux"]

FROM builder AS live_env
RUN cp /app/configs/files/live.yml /app/config.yml
CMD ["/app/newsfeed_service.linux"]