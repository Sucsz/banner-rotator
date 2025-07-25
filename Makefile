BINARY_NAME=banner-rotator
MAIN_FILE=cmd/main.go

ENV_FILE=.env
ifneq (,$(wildcard $(ENV_FILE)))
  include $(ENV_FILE)
  # Экспортим все ключи из .env в окружение Make
  export $(shell sed 's/=.*//' $(ENV_FILE))
endif

DOCKER_COMPOSE=docker compose --env-file $(ENV_FILE)

## Запускает контейнеры Kafka, PostgreSQL и т.д.
run:
	$(DOCKER_COMPOSE) up -d

## Останавливает и удаляет контейнеры + сеть
stop:
	$(DOCKER_COMPOSE) down

## Перезапускает инфраструктуру (stop + run)
restart:
	$(MAKE) stop
	$(MAKE) run
	@echo "Ensuring Kafka topic '$(APP_KAFKA_TOPIC)' exists…"
	$(DOCKER_COMPOSE) exec kafka \
	  kafka-topics --bootstrap-server $(APP_KAFKA_BROKERS) \
	    --create --if-not-exists \
	    --topic $(APP_KAFKA_TOPIC) \
	    --partitions 1 \
	    --replication-factor 1 || true

## Пересобирает образ приложения без кэша
rebuild:
	$(DOCKER_COMPOSE) build --no-cache app

## Пересобирает образ и перезапускает сервисы
full-restart:
	$(MAKE) rebuild
	$(MAKE) restart

## Показывает логи всех контейнеров
logs:
	$(DOCKER_COMPOSE) logs -f

## Собирает Go-бинарник из cmd/main.go
build:
	go build -o $(BINARY_NAME) $(MAIN_FILE)

## Запускает все юнит-тесты с -race
test:
	go test -race -count=1 ./...

## Локальный запуск приложения с загрузкой ENV из .env
local:
	bash -c '\
	  set -o allexport; source $(ENV_FILE); set +o allexport; \
	  export APP_POSTGRES_HOST=localhost \
	  		 APP_POSTGRES_PORT=$$HOST_POSTGRES_PORT \
	  		 APP_KAFKA_BROKERS=localhost:9092; \
	  go run $(MAIN_FILE) \
	'
## Запускает golangci-lint
lint:
	golangci-lint run ./...

## Показывает список всех команд с описанием
help:
	@echo "Доступные команды:"
	@awk 'BEGIN {FS = ":.*?## "} \
		/^##/ { help = substr($$0, 4); next } \
		/^[a-zA-Z0-9_-]+:/ && help { \
			sub(/:.*/, "", $$1); \
			printf "  \033[36m%-12s\033[0m %s\n", $$1, help; \
			help = "" \
		}' Makefile

## Интеграционные смок-тесты на основные хендлеры апи
smoke:
	./tests/http/smoke.sh



.PHONY: run stop restart logs build test lint help
