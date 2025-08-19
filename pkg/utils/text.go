package utils

import (
	"regexp"
	"strings"
)

// parseLanguageNames converts language names to locale codes
func ParseLanguageNames(langStr string) []string {
	// Language name to locale code mapping
	languageMap := map[string]string{
		"english":    "en",
		"spanish":    "es",
		"french":     "fr",
		"german":     "de",
		"italian":    "it",
		"portuguese": "pt",
		"russian":    "ru",
		"chinese":    "zh",
		"japanese":   "ja",
		"korean":     "ko",
		"arabic":     "ar",
		"hindi":      "hi",
		"dutch":      "nl",
		"swedish":    "sv",
		"danish":     "da",
		"norwegian":  "no",
		"finnish":    "fi",
		"polish":     "pl",
		"czech":      "cs",
		"hungarian":  "hu",
		"turkish":    "tr",
		"hebrew":     "he",
		"thai":       "th",
		"vietnamese": "vi",
		"indonesian": "id",
		"malay":      "ms",
		"tagalog":    "tl",
		"swahili":    "sw",
	}

	var locales []string
	langStr = strings.ToLower(langStr)

	// Split by common delimiters
	languages := regexp.MustCompile(`[,;]|\s+and\s+`).Split(langStr, -1)

	for _, lang := range languages {
		lang = strings.TrimSpace(lang)
		lang = strings.Trim(lang, ".,")

		if locale, exists := languageMap[lang]; exists {
			locales = append(locales, locale)
		}
	}

	return locales
}

// generateDescriptionFromModelName creates a readable description from a model name
func GenerateDescriptionFromModelName(modelName string) string {
	if modelName == "" {
		return ""
	}

	// Remove common prefixes
	cleaned := modelName
	cleaned = strings.TrimPrefix(cleaned, "RedHatAI/")
	cleaned = strings.TrimPrefix(cleaned, "meta-llama/")
	cleaned = strings.TrimPrefix(cleaned, "microsoft/")
	cleaned = strings.TrimPrefix(cleaned, "mistralai/")
	cleaned = strings.TrimPrefix(cleaned, "Qwen/")
	cleaned = strings.TrimPrefix(cleaned, "ibm-granite/")

	// Replace hyphens and underscores with spaces
	cleaned = strings.ReplaceAll(cleaned, "-", " ")
	cleaned = strings.ReplaceAll(cleaned, "_", " ")

	// Handle specific model patterns
	cleaned = regexp.MustCompile(`(?i)\bLlama\b`).ReplaceAllString(cleaned, "Llama")
	cleaned = regexp.MustCompile(`(?i)\bMistral\b`).ReplaceAllString(cleaned, "Mistral")
	cleaned = regexp.MustCompile(`(?i)\bGranite\b`).ReplaceAllString(cleaned, "Granite")
	cleaned = regexp.MustCompile(`(?i)\bPhi\b`).ReplaceAllString(cleaned, "Phi")
	cleaned = regexp.MustCompile(`(?i)\bQwen\b`).ReplaceAllString(cleaned, "Qwen")
	cleaned = regexp.MustCompile(`(?i)\bWhisper\b`).ReplaceAllString(cleaned, "Whisper")

	// Handle version patterns (e.g., "3.3" -> "3.3", "70B" -> "70B")
	cleaned = regexp.MustCompile(`\b(\d+)\.(\d+)\b`).ReplaceAllString(cleaned, "$1.$2")
	cleaned = regexp.MustCompile(`\b(\d+)([BM])\b`).ReplaceAllString(cleaned, "$1$2")

	// Handle common suffixes
	cleaned = regexp.MustCompile(`(?i)\bInstruct\b`).ReplaceAllString(cleaned, "Instruct")
	cleaned = regexp.MustCompile(`(?i)\bBase\b`).ReplaceAllString(cleaned, "Base")
	cleaned = regexp.MustCompile(`(?i)\bChat\b`).ReplaceAllString(cleaned, "Chat")
	cleaned = regexp.MustCompile(`(?i)\bCode\b`).ReplaceAllString(cleaned, "Code")

	// Handle quantization suffixes
	cleaned = regexp.MustCompile(`(?i)\bquantized\.(w\d+a\d+)\b`).ReplaceAllString(cleaned, "($1 quantized)")
	cleaned = regexp.MustCompile(`(?i)\bfp8 dynamic\b`).ReplaceAllString(cleaned, "(FP8 dynamic)")
	cleaned = regexp.MustCompile(`(?i)\bfp8\b`).ReplaceAllString(cleaned, "(FP8)")

	// Clean up multiple spaces
	cleaned = regexp.MustCompile(`\s+`).ReplaceAllString(cleaned, " ")
	cleaned = strings.TrimSpace(cleaned)

	// Capitalize first letter
	if len(cleaned) > 0 {
		cleaned = strings.ToUpper(string(cleaned[0])) + cleaned[1:]
	}

	return cleaned
}

