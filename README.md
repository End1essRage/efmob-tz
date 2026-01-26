## run
1. устанавливаем go v1.24
2. устаналиваем docker
3. устанавливаем taskfile
4. task up

## deploy definition - что поднимается при "task up"
r-proxy:
	traefik
observability:
	loki, grafana, alloy, prometheus
infra:
	postgres
other:
	swagger

## tasks
e2e - перезапустит сервис и запустит hurl тесты(для локального использования в dev окружении без внешних зависимостей)
ci:test-integration - запуск интеграционных тестов с поднятием сервисов через testcontainers


## Environment:
dev:
	loglevel: debug
	db: inmemory
test:
	loglevel: debug
	db: postgres


## links
[[[logs](http://localhost:3000/a/grafana-lokiexplore-app/explore/service/unknown_service/logs?from=now-15m&to=now&var-ds=P8E80F9AEF21F6940&var-filters=service_name%7C%3D%7Cunknown_service&patterns=%5B%5D&var-lineFormat=&var-fields=service%7C%3D%7C%7B%22parser%22:%22json%22__gfc__%22value%22:%22subs%22%7D,subs&var-levels=&var-metadata=&var-jsonFields=&var-patterns=&var-lineFilterV2=&var-lineFilters=&timezone=browser&var-all-fields=service%7C%3D%7C%7B%22parser%22:%22json%22__gfc__%22value%22:%22subs%22%7D,subs&displayedFields=%5B%22_caller%22,%22_message%22,%22package%22,%22service%22%5D&urlColumns=%5B%5D&visualizationType=%22logs%22&prettifyLogMessage=false&sortOrder=%22Descending%22&wrapLogMessage=false)]]
