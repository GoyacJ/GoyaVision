package port

import "context"

type InferenceRequest struct {
	Body []byte
}

type InferenceResponse struct {
	Body []byte
}

type Inference interface {
	Post(ctx context.Context, endpoint string, req InferenceRequest) (*InferenceResponse, error)
}
