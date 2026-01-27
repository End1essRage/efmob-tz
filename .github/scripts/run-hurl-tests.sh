#!/bin/bash
# .github/scripts/run-hurl-tests.sh
# Запуск HURL тестов

set -e

echo "Running HURL tests..."

# Вспомогательные функции
source ./.github/scripts/generate-compose-flags.sh

COMPOSE_FLAGS=$(generate_compose_flags "$COMPOSE_FILES")
echo "Compose flags: $COMPOSE_FLAGS"

# Проверяем наличие сервиса api-tests в docker-compose
if docker compose $COMPOSE_FLAGS config --services | grep -q "api-tests"; then
    echo "Using api-tests service from docker-compose..."
    docker compose $COMPOSE_FLAGS run --rm api-tests
else
    echo "Creating and running HURL tests container..."

    # Получаем имя сети
    NETWORK_NAME=$(docker compose $COMPOSE_FLAGS ps -q subs | head -1 | xargs docker inspect -f '{{range .NetworkSettings.Networks}}{{.NetworkID}}{{end}}')

    docker run --rm \
        --network "$NETWORK_NAME" \
        -v "$PWD/tests:/tests:ro" \
        ghcr.io/orange-opensource/hurl:4.2.0 \
        sh -c "
            echo 'Waiting for subs service to respond...'
            for i in {1..10}; do
                if wget -qO- http://subs:8080/health > /dev/null 2>&1; then
                    echo 'Service is responding!'
                    break
                fi
                echo 'Waiting... (\$i/10)'
                sleep 2
            done

            echo 'Running HURL tests...'
            find /tests -name '*.hurl' -type f

            # Запускаем все тесты
            for file in /tests/**/*.hurl; do
                if [ -f \"\$file\" ]; then
                    echo \"Running test: \$file\"
                    hurl --test \"\$file\"
                fi
            done
        "
fi