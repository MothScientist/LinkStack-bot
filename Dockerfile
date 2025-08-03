FROM golang:1.24-alpine

WORKDIR /app

COPY . .

WORKDIR /app/database-init
RUN go run main.go

WORKDIR /app
RUN go build -o bot .

CMD ["./bot"]