package engine

import (
	"context"

	"goyavision/internal/domain/ai_model"
)

// ChatMessage represents a single message in an AI conversation.
type ChatMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

// ContentPart represents a multimodal content part (text or image).
type ContentPart struct {
	Type     string    `json:"type"`
	Text     string    `json:"text,omitempty"`
	ImageURL *ImageURL `json:"image_url,omitempty"`
}

// ImageURL wraps an image URL for vision requests.
type ImageURL struct {
	URL string `json:"url"`
}

// ChatRequest is the unified request for all providers.
type ChatRequest struct {
	Messages       []ChatMessage
	Temperature    *float64
	MaxTokens      *int
	TopP           *float64
	ResponseFormat string
}

// ChatResponse is the unified response from all providers.
type ChatResponse struct {
	Content    string
	TokensUsed int
	Model      string
}

// AIProvider defines the interface for AI model providers.
type AIProvider interface {
	Chat(ctx context.Context, model *ai_model.AIModel, req *ChatRequest) (*ChatResponse, error)
	HealthCheck(ctx context.Context, model *ai_model.AIModel) error
}
