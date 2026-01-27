#!/bin/bash
# .github/scripts/generate-compose-flags.sh
# Генерация флагов для docker compose

generate_compose_flags() {
    local compose_files="$1"
    local flags=""

    while IFS= read -r line || [[ -n "$line" ]]; do
        trimmed_line=$(echo "$line" | xargs)
        if [[ -n "$trimmed_line" ]]; then
            flags="$flags -f $trimmed_line"
        fi
    done <<< "$compose_files"

    echo "${flags# }"
}

# Экспортируем функцию для использования в других скриптах
export -f generate_compose_flags