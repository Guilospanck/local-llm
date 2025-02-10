package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

const DEFAULT_OLLAMA_SERVER = "http://localhost:7869"
const SERVER_PORT = 4444
const MODEL_TEMPERATURE = 0.1

var modelName *string

type QueryRequestData struct {
	Query string `json:"query"`
}

func main() {
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	modelName = flag.String("model", "deepseek-r1", "Deep Seek R1")
	flag.Parse()

	// Routes
	e.GET("/healthz", func(c echo.Context) error { return c.String(http.StatusOK, "I'm healthy") })
	e.POST("/", callStream)

	// Start server
	if err := e.Start(fmt.Sprintf(":%d", SERVER_PORT)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}

func callStream(c echo.Context) error {
	var body QueryRequestData

	// parse the body into a json structure
	if err := c.Bind(&body); err != nil {
		fmt.Print(err)
		return c.JSON(http.StatusBadRequest, echo.Map{"error": "Invalid JSON"})
	}

	ollama_server, exists := os.LookupEnv("OLLAMA_SERVER")
	if !exists {
		ollama_server = DEFAULT_OLLAMA_SERVER
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
		llms.WithTemperature(MODEL_TEMPERATURE),
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

func callSimple(c echo.Context) error {
	ollama_server, exists := os.LookupEnv("OLLAMA_SERVER")
	if !exists {
		ollama_server = DEFAULT_OLLAMA_SERVER
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
