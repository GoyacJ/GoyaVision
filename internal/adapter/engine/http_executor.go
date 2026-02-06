package engine

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"goyavision/internal/domain/operator"
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

// Execute 执行算子版本
func (e *HTTPOperatorExecutor) Execute(ctx context.Context, version *operator.OperatorVersion, input *operator.Input) (*operator.Output, error) {
	if version == nil {
		return nil, fmt.Errorf("operator version is nil")
	}

	if input == nil {
		input = &operator.Input{}
	}

	if version.ExecMode != operator.ExecModeHTTP {
		return nil, fmt.Errorf("http executor does not support exec mode: %s", version.ExecMode)
	}

	if version.ExecConfig == nil || version.ExecConfig.HTTP == nil {
		return nil, fmt.Errorf("http exec config is required")
	}

	httpCfg := version.ExecConfig.HTTP
	if httpCfg.Endpoint == "" {
		return nil, fmt.Errorf("http endpoint is required")
	}

	method := httpCfg.Method
	if method == "" {
		method = http.MethodPost
	}

	requestBody, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal input: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, method, httpCfg.Endpoint, bytes.NewReader(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	for k, v := range httpCfg.Headers {
		req.Header.Set(k, v)
	}

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

	var output operator.Output
	if err := json.Unmarshal(body, &output); err != nil {
		return nil, fmt.Errorf("failed to unmarshal output: %w", err)
	}

	return &output, nil
}

// Mode 返回执行模式
func (e *HTTPOperatorExecutor) Mode() operator.ExecMode {
	return operator.ExecModeHTTP
}

// HealthCheck 检查执行器配置有效性
func (e *HTTPOperatorExecutor) HealthCheck(ctx context.Context, version *operator.OperatorVersion) error {
	if version == nil {
		return fmt.Errorf("operator version is nil")
	}
	if version.ExecConfig == nil || version.ExecConfig.HTTP == nil {
		return fmt.Errorf("http exec config is required")
	}
	if version.ExecConfig.HTTP.Endpoint == "" {
		return fmt.Errorf("http endpoint is required")
	}
	return nil
}
