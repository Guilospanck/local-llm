package main

import (
	"base/pkg/domain"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strings"

	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

const DEFAULT_OLLAMA_SERVER = "http://localhost:7869"
const SERVER_PORT = 4444
const MODEL_TEMPERATURE = 0

var modelName *string

type QueryRequestData struct {
	Query string `json:"query"`
}

var db *domain.Database

func main() {
	e := echo.New()

	db = domain.NewDb()
	db.Connect()

	defer db.Close()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	modelName = flag.String("model", "deepseek-r1", "Deep Seek R1")
	flag.Parse()

	// Routes
	e.GET("/healthz", func(c echo.Context) error { return c.String(http.StatusOK, "I'm healthy") })
	e.POST("/", callStream)
	e.POST("/call", callSimple)
	e.POST("/extract", extract)

	// Start server
	if err := e.Start(fmt.Sprintf(":%d", SERVER_PORT)); err != nil && !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed to start server", "error", err)
	}
}

var extractPrompt string = `
	You are an AI that only provides direct answers. Do not include <think> or reasoning steps.
	Be brief and direct to the point.

	Given the input "%s", extract info that might fall into some house categories,
	like "views", "size", "color", "priceMin", "priceMax" and so on. Answer giving a JSON object
	with the keys and their values.

	The prices should be type number (int).
	The views can be an array of strings.
	Other keys should be strings.

	If the price is given in some natural language, like 'not expensive', try to
	fit it into a range of price that would make sense considering the current
	house market.

	A good range of price, probably:
	- cheap: 0-100000
	- medium: 101000-500000
	- expensive: +500000

	In the case that the characteristic of the property falls into "expensive", which doesn't have a maximum price,
	only minimum price (500000), we should set priceMax to null. The other categories should set the adequate priceMin
	and priceMax based on the information given above (cheap, medium).

	But if the price is given to you in numbers, like "will spend until 300000", you should set the priceMin to 0
	and priceMax to that value. If it is something like "will spend minimum of 200000", you should set the
	priceMin to that value and set the priceMax to null (because no maximum price was given). If, on the other hand,
	the user gives you "will spend between 1000 and 2000" (or something like that), you should set the priceMin and
	priceMax to those boundaries: priceMin: 1000 and priceMax: 2000.

	if you don't know about what the price should be, set priceMin to 0 and leave priceMax as null.

	If the size is given in some natural language, like "a mansion" or "a big house" or "a small apartment"
	or anything that could resemble sizes, try to fit it into a range of sizes that would make sense
	considering the current house market.

	A good range of size, probably:
	- small: 0-50
	- medium: 51-300
	- big: +300

	In the case that the characteristic of the property falls into the "big" category, which doesn't have a maximum size,
	only minimum size (sizeMin: 300), we should set sizeMax to null. Other categories for size of property should
	follow the sizeMin and sizeMax from the specified values above (small, medium).

	if you don't know about what the size should be (because you don't know which size characteristic the house should have),
	set sizeMin to 0 and leave sizeMax as null.

	If the color can be any (or is not specified), set it to null.

	If no views specified, just leave it an empty array like this: [].

	You must respond **only** in JSON format. Do not include explanations, greetings, or extra text.
	Your response must be valid JSON. Go straight to the answer. Do NOT hallucinate.

	Example of input:
	User: "I want a big house, close to the sea and to the mountains. Not very expensive. Maybe marble colored"

	Your response (a valid JSON, and nothing more than it):

	{
		"sizeMin": 300, // big house
		"sizeMax": null, // big house has no max limit for size
		"priceMin": 0, // not very expensive = cheap category
		"priceMax": 100000, // cheap category price max
		"views": ["sea", "mountains"],
		"color": "marble"
	}
`

// cleanLLMResponse removes <think>...</think> tags
func cleanLLMResponse(text string) string {
	// Remove <think>...</think> content
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	cleanedText := re.ReplaceAllString(text, "")
	// Remove ```json```
	cleanedText = strings.ReplaceAll(cleanedText, "```json", "")
	cleanedText = strings.ReplaceAll(cleanedText, "```", "")

	// Trim any extra spaces or newlines
	return strings.TrimSpace(cleanedText)
}

type QueryProperty struct {
	Color    string   `json:"color"`
	SizeMin  float32  `json:"sizeMin"`
	SizeMax  float32  `json:"sizeMax"`
	PriceMin float64  `json:"priceMin"`
	PriceMax float64  `json:"priceMax"`
	Views    []string `json:"views"`
}

var DEFAULT_QUERY = QueryProperty{
	Color:    "",
	SizeMin:  0,
	SizeMax:  math.MaxFloat32,
	PriceMin: 0,
	PriceMax: math.MaxFloat64,
	Views:    []string{},
}

// TODO: use streams...simple call too slow
func extract(c echo.Context) error {
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

	prompt := fmt.Sprintf(extractPrompt, body.Query)

	ctx := context.Background()
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		log.Fatal(err)
	}

	limit := false

	// Extract JSON response
	cleanedResponse := cleanLLMResponse(completion)
	var result QueryProperty
	err = json.Unmarshal([]byte(cleanedResponse), &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		// Set the query to query all data
		result = DEFAULT_QUERY
		limit = true
	}

	// Query DB with what we found
	properties := db.QueryByCharacteristics(result.Color, result.PriceMin, result.PriceMax, result.SizeMin, result.SizeMax, limit)

	c.Response().Header().Set("Content-Type", "application/json")
	return c.JSON(http.StatusOK, properties)
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
