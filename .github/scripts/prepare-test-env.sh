#!/bin/bash
# .github/scripts/prepare-test-env.sh
# Подготовка тестового окружения

set -e

echo "Preparing test environment..."

# Создаем .env файл для тестовой среды
cat > .env << EOF
# Application
ENV=test
SERVICE_NAME=subs
PORT=8080

# PostgreSQL DSN для тестов
POSTGRES_DSN=host=postgres user=${POSTGRES_USER} password=${POSTGRES_PASSWORD} dbname=${POSTGRES_DB} port=5432 sslmode=disable

# Отключаем ненужные сервисы для CI
DISABLE_TRAEFIK=true
DISABLE_SWAGGER_UI=true
EOF

echo "Test environment prepared successfully."