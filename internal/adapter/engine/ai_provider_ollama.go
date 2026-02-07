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

// OllamaProvider handles Ollama API calls.
type OllamaProvider struct {
	client *http.Client
}

func NewOllamaProvider() *OllamaProvider {
	return &OllamaProvider{
		client: &http.Client{Timeout: 10 * time.Minute},
	}
}

type ollamaRequest struct {
	Model    string         `json:"model"`
	Messages []ChatMessage  `json:"messages"`
	Stream   bool           `json:"stream"`
	Options  *ollamaOptions `json:"options,omitempty"`
	Format   string         `json:"format,omitempty"`
}

type ollamaOptions struct {
	Temperature *float64 `json:"temperature,omitempty"`
	NumPredict  *int     `json:"num_predict,omitempty"`
	TopP        *float64 `json:"top_p,omitempty"`
}

type ollamaResponse struct {
	Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	} `json:"message"`
	Model              string `json:"model"`
	PromptEvalCount    int    `json:"prompt_eval_count"`
	EvalCount          int    `json:"eval_count"`
}

func (p *OllamaProvider) Chat(ctx context.Context, model *ai_model.AIModel, req *ChatRequest) (*ChatResponse, error) {
	endpoint := strings.TrimRight(model.Endpoint, "/") + "/api/chat"

	body := ollamaRequest{
		Model:    model.ModelName,
		Messages: req.Messages,
		Stream:   false,
	}

	if req.Temperature != nil || req.MaxTokens != nil || req.TopP != nil {
		body.Options = &ollamaOptions{
			Temperature: req.Temperature,
			NumPredict:  req.MaxTokens,
			TopP:        req.TopP,
		}
	}

	if req.ResponseFormat == "json" {
		body.Format = "json"
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

	var result ollamaResponse
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &ChatResponse{
		Content:    result.Message.Content,
		TokensUsed: result.PromptEvalCount + result.EvalCount,
		Model:      result.Model,
	}, nil
}

func (p *OllamaProvider) HealthCheck(ctx context.Context, model *ai_model.AIModel) error {
	endpoint := strings.TrimRight(model.Endpoint, "/") + "/api/tags"

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
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
