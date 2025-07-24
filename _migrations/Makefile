# Makefile для управления проектом на Docker.
#
# Инструкции:
# - Для добавления новой команды определите новую цель и добавьте комментарий перед ней,
#   чтобы описание было видно при вызове 'make help'.
# - Для использования docker-compose с определенными настройками используйте функцию DOCKER_COMPOSE_UP.
#
# Пример: 
#
# # my-command: Мое описание для команды
# my-command: 
#   <ваш код команды>

## migration: Запуск миграции базы данных
migration:
	cd versions && migrate -database ${POSTGRESQL_URL} -path db/migrations up

## migration-down: Откат миграции базы данных
migration-down:
	cd versions && migrate -database ${POSTGRESQL_URL} -path db/migrations down

## migration-generate: Генерация миграции, пример: make migration-generate name=test_migration_name
migration-generate:
	cd versions && migrate create -ext sql -dir db/migrations $(name)

help:
	@echo "Доступные команды:"
	@echo ""
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'
	@echo ""
