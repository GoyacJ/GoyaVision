package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"goyavision/internal/api/dto"
	appdto "goyavision/internal/app/dto"
	"goyavision/internal/domain/operator"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterOperator(g *echo.Group, h *Handlers) {
	handler := &operatorHandler{h: h}
	g.GET("/operators", handler.List)
	g.POST("/operators", handler.Create)
	g.GET("/operators/:id", handler.Get)
	g.PUT("/operators/:id", handler.Update)
	g.DELETE("/operators/:id", handler.Delete)
	g.GET("/operators/:id/versions", handler.ListVersions)
	g.POST("/operators/:id/versions", handler.CreateVersion)
	g.GET("/operators/:id/versions/:version_id", handler.GetVersion)
	g.POST("/operators/:id/versions/activate", handler.ActivateVersion)
	g.POST("/operators/:id/versions/rollback", handler.RollbackVersion)
	g.POST("/operators/:id/versions/archive", handler.ArchiveVersion)
	g.POST("/operators/validate-schema", handler.ValidateSchema)
	g.POST("/operators/validate-connection", handler.ValidateConnection)
	g.GET("/operators/templates", handler.ListTemplates)
	g.GET("/operators/templates/:template_id", handler.GetTemplate)
	g.POST("/operators/templates/install", handler.InstallTemplate)
	g.GET("/operators/:id/dependencies", handler.ListDependencies)
	g.PUT("/operators/:id/dependencies", handler.SetDependencies)
	g.GET("/operators/:id/dependencies/check", handler.CheckDependencies)
	g.POST("/operators/:id/publish", handler.Publish)
	g.POST("/operators/:id/deprecate", handler.Deprecate)
	g.POST("/operators/:id/test", handler.Test)
	g.GET("/operators/mcp/servers", handler.ListMCPServers)
	g.GET("/operators/mcp/servers/:id/tools", handler.ListMCPTools)
	g.GET("/operators/mcp/servers/:id/tools/:tool/preview", handler.PreviewMCPTool)
	g.POST("/operators/mcp/install", handler.InstallMCPOperator)
	g.POST("/operators/mcp/sync-templates", handler.SyncMCPTemplates)
	g.GET("/operators/category/:category", handler.ListByCategory)
}

type operatorHandler struct {
	h *Handlers
}

func (h *operatorHandler) List(c echo.Context) error {
	var query dto.OperatorListQuery
	if err := c.Bind(&query); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameters")
	}

	q := appdto.ListOperatorsQuery{
		Pagination: appdto.Pagination{
			Limit:  query.Limit,
			Offset: query.Offset,
		},
	}

	if query.Category != nil {
		cat := operator.Category(*query.Category)
		q.Category = &cat
	}

	if query.Type != nil {
		t := operator.Type(*query.Type)
		q.Type = &t
	}

	if query.Status != nil {
		s := operator.Status(*query.Status)
		q.Status = &s
	}

	if query.Origin != nil {
		o := operator.Origin(*query.Origin)
		q.Origin = &o
	}

	if query.ExecMode != nil {
		m := operator.ExecMode(*query.ExecMode)
		q.ExecMode = &m
	}

	if query.Tags != nil && *query.Tags != "" {
		q.Tags = strings.Split(*query.Tags, ",")
	}

	if query.Keyword != nil {
		q.Keyword = *query.Keyword
	}

	result, err := h.h.ListOperators.Handle(c.Request().Context(), q)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.OperatorListResponse{
		Items: dto.OperatorsToResponse(result.Items),
		Total: result.Total,
	})
}

func (h *operatorHandler) Create(c echo.Context) error {
	var req dto.OperatorCreateReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	status := operator.StatusDraft
	if req.Status != "" {
		status = operator.Status(req.Status)
	}

	cmd := appdto.CreateOperatorCommand{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Category:    operator.Category(req.Category),
		Type:        operator.Type(req.Type),
		Origin:      operator.Origin(req.Origin),
		ExecMode:    operator.ExecMode(req.ExecMode),
		Status:      status,
		Tags:        req.Tags,
	}

	if len(req.ExecConfig) > 0 {
		var cfg operator.ExecConfig
		b, _ := json.Marshal(req.ExecConfig)
		if err := json.Unmarshal(b, &cfg); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid exec_config")
		}
		cmd.ExecConfig = &cfg
	}

	op, err := h.h.CreateOperator.Handle(c.Request().Context(), cmd)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, dto.OperatorToResponse(op))
}

