package utils

import (
	"base/pkg/types"
	"math"
)

const DEFAULT_OLLAMA_SERVER = "http://localhost:7869"
const SERVER_PORT = 4444
const MODEL_TEMPERATURE = 0

var DEFAULT_QUERY = types.QueryProperty{
	Color:    "",
	SizeMin:  0,
	SizeMax:  math.MaxFloat32,
	PriceMin: 0,
	PriceMax: math.MaxFloat64,
	Views:    []string{},
}
