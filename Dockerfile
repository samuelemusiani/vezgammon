FROM node:alpine as build-client

WORKDIR /app

RUN apk add just
COPY justfile justfile

COPY client client

RUN just build-client

FROM golang:1.23.2-alpine as build-server

WORKDIR /app    

RUN apk add just
COPY justfile justfile

COPY server server
COPY --from=build-client /app/client/dist server/dist
COPY go.mod go.mod
COPY go.sum go.sum

RUN just build-server

FROM alpine

WORKDIR /app

COPY --from=build-server /app/build .
COPY server/config/release-config.toml .

CMD [ "sh", "-c", "GIN_MODE=release ./vezgammon release-config.toml" ]
