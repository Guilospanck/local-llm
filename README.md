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

## Extract

There's a simple application via the `/extract` endpoint that will query a database via natural language.

### Preparing the data


