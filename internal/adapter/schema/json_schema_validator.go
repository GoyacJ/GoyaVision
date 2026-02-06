package schema

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"goyavision/pkg/apperr"

	jsonschema "github.com/santhosh-tekuri/jsonschema/v5"
)

type JSONSchemaValidator struct{}

func NewJSONSchemaValidator() *JSONSchemaValidator {
	return &JSONSchemaValidator{}
}

func (v *JSONSchemaValidator) IsValidJSONSchema(_ context.Context, schema map[string]interface{}) error {
	if schema == nil {
		return apperr.InvalidInput("schema is required")
	}

	if _, err := compileSchema(schema); err != nil {
		return err
	}

	return nil
}

func (v *JSONSchemaValidator) ValidateInput(ctx context.Context, schema map[string]interface{}, input map[string]interface{}) error {
	compiled, err := compileSchema(schema)
	if err != nil {
		return err
	}
	if input == nil {
		input = map[string]interface{}{}
	}
	if err := compiled.Validate(input); err != nil {
		return apperr.Wrap(err, apperr.CodeInvalidInput, "input does not match schema")
	}
	return validateRequiredFields(schema, input, "input")
}

func (v *JSONSchemaValidator) ValidateOutput(ctx context.Context, schema map[string]interface{}, output map[string]interface{}) error {
	compiled, err := compileSchema(schema)
	if err != nil {
		return err
	}
	if output == nil {
		output = map[string]interface{}{}
	}
	if err := compiled.Validate(output); err != nil {
		return apperr.Wrap(err, apperr.CodeInvalidInput, "output does not match schema")
	}
	return validateRequiredFields(schema, output, "output")
}

func (v *JSONSchemaValidator) ValidateConnection(ctx context.Context, upstreamOutputSpec map[string]interface{}, downstreamInputSchema map[string]interface{}) error {
	if _, err := compileSchema(upstreamOutputSpec); err != nil {
		return apperr.Wrap(err, apperr.CodeInvalidInput, "invalid upstream output spec")
	}
	if _, err := compileSchema(downstreamInputSchema); err != nil {
		return apperr.Wrap(err, apperr.CodeInvalidInput, "invalid downstream input schema")
	}

	upstreamProps := extractProperties(upstreamOutputSpec)
	downstreamProps := extractProperties(downstreamInputSchema)
	requiredFields := extractRequired(downstreamInputSchema)
	for _, field := range requiredFields {
		up, ok := upstreamProps[field]
		if !ok {
			return apperr.InvalidInput(fmt.Sprintf("connection invalid: missing required field '%s'", field))
		}

		down, ok := downstreamProps[field]
		if !ok {
			continue
		}

		if !isPropertyTypeCompatible(up, down) {
			return apperr.InvalidInput(fmt.Sprintf("connection invalid: field '%s' type is incompatible", field))
		}
	}

	return nil
}

func compileSchema(schema map[string]interface{}) (*jsonschema.Schema, error) {
	b, err := json.Marshal(schema)
	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeInvalidInput, "schema is not valid json")
	}

	var root interface{}
	if err := json.Unmarshal(b, &root); err != nil {
		return nil, apperr.Wrap(err, apperr.CodeInvalidInput, "schema is not valid json")
	}
	if _, ok := root.(map[string]interface{}); !ok {
		return nil, apperr.InvalidInput("schema root must be object")
	}

	compiler := jsonschema.NewCompiler()
	const schemaURL = "memory://schema.json"
	if err := compiler.AddResource(schemaURL, bytes.NewReader(b)); err != nil {
		return nil, apperr.Wrap(err, apperr.CodeInvalidInput, "failed to load schema resource")
	}
	compiled, err := compiler.Compile(schemaURL)
	if err != nil {
		return nil, apperr.Wrap(err, apperr.CodeInvalidInput, "invalid json schema")
	}

	return compiled, nil
}

func validateRequiredFields(schema map[string]interface{}, data map[string]interface{}, label string) error {
	if data == nil {
		data = map[string]interface{}{}
	}

	required := extractRequired(schema)
	for _, key := range required {
		if _, ok := data[key]; !ok {
			return apperr.InvalidInput(fmt.Sprintf("%s missing required field '%s'", label, key))
		}
	}
	return nil
}

func extractRequired(schema map[string]interface{}) []string {
	arr, ok := schema["required"].([]interface{})
	if !ok {
		return nil
	}

	result := make([]string, 0, len(arr))
	for _, item := range arr {
		if s, ok := item.(string); ok && s != "" {
			result = append(result, s)
		}
	}
	return result
}

func extractProperties(schema map[string]interface{}) map[string]interface{} {
	props, ok := schema["properties"].(map[string]interface{})
	if !ok {
		return map[string]interface{}{}
	}
	return props
}

func isPropertyTypeCompatible(upstreamProp, downstreamProp interface{}) bool {
	upstreamTypes := extractSchemaTypes(upstreamProp)
	downstreamTypes := extractSchemaTypes(downstreamProp)

	if len(upstreamTypes) == 0 || len(downstreamTypes) == 0 {
		return true
	}

	for _, dt := range downstreamTypes {
		for _, ut := range upstreamTypes {
			if schemaTypeCompatible(ut, dt) {
				return true
			}
		}
	}

	return false
}

func extractSchemaTypes(prop interface{}) []string {
	m, ok := prop.(map[string]interface{})
	if !ok {
		return nil
	}

	v, ok := m["type"]
	if !ok {
		return nil
	}

	switch t := v.(type) {
	case string:
		if t == "" {
			return nil
		}
		return []string{t}
	case []interface{}:
		res := make([]string, 0, len(t))
		for i := range t {
			ts, ok := t[i].(string)
			if ok && ts != "" {
				res = append(res, ts)
			}
		}
		return res
	default:
		return nil
	}
}

func schemaTypeCompatible(upstream, downstream string) bool {
	if upstream == downstream {
		return true
	}
	// JSON Schema 中 integer 是 number 的子类型
	if upstream == "integer" && downstream == "number" {
		return true
	}
	return false
}
