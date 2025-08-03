FROM golang:1.24.1-alpine

WORKDIR /app

COPY . .
RUN cd database-init && go run main.go

WORKDIR /app
RUN go build -o bot .

CMD ["./bot"]