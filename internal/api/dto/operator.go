package dto

import (
	"time"

	"goyavision/internal/domain/operator"

	"github.com/google/uuid"
)

// OperatorListQuery 列出算子查询参数
type OperatorListQuery struct {
	Category  *string `query:"category"`
	Type      *string `query:"type"`
	Status    *string `query:"status"`
	Origin    *string `query:"origin"`
	ExecMode  *string `query:"exec_mode"`
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
	Origin      string                 `json:"origin,omitempty"`
	ExecMode    string                 `json:"exec_mode,omitempty"`
	ExecConfig  map[string]interface{} `json:"exec_config,omitempty"`
	Version     string                 `json:"version,omitempty"`
	Endpoint    string                 `json:"endpoint,omitempty"`
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
	Category    *string                `json:"category,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
}

type TestOperatorReq struct {
	AssetID *uuid.UUID             `json:"asset_id,omitempty"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

type TestOperatorResponse struct {
	Success     bool                   `json:"success"`
	Message     string                 `json:"message"`
	Diagnostics map[string]interface{} `json:"diagnostics,omitempty"`
}

type MCPServerResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}

type MCPToolResponse struct {
	Name         string                 `json:"name"`
	Description  string                 `json:"description,omitempty"`
	Version      string                 `json:"version,omitempty"`
	InputSchema  map[string]interface{} `json:"input_schema,omitempty"`
	OutputSchema map[string]interface{} `json:"output_schema,omitempty"`
}

type MCPInstallReq struct {
	ServerID     string   `json:"server_id" validate:"required"`
	ToolName     string   `json:"tool_name" validate:"required"`
	OperatorCode string   `json:"operator_code" validate:"required"`
	OperatorName string   `json:"operator_name" validate:"required"`
	Category     *string  `json:"category,omitempty"`
	Type         *string  `json:"type,omitempty"`
	TimeoutSec   int      `json:"timeout_sec,omitempty"`
	Tags         []string `json:"tags,omitempty"`
}

type SyncMCPTemplatesReq struct {
	ServerID string `json:"server_id" validate:"required"`
}

type SyncMCPTemplatesResponse struct {
	ServerID string `json:"server_id"`
	Total    int    `json:"total"`
	Created  int    `json:"created"`
	Updated  int    `json:"updated"`
}

type OperatorVersionResponse struct {
	ID         uuid.UUID              `json:"id"`
	Version    string                 `json:"version"`
	ExecMode   string                 `json:"exec_mode"`
	ExecConfig map[string]interface{} `json:"exec_config,omitempty"`
	Status     string                 `json:"status"`
}

// OperatorResponse 算子响应
type OperatorResponse struct {
	ID          uuid.UUID              `json:"id"`
	Code        string                 `json:"code"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Category    string                 `json:"category"`
	Type        string                 `json:"type"`
	Origin      string                 `json:"origin,omitempty"`
	ActiveVersionID *uuid.UUID         `json:"active_version_id,omitempty"`
	ExecMode    string                 `json:"exec_mode,omitempty"`
	ActiveVersion *OperatorVersionResponse `json:"active_version,omitempty"`
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
func OperatorToResponse(o *operator.Operator) *OperatorResponse {
	if o == nil {
		return nil
	}

	inputSchema := o.InputSchema
	if inputSchema == nil {
		inputSchema = make(map[string]interface{})
	}

	outputSpec := o.OutputSpec
	if outputSpec == nil {
		outputSpec = make(map[string]interface{})
	}

	config := o.Config
	if config == nil {
		config = make(map[string]interface{})
	}

	tags := o.Tags
	if tags == nil {
		tags = []string{}
	}

	var activeVersion *OperatorVersionResponse
	execMode := ""
	if o.ActiveVersion != nil {
		execMode = string(o.ActiveVersion.ExecMode)
		activeVersion = &OperatorVersionResponse{
			ID:      o.ActiveVersion.ID,
			Version: o.ActiveVersion.Version,
			ExecMode: string(o.ActiveVersion.ExecMode),
			Status:  string(o.ActiveVersion.Status),
		}
	}

	return &OperatorResponse{
		ID:          o.ID,
		Code:        o.Code,
		Name:        o.Name,
		Description: o.Description,
		Category:    string(o.Category),
		Type:        string(o.Type),
		Origin:      string(o.Origin),
		ActiveVersionID: o.ActiveVersionID,
		ExecMode:    execMode,
		ActiveVersion: activeVersion,
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
func OperatorsToResponse(operators []*operator.Operator) []*OperatorResponse {
	result := make([]*OperatorResponse, len(operators))
	for i, o := range operators {
		result[i] = OperatorToResponse(o)
	}
	return result
}
