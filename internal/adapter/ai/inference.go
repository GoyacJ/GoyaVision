package ai

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"goyavision/internal/port"
)

type inferenceAdapter struct {
	client  *http.Client
	timeout time.Duration
	retry   int
}

func NewInferenceAdapter(timeout time.Duration, retry int) port.Inference {
	return &inferenceAdapter{
		client: &http.Client{
			Timeout: timeout,
		},
		timeout: timeout,
		retry:   retry,
	}
}

func (a *inferenceAdapter) Post(ctx context.Context, endpoint string, req port.InferenceRequest) (*port.InferenceResponse, error) {
	var lastErr error
	for i := 0; i <= a.retry; i++ {
		if i > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(time.Duration(i) * time.Second):
			}
		}

		resp, err := a.doPost(ctx, endpoint, req)
		if err == nil {
			return resp, nil
		}
		lastErr = err
	}
	return nil, fmt.Errorf("inference failed after %d retries: %w", a.retry, lastErr)
}

func (a *inferenceAdapter) doPost(ctx context.Context, endpoint string, req port.InferenceRequest) (*port.InferenceResponse, error) {
	httpReq, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewReader(req.Body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := a.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("inference returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	return &port.InferenceResponse{
		Body: body,
	}, nil
}
