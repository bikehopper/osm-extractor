FROM node:current-slim
RUN apt-get update && \
    apt-get install -y osmium-tool dumb-init && \
    mkdir /app

VOLUME ["/mnt/input", "/mnt/output"]
WORKDIR /app
COPY package.json /app
COPY package-lock.json /app
RUN npm install

COPY tsconfig.json ./tsconfig.json
COPY ./src ./src
RUN npm run build
COPY ./polygons ./polygons

ENTRYPOINT ["/usr/bin/dumb-init", "--"]