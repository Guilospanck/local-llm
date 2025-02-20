package handlers

import (
	"base/pkg/domain"
	"base/pkg/types"
	"base/pkg/utils"
	"context"
	"fmt"
	"net/http"
	"os"

	"encoding/json"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/ollama"
)

// TODO: the call is too slow. Find a better way.
func Extract(c echo.Context, db *domain.Database, modelName *string) error {
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

	extractPrompt := utils.GetPromptBasedOnModel(utils.Model(*modelName))

	prompt := fmt.Sprintf(extractPrompt, body.Query)

	ctx := context.Background()
	completion, err := llms.GenerateFromSinglePrompt(ctx, llm, prompt)
	if err != nil {
		log.Fatal(err)
	}

	limit := false

	fmt.Println()
	fmt.Println(completion)

	// Extract JSON response
	cleanedResponse := utils.CleanLLMResponse(completion)
	var result types.QueryProperty
	err = json.Unmarshal([]byte(cleanedResponse), &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		fmt.Println("Using default query...")
		// Set the query to query all data
		result = utils.DEFAULT_QUERY
		limit = true
	}

	// Query DB with what we found
	properties := db.QueryByCharacteristics(result.Color, result.PriceMin, result.PriceMax, result.SizeMin, result.SizeMax, result.Views, limit)

	c.Response().Header().Set("Content-Type", "application/json")
	return c.JSON(http.StatusOK, properties)
}
