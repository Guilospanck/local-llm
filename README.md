# Local LLM

Easiest way to have your own local setup running LLM.
It uses [Ollama](https://ollama.com/) with the [Deep Seek R1](https://ollama.com/library/deepseek-r1) model.

## Installation

Clone the repository:

```shell
git clone https://github.com/Guilospanck/local-llm.git
cd local-llm/
```
And be sure to have [Go](https://go.dev/doc/install), [pnpm](https://pnpm.io/installation#using-npm), [Docker](https://docs.docker.com/get-started/get-docker/) and [Docker compose](https://docs.docker.com/compose/install/) installed.

### Optional

Install [just](https://github.com/casey/just) and [air](https://github.com/air-verse/air) to simplify the development in dev mode.


## Running

### Docker compose

If you want to just have it running to try it out, do:

```shell
docker-compose up -d --build
```

You should be able able now to see it running at http://localhost:3000/.

To bring it down after using it:

```shell
docker-compose down
```

Otherwise, if you want to run things in Dev mode, check next section.

### Dev mode

Be sure to have Ollama running with the `Deep Seek R1` model:

```shell
docker-compose up -d ollama model-puller
```

To bring it down after using it:

```shell
docker-compose down ollama
```

Now you have two ways of initialising the other components:

#### With `just`

##### Sequentially

In one terminal:

```shell
just front-init
```

In another:

```shell
just back-init
```

##### Concurrently

Non-watch mode:

```shell
just dev
```

Watch mode (requires [air](https://github.com/air-verse/air)):

```shell
just dev-watch
```

#### Manual way

To have the React frontend running:

In one terminal:

1 - Be sure to have `pnpm` installed:

```shell
npm install -g pnpm
```

2 - Then:

```shell
cd front/
pnpm install
pnpm dev
```

To have the Golang backend running:

In another terminal:
```shell
cd back/
go mod tidy
go run .
```

## Database

We want to spin up a Postgres database with some information. The next commands will create a databsae called `local-ai` that will be available at default port `5432` with three tables: `property`, `view` and `property_view` (Many-2-Many table between property and view).

> [!TIP]
> You can easily create and run migrations by using the `just` formulas. See [Migrations](#migrations).

### Docker-compose way

Just run

```shell
docker-compose up -d postgres

```

To bring it down after using it:

```shell
docker-compose down postgres
```

### Docker way

Just run

```shell
just start-postgres
```

### Migrations

Create a migration file by running:

```shell
just create-migration MIGRATION_NAME
```

where `MIGRATION_NAME` is whatever name you want to give to it.

After you create it and add the SQL information there, you can run migrations via (requires [migrate](https://github.com/golang-migrate/migrate)):

```shell
just migration-up
# or to rollback
just migration-down
```


## Possible calls

There are two types of calls for the application: [Chat](#chat) and [Extract](#extract).

### Chat

This is the normal GPT-like chat with the LLM. Just call the `/` endpoint with a query.

Example:

```shell
curl 'http://localhost:4444/' \
  -H 'Accept: */*' \
  -H 'Connection: keep-alive' \
  -H 'Content-Type: application/json' \
  --data-raw $'{"query":"Hi"}'
```

### Extract

It's used to extract some characteristics from a natural language input. This is used to query the database.

There's a simple application via the `/extract` endpoint that will query a database via natural language. Before running it, be sure to have the [Database](#database) up and running.


#### Running the extract 

```shell
curl 'http://localhost:4444/extract' \
  -H 'Content-Type: application/json' \
  --data-raw '{"query":"I want a big house, close to the sea and to the mountains. Not very expensive. Maybe marble colored"}'
```

It should return a JSON object with some characteristics of the property you wanna buy.

