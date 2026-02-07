package engine

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	appport "goyavision/internal/app/port"
	"goyavision/internal/domain/ai_model"
	"goyavision/internal/domain/operator"
	"goyavision/internal/port"
)

var _ port.OperatorExecutor = (*AIModelExecutor)(nil)

// AIModelExecutor executes operators via AI model inference.
type AIModelExecutor struct {
	repo      port.Repository
	crypto    appport.CryptoService
	providers map[ai_model.Provider]AIProvider
}

// NewAIModelExecutor creates an AI model executor with all supported providers.
func NewAIModelExecutor(repo port.Repository, crypto appport.CryptoService) *AIModelExecutor {
	openai := NewOpenAIProvider()
	return &AIModelExecutor{
		repo:   repo,
		crypto: crypto,
		providers: map[ai_model.Provider]AIProvider{
			ai_model.ProviderOpenAI:    openai,
			ai_model.ProviderAnthropic: NewAnthropicProvider(),
			ai_model.ProviderOllama:    NewOllamaProvider(),
			ai_model.ProviderLocal:     openai,
			ai_model.ProviderCustom:    openai,
			ai_model.ProviderQwen:      NewQwenProvider(),
			ai_model.ProviderDoubao:    NewDoubaoProvider(),
			ai_model.ProviderZhipu:     NewZhipuProvider(),
			ai_model.ProviderVLLM:      NewVLLMProvider(),
		},
	}
}

func (e *AIModelExecutor) Execute(ctx context.Context, version *operator.OperatorVersion, input *operator.Input) (*operator.Output, error) {
	if version == nil {
		return nil, fmt.Errorf("operator version is nil")
	}
	if version.ExecMode != operator.ExecModeAIModel {
		return nil, fmt.Errorf("ai_model executor does not support exec mode: %s", version.ExecMode)
	}
	if version.ExecConfig == nil || version.ExecConfig.AIModel == nil {
		return nil, fmt.Errorf("ai_model exec config is required")
	}

	cfg := version.ExecConfig.AIModel

	model, err := e.repo.GetAIModel(ctx, cfg.ModelID)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve AI model: %w", err)
	}
	if !model.IsActive() {
		return nil, fmt.Errorf("AI model %s is disabled", model.Name)
	}

	if model.APIKey != "" && e.crypto != nil {
		if decrypted, decErr := e.crypto.Decrypt(model.APIKey); decErr == nil {
			model.APIKey = decrypted
		}
	}

	provider, ok := e.providers[model.Provider]
	if !ok {
		return nil, fmt.Errorf("unsupported AI provider: %s", model.Provider)
	}

	vars := BuildTemplateVars(input.AssetID.String(), input.Params, nil)
	systemPrompt := RenderPromptTemplate(cfg.SystemPrompt, vars)
	userPrompt := RenderPromptTemplate(cfg.UserPromptTemplate, vars)

	var messages []ChatMessage
	if systemPrompt != "" {
		messages = append(messages, ChatMessage{Role: "system", Content: systemPrompt})
	}

	if cfg.InteractionMode == "vision" {
		content := []ContentPart{
			{Type: "text", Text: userPrompt},
		}
		if assetPath, ok := input.Params["image_url"]; ok {
			content = append(content, ContentPart{
				Type:     "image_url",
				ImageURL: &ImageURL{URL: fmt.Sprintf("%v", assetPath)},
			})
		}
		messages = append(messages, ChatMessage{Role: "user", Content: content})
	} else {
		messages = append(messages, ChatMessage{Role: "user", Content: userPrompt})
	}

	if cfg.TimeoutSec > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Duration(cfg.TimeoutSec)*time.Second)
		defer cancel()
	}

	chatReq := &ChatRequest{
		Messages:       messages,
		Temperature:    cfg.Temperature,
		MaxTokens:      cfg.MaxTokens,
		TopP:           cfg.TopP,
		ResponseFormat: cfg.ResponseFormat,
	}

	chatResp, err := provider.Chat(ctx, model, chatReq)
	if err != nil {
		return nil, fmt.Errorf("AI model execution failed: %w", err)
	}

	output := e.mapResponse(chatResp, cfg.OutputMapping)
	if output.Diagnostics == nil {
		output.Diagnostics = make(map[string]interface{})
	}
	output.Diagnostics["tokens_used"] = chatResp.TokensUsed
	output.Diagnostics["model"] = chatResp.Model
	output.Diagnostics["provider"] = string(model.Provider)

	return output, nil
}

func (e *AIModelExecutor) Mode() operator.ExecMode {
	return operator.ExecModeAIModel
}

func (e *AIModelExecutor) HealthCheck(ctx context.Context, version *operator.OperatorVersion) error {
	if version == nil {
		return fmt.Errorf("operator version is nil")
	}
	if version.ExecConfig == nil || version.ExecConfig.AIModel == nil {
		return fmt.Errorf("ai_model exec config is required")
	}

	model, err := e.repo.GetAIModel(ctx, version.ExecConfig.AIModel.ModelID)
	if err != nil {
		return fmt.Errorf("failed to resolve AI model: %w", err)
	}

	if model.APIKey != "" && e.crypto != nil {
		if decrypted, decErr := e.crypto.Decrypt(model.APIKey); decErr == nil {
			model.APIKey = decrypted
		}
	}

	provider, ok := e.providers[model.Provider]
	if !ok {
		return fmt.Errorf("unsupported AI provider: %s", model.Provider)
	}

	return provider.HealthCheck(ctx, model)
}

func (e *AIModelExecutor) mapResponse(resp *ChatResponse, mapping map[string]interface{}) *operator.Output {
	output := &operator.Output{}

	if mapping != nil && len(mapping) > 0 {
		var parsed map[string]interface{}
		if err := json.Unmarshal([]byte(resp.Content), &parsed); err == nil {
			output.Results = []operator.Result{
				{
					Type: "ai_response",
					Data: parsed,
				},
			}
			return output
		}
	}

	output.Results = []operator.Result{
		{
			Type: "ai_response",
			Data: map[string]interface{}{
				"content": resp.Content,
			},
		},
	}
	return output
}
