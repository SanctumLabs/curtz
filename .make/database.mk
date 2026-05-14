# Database defaults

DATABASE_MIGRATION_PATH ?= ./internal/infra/database/postgresql/migrations
DATABASE_USER ?= curtz-svc-user
DATABASE_PASSWORD ?= curtz-svc-pass
DATABASE_NAME ?= curtz_system_db
DATABASE_HOST ?= localhost
DATABASE_PORT ?= 5433
DATABASE_SCHEMA ?= curtz_system
DATABASE_SCHEMA_MIGRATION ?= curtz_schema_migrations
MIGRATE_DIRECTION ?= up

.PHONY: sqlc.generate
sqlc.generate: # generate SQLC code
	@echo "${GREEN} Generating SQLC code ${NC}"
	sqlc generate -f sqlc.yaml
	@echo "${GREEN} Done generating SQLC code ${NC}"

.PHONY: migrate
migrate: createDockerEnvFile ## Runs the migrations. Defaults to up, usage: make migrate MIGRATE_DIRECTION=down
	@echo "${GREEN} >>>>> Running migrations from directory $(DATABASE_MIGRATION_PATH) ${NC}"
	docker run -v $(DATABASE_MIGRATION_PATH):/migrations --network host migrate/migrate -path=/migrations/ -database 'postgres://$(DATABASE_USER):$(DATABASE_PASSWORD)@$(DATABASE_HOST):$(DATABASE_PORT)/$(DATABASE_NAME)?sslmode=disable&x-migrations-table=$(DATABASE_SCHEMA_MIGRATION)' $(MIGRATE_DIRECTION)
	@echo "${GREEN} >>>> Done running migrations ${NC}"

.PHONY: migrateCreate
migrateCreate: createDockerEnvFile ## Creates a migration up and down sql script, usage: make migrateCreate MIGRATE_NAME=20220101_create_parking
	@echo "${GREEN} >>>>> Creating migration script $(MIGRATION_NAME) in $(DATABASE_MIGRATION_PATH) ${NC}"
	docker run --rm -v $(shell pwd)/$(DATABASE_MIGRATION_PATH):/migrations --network host migrate/migrate create -ext sql -dir /migrations -seq $(MIGRATION_NAME)
	@echo "${GREEN} >>>> Done creating migration script $(MIGRATION_NAME) ${NC}"
