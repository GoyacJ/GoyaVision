package engine

import (
	"fmt"
	"strings"
)

// RenderPromptTemplate renders a prompt template by replacing variables.
// Supported: {{asset_id}}, {{params.key}}, {{asset.path}}, {{asset.type}}
func RenderPromptTemplate(tmpl string, vars map[string]interface{}) string {
	if tmpl == "" {
		return ""
	}

	result := tmpl
	for key, val := range vars {
		placeholder := "{{" + key + "}}"
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", val))
	}

	if params, ok := vars["params"]; ok {
		if paramsMap, ok := params.(map[string]interface{}); ok {
			for k, v := range paramsMap {
				placeholder := "{{params." + k + "}}"
				result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", v))
			}
		}
	}

	if asset, ok := vars["asset"]; ok {
		if assetMap, ok := asset.(map[string]interface{}); ok {
			for k, v := range assetMap {
				placeholder := "{{asset." + k + "}}"
				result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", v))
			}
		}
	}

	return result
}

// BuildTemplateVars builds the template variable map from operator input.
func BuildTemplateVars(assetID string, params map[string]interface{}, assetInfo map[string]interface{}) map[string]interface{} {
	vars := map[string]interface{}{
		"asset_id": assetID,
	}
	if params != nil {
		vars["params"] = params
		for k, v := range params {
			vars["params."+k] = v
		}
	}
	if assetInfo != nil {
		vars["asset"] = assetInfo
		for k, v := range assetInfo {
			vars["asset."+k] = v
		}
	}
	return vars
}
