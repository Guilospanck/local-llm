front-init:
	cd front/ && pnpm i && pnpm dev

front-dev:
	cd front/ && pnpm dev

back-init:
	cd back/ && go mod tidy && go run .

back-dev:
	cd back/ && air

start-postgres:
	docker run --rm -d --name local-postgres -v ./postgres/data:/var/lib/postgresql/data -e POSTGRES_PASSWORD=postgres -p 5432:5432 -d postgres

MIGRATIONS_PATH := "./back/pkg/domain/migrations"
SEEDS_PATH := "./back/pkg/domain/seeds/*"
POSTGRES_URL := "postgres://postgres:postgres@localhost:5432/local-ai?sslmode=disable"
POSTGRES_USERNAME := "postgres"
POSTGRES_PASSWORD := "postgres"
POSTGRES_DB := "local-ai"

# requires golang-migrate
create-migration MIGRATION_NAME:
	migrate create -ext sql -dir {{MIGRATIONS_PATH}} {{MIGRATION_NAME}}

create-seed SEED_NAME:
	migrate create -ext sql -dir {{SEEDS_PATH}} {{SEED_NAME}}

migration-up:
	migrate -database "{{POSTGRES_URL}}" -path {{MIGRATIONS_PATH}} up

migration-down:
	migrate -database "{{POSTGRES_URL}}" -path {{MIGRATIONS_PATH}} down -all

seed-up:
	for seed in {{SEEDS_PATH}}; do PGPASSWORD={{POSTGRES_PASSWORD}} psql -h localhost -U {{POSTGRES_USERNAME}} --no-password -d {{POSTGRES_DB}} -f $seed; done
