# syntax=docker/dockerfile:1

FROM golang:1.16-alpine

WORKDIR /simple_json_api

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . /simple_json_api/

RUN go build

EXPOSE 8080

# default 8080
ARG PORTNO_ARG="8080"

CMD [ "./simple_json_api --PORT=${PORTNO_ARG} --MODE=release" ]