# Stage 1: Builder - собирает бинарник
FROM golang:1.24.4-alpine AS builder

WORKDIR /app

# Копируем зависимости для кэширования
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .

# ARG для имени сервиса (передается при сборке)
ARG SERVICE_NAME
ENV SERVICE_NAME=${SERVICE_NAME}

# Собираем только указанный сервис
RUN if [ -n "$SERVICE_NAME" ] && [ -f "cmd/${SERVICE_NAME}/main.go" ]; then \
        echo "Building ${SERVICE_NAME}..." && \
        CGO_ENABLED=0 GOOS=linux go build \
            -ldflags="-s -w" \
            -o /app/bin/main \
            ./cmd/${SERVICE_NAME}/; \
    else \
        echo "ERROR: SERVICE_NAME not specified or main.go not found" && \
        echo "Available services:" && ls cmd/ && exit 1; \
    fi

# Stage 2: минимальный образ
FROM alpine:latest

# Устанавливаем необходимые пакеты
RUN apk --no-cache add ca-certificates tzdata

# Создаем непривилегированного пользователя
RUN addgroup -g 1000 -S appuser && \
    adduser -u 1000 -S appuser -G appuser

WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder --chown=appuser:appuser /app/bin/ /app/

# Переключаемся на непривилегированного пользователя
USER appuser

#HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
#    CMD wget --no-verbose --tries=1 --spider http://localhost:${PORT:-8080}/health || exit 1

# Запускаем бинарник
ENTRYPOINT ["./main"]
CMD []