func (h *operatorHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid operator id")
	}

	op, err := h.h.GetOperator.Handle(c.Request().Context(), appdto.GetOperatorQuery{ID: id})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.OperatorToResponse(op))
}

func (h *operatorHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid operator id")
	}

	var req dto.OperatorUpdateReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	cmd := appdto.UpdateOperatorCommand{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Tags:        req.Tags,
	}

	if req.Category != nil {
		category := operator.Category(*req.Category)
		cmd.Category = &category
	}

	op, err := h.h.UpdateOperator.Handle(c.Request().Context(), cmd)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.OperatorToResponse(op))
}

func (h *operatorHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid operator id")
	}

	err = h.h.DeleteOperator.Handle(c.Request().Context(), appdto.DeleteOperatorCommand{ID: id})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *operatorHandler) ListVersions(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid operator id")
	}

	var query dto.OperatorVersionListQuery
	if err := c.Bind(&query); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameters")
	}

	result, err := h.h.ListOperatorVersions.Handle(c.Request().Context(), appdto.ListOperatorVersionsQuery{
		OperatorID: id,
		Pagination: appdto.Pagination{Limit: query.Limit, Offset: query.Offset},
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.OperatorVersionListResponse{
		Items: dto.OperatorVersionsToResponse(result.Items),
		Total: result.Total,
	})
}

func (h *operatorHandler) GetVersion(c echo.Context) error {
	operatorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid operator id")
	}
	versionID, err := uuid.Parse(c.Param("version_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid version id")
	}

	v, err := h.h.GetOperatorVersion.Handle(c.Request().Context(), appdto.GetOperatorVersionQuery{
		OperatorID: operatorID,
		VersionID:  versionID,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.OperatorVersionToResponse(v))
}

func (h *operatorHandler) CreateVersion(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid operator id")
	}

	var req dto.OperatorVersionCreateReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	cmd := appdto.CreateOperatorVersionCommand{
		OperatorID:  id,
		Version:     req.Version,
		ExecMode:    operator.ExecMode(req.ExecMode),
		InputSchema: req.InputSchema,
		OutputSpec:  req.OutputSpec,
		Config:      req.Config,
		Changelog:   req.Changelog,
		Status:      operator.VersionStatus(req.Status),
	}

	if len(req.ExecConfig) > 0 {
		var cfg operator.ExecConfig
		b, _ := json.Marshal(req.ExecConfig)
		if err := json.Unmarshal(b, &cfg); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid exec_config")
		}
		cmd.ExecConfig = &cfg
	}

	v, err := h.h.CreateOperatorVersion.Handle(c.Request().Context(), cmd)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, dto.OperatorVersionToResponse(v))
}

func (h *operatorHandler) ActivateVersion(c echo.Context) error {
	operatorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid operator id")
	}

	var req dto.OperatorVersionActionReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	op, err := h.h.ActivateVersion.Handle(c.Request().Context(), appdto.ActivateVersionCommand{
		OperatorID: operatorID,
		VersionID:  req.VersionID,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.OperatorToResponse(op))
}

func (h *operatorHandler) RollbackVersion(c echo.Context) error {
	operatorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid operator id")
	}

	var req dto.OperatorVersionActionReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	op, err := h.h.RollbackVersion.Handle(c.Request().Context(), appdto.RollbackVersionCommand{
		OperatorID: operatorID,
		VersionID:  req.VersionID,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.OperatorToResponse(op))
}

func (h *operatorHandler) ArchiveVersion(c echo.Context) error {
	operatorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid operator id")
	}

	var req dto.OperatorVersionActionReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	v, err := h.h.ArchiveVersion.Handle(c.Request().Context(), appdto.ArchiveVersionCommand{
		OperatorID: operatorID,
		VersionID:  req.VersionID,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.OperatorVersionToResponse(v))
}

