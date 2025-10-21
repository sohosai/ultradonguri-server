FROM golang:1.25-alpine

WORKDIR /app

RUN apk add --no-cache git inotify-tools
RUN go install github.com/air-verse/air@latest

EXPOSE 8080
