FROM golang:1.19-alpine AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /app/main cmd/main.go

FROM alpine:latest

COPY config/config.json /config/config.json

COPY docs ./docs

COPY --from=builder /app/main /main

ENTRYPOINT ["/main"]