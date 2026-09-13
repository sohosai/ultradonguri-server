FROM golang:1.25-alpine

WORKDIR /app

RUN apk add --no-cache git inotify-tools
RUN go install github.com/air-verse/air@v1.64.5

EXPOSE 8080
