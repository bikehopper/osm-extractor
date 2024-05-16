FROM debian:stable-slim
RUN apt-get update && \
    apt-get install -y osmium-tool dumb-init && \
    mkdir /app

VOLUME ["/mnt/input", "/mnt/output"]
WORKDIR /app
COPY polygons .
COPY config.json .
ENTRYPOINT ["/usr/bin/dumb-init", "--"]
CMD osmium extract -d /mnt/output -c ./config.json /mnt/input/norcal-latest.osm.pbf