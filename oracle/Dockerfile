# Build: /auxo > `docker build -f oracle/Dockerfile . -t auxo/oracle:v1.0`
# Run: /auxo > `docker run -it --rm -p 3000:3000 --name oracle-instance auxo/oracle:v1.0`
FROM golang:1.16.12-alpine3.15 AS build
ENV APP_HOME /go/src/github.com/project-auxo/oracle/src
WORKDIR ${APP_HOME}
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build \
CGO_ENABLED=0 GOOS=linux go build -o /out/oracle-app oracle/main.go

FROM scratch
COPY oracle/internal/config/config.yml config/
COPY oracle/services/ oracle/services/
COPY --from=build /out/oracle-app .
EXPOSE 3000
CMD ["./oracle-app", "--config", "config/config.yml"]
