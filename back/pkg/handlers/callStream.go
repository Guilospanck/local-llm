package handlers

import (
	"base/pkg/types"
	"base/pkg/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

func CallStream(c echo.Context, modelName *string) error {
	var body types.QueryRequestData

	// parse the body into a json structure
	if err := c.Bind(&body); err != nil {
		fmt.Print(err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid JSON"})
	}

	ollama_server, exists := os.LookupEnv("OLLAMA_SERVER")
	if !exists {
		ollama_server = utils.DEFAULT_OLLAMA_SERVER
	}

	llm, err := ollama.New(ollama.WithModel(*modelName), ollama.WithServerURL(ollama_server))
	if err != nil {
		log.Fatal(err)
	}

	// Set headers for streaming
	c.Response().Header().Set("Content-Type", "text/plain")
	c.Response().Header().Set("Transfer-Encoding", "chunked")

	// Get the response writer and flusher
	w := c.Response().Writer
	flusher, ok := w.(http.Flusher)
	if !ok {
		return c.String(http.StatusInternalServerError, "Streaming unsupported")
	}

	ctx := context.Background()

	// "Human: %s\nAssistant:", body.Query
	prompt := body.Query

	// INFO: the `llm.Call` method does something like this:
	// messages := []llms.MessageContent{
	// 	{
	// 		Role:  llms.ChatMessageTypeSystem,
	// 		Parts: []llms.ContentPart{llms.TextContent{Text: prompt}},
	// 	},
	// }
	// completion, err := llm.GenerateContent(ctx, messages,
	//
	// You could use it to pass different types of data (like binary data, for example)

	completion, err := llm.Call(ctx, prompt,
		llms.WithTemperature(utils.MODEL_TEMPERATURE),
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			fmt.Fprint(w, string(chunk))
			flusher.Flush() // Flush data immediately

			return nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Prevents the compiler from saying that `completion` is not being used
	_ = completion

	return nil
}
