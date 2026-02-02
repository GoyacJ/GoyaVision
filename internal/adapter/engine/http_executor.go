package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"goyavision/internal/domain"
	"goyavision/internal/port"
)

var _ port.OperatorExecutor = (*HTTPOperatorExecutor)(nil)

// HTTPOperatorExecutor HTTP 算子执行器
type HTTPOperatorExecutor struct {
	client *http.Client
}

// NewHTTPOperatorExecutor 创建 HTTP 算子执行器
func NewHTTPOperatorExecutor() *HTTPOperatorExecutor {
	return &HTTPOperatorExecutor{
		client: &http.Client{
			Timeout: 5 * time.Minute,
		},
	}
}

// Execute 执行算子
func (e *HTTPOperatorExecutor) Execute(ctx context.Context, operator *domain.Operator, input *domain.OperatorInput) (*domain.OperatorOutput, error) {
	if operator == nil {
		return nil, fmt.Errorf("operator is nil")
	}

	if input == nil {
		input = &domain.OperatorInput{}
	}

	requestBody, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, operator.Method, operator.Endpoint, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("operator returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var output domain.OperatorOutput
	if err := json.Unmarshal(body, &output); err != nil {
		return nil, fmt.Errorf("failed to unmarshal output: %w", err)
	}

	return &output, nil
}
