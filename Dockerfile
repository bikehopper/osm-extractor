FROM golang:1.22 as build
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY workflow ./workflow
RUN go build -v -o /usr/local/bin/create /usr/src/app/workflow/create/
RUN go build -v -o /usr/local/bin/worker /usr/src/app/workflow/worker/

FROM debian:12-slim
RUN apt-get update && \
    apt-get install -y osmium-tool dumb-init && \
    mkdir /app
VOLUME ["/mnt/input", "/mnt/output"]
WORKDIR /app
COPY --from=build /usr/local/bin/create /usr/local/bin/worker .
COPY polygons ./polygons
COPY config.json .
ENTRYPOINT ["/usr/bin/dumb-init", "--"]