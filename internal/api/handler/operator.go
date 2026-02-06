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

	if query.IsBuiltin != nil {
		q.IsBuiltin = query.IsBuiltin
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
		Version:     req.Version,
		Endpoint:    req.Endpoint,
		Method:      req.Method,
		InputSchema: req.InputSchema,
		OutputSpec:  req.OutputSpec,
		Config:      req.Config,
		Status:      status,
		IsBuiltin:   req.IsBuiltin,
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
