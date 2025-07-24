FROM golang:1.23.0-alpine

WORKDIR /var/www/app

COPY . .

EXPOSE 8080

# RUN go mod download
# RUN go install github.com/pressly/goose/v3/cmd/goose@latest

RUN chmod +x entrypoint.sh

CMD ["./entrypoint.sh"]
