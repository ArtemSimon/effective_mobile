# ========== Build Stage ==========
FROM golang:1.23-alpine AS builder

# Устанавливаем зависимости для сборки
RUN apk add --no-cache git make gcc musl-dev

# Создаем рабочую директорию
WORKDIR /app

# Копируем файлы модулей (для лучшего кэширования)
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение с отключенным CGO и флагами оптимизации
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-w -s" \
    -o /app/bin/effective-mobile \
    ./app/main.go 
# Собираем тестовый бинарник (дополнительный шаг)
RUN go test -c -o /app/bin/effective-mobile-test ./internal/api

# ========== Runtime Stage ==========
FROM alpine:3.18

# Устанавливаем зависимости для runtime
RUN apk add --no-cache ca-certificates tzdata

# Копируем бинарник из builder stage
COPY --from=builder --chown=appuser:appgroup /app/bin/effective-mobile /app/

# Копируем тестовый бинарник 
COPY --from=builder /app/bin/effective-mobile-test /app/

# Копируем файлы миграций 
COPY --from=builder /app/migrations /app/migrations

# В продакш так не делают это не безопасно можно использовать секреты или тот же Vault
# COPY .env /app/.env 
# Настраиваем рабочую директорию
WORKDIR /app

# Экспонируем порт (должен совпадать с конфигом приложения)
EXPOSE 8080

# Запускаем приложение
CMD ["./effective-mobile"]