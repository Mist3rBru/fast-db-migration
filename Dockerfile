FROM golang:1.22 as base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY /cmd ./cmd
COPY /config ./config
COPY /internal ./internal

FROM base as build
RUN GOOS=linux go build -a -o main cmd/app/main.go