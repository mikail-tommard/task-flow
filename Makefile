DB_DSN ?= postgres://tasksflowuser:tasksflowpass@localhost:5433/taskflowdb?sslmode=disable
MIGRATION_DIR := $(CURDIR)/migrations

migrate-up:
	migrate -path $(MIGRATION_DIR) -database "$(DB_DSN)" up

migrate-down:
	migrate -path $(MIGRATION_DIR) -database "$(DB_DSN)" down

migrate-force:
	migrate -path $(MIGRATIONS_DIR) -database "$(DB_DSN)" force $(V)

fmt:
	@go fmt ./...