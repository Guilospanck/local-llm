front-init:
	cd front/ && pnpm i && pnpm dev

front-dev:
	cd front/ && pnpm dev

back-init MODEL='deepseek-r1:1.5b':
	bash -c 'source ./scripts/valid_models.sh && validate_model "{{MODEL}}"'
	cd back/ && go mod tidy && OLLAMA_MODEL={{MODEL}} go run .

back-dev MODEL='deepseek-r1:1.5b':
	bash -c 'source ./scripts/valid_models.sh && validate_model "{{MODEL}}"'
	cd back/ && OLLAMA_MODEL={{MODEL}} air

dev-watch MODEL='deepseek-r1:1.5b':
	bash -c 'source ./scripts/valid_models.sh && validate_model "{{MODEL}}"'
	OLLAMA_MODEL={{MODEL}} ./scripts/dev.sh --watch

dev MODEL='deepseek-r1:1.5b':
	bash -c 'source ./scripts/valid_models.sh && validate_model "{{MODEL}}"'
	OLLAMA_MODEL={{MODEL}} ./scripts/dev.sh 

dcup:
	docker-compose up -d --build --remove-orphans

dc-ollama MODEL='deepseek-r1:1.5b':
	bash -c 'source ./scripts/valid_models.sh && validate_model "{{MODEL}}"'
	OLLAMA_MODEL={{MODEL}} docker-compose up -d --build ollama model-puller

dcdown:
	docker-compose down

start-postgres:
	docker rm -f local-postgres
	docker run --rm -d \
		--name local-postgres \
		-v $(pwd)/postgres/data:/var/lib/postgresql/data \
		-v $(pwd)/back/pkg/domain/data/init.sql:/docker-entrypoint-initdb.d/init.sql \
		-e POSTGRES_USERNAME=postgres \
		-e POSTGRES_PASSWORD=postgres \
		-e POSTGRES_DB=local-ai \
		-p 5432:5432 \
		-d postgres:latest

MIGRATIONS_PATH := "./back/pkg/domain/data/migrations/*"
SEEDS_PATH := "./back/pkg/domain/data/seeds"
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
