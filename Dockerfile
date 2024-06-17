FROM golang:1.22 as build
WORKDIR /usr/src/app
COPY go.mod go.sum Makefile ./
RUN make install

COPY workflow ./workflow
RUN make build

FROM debian:12-slim
RUN apt-get update && \
    apt-get install -y ca-certificates osmium-tool dumb-init && \
    mkdir /app
VOLUME ["/mnt/input", "/mnt/output"]
WORKDIR /app
COPY --from=build /usr/src/app/bin/osm-extractor-workflow .
COPY polygons ./polygons
COPY config.json .
ENTRYPOINT ["/usr/bin/dumb-init", "--"]