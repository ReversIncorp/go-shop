FROM golang:1.23.0-alpine

WORKDIR /var/www/app
RUN apk --no-cache add postgresql-client

COPY go.mod go.sum ./
RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .
RUN go build -o main app/main.go
RUN chmod +x main

EXPOSE 8080

RUN chmod +x entrypoint.sh

CMD ["./entrypoint.sh"]
