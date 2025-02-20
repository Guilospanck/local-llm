package utils

import (
	"base/pkg/types"
	"math"
)

type Model string

const (
	DEEPSEEK_R1_1_5B Model = "deepseek-r1:1.5b"
	GEMMA_2B               = "gemma:2b"
	LLAMA_3_2              = "llama3.2"
)

const DEFAULT_OLLAMA_MODEL Model = DEEPSEEK_R1_1_5B

const DEFAULT_OLLAMA_SERVER = "http://localhost:7869"
const SERVER_PORT = 4444
const MODEL_TEMPERATURE = 0

const DB_USER = "postgres"
const DB_PASSWORD = "postgres"
const DB_HOSTNAME = "localhost"
const DB_DBNAME = "local-ai"

const MAX_ITEMS_TO_QUERY = 3

var DEFAULT_QUERY = types.QueryProperty{
	Color:    "",
	SizeMin:  0,
	SizeMax:  math.MaxFloat32,
	PriceMin: 0,
	PriceMax: math.MaxFloat64,
	Views:    []string{},
}
