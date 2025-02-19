package utils

import (
	"base/pkg/types"
	"math"
)

const DEFAULT_OLLAMA_SERVER = "http://localhost:7869"
const SERVER_PORT = 4444
const MODEL_TEMPERATURE = 0

const DB_USER = "postgres"
const DB_PASSWORD = "postgres"
const DB_HOSTNAME = "localhost"
const DB_DBNAME = "local-ai"

var DEFAULT_QUERY = types.QueryProperty{
	Color:    "",
	SizeMin:  0,
	SizeMax:  math.MaxFloat32,
	PriceMin: 0,
	PriceMax: math.MaxFloat64,
	Views:    []string{},
}
