#!/bin/bash
# .github/scripts/docker-compose.sh
# Управление Docker Compose в CI

set -e

# Вспомогательные функции
source ./.github/scripts/generate-compose-flags.sh

setup() {
    echo "Setting up Docker Compose..."
    sudo apt-get update
    sudo apt-get install -y docker-compose-v2
    docker compose version
}

start-services() {
    echo "Building and starting services..."

    COMPOSE_FLAGS=$(generate_compose_flags "$COMPOSE_FILES")
    echo "Using docker compose flags: $COMPOSE_FLAGS"

    # Собираем сервис
    docker compose $COMPOSE_FLAGS build subs

    # Запускаем зависимости
    docker compose $COMPOSE_FLAGS up -d postgres

    # Ждем PostgreSQL
    echo "Waiting for PostgreSQL to start..."
    for i in {1..30}; do
        if docker compose $COMPOSE_FLAGS exec -T postgres pg_isready -U $POSTGRES_USER; then
            echo "PostgreSQL is ready!"
            break
        fi
        echo "Waiting for PostgreSQL... ($i/30)"
        sleep 2
    done

    # Запускаем основной сервис
    docker compose $COMPOSE_FLAGS up -d subs

    # Ждем сервис
    echo "Waiting for subs service to start..."
    for i in {1..30}; do
        if docker compose $COMPOSE_FLAGS exec -T subs wget -qO- http://localhost:8080/health > /dev/null 2>&1; then
            echo "Subs service is ready!"
            break
        fi
        echo "Waiting for subs service... ($i/30)"
        sleep 2
    done
}

logs() {
    COMPOSE_FLAGS=$(generate_compose_flags "$COMPOSE_FILES")
    echo "========== Service logs =========="
    docker compose $COMPOSE_FLAGS logs --tail=100
    echo "========== End logs =========="
}

cleanup() {
    COMPOSE_FLAGS=$(generate_compose_flags "$COMPOSE_FILES")
    echo "Cleaning up Docker Compose..."
    docker compose $COMPOSE_FLAGS down -v --remove-orphans
}

# Обработка команд
case "$1" in
    setup)
        setup
        ;;
    start-services)
        start-services
        ;;
    logs)
        logs
        ;;
    cleanup)
        cleanup
        ;;
    *)
        echo "Usage: $0 {setup|start-services|logs|cleanup}"
        exit 1
        ;;
esac