// normalizeModelName normalizes model names for comparison
func NormalizeModelName(name string) string {
	// Remove common prefixes and suffixes
	normalized := strings.ToLower(name)

	// Remove registry prefix
	normalized = strings.TrimPrefix(normalized, "registry.redhat.io/rhelai1/modelcar-")

	// Remove RedHatAI prefix
	normalized = strings.TrimPrefix(normalized, "redhatai/")

	// Remove version tags
	if idx := strings.LastIndex(normalized, ":"); idx != -1 {
		normalized = normalized[:idx]
	}

	// Replace various separators with hyphens for consistency
	normalized = strings.ReplaceAll(normalized, "_", "-")
	normalized = strings.ReplaceAll(normalized, ".", "-")
	normalized = strings.ReplaceAll(normalized, " ", "-")

	// Remove duplicate hyphens
	for strings.Contains(normalized, "--") {
		normalized = strings.ReplaceAll(normalized, "--", "-")
	}

	return normalized
}

// calculateSimilarity calculates a simple similarity score between two strings
func CalculateSimilarity(s1, s2 string) float64 {
	s1Norm := NormalizeModelName(s1)
	s2Norm := NormalizeModelName(s2)

	// Exact match
	if s1Norm == s2Norm {
		return 1.0
	}

	// Check if one contains the other
	if strings.Contains(s1Norm, s2Norm) || strings.Contains(s2Norm, s1Norm) {
		return 0.8
	}

	// Count common words/tokens
	s1Tokens := strings.Split(s1Norm, "-")
	s2Tokens := strings.Split(s2Norm, "-")

	commonTokens := 0
	for _, token1 := range s1Tokens {
		if token1 == "" {
			continue
		}
		for _, token2 := range s2Tokens {
			if token1 == token2 {
				commonTokens++
				break
			}
		}
	}

	maxTokens := len(s1Tokens)
	if len(s2Tokens) > maxTokens {
		maxTokens = len(s2Tokens)
	}

	if maxTokens == 0 {
		return 0.0
	}

	return float64(commonTokens) / float64(maxTokens)
}

// GenerateReadableDescription creates a human-readable description from a model name
func GenerateReadableDescription(modelName string) string {
	if modelName == "" {
		return ""
	}

	// Remove common prefixes/suffixes and make it more readable
	cleaned := modelName

	// Remove registry paths
	if strings.Contains(cleaned, "/") {
		parts := strings.Split(cleaned, "/")
		cleaned = parts[len(parts)-1] // Take the last part
	}

	// Remove version tags
	if strings.Contains(cleaned, ":") {
		parts := strings.Split(cleaned, ":")
		cleaned = parts[0]
	}

	// Replace hyphens and underscores with spaces
	cleaned = strings.ReplaceAll(cleaned, "-", " ")
	cleaned = strings.ReplaceAll(cleaned, "_", " ")

	// Remove common prefixes
	prefixes := []string{"modelcar", "rhelai1", "model", "ai"}
	for _, prefix := range prefixes {
		if strings.HasPrefix(strings.ToLower(cleaned), prefix+" ") {
			cleaned = cleaned[len(prefix)+1:]
		}
	}

	// Capitalize words appropriately
	words := strings.Fields(cleaned)
	var result []string

	for _, word := range words {
		if word == "" {
			continue
		}

		// Special handling for known model names and companies
		switch strings.ToLower(word) {
		case "granite", "llama", "mistral", "qwen", "phi", "gemma":
			result = append(result, strings.Title(word))
		case "instruct":
			result = append(result, "Instruct")
		case "base":
			result = append(result, "Base")
		case "chat":
			result = append(result, "Chat")
		case "quantized":
			result = append(result, "Quantized")
		case "ibm":
			result = append(result, "IBM")
		case "microsoft":
			result = append(result, "Microsoft")
		case "meta":
			result = append(result, "Meta")
		case "redhat", "redhatai":
			result = append(result, "Red Hat AI")
		default:
			// Check if it's a version number
			if matched, _ := regexp.MatchString(`^\d+(\.\d+)*[a-z]*$`, word); matched {
				result = append(result, word)
			} else {
				result = append(result, strings.Title(word))
			}
		}
	}

	if len(result) == 0 {
		return ""
	}

	description := strings.Join(result, " ")

	// Add appropriate suffix based on model type
	if strings.Contains(strings.ToLower(description), "instruct") {
		description += " - An instruction-tuned language model"
	} else if strings.Contains(strings.ToLower(description), "chat") {
		description += " - A conversational AI model"
	} else if strings.Contains(strings.ToLower(description), "base") {
		description += " - A foundation language model"
	} else {
		description += " - A large language model"
	}

	return description
}

