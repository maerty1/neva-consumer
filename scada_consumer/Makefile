# # Makefile для управления проектом на Docker.
# #
# # Инструкции:
# # - Для добавления новой команды определите новую цель и добавьте комментарий перед ней,
# #   чтобы описание было видно при вызове 'make help'.
# # - Для использования docker-compose с определенными настройками используйте функцию DOCKER_COMPOSE_UP.
# #
# # Пример: 
# # Test build
# # ## my-command: Мое описание для команды 
# # my-command: 
# #   <ваш код команды>

# # Задаем переменные с именами файлов и переменных окружения
# LOCAL_ENV_FILE := env/local.sh
# PROD_ENV_FILE := env/prod.sh
# TEST_ENV_FILE := env/test.sh

# LOCAL_COMPOSE_FILE := docker/local-local-docker-compose.yml
# PROD_COMPOSE_FILE := docker/prod-local-docker-compose.yml
# TEST_COMPOSE_FILE := docker/test-local-docker-compose.yml

# # Цвета для вывода сообщений
# RED=\033[0;31m
# YELLOW=\033[0;33m
# GREEN=\033[0;32m
# NC=\033[0m # No Color

# # Определяем функции
# define DOCKER_COMPOSE_UP
# source $(1) && docker compose -f $(2) up $(3) -d
# endef

# # Базовые команды

# ## migration: Запуск миграции базы данных
# migration: check-psql_url
# 	cd ../_migrations && migrate -database ${POSTGRESQL_URL} -path versions up

# ## migration-down: Откат миграции базы данных
# migration-down: check-psql_url
# 	cd ../_migrations && migrate -database ${POSTGRESQL_URL} -path versions down

# ## migration-generate: Генерация миграции, пример: make migration-generate name=test_migration_name
# migration-generate: check-psql_url
# 	cd ../_migrations && migrate create -ext sql -dir versions $(name)

# ## local: Запуск локальной среды
# local:
# 	$(call DOCKER_COMPOSE_UP,$(LOCAL_ENV_FILE),$(LOCAL_COMPOSE_FILE),)

# ## prod: Запуск prod среды
# prod:
# 	$(call DOCKER_COMPOSE_UP,$(PROD_ENV_FILE),$(PROD_COMPOSE_FILE),)

# ## test: Запуск тестов
# test:
# 	$(call DOCKER_COMPOSE_UP,$(TEST_ENV_FILE),$(TEST_COMPOSE_FILE),--build)

# # Дополнительные команды для сборки и запуска

# ## local-build: Запуск локальной среды с опцией сборки
# local-build:
# 	$(call DOCKER_COMPOSE_UP,$(LOCAL_ENV_FILE),$(LOCAL_COMPOSE_FILE),--build)

# ## prod-build: Запуск prod среды с опцией сборки
# prod-build:
# 	$(call DOCKER_COMPOSE_UP,$(PROD_ENV_FILE),$(PROD_COMPOSE_FILE),--build)

# help:
# 	@echo "${GREEN}Доступные команды:${NC}"
# 	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//' | awk -F ': ' '{printf "${YELLOW}%-20s${NC} %s\n", $$1, $$2}'
# 	@echo ""

# check-psql_url:
# 	@if [ -z "${POSTGRESQL_URL}" ]; then \
# 		echo "${RED}Ошибка: POSTGRESQL_URL не установлен.${NC}"; \
# 		echo "${YELLOW}Решение: выполните команду 'source env/out_docker.sh' для установки переменных окружения${NC}"; \
# 		echo "${YELLOW}или используйте export POSTGRESQL_URL=...${NC}"; \
# 		exit 1; \
# 	fi
