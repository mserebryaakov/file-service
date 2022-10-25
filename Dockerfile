FROM golang:1.19-alpine

WORKDIR /app
COPY . .

RUN go build -o /main cmd/main.go

ENTRYPOINT ["/main"]