// NormalizeTask normalizes a task description to standard task categories
func NormalizeTask(task string) string {
	if task == "" {
		return ""
	}

	taskLower := strings.ToLower(strings.TrimSpace(task))

	// Task normalization mapping
	taskMap := map[string]string{
		// Text generation tasks
		"text generation":   "text-generation",
		"text-generation":   "text-generation",
		"language modeling": "text-generation",
		"conversation":      "text-generation",
		"conversational":    "text-generation",
		"chat":              "text-generation",
		"chatbot":           "text-generation",
		"dialogue":          "text-generation",
		"code generation":   "text-generation",
		"coding":            "text-generation",
		"programming":       "text-generation",
		"completion":        "text-generation",
		"writing":           "text-generation",
		"creative writing":  "text-generation",
		"storytelling":      "text-generation",

		// Classification tasks
		"text classification": "text-classification",
		"text-classification": "text-classification",
		"classification":      "text-classification",
		"sentiment analysis":  "text-classification",
		"sentiment":           "text-classification",
		"categorization":      "text-classification",
		"labeling":            "text-classification",

		// Question answering
		"question answering":    "question-answering",
		"question-answering":    "question-answering",
		"qa":                    "question-answering",
		"q&a":                   "question-answering",
		"question and answer":   "question-answering",
		"information retrieval": "question-answering",
		"search":                "question-answering",

		// Image tasks
		"image classification":      "image-classification",
		"image-classification":      "image-classification",
		"image captioning":          "image-to-text",
		"image-to-text":             "image-to-text",
		"image description":         "image-to-text",
		"visual question answering": "image-text-to-text",
		"image-text-to-text":        "image-text-to-text",
		"image-to-image":            "image-to-image",

		// Other specific tasks
		"sentence similarity": "sentence-similarity",
		"sentence-similarity": "sentence-similarity",
		"text ranking":        "text-ranking",
		"text-ranking":        "text-ranking",
		"ranking":             "text-ranking",
		"any-to-any":          "any-to-any",
		"text-to-video":       "text-to-video",
		"video-to-video":      "video-to-video",
	}

	// Check exact matches first
	if normalized, exists := taskMap[taskLower]; exists {
		return normalized
	}

	// Check for partial matches
	for pattern, normalized := range taskMap {
		if strings.Contains(taskLower, pattern) {
			return normalized
		}
	}

	// If no match found, try to infer from keywords
	if strings.Contains(taskLower, "question") && strings.Contains(taskLower, "answer") {
		return "question-answering"
	}
	if strings.Contains(taskLower, "image") && strings.Contains(taskLower, "text") {
		return "image-text-to-text"
	}
	if strings.Contains(taskLower, "image") {
		return "image-classification"
	}
	if strings.Contains(taskLower, "generat") {
		return "text-generation"
	}
	if strings.Contains(taskLower, "classif") {
		return "text-classification"
	}

	// Return original if no normalization possible
	return task
}
