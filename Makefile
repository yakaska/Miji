include .env
export

export PROJECT_ROOT = $(shell pwd)

env-up:
	@docker compose up -d

env-down:
	@docker compose down
  
env-cleanup:
	read -p "Clear all env volumes? Possible data loss. [y/N]: " ans; \
	if [ "$$ans" = "y" ]; then \
		docker compose down && \
		rm -rf out/pgdata && \
		echo "Done."; \
	else \
		echo "Cleanup canceled."; \
	fi
	docker compose down

migrate-create:
	@if [-z "$(seq)"]; then \
		echo "No var seq. Ex: make migrate-create seq=init"; \
		exit 1; \
	fi; \

	docker compose run --rm miji-postgres-migrate \
	  create \
	  -ext sql \
	  -dir /migrations \
	  -seq "$(seq)"

migrate:
	@if [-z "$(action)"]; then \
		echo "No var action. Ex: make migrate-create action=up 1"; \
		exit 1; \
	fi; \

	docker compose run --rm miji-postgres-migrate \
		-path /migrations \
		-database postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@miji-postgres:${POSTGRES_PORT}/${POSTGRES_DB}?sslmode=disable \
		$(action)