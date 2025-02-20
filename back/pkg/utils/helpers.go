package utils

import (
	"fmt"
	"regexp"
	"strings"
)

// Helper function to remove <think>...</think> tags
func CleanLLMResponse(text string) string {
	// Remove <think>...</think> content
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	cleanedText := re.ReplaceAllString(text, "")
	// Remove ```json```
	cleanedText = strings.ReplaceAll(cleanedText, "```json", "")
	cleanedText = strings.ReplaceAll(cleanedText, "```", "")

	// Trim any extra spaces or newlines
	return strings.TrimSpace(cleanedText)
}

// Helper function to join the elements of a slice
func JoinSlice[T []E, E any](slice T, separator string) string {
	strArray := make([]string, len(slice))
	for i, element := range slice {
		strArray[i] = fmt.Sprint(element)
	}

	return strings.Join(strArray, separator)
}

func GetPromptBasedOnModel(modelName Model) string {

	switch modelName {
	case DEEPSEEK_R1_1_5B:
		println("Using extract for Deepseek-r1:1.5b")
		return ExtractPromptDeepseek
	case GEMMA_2B:
		println("Using extract for Gemma:2b")
		return ExtractPromptGemma2b
	case LLAMA_3_2:
		println("Using extract for Llama3.2")
		return ExtractPromptGemma2b
	default:
		println("Using default extract")
		return ExtractPromptDeepseek
	}

}
