## Описание решения
 Тестирование
	repo - покрыт интеграционными тестами через test-containers
	e2e - реализованы с помощью hurl
 Docs
	настроена генерация swagger документации
	деплоится page с swagger документацией [[[link](https://end1essrage.github.io/efmob-tz/)]]
	поднимается swagger сервис с доступом к мс
 CI
	настроен ci пайплайн - linter(golangci), tests, e2e(hurl), SAST(semgrep)
 API
	настроен rate_limit для всех входящих запросов 1m:100 burst:30()
	логируется информация о выполнении запроса- длительность, путь, user_agent, etc.
 Логирование
	настроена генерация request_id для трейсинга запроса в рамках сервиса
 БД
 	для тест окружения настроена автомиграция
 	retry логика для mutable запросов - экспонециально + jitter
	оптимистичная блокировка + обработка конкурентной попытки записи как retryable ошибки(покрыто тестами)
 Observability
	настроен сбор логов и метрик
	отображение в графане + дашборды

## Business logic annotations
1) Обновлять в данных подписки можно только цену и период(начало и конец)

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







