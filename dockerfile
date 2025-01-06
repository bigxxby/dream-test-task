# Этап сборки приложения
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Устанавливаем зависимости для компиляции в Alpine
RUN apk add --no-cache git

# Копируем все исходные файлы в контейнер
COPY . .

# Загружаем зависимости и собираем приложение
RUN go mod tidy
RUN go build -o app ./cmd

# Этап с конечным образом
FROM alpine:latest

# Устанавливаем сертификаты
RUN apk update && apk add --no-cache ca-certificates libc6-compat

# Копируем собранный бинарный файл из предыдущего этапа
COPY --from=builder /app/app /usr/local/bin/app

# Копируем .env файл в контейнер
COPY .env /app/.env

# Устанавливаем рабочую директорию
WORKDIR /app

# Открываем порт для приложения
EXPOSE 8081

# Устанавливаем права на бинарник
RUN chmod +x /usr/local/bin/app

# Команда для запуска приложения
CMD ["/usr/local/bin/app"]
