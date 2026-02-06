package port

import "context"

type SchemaValidator interface {
	IsValidJSONSchema(ctx context.Context, schema map[string]interface{}) error
	ValidateInput(ctx context.Context, schema map[string]interface{}, input map[string]interface{}) error
	ValidateOutput(ctx context.Context, schema map[string]interface{}, output map[string]interface{}) error
	ValidateConnection(ctx context.Context, upstreamOutputSpec map[string]interface{}, downstreamInputSchema map[string]interface{}) error
}
