package dto

import (
	"time"

	"goyavision/internal/domain"

	"github.com/google/uuid"
)

// OperatorListQuery 列出算子查询参数
type OperatorListQuery struct {
	Category  *string `query:"category"`
	Type      *string `query:"type"`
	Status    *string `query:"status"`
	IsBuiltin *bool   `query:"is_builtin"`
	Tags      *string `query:"tags"`
	Keyword   *string `query:"keyword"`
	Limit     int     `query:"limit"`
	Offset    int     `query:"offset"`
}

// OperatorCreateReq 创建算子请求
type OperatorCreateReq struct {
	Code        string                 `json:"code" validate:"required"`
	Name        string                 `json:"name" validate:"required"`
	Description string                 `json:"description,omitempty"`
	Category    string                 `json:"category" validate:"required"`
	Type        string                 `json:"type" validate:"required"`
	Version     string                 `json:"version,omitempty"`
	Endpoint    string                 `json:"endpoint" validate:"required"`
	Method      string                 `json:"method,omitempty"`
	InputSchema map[string]interface{} `json:"input_schema,omitempty"`
	OutputSpec  map[string]interface{} `json:"output_spec,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Status      string                 `json:"status,omitempty"`
	IsBuiltin   bool                   `json:"is_builtin,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
}

// OperatorUpdateReq 更新算子请求
type OperatorUpdateReq struct {
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	Endpoint    *string                `json:"endpoint,omitempty"`
	Method      *string                `json:"method,omitempty"`
	InputSchema map[string]interface{} `json:"input_schema,omitempty"`
	OutputSpec  map[string]interface{} `json:"output_spec,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Status      *string                `json:"status,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
}

// OperatorResponse 算子响应
type OperatorResponse struct {
	ID          uuid.UUID              `json:"id"`
	Code        string                 `json:"code"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Category    string                 `json:"category"`
	Type        string                 `json:"type"`
	Version     string                 `json:"version"`
	Endpoint    string                 `json:"endpoint"`
	Method      string                 `json:"method"`
	InputSchema map[string]interface{} `json:"input_schema,omitempty"`
	OutputSpec  map[string]interface{} `json:"output_spec,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Status      string                 `json:"status"`
	IsBuiltin   bool                   `json:"is_builtin"`
	Tags        []string               `json:"tags,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// OperatorListResponse 算子列表响应
type OperatorListResponse struct {
	Items []*OperatorResponse `json:"items"`
	Total int64               `json:"total"`
}

// OperatorToResponse 转换为响应
func OperatorToResponse(o *domain.Operator) *OperatorResponse {
	if o == nil {
		return nil
	}

	var inputSchema map[string]interface{}
	if o.InputSchema != nil {
		if err := o.InputSchema.Unmarshal(&inputSchema); err == nil {
		}
	}

	var outputSpec map[string]interface{}
	if o.OutputSpec != nil {
		if err := o.OutputSpec.Unmarshal(&outputSpec); err == nil {
		}
	}

	var config map[string]interface{}
	if o.Config != nil {
		if err := o.Config.Unmarshal(&config); err == nil {
		}
	}

	var tags []string
	if o.Tags != nil {
		if err := o.Tags.Unmarshal(&tags); err == nil {
		}
	}

	return &OperatorResponse{
		ID:          o.ID,
		Code:        o.Code,
		Name:        o.Name,
		Description: o.Description,
		Category:    string(o.Category),
		Type:        string(o.Type),
		Version:     o.Version,
		Endpoint:    o.Endpoint,
		Method:      o.Method,
		InputSchema: inputSchema,
		OutputSpec:  outputSpec,
		Config:      config,
		Status:      string(o.Status),
		IsBuiltin:   o.IsBuiltin,
		Tags:        tags,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
}

// OperatorsToResponse 转换为响应列表
func OperatorsToResponse(operators []*domain.Operator) []*OperatorResponse {
	result := make([]*OperatorResponse, len(operators))
	for i, o := range operators {
		result[i] = OperatorToResponse(o)
	}
	return result
}