func (h *operatorHandler) ValidateSchema(c echo.Context) error {
	var req dto.ValidateSchemaReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	err := h.h.ValidateSchema.Handle(c.Request().Context(), appdto.ValidateSchemaQuery{Schema: req.Schema})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.ValidateResultResponse{Valid: true})
}

func (h *operatorHandler) ValidateConnection(c echo.Context) error {
	var req dto.ValidateConnectionReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	err := h.h.ValidateConnection.Handle(c.Request().Context(), appdto.ValidateConnectionQuery{
		UpstreamOutputSpec:    req.UpstreamOutputSpec,
		DownstreamInputSchema: req.DownstreamInputSchema,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.ValidateResultResponse{Valid: true})
}

func (h *operatorHandler) ListTemplates(c echo.Context) error {
	var query dto.TemplateListQuery
	if err := c.Bind(&query); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameters")
	}

	q := appdto.ListTemplatesQuery{
		Pagination: appdto.Pagination{Limit: query.Limit, Offset: query.Offset},
	}
	if query.Category != nil {
		v := operator.Category(*query.Category)
		q.Category = &v
	}
	if query.Type != nil {
		v := operator.Type(*query.Type)
		q.Type = &v
	}
	if query.ExecMode != nil {
		v := operator.ExecMode(*query.ExecMode)
		q.ExecMode = &v
	}
	if query.Tags != nil && *query.Tags != "" {
		q.Tags = strings.Split(*query.Tags, ",")
	}
	if query.Keyword != nil {
		q.Keyword = *query.Keyword
	}

	result, err := h.h.ListTemplates.Handle(c.Request().Context(), q)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.OperatorTemplateListResponse{
		Items: dto.TemplatesToResponse(result.Items),
		Total: result.Total,
	})
}

func (h *operatorHandler) GetTemplate(c echo.Context) error {
	id, err := uuid.Parse(c.Param("template_id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid template id")
	}

	tpl, err := h.h.GetTemplate.Handle(c.Request().Context(), appdto.GetTemplateQuery{ID: id})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.TemplateToResponse(tpl))
}

func (h *operatorHandler) InstallTemplate(c echo.Context) error {
	var req dto.InstallTemplateReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	op, err := h.h.InstallTemplate.Handle(c.Request().Context(), appdto.InstallTemplateCommand{
		TemplateID:   req.TemplateID,
		OperatorCode: req.OperatorCode,
		OperatorName: req.OperatorName,
		Tags:         req.Tags,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, dto.OperatorToResponse(op))
}

func (h *operatorHandler) ListDependencies(c echo.Context) error {
	operatorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid operator id")
	}

	deps, err := h.h.ListOperatorDependencies.Handle(c.Request().Context(), appdto.ListOperatorDependenciesQuery{OperatorID: operatorID})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.DependenciesToResponse(deps))
}

func (h *operatorHandler) SetDependencies(c echo.Context) error {
	operatorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid operator id")
	}

	var req dto.SetDependenciesReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	deps := make([]appdto.DependencyItemInput, 0, len(req.Dependencies))
	for i := range req.Dependencies {
		deps = append(deps, appdto.DependencyItemInput{
			DependsOnID: req.Dependencies[i].DependsOnID,
			MinVersion:  req.Dependencies[i].MinVersion,
			IsOptional:  req.Dependencies[i].IsOptional,
		})
	}

	err = h.h.SetOperatorDependencies.Handle(c.Request().Context(), appdto.SetOperatorDependenciesCommand{
		OperatorID:   operatorID,
		Dependencies: deps,
	})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *operatorHandler) CheckDependencies(c echo.Context) error {
	operatorID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid operator id")
	}

	ok, unmet, err := h.h.CheckDependencies.Handle(c.Request().Context(), appdto.CheckDependenciesQuery{OperatorID: operatorID})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.DependencyCheckResponse{Satisfied: ok, Unmet: unmet})
}

func (h *operatorHandler) Publish(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid operator id")
	}

	op, err := h.h.PublishOperator.Handle(c.Request().Context(), appdto.PublishOperatorCommand{ID: id})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.OperatorToResponse(op))
}

