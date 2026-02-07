package command

import (
	"context"
	"errors"
	"time"

	"goyavision/internal/adapter/engine"
	"goyavision/internal/app/dto"
	appport "goyavision/internal/app/port"
	"goyavision/internal/domain/ai_model"
	"goyavision/internal/port"
	"goyavision/pkg/apperr"

	"gorm.io/gorm"
)

type TestAIModelHandler struct {
	repo   port.Repository
	crypto appport.CryptoService
}

func NewTestAIModelHandler(repo port.Repository, crypto appport.CryptoService) *TestAIModelHandler {
	return &TestAIModelHandler{repo: repo, crypto: crypto}
}

func (h *TestAIModelHandler) Handle(ctx context.Context, cmd dto.TestAIModelCommand) (*dto.TestAIModelResult, error) {
	model, err := h.repo.GetAIModel(ctx, cmd.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("ai model", cmd.ID.String())
		}
		return nil, apperr.Wrap(err, apperr.CodeDBError, "failed to get ai model")
	}

	if model.APIKey != "" && h.crypto != nil {
		decrypted, decErr := h.crypto.Decrypt(model.APIKey)
		if decErr == nil {
			model.APIKey = decrypted
		}
	}

	var provider engine.AIProvider
	switch model.Provider {
	case ai_model.ProviderOpenAI, ai_model.ProviderLocal, ai_model.ProviderCustom:
		provider = engine.NewOpenAIProvider()
	case ai_model.ProviderAnthropic:
		provider = engine.NewAnthropicProvider()
	case ai_model.ProviderOllama:
		provider = engine.NewOllamaProvider()
	default:
		return &dto.TestAIModelResult{
			Success: false,
			Message: "unsupported provider: " + string(model.Provider),
		}, nil
	}

	testCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	if err := provider.HealthCheck(testCtx, model); err != nil {
		return &dto.TestAIModelResult{
			Success: false,
			Message: "connection failed: " + err.Error(),
		}, nil
	}

	return &dto.TestAIModelResult{
		Success: true,
		Message: "connection successful",
	}, nil
}
