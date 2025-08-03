FROM golang:1.24.1-alpine

WORKDIR /app

COPY . /app

WORKDIR /app/database-init
RUN go run main.go

WORKDIR /app
RUN go build -o bot .

CMD ["./bot"]