func (h *operatorHandler) Deprecate(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid operator id")
	}

	op, err := h.h.DeprecateOperator.Handle(c.Request().Context(), appdto.DeprecateOperatorCommand{ID: id})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.OperatorToResponse(op))
}

func (h *operatorHandler) Test(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid operator id")
	}

	var req dto.TestOperatorReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	res, err := h.h.TestOperator.Handle(c.Request().Context(), appdto.TestOperatorCommand{
		ID:      id,
		AssetID: req.AssetID,
		Params:  req.Params,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.TestOperatorResponse{
		Success:     res.Success,
		Message:     res.Message,
		Diagnostics: res.Diagnostics,
	})
}

func (h *operatorHandler) ListByCategory(c echo.Context) error {
	categoryStr := c.Param("category")
	category := operator.Category(categoryStr)

	q := appdto.ListOperatorsQuery{
		Category: &category,
		Pagination: appdto.Pagination{
			Limit:  100,
			Offset: 0,
		},
	}

	result, err := h.h.ListOperators.Handle(c.Request().Context(), q)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.OperatorsToResponse(result.Items))
}

func (h *operatorHandler) ListMCPServers(c echo.Context) error {
	servers, err := h.h.ListMCPServers.Handle(c.Request().Context(), appdto.ListMCPServersQuery{})
	if err != nil {
		return err
	}

	resp := make([]dto.MCPServerResponse, 0, len(servers))
	for i := range servers {
		resp = append(resp, dto.MCPServerResponse{
			ID:          servers[i].ID,
			Name:        servers[i].Name,
			Description: servers[i].Description,
			Status:      servers[i].Status,
		})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *operatorHandler) ListMCPTools(c echo.Context) error {
	serverID := c.Param("id")
	tools, err := h.h.ListMCPTools.Handle(c.Request().Context(), appdto.ListMCPToolsQuery{ServerID: serverID})
	if err != nil {
		return err
	}

	resp := make([]dto.MCPToolResponse, 0, len(tools))
	for i := range tools {
		resp = append(resp, dto.MCPToolResponse{
			Name:         tools[i].Name,
			Description:  tools[i].Description,
			Version:      tools[i].Version,
			InputSchema:  tools[i].InputSchema,
			OutputSchema: tools[i].OutputSchema,
		})
	}

	return c.JSON(http.StatusOK, resp)
}

func (h *operatorHandler) PreviewMCPTool(c echo.Context) error {
	serverID := c.Param("id")
	toolName := c.Param("tool")

	tool, err := h.h.PreviewMCPTool.Handle(c.Request().Context(), appdto.PreviewMCPToolQuery{
		ServerID: serverID,
		ToolName: toolName,
	})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.MCPToolResponse{
		Name:         tool.Name,
		Description:  tool.Description,
		Version:      tool.Version,
		InputSchema:  tool.InputSchema,
		OutputSchema: tool.OutputSchema,
	})
}

func (h *operatorHandler) InstallMCPOperator(c echo.Context) error {
	var req dto.MCPInstallReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	cmd := appdto.InstallMCPOperatorCommand{
		ServerID:     req.ServerID,
		ToolName:     req.ToolName,
		OperatorCode: req.OperatorCode,
		OperatorName: req.OperatorName,
		TimeoutSec:   req.TimeoutSec,
		Tags:         req.Tags,
	}
	if req.Category != nil {
		cat := operator.Category(*req.Category)
		cmd.Category = &cat
	}
	if req.Type != nil {
		t := operator.Type(*req.Type)
		cmd.Type = &t
	}

	op, err := h.h.InstallMCPOperator.Handle(c.Request().Context(), cmd)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, dto.OperatorToResponse(op))
}

func (h *operatorHandler) SyncMCPTemplates(c echo.Context) error {
	var req dto.SyncMCPTemplatesReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	res, err := h.h.SyncMCPTemplates.Handle(c.Request().Context(), appdto.SyncMCPTemplatesCommand{ServerID: req.ServerID})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.SyncMCPTemplatesResponse{
		ServerID: res.ServerID,
		Total:    res.Total,
		Created:  res.Created,
		Updated:  res.Updated,
	})
}
