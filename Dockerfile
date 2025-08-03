FROM golang:1.24.1-alpine

WORKDIR /app

COPY . .

RUN cd database-init && go run main.go

RUN go build -o bot .

CMD ["./bot"]