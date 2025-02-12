front-init:
	cd front/ && pnpm i && pnpm dev

front-dev:
	cd front/ && pnpm dev

back-init:
	cd back/ && go mod tidy && go run .

back-dev:
	cd back/ && air

############## FIXME ##################################
MIGRATIONS_PATH := "pkg/domain/migrations"
POSTGRES_URL := "postgres://postgres:postgres@localhost:5432/local-ai?sslmode=disable"

# requires golang-migrate
create-migration MIGRATION_NAME:
	migrate create -ext sql -dir {{MIGRATIONS_PATH}} -seq MIGRATION_NAME

migration-up:
	migrate -database "{{POSTGRES_URL}}" -path {{MIGRATIONS_PATH}} up

migration-down:
	migrate -database "{{POSTGRES_URL}}" -path {{MIGRATIONS_PATH}} down -all
######################################################33333
