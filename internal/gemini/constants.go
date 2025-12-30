package gemini

import "google.golang.org/genai"

// Model defaults
const (
	DefaultModelID      = "gemini-2.5-pro"
	DefaultFlashModelID = "gemini-2.5-flash"
)

// Generation parameter defaults
const (
	DefaultModelTemperature = float32(0.7)
	DefaultModelTopP        = float32(0.95)
	DefaultModelTopK        = 64
)

// Buffer sizes for streaming channels
const (
	DefaultHistoryCapacity  = 20
	StreamChannelBufferSize = 20
	ChunkChannelBufferSize  = 10
	ErrorChannelBufferSize  = 1
	DoneChannelBufferSize   = 1
)

// Safety threshold defaults
const DefaultSafetyThreshold = "HIGH"

// mapSafetyThreshold maps a string threshold to genai.HarmBlockThreshold
// Supported values: "HIGH", "MEDIUM_AND_ABOVE", "LOW_AND_ABOVE", "NONE", "OFF"
func mapSafetyThreshold(threshold string) genai.HarmBlockThreshold {
	switch threshold {
	case "HIGH":
		return genai.HarmBlockThresholdBlockOnlyHigh
	case "MEDIUM_AND_ABOVE":
		return genai.HarmBlockThresholdBlockMediumAndAbove
	case "LOW_AND_ABOVE":
		return genai.HarmBlockThresholdBlockLowAndAbove
	case "NONE":
		return genai.HarmBlockThresholdBlockNone
	case "OFF":
		return genai.HarmBlockThresholdOff
	default:
		return genai.HarmBlockThresholdBlockOnlyHigh
	}
}

// getAllHarmCategories returns all harm categories for safety settings
func getAllHarmCategories() []genai.HarmCategory {
	return []genai.HarmCategory{
		genai.HarmCategoryHarassment,
		genai.HarmCategoryHateSpeech,
		genai.HarmCategorySexuallyExplicit,
		genai.HarmCategoryDangerousContent,
		genai.HarmCategoryCivicIntegrity,
	}
}
