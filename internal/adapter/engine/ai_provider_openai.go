package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"goyavision/internal/domain/ai_model"
)

// OpenAICompatProvider handles OpenAI-compatible APIs.
// Covers OpenAI, Qwen (DashScope), Doubao (Volcano Engine), Zhipu (GLM), vLLM, and local/custom providers.
type OpenAICompatProvider struct {
	chatPath   string
	healthPath string
	client     *http.Client
}

func NewOpenAICompatProvider(chatPath, healthPath string, timeout time.Duration) *OpenAICompatProvider {
	return &OpenAICompatProvider{
		chatPath:   chatPath,
		healthPath: healthPath,
		client:     &http.Client{Timeout: timeout},
	}
}

func NewOpenAIProvider() *OpenAICompatProvider {
	return NewOpenAICompatProvider("/v1/chat/completions", "/v1/models", 5*time.Minute)
}

func NewQwenProvider() *OpenAICompatProvider {
	return NewOpenAICompatProvider("/chat/completions", "/models", 5*time.Minute)
}

func NewDoubaoProvider() *OpenAICompatProvider {
	return NewOpenAICompatProvider("/chat/completions", "/models", 5*time.Minute)
}

func NewZhipuProvider() *OpenAICompatProvider {
	return NewOpenAICompatProvider("/chat/completions", "/models", 5*time.Minute)
}

func NewVLLMProvider() *OpenAICompatProvider {
	return NewOpenAICompatProvider("/v1/chat/completions", "/v1/models", 10*time.Minute)
}

type openAIRequest struct {
	Model          string         `json:"model"`
	Messages       []ChatMessage  `json:"messages"`
	Temperature    *float64       `json:"temperature,omitempty"`
	MaxTokens      *int           `json:"max_tokens,omitempty"`
	TopP           *float64       `json:"top_p,omitempty"`
	ResponseFormat *openAIRespFmt `json:"response_format,omitempty"`
}

type openAIRespFmt struct {
	Type string `json:"type"`
}

type openAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
	Model string `json:"model"`
}

func (p *OpenAICompatProvider) Chat(ctx context.Context, model *ai_model.AIModel, req *ChatRequest) (*ChatResponse, error) {
	endpoint := strings.TrimRight(model.Endpoint, "/") + p.chatPath

	body := openAIRequest{
		Model:       model.ModelName,
		Messages:    req.Messages,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		TopP:        req.TopP,
	}
	if req.ResponseFormat == "json" {
		body.ResponseFormat = &openAIRespFmt{Type: "json_object"}
	}

	data, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")
	if model.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+model.APIKey)
	}

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	var result openAIResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no choices in response")
	}

	return &ChatResponse{
		Content:    result.Choices[0].Message.Content,
		TokensUsed: result.Usage.TotalTokens,
		Model:      result.Model,
	}, nil
}

func (p *OpenAICompatProvider) HealthCheck(ctx context.Context, model *ai_model.AIModel) error {
	endpoint := strings.TrimRight(model.Endpoint, "/") + p.healthPath

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	if model.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+model.APIKey)
	}

	resp, err := p.client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("health check returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}
