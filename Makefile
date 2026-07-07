include .env
export

export PROJECT_ROOT=${shell pwd}

env-up:
	@docker compose up -d advertapp-postgres

env-down:
	@docker compose down advertapp-postgres

env-cleanup:
	@read -p "Очистить все volume файлы окружения? Опасность утери данных. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down advertapp-postgres port-forwarder && \
		sudo rm -rf "${PROJECT_ROOT}/out/pgdata" && \
		echo "Файлы окружения очищены"; \
	else \
		echo "Очистка окружения отменена"; \
	fi

env-port-forward:
	@docker compose up -d port-forwarder

env-port-close:
	@docker compose down port-forwarder

migrate-create:
	@if [ -z "$(seq)" ]; then \
		echo "Отсутствует необходимый параметр seq. Пример: make migrate-create seq=init"; \
		exit 1; \
	fi; \
	docker compose run --rm advertapp-postgres-migrate \
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
	docker compose run --rm advertapp-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@advertapp-postgres:5432/${POSTGRES_DB}?sslmode=disable \
		"$(action)"

grpc-gen-advert:
	@protoc -I internal/protos \
	--go_out=internal/gen \
	--go_opt=paths=source_relative \
	--go-grpc_out=internal/gen \
	--go-grpc_opt=paths=source_relative \
	advert/advert.proto

grpc-gen-user:
	@protoc -I internal/protos \
	--go_out=internal/gen \
	--go_opt=paths=source_relative \
	--go-grpc_out=internal/gen \
	--go-grpc_opt=paths=source_relative \
	user/user.proto

grpc-gen-category:
	@protoc -I internal/protos \
	--go_out=internal/gen \
	--go_opt=paths=source_relative \
	--go-grpc_out=internal/gen \
	--go-grpc_opt=paths=source_relative \
	category/category.proto

grpc-gen-favourite:
	@protoc -I internal/protos \
	--go_out=internal/gen \
	--go_opt=paths=source_relative \
	--go-grpc_out=internal/gen \
	--go-grpc_opt=paths=source_relative \
	favourite/favourite.proto

app-run:
	@export LOGGER_FOLDER="${PROJECT_ROOT}/out/logs" && \
	export POSTGRES_HOST=localhost && \
	go run "${PROJECT_ROOT}/cmd/app/main.go"