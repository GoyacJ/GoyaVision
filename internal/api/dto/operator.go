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
	Version     string                 `json:"version,omitempty"`      // Deprecated: 兼容字段，优先通过版本 API 维护
	Endpoint    string                 `json:"endpoint,omitempty"`     // Deprecated: 兼容字段，建议改用 exec_config.http.endpoint
	Method      string                 `json:"method,omitempty"`       // Deprecated: 兼容字段，建议改用 exec_config.http.method
	InputSchema map[string]interface{} `json:"input_schema,omitempty"` // Deprecated: 兼容字段，建议改用版本 API
	OutputSpec  map[string]interface{} `json:"output_spec,omitempty"`  // Deprecated: 兼容字段，建议改用版本 API
	Config      map[string]interface{} `json:"config,omitempty"`       // Deprecated: 兼容字段，建议改用版本 API
	Status      string                 `json:"status,omitempty"`
	IsBuiltin   bool                   `json:"is_builtin,omitempty"` // Deprecated: 兼容字段，建议改用 origin
	Tags        []string               `json:"tags,omitempty"`
}

// OperatorUpdateReq 更新算子请求
type OperatorUpdateReq struct {
	Name        *string  `json:"name,omitempty"`
	Description *string  `json:"description,omitempty"`
	Category    *string  `json:"category,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

type TestOperatorReq struct {
	AssetID *uuid.UUID             `json:"asset_id,omitempty"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

type OperatorVersionListQuery struct {
	Limit  int `query:"limit"`
	Offset int `query:"offset"`
}

type OperatorVersionCreateReq struct {
	Version     string                 `json:"version" validate:"required"`
	ExecMode    string                 `json:"exec_mode" validate:"required"`
	ExecConfig  map[string]interface{} `json:"exec_config,omitempty"`
	InputSchema map[string]interface{} `json:"input_schema,omitempty"`
	OutputSpec  map[string]interface{} `json:"output_spec,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Changelog   string                 `json:"changelog,omitempty"`
	Status      string                 `json:"status,omitempty"`
}

type OperatorVersionActionReq struct {
	VersionID uuid.UUID `json:"version_id" validate:"required"`
}

type OperatorVersionListResponse struct {
	Items []*OperatorVersionResponse `json:"items"`
	Total int64                      `json:"total"`
}

type ValidateSchemaReq struct {
	Schema map[string]interface{} `json:"schema"`
}

type ValidateConnectionReq struct {
	UpstreamOutputSpec    map[string]interface{} `json:"upstream_output_spec"`
	DownstreamInputSchema map[string]interface{} `json:"downstream_input_schema"`
}

type ValidateResultResponse struct {
	Valid   bool   `json:"valid"`
	Message string `json:"message,omitempty"`
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

type TemplateListQuery struct {
	Category *string `query:"category"`
	Type     *string `query:"type"`
	ExecMode *string `query:"exec_mode"`
	Tags     *string `query:"tags"`
	Keyword  *string `query:"keyword"`
	Limit    int     `query:"limit"`
	Offset   int     `query:"offset"`
}

type OperatorTemplateResponse struct {
	ID          uuid.UUID              `json:"id"`
	Code        string                 `json:"code"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Category    string                 `json:"category"`
	Type        string                 `json:"type"`
	ExecMode    string                 `json:"exec_mode"`
	ExecConfig  map[string]interface{} `json:"exec_config,omitempty"`
	InputSchema map[string]interface{} `json:"input_schema,omitempty"`
	OutputSpec  map[string]interface{} `json:"output_spec,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Author      string                 `json:"author,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
	IconURL     string                 `json:"icon_url,omitempty"`
	Downloads   int64                  `json:"downloads"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type OperatorTemplateListResponse struct {
	Items []*OperatorTemplateResponse `json:"items"`
	Total int64                       `json:"total"`
}

type InstallTemplateReq struct {
	TemplateID   uuid.UUID `json:"template_id" validate:"required"`
	OperatorCode string    `json:"operator_code" validate:"required"`
	OperatorName string    `json:"operator_name" validate:"required"`
	Tags         []string  `json:"tags,omitempty"`
}

type OperatorDependencyItemReq struct {
	DependsOnID uuid.UUID `json:"depends_on_id" validate:"required"`
	MinVersion  string    `json:"min_version,omitempty"`
	IsOptional  bool      `json:"is_optional,omitempty"`
}

type SetDependenciesReq struct {
	Dependencies []OperatorDependencyItemReq `json:"dependencies"`
}

type OperatorDependencyResponse struct {
	ID          uuid.UUID `json:"id"`
	OperatorID  uuid.UUID `json:"operator_id"`
	DependsOnID uuid.UUID `json:"depends_on_id"`
	MinVersion  string    `json:"min_version,omitempty"`
	IsOptional  bool      `json:"is_optional"`
	CreatedAt   time.Time `json:"created_at"`
}

type DependencyCheckResponse struct {
	Satisfied bool     `json:"satisfied"`
	Unmet     []string `json:"unmet,omitempty"`
}

type OperatorVersionResponse struct {
	ID         uuid.UUID              `json:"id"`
	Version    string                 `json:"version"`
	ExecMode   string                 `json:"exec_mode"`
	ExecConfig map[string]interface{} `json:"exec_config,omitempty"`
	Status     string                 `json:"status"`
}

func OperatorVersionToResponse(v *operator.OperatorVersion) *OperatorVersionResponse {
	if v == nil {
		return nil
	}

	var execConfig map[string]interface{}
	if v.ExecConfig != nil {
		execConfig = map[string]interface{}{}
		if v.ExecConfig.HTTP != nil {
			execConfig["http"] = v.ExecConfig.HTTP
		}
		if v.ExecConfig.CLI != nil {
			execConfig["cli"] = v.ExecConfig.CLI
		}
		if v.ExecConfig.MCP != nil {
			execConfig["mcp"] = v.ExecConfig.MCP
		}
	}

	return &OperatorVersionResponse{
		ID:         v.ID,
		Version:    v.Version,
		ExecMode:   string(v.ExecMode),
		ExecConfig: execConfig,
		Status:     string(v.Status),
	}
}

func OperatorVersionsToResponse(items []*operator.OperatorVersion) []*OperatorVersionResponse {
	result := make([]*OperatorVersionResponse, len(items))
	for i := range items {
		result[i] = OperatorVersionToResponse(items[i])
	}
	return result
}

func TemplateToResponse(t *operator.OperatorTemplate) *OperatorTemplateResponse {
	if t == nil {
		return nil
	}

	var execConfig map[string]interface{}
	if t.ExecConfig != nil {
		execConfig = map[string]interface{}{}
		if t.ExecConfig.HTTP != nil {
			execConfig["http"] = t.ExecConfig.HTTP
		}
		if t.ExecConfig.CLI != nil {
			execConfig["cli"] = t.ExecConfig.CLI
		}
		if t.ExecConfig.MCP != nil {
			execConfig["mcp"] = t.ExecConfig.MCP
		}
	}

	tags := t.Tags
	if tags == nil {
		tags = []string{}
	}

	return &OperatorTemplateResponse{
		ID:          t.ID,
		Code:        t.Code,
		Name:        t.Name,
		Description: t.Description,
		Category:    string(t.Category),
		Type:        string(t.Type),
		ExecMode:    string(t.ExecMode),
		ExecConfig:  execConfig,
		InputSchema: t.InputSchema,
		OutputSpec:  t.OutputSpec,
		Config:      t.Config,
		Author:      t.Author,
		Tags:        tags,
		IconURL:     t.IconURL,
		Downloads:   t.Downloads,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
}

func TemplatesToResponse(items []*operator.OperatorTemplate) []*OperatorTemplateResponse {
	res := make([]*OperatorTemplateResponse, len(items))
	for i := range items {
		res[i] = TemplateToResponse(items[i])
	}
	return res
}

func DependencyToResponse(dep *operator.OperatorDependency) *OperatorDependencyResponse {
	if dep == nil {
		return nil
	}
	return &OperatorDependencyResponse{
		ID:          dep.ID,
		OperatorID:  dep.OperatorID,
		DependsOnID: dep.DependsOnID,
		MinVersion:  dep.MinVersion,
		IsOptional:  dep.IsOptional,
		CreatedAt:   dep.CreatedAt,
	}
}

func DependenciesToResponse(items []*operator.OperatorDependency) []*OperatorDependencyResponse {
	res := make([]*OperatorDependencyResponse, len(items))
	for i := range items {
		res[i] = DependencyToResponse(items[i])
	}
	return res
}

// OperatorResponse 算子响应
type OperatorResponse struct {
	ID              uuid.UUID                `json:"id"`
	Code            string                   `json:"code"`
	Name            string                   `json:"name"`
	Description     string                   `json:"description,omitempty"`
	Category        string                   `json:"category"`
	Type            string                   `json:"type"`
	Origin          string                   `json:"origin,omitempty"`
	ActiveVersionID *uuid.UUID               `json:"active_version_id,omitempty"`
	ExecMode        string                   `json:"exec_mode,omitempty"`
	ActiveVersion   *OperatorVersionResponse `json:"active_version,omitempty"`
	Version         string                   `json:"version,omitempty"`      // Deprecated: 兼容字段，建议使用 active_version.version
	Endpoint        string                   `json:"endpoint,omitempty"`     // Deprecated: 兼容字段，建议使用 active_version.exec_config.http.endpoint
	Method          string                   `json:"method,omitempty"`       // Deprecated: 兼容字段，建议使用 active_version.exec_config.http.method
	InputSchema     map[string]interface{}   `json:"input_schema,omitempty"` // Deprecated: 兼容字段，建议使用版本级 schema
	OutputSpec      map[string]interface{}   `json:"output_spec,omitempty"`  // Deprecated: 兼容字段，建议使用版本级 schema
	Config          map[string]interface{}   `json:"config,omitempty"`       // Deprecated: 兼容字段，建议使用版本级 config
	Status          string                   `json:"status"`
	IsBuiltin       bool                     `json:"is_builtin"` // Deprecated: 兼容字段，建议使用 origin
	Tags            []string                 `json:"tags,omitempty"`
	CreatedAt       time.Time                `json:"created_at"`
	UpdatedAt       time.Time                `json:"updated_at"`
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

	inputSchema := make(map[string]interface{})
	outputSpec := make(map[string]interface{})
	config := make(map[string]interface{})
	version := ""
	endpoint := ""
	method := ""

	tags := o.Tags
	if tags == nil {
		tags = []string{}
	}

	var activeVersion *OperatorVersionResponse
	execMode := ""
	if o.ActiveVersion != nil {
		version = o.ActiveVersion.Version
		execMode = string(o.ActiveVersion.ExecMode)
		activeVersion = OperatorVersionToResponse(o.ActiveVersion)

		if o.ActiveVersion.InputSchema != nil {
			inputSchema = o.ActiveVersion.InputSchema
		}
		if o.ActiveVersion.OutputSpec != nil {
			outputSpec = o.ActiveVersion.OutputSpec
		}
		if o.ActiveVersion.Config != nil {
			config = o.ActiveVersion.Config
		}
		if o.ActiveVersion.ExecMode == operator.ExecModeHTTP && o.ActiveVersion.ExecConfig != nil && o.ActiveVersion.ExecConfig.HTTP != nil {
			endpoint = o.ActiveVersion.ExecConfig.HTTP.Endpoint
			method = o.ActiveVersion.ExecConfig.HTTP.Method
		}
	}

	return &OperatorResponse{
		ID:              o.ID,
		Code:            o.Code,
		Name:            o.Name,
		Description:     o.Description,
		Category:        string(o.Category),
		Type:            string(o.Type),
		Origin:          string(o.Origin),
		ActiveVersionID: o.ActiveVersionID,
		ExecMode:        execMode,
		ActiveVersion:   activeVersion,
		Version:         version,
		Endpoint:        endpoint,
		Method:          method,
		InputSchema:     inputSchema,
		OutputSpec:      outputSpec,
		Config:          config,
		Status:          string(o.Status),
		IsBuiltin:       o.Origin == operator.OriginBuiltin,
		Tags:            tags,
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
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
