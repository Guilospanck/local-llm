package handlers

import (
	"base/pkg/utils"
	"context"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func CallSimple(c echo.Context, modelName *string) error {
	ollama_server, exists := os.LookupEnv("OLLAMA_SERVER")
	if !exists {
		ollama_server = utils.DEFAULT_OLLAMA_SERVER
	}

	llm, err := ollama.New(ollama.WithModel(*modelName), ollama.WithServerURL(ollama_server))
	if err != nil {
		log.Fatal(err)
	}

	query := "very briefly, tell me the difference between a comet and a meteor"

	ctx := context.Background()
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, query)
	if err != nil {
		log.Fatal(err)
	}

	return c.String(http.StatusOK, completion)
}
