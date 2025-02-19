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
