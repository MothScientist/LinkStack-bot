# Базовый образ
FROM golang:1.24-alpine

# Рабочая директория
WORKDIR /app

# Копируем ВСЕ файлы проекта с учётом .dockerignore
COPY . .

WORKDIR /app/database-init
RUN go run main.go

WORKDIR /app
RUN go build -o bot .

# 3. Запускаем приложение
CMD ["./bot"]