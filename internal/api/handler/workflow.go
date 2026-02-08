package handler

import (
	"net/http"
	"strings"

	"goyavision/internal/api/dto"
	authmiddleware "goyavision/internal/api/middleware"
	appdto "goyavision/internal/app/dto"
	"goyavision/internal/domain/workflow"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterWorkflowRoutes(public *echo.Group, protected *echo.Group, h *Handlers) {
	handler := &workflowHandler{h: h}
	// Public
	public.GET("/workflows", handler.List)
	public.GET("/workflows/:id", handler.Get)

	// Protected
	protected.POST("/workflows", handler.Create)
	protected.PUT("/workflows/:id", handler.Update)
	protected.DELETE("/workflows/:id", handler.Delete)
	protected.POST("/workflows/:id/enable", handler.Enable)
	protected.POST("/workflows/:id/disable", handler.Disable)
	protected.POST("/workflows/:id/trigger", handler.Trigger)
}

type workflowHandler struct {
	h *Handlers
}

func (h *workflowHandler) List(c echo.Context) error {
	var query dto.WorkflowListQuery
	if err := c.Bind(&query); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameters")
	}

	q := appdto.ListWorkflowsQuery{
		Pagination: appdto.Pagination{
			Limit:  query.Limit,
			Offset: query.Offset,
		},
	}

	if query.Status != nil {
		s := workflow.Status(*query.Status)
		q.Status = &s
	}

	if query.TriggerType != nil {
		t := workflow.TriggerType(*query.TriggerType)
		q.TriggerType = &t
	}

	if query.Tags != nil && *query.Tags != "" {
		q.Tags = strings.Split(*query.Tags, ",")
	}

	if query.Keyword != nil {
		q.Keyword = *query.Keyword
	}

	result, err := h.h.ListWorkflows.Handle(c.Request().Context(), q)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.WorkflowListResponse{
		Items: dto.WorkflowsToResponse(result.Items),
		Total: result.Total,
	})
}

func (h *workflowHandler) Create(c echo.Context) error {
	var req dto.WorkflowCreateReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	status := workflow.StatusDraft
	if req.Status != "" {
		status = workflow.Status(req.Status)
	}

	nodes := make([]appdto.WorkflowNodeInput, len(req.Nodes))
	for i, n := range req.Nodes {
		nodes[i] = appdto.WorkflowNodeInput{
			NodeKey:    n.NodeKey,
			NodeType:   n.NodeType,
			OperatorID: n.OperatorID,
			Config:     n.Config,
			Position:   n.Position,
		}
	}

	edges := make([]appdto.WorkflowEdgeInput, len(req.Edges))
	for i, e := range req.Edges {
		edges[i] = appdto.WorkflowEdgeInput{
			SourceKey: e.SourceKey,
			TargetKey: e.TargetKey,
			Condition: e.Condition,
		}
	}

	visibility := workflow.VisibilityPrivate
	if req.Visibility != nil {
		visibility = workflow.Visibility(*req.Visibility)
	}

	visibleRoleIDs := req.VisibleRoleIDs
	if visibility == workflow.VisibilityRole && len(visibleRoleIDs) == 0 {
		if ids, ok := c.Get(authmiddleware.ContextKeyRoleIDs).([]uuid.UUID); ok {
			for _, id := range ids {
				visibleRoleIDs = append(visibleRoleIDs, id.String())
			}
		}
	}

	cmd := appdto.CreateWorkflowCommand{
		Code:           req.Code,
		Name:           req.Name,
		Description:    req.Description,
		Version:        req.Version,
		TriggerType:    workflow.TriggerType(req.TriggerType),
		TriggerConf:    req.TriggerConf,
		Status:         status,
		Tags:           req.Tags,
		Nodes:          nodes,
		Edges:          edges,
		Visibility:     visibility,
		VisibleRoleIDs: visibleRoleIDs,
	}

	wf, err := h.h.CreateWorkflow.Handle(c.Request().Context(), cmd)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, dto.WorkflowToResponseWithNodes(wf))
}

func (h *workflowHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid workflow id")
	}

	withNodes := c.QueryParam("with_nodes") == "true"
	if withNodes {
		wf, err := h.h.GetWorkflowWithNodes.Handle(c.Request().Context(), appdto.GetWorkflowWithNodesQuery{ID: id})
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, dto.WorkflowToResponseWithNodes(wf))
	}

	wf, err := h.h.GetWorkflow.Handle(c.Request().Context(), appdto.GetWorkflowQuery{ID: id})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.WorkflowToResponse(wf))
}

func (h *workflowHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid workflow id")
	}

	var req dto.WorkflowUpdateReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	nodes := make([]appdto.WorkflowNodeInput, len(req.Nodes))
	for i, n := range req.Nodes {
		nodes[i] = appdto.WorkflowNodeInput{
			NodeKey:    n.NodeKey,
			NodeType:   n.NodeType,
			OperatorID: n.OperatorID,
			Config:     n.Config,
			Position:   n.Position,
		}
	}

	edges := make([]appdto.WorkflowEdgeInput, len(req.Edges))
	for i, e := range req.Edges {
		edges[i] = appdto.WorkflowEdgeInput{
			SourceKey: e.SourceKey,
			TargetKey: e.TargetKey,
			Condition: e.Condition,
		}
	}

	cmd := appdto.UpdateWorkflowCommand{
		ID:             id,
		Name:           req.Name,
		Description:    req.Description,
		TriggerConf:    req.TriggerConf,
		Tags:           req.Tags,
		Nodes:          nodes,
		Edges:          edges,
		VisibleRoleIDs: req.VisibleRoleIDs,
	}

	if req.Visibility != nil {
		v := workflow.Visibility(*req.Visibility)
		cmd.Visibility = &v
	}

	if req.Status != nil {
		s := workflow.Status(*req.Status)
		cmd.Status = &s
	}

	result, err := h.h.UpdateWorkflow.Handle(c.Request().Context(), cmd)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.WorkflowToResponseWithNodes(result))
}

func (h *workflowHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid workflow id")
	}

	err = h.h.DeleteWorkflow.Handle(c.Request().Context(), appdto.DeleteWorkflowCommand{ID: id})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *workflowHandler) Enable(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid workflow id")
	}

	result, err := h.h.EnableWorkflow.Handle(c.Request().Context(), appdto.EnableWorkflowCommand{ID: id, Enabled: true})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.WorkflowToResponseWithNodes(result))
}

func (h *workflowHandler) Disable(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid workflow id")
	}

	result, err := h.h.EnableWorkflow.Handle(c.Request().Context(), appdto.EnableWorkflowCommand{ID: id, Enabled: false})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.WorkflowToResponse(result))
}

func (h *workflowHandler) Trigger(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid workflow id")
	}

	var req struct {
		AssetID *uuid.UUID `json:"asset_id,omitempty"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	if h.h.WorkflowScheduler == nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "scheduler not available")
	}

	task, err := h.h.WorkflowScheduler.TriggerWorkflow(c.Request().Context(), id, req.AssetID)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusAccepted, dto.TaskToResponse(task))
}
