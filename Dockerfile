FROM golang:1.22 as build
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY workflow ./workflow
RUN go build -v -o /usr/local/bin/osm-extractor-workflow /usr/src/app/workflow/cmd

FROM debian:12-slim
RUN apt-get update && \
    apt-get install -y osmium-tool dumb-init && \
    mkdir /app
VOLUME ["/mnt/input", "/mnt/output"]
WORKDIR /app
COPY --from=build /usr/local/bin/osm-extractor-workflow .
COPY polygons ./polygons
COPY config.json .
ENTRYPOINT ["/usr/bin/dumb-init", "--"]