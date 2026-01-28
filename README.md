## Описание решения
### Тестирование
- **Репозиторий** - покрыт интеграционными тестами через TestContainers
- **E2E тесты** - реализованы с помощью Hurl
### Docs
- Настроена генерация Swagger документации
- Автоматический деплой документации на GitHub Pages
  [Документация API](https://end1essrage.github.io/efmob-tz/)
- Локальный Swagger сервис с доступом к микросервису
### CI
Настроен CI пайплайн со следующими этапами:
1. **Linter** - golangci-lint
2. **Unit-тесты** - go test
3. **E2E тесты** - Hurl
4. **SAST** - Semgrep (статический анализ безопасности)
### API
- **Rate Limiting** - настроен для всех входящих запросов: 100 запросов/мин, burst: 30
- **Логирование** - записывается информация о выполнении запроса:
  - Длительность выполнения
  - Путь запроса
  - User-Agent
  - Другие метаданные
### Логирование
- Генерация уникального `request_id` для трейсинга запросов в рамках сервиса
- Структурированные логи в формате JSON
### БД
- **Тестовое окружение** - настроена автоматическая миграция
- **Retry логика** для mutable запросов:
  - Экспоненциальная backoff стратегия
  - Добавлен jitter для предотвращения thundering herd
- **Оптимистичная блокировка** с обработкой конкурентных запросов как retryable ошибок (покрыто тестами)
### Observability
- Сбор логов и метрик в реальном времени
- Визуализация в Grafana с готовыми дашбордами
- Интеграция с:
  - Loki (хранение логов)
  - Prometheus (хранение метрик)
  - Alloy (сбор данных)

## Бизнес-логика

### Правила обновления подписок
При обновлении данных подписки можно изменять только:
1. **Цену** подписки
2. **Период** подписки (дата начала и окончания)

## Run
1. устанавливаем go v1.24
2. устаналиваем docker
3. устанавливаем taskfile
4. task up

## Deploy definition - что поднимается при "task up"
r-proxy:
	traefik
observability:
	loki, grafana, alloy, prometheus
infra:
	postgres
other:
	swagger




## Links
[[[logs](http://localhost:3000/a/grafana-lokiexplore-app/explore/service/unknown_service/logs?from=now-15m&to=now&var-ds=P8E80F9AEF21F6940&var-filters=service_name%7C%3D%7Cunknown_service&patterns=%5B%5D&var-lineFormat=&var-fields=service%7C%3D%7C%7B%22parser%22:%22json%22__gfc__%22value%22:%22subs%22%7D,subs&var-levels=&var-metadata=&var-jsonFields=&var-patterns=&var-lineFilterV2=&var-lineFilters=&timezone=browser&var-all-fields=service%7C%3D%7C%7B%22parser%22:%22json%22__gfc__%22value%22:%22subs%22%7D,subs&displayedFields=%5B%22_caller%22,%22_message%22,%22package%22,%22service%22%5D&urlColumns=%5B%5D&visualizationType=%22logs%22&prettifyLogMessage=false&sortOrder=%22Descending%22&wrapLogMessage=false)]]

[[[swagger-local-deploy](http://localhost/swagger)]]







