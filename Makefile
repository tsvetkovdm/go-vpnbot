include .env
export

export PROJECT_ROOT=$(shell pwd)
export DOCKER_USER=$(shell id -u):$(shell id -g)

env-up:
	@docker compose up -d vpnbot-postgres

env-down:
	@docker compose down vpnbot-postgres

env-cleanup:
	@read -p "Очистить все volume файлы окружения? Опасность утери данных. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down vpnbot-postgres port-forwarder && \
		docker run --rm -v $(PROJECT_ROOT)/out:/out alpine rm -rf /out/.pgdata && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует необходимый параметр seq. Пример: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	mkdir -p ./migrations
	docker compose run --rm vpnbot-postgres-migrate \
		create \
		-ext sql \
		-dir /migrations \
		-seq "$(seq)"

migrate-up:
	@make migrate-action action=up

migrate-down:
	@make migrate-action action=down

migrate-action:
	@if [ -z "$(action)" ]; then \
		echo "Отсутствует необходимый параметр action. Пример: make migrate-action action=up"; \
		exit 1; \
	fi; \
	docker compose run --rm vpnbot-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@vpnbot-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		$(action)

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder