FROM golang:1.25-alpine

WORKDIR /app

RUN apk add --no-cache git inotify-tools

EXPOSE 8080

COPY src /app/src
RUN cd ./src && go mod download && go build -o ../server main.go

CMD ./server
