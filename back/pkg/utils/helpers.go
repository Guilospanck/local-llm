package utils

import (
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
