package handler

import (
	"strings"

	"goyavision/internal/api/dto"
	"goyavision/internal/app"
	"goyavision/internal/domain"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterWorkflow(g *echo.Group, d Deps) {
	svc := app.NewWorkflowService(d.Repo)
	h := workflowHandler{
		svc:       svc,
		scheduler: d.WorkflowScheduler,
	}
	g.GET("/workflows", h.List)
	g.POST("/workflows", h.Create)
	g.GET("/workflows/:id", h.Get)
	g.PUT("/workflows/:id", h.Update)
	g.DELETE("/workflows/:id", h.Delete)
	g.POST("/workflows/:id/enable", h.Enable)
	g.POST("/workflows/:id/disable", h.Disable)
	g.POST("/workflows/:id/trigger", h.Trigger)
}

type workflowHandler struct {
	svc       *app.WorkflowService
	scheduler *app.WorkflowScheduler
}

func (h *workflowHandler) List(c echo.Context) error {
	var query dto.WorkflowListQuery
	if err := c.Bind(&query); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid query parameters",
		})
	}

	req := &app.ListWorkflowsRequest{
		Limit:  query.Limit,
		Offset: query.Offset,
	}

	if query.Status != nil {
		s := domain.WorkflowStatus(*query.Status)
		req.Status = &s
	}

	if query.TriggerType != nil {
		t := domain.TriggerType(*query.TriggerType)
		req.TriggerType = &t
	}

	if query.Tags != nil && *query.Tags != "" {
		req.Tags = strings.Split(*query.Tags, ",")
	}

	if query.Keyword != nil {
		req.Keyword = *query.Keyword
	}

	workflows, total, err := h.svc.List(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.WorkflowListResponse{
		Items: dto.WorkflowsToResponse(workflows),
		Total: total,
	})
}

func (h *workflowHandler) Create(c echo.Context) error {
	var req dto.WorkflowCreateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	status := domain.WorkflowStatusDraft
	if req.Status != "" {
		status = domain.WorkflowStatus(req.Status)
	}

	nodes := make([]app.WorkflowNodeInput, len(req.Nodes))
	for i, n := range req.Nodes {
		nodes[i] = app.WorkflowNodeInput{
			NodeKey:    n.NodeKey,
			NodeType:   n.NodeType,
			OperatorID: n.OperatorID,
			Config:     n.Config,
			Position:   n.Position,
		}
	}

	edges := make([]app.WorkflowEdgeInput, len(req.Edges))
	for i, e := range req.Edges {
		edges[i] = app.WorkflowEdgeInput{
			SourceKey: e.SourceKey,
			TargetKey: e.TargetKey,
			Condition: e.Condition,
		}
	}

	createReq := &app.CreateWorkflowRequest{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Version:     req.Version,
		TriggerType: domain.TriggerType(req.TriggerType),
		TriggerConf: req.TriggerConf,
		Status:      status,
		Tags:        req.Tags,
		Nodes:       nodes,
		Edges:       edges,
	}

	workflow, err := h.svc.Create(c.Request().Context(), createReq)
	if err != nil {
		return err
	}

	return c.JSON(201, dto.WorkflowToResponseWithNodes(workflow))
}

func (h *workflowHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid workflow id",
		})
	}

	withNodes := c.QueryParam("with_nodes") == "true"
	if withNodes {
		workflow, err := h.svc.GetWithNodes(c.Request().Context(), id)
		if err != nil {
			return err
		}
		return c.JSON(200, dto.WorkflowToResponseWithNodes(workflow))
	}

	workflow, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.WorkflowToResponse(workflow))
}

func (h *workflowHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid workflow id",
		})
	}

	var req dto.WorkflowUpdateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	nodes := make([]app.WorkflowNodeInput, len(req.Nodes))
	for i, n := range req.Nodes {
		nodes[i] = app.WorkflowNodeInput{
			NodeKey:    n.NodeKey,
			NodeType:   n.NodeType,
			OperatorID: n.OperatorID,
			Config:     n.Config,
			Position:   n.Position,
		}
	}

	edges := make([]app.WorkflowEdgeInput, len(req.Edges))
	for i, e := range req.Edges {
		edges[i] = app.WorkflowEdgeInput{
			SourceKey: e.SourceKey,
			TargetKey: e.TargetKey,
			Condition: e.Condition,
		}
	}

	updateReq := &app.UpdateWorkflowRequest{
		Name:        req.Name,
		Description: req.Description,
		TriggerConf: req.TriggerConf,
		Tags:        req.Tags,
		Nodes:       nodes,
		Edges:       edges,
	}

	if req.Status != nil {
		s := domain.WorkflowStatus(*req.Status)
		updateReq.Status = &s
	}

	workflow, err := h.svc.Update(c.Request().Context(), id, updateReq)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.WorkflowToResponseWithNodes(workflow))
}

func (h *workflowHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid workflow id",
		})
	}

	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		return err
	}

	return c.NoContent(204)
}

func (h *workflowHandler) Enable(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid workflow id",
		})
	}

	workflow, err := h.svc.Enable(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.WorkflowToResponseWithNodes(workflow))
}

func (h *workflowHandler) Disable(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid workflow id",
		})
	}

	workflow, err := h.svc.Disable(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.WorkflowToResponse(workflow))
}

func (h *workflowHandler) Trigger(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid workflow id",
		})
	}

	var req struct {
		AssetID *uuid.UUID `json:"asset_id,omitempty"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	if h.scheduler == nil {
		return c.JSON(500, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "scheduler not available",
		})
	}

	task, err := h.scheduler.TriggerWorkflow(c.Request().Context(), id, req.AssetID)
	if err != nil {
		return err
	}

	return c.JSON(202, dto.TaskToResponse(task))
}
