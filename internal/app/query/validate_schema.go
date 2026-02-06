package query

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/pkg/apperr"
)

type ValidateSchemaHandler struct {
	validator port.SchemaValidator
}

func NewValidateSchemaHandler(validator port.SchemaValidator) *ValidateSchemaHandler {
	return &ValidateSchemaHandler{validator: validator}
}

func (h *ValidateSchemaHandler) Handle(ctx context.Context, q dto.ValidateSchemaQuery) error {
	if h.validator == nil {
		return apperr.ServiceUnavailable("schema validator is not configured")
	}
	return h.validator.IsValidJSONSchema(ctx, q.Schema)
}
