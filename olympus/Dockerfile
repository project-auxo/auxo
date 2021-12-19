# Build: /auxo > `docker build -f olympus/Dockerfile . -t auxo/olympus:v1.0`
# Run: /auxo > `docker run -it --rm -p 5555:5555 -p 5556:5556 -p 3000:3000 --name olympus-instance auxo/olympus:v1.0`
FROM golang:1.16.12-alpine3.15 AS build
ENV APP_HOME /go/src/github.com/project-auxo/olympus/src
WORKDIR ${APP_HOME}
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY . .

# Install extra dependencies
RUN apk add --no-cache zeromq-dev musl-dev pkgconfig alpine-sdk libsodium-dev libzmq-static libsodium-static

RUN --mount=type=cache,target=/root/.cache/go-build \
  CGO_LDFLAGS="$CGO_LDFLAGS -lstdc++ -lm -lsodium" \
  CGO_ENABLED=1 \
  GOOS=linux \
  go build -v -a --ldflags '-extldflags "-static" -v' -o /out/olympus-app olympus/main.go

FROM scratch
COPY olympus/internal/config/config.yml config/
COPY --from=build /out/olympus-app .
EXPOSE 5555 5556 3000
CMD ["./olympus-app", "--config", "config/config.yml"]
