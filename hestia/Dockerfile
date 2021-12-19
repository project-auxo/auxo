# Build: /auxo > `docker build -f hestia/Dockerfile . -t auxo/hestia:v1.0`
# Run: /auxo > `docker run -it --rm -p 8080:8080 --name hestia-instance auxo/hestia:v1.0`
FROM golang:1.16.12-alpine3.15 AS build
ENV APP_HOME /go/src/github.com/project-auxo/hestia/src
WORKDIR ${APP_HOME}
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN apk add --no-cache ca-certificates
RUN --mount=type=cache,target=/root/.cache/go-build \
CGO_ENABLED=0 GOOS=linux go build -o /out/hestia-app hestia/main.go

FROM scratch
COPY hestia/internal/config/config.yml config/
COPY hestia/.env .
COPY hestia/web hestia/web/
COPY --from=build /out/hestia-app .
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
EXPOSE 8080
CMD ["./hestia-app", "--config", "config/config.yml", "-env", ".env"]
