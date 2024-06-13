FROM node:current-slim as build
RUN mkdir /app
WORKDIR /app
COPY package.json package-lock.json tsconfig.json /app
RUN npm install
COPY ./src ./src
RUN npm run build

FROM node:current-slim
RUN apt-get update && \
    apt-get install -y osmium-tool dumb-init && \
    mkdir /app
VOLUME ["/mnt/input", "/mnt/output"]
WORKDIR /app
COPY package.json /app
RUN npm install --production

COPY --from=build /app/lib ./lib
COPY ./polygons ./polygons
ENTRYPOINT ["/usr/bin/dumb-init", "--"]