# Makefile для управления проектом на Docker.
#
# Инструкции:
# - Для добавления новой команды определите новую цель и добавьте комментарий перед ней,
#   чтобы описание было видно при вызове 'make help'.
# - Для использования docker-compose с определенными настройками используйте функцию DOCKER_COMPOSE_UP.
#
# Пример: 
#
# ## my-command: Мое описание для команды
# my-command: 
#   <ваш код команды>

# Задаем переменные с именами файлов и переменных окружения
LOCAL_ENV_FILE := env/local.sh
OUT_DOCKER_ENV := env/out_docker.sh
PROD_ENV_FILE := env/prod.sh
TEST_ENV_FILE := env/test.sh

LOCAL_COMPOSE_FILE := docker/local-docker-compose.yml
PROD_COMPOSE_FILE := docker/prod-docker-compose.yml
TEST_COMPOSE_FILE := docker/test-docker-compose.yml

# Цвета для вывода сообщений
RED=\033[0;31m
YELLOW=\033[0;33m
GREEN=\033[0;32m
NC=\033[0m # No Color

# Определяем функции
define DOCKER_COMPOSE_UP
source $(1) && docker compose -f $(2) up $(3) -d
endef

# Базовые команды

## migration: Запуск миграции базы данных
migration: check-psql_url
	cd ../_migrations && migrate -database ${POSTGRESQL_URL} -path versions up

## migration-down: Откат миграции базы данных
migration-down: check-psql_url
	cd ../_migrations && migrate -database ${POSTGRESQL_URL} -path versions down

## migration-generate: Генерация миграции, пример: make migration-generate name=test_migration_name
migration-generate: check-psql_url
	cd ../_migrations && migrate create -ext sql -dir versions $(name)

## local: Запуск локальной среды
local:
	source $(LOCAL_ENV_FILE) && docker compose -f $(LOCAL_COMPOSE_FILE) --profile main up -d

## local-retryer: Запуск локальной среды с retryer
local-retryer:
	source $(LOCAL_ENV_FILE) && docker compose -f $(LOCAL_COMPOSE_FILE) --profile retryer up -d

## local-poller: Запуск локальной среды с poller
local-poller:
	source $(LOCAL_ENV_FILE) && docker compose -f $(LOCAL_COMPOSE_FILE) --profile poller up -d

## local-sync: Запуск локальной среды с sync
local-sync:
	source $(LOCAL_ENV_FILE) && docker compose -f $(LOCAL_COMPOSE_FILE) --profile sync up -d

## prod: Запуск prod среды
prod:
	source $(PROD_ENV_FILE) && docker compose -f $(PROD_COMPOSE_FILE) up -d

## test: Запуск тестов
test:
	$(call DOCKER_COMPOSE_UP,$(TEST_ENV_FILE),$(TEST_COMPOSE_FILE),--build)

# Дополнительные команды для сборки и запуска

## local-build: Запуск poller с опцией сборки
local-poller-build:
	source $(LOCAL_ENV_FILE) && docker compose -f $(LOCAL_COMPOSE_FILE) --profile poller up -d --build

## local-sync-build: Запуск syncer с опцией сборки
local-sync-build:
	source $(LOCAL_ENV_FILE) && docker compose -f $(LOCAL_COMPOSE_FILE) --profile sync up -d --build

## local-retryer-build: Запуск retryer среды с опцией сброки
local-retryer-build:
	source $(LOCAL_ENV_FILE) && docker compose -f $(LOCAL_COMPOSE_FILE) --profile retryer up -d --build

## local-build: Запуск локальной среды с опцией сброки
local-build:
	source $(LOCAL_ENV_FILE) && docker compose -f $(LOCAL_COMPOSE_FILE) --profile main up -d --build

## prod-build: Запуск prod среды с опцией сброки
prod-build:
	source $(PROD_ENV_FILE) && docker compose -f $(PROD_COMPOSE_FILE) --profile main up -d --build

## prod-sync-build: Запуск prod среды с опцией сброки
prod-sync-build:
	source $(PROD_ENV_FILE) && docker compose -f $(PROD_COMPOSE_FILE) --profile sync up -d --build

## mock-gen:
mock-gen:
	mockgen -source=internal/repositories/measure_points/repository.go \
    -destination=internal/repositories/measure_points/mock_repository.go \
    -package=measure_points

	mockgen -source=internal/api_clients/lers/api_client.go \
    -destination=internal/api_clients/lers/mock_api_client.go \
    -package=lers

help:
	@echo "${GREEN}Доступные команды:${NC}"
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//' | awk -F ': ' '{printf "${YELLOW}%-20s${NC} %s\n", $$1, $$2}'
	@echo ""

check-psql_url:
	@if [ -z "${POSTGRESQL_URL}" ]; then \
		echo "${RED}Ошибка: POSTGRESQL_URL не установлен.${NC}"; \
		echo "${YELLOW}Решение: выполните команду 'source env/out_docker.sh' для установки переменных окружения${NC}"; \
		echo "${YELLOW}или используйте export POSTGRESQL_URL=...${NC}"; \
		exit 1; \
	fi