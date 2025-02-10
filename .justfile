front-init:
	cd front/ && pnpm i && pnpm dev

front-dev:
	cd front/ && pnpm dev

back-init:
	cd back/ && go mod tidy && go run .

back-dev:
	cd back/ && air
