package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/pkg/apperr"
)

type ValidateConnectionHandler struct {
	validator port.SchemaValidator
}

func NewValidateConnectionHandler(validator port.SchemaValidator) *ValidateConnectionHandler {
	return &ValidateConnectionHandler{validator: validator}
}

func (h *ValidateConnectionHandler) Handle(ctx context.Context, q dto.ValidateConnectionQuery) error {
	if h.validator == nil {
		return apperr.ServiceUnavailable("schema validator is not configured")
	}
	return h.validator.ValidateConnection(ctx, q.UpstreamOutputSpec, q.DownstreamInputSchema)
}
