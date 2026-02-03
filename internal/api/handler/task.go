package handler

import (
	"time"

	"goyavision/config"
	"goyavision/internal/api/dto"
	"goyavision/internal/app"
	"goyavision/internal/domain"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterTask(g *echo.Group, d Deps) {
	svc := app.NewTaskService(d.Repo)
	h := taskHandler{
		svc: svc,
		cfg: d.Cfg,
	}
	g.GET("/tasks", h.List)
	g.POST("/tasks", h.Create)
	g.GET("/tasks/:id", h.Get)
	g.PUT("/tasks/:id", h.Update)
	g.DELETE("/tasks/:id", h.Delete)
	g.POST("/tasks/:id/start", h.Start)
	g.POST("/tasks/:id/complete", h.Complete)
	g.POST("/tasks/:id/fail", h.Fail)
	g.POST("/tasks/:id/cancel", h.Cancel)
	g.GET("/tasks/stats", h.Stats)
}

type taskHandler struct {
	svc *app.TaskService
	cfg *config.Config
}

func (h *taskHandler) List(c echo.Context) error {
	var query dto.TaskListQuery
	if err := c.Bind(&query); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid query parameters",
		})
	}

	req := &app.ListTasksRequest{
		WorkflowID: query.WorkflowID,
		AssetID:    query.AssetID,
		Limit:      query.Limit,
		Offset:     query.Offset,
	}

	if query.Status != nil {
		s := domain.TaskStatus(*query.Status)
		req.Status = &s
	}

	if query.From != nil {
		t := time.Unix(*query.From, 0)
		req.From = &t
	}

	if query.To != nil {
		t := time.Unix(*query.To, 0)
		req.To = &t
	}

	tasks, total, err := h.svc.List(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.TaskListResponse{
		Items: dto.TasksToResponse(tasks),
		Total: total,
	})
}

func (h *taskHandler) Create(c echo.Context) error {
	var req dto.TaskCreateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	createReq := &app.CreateTaskRequest{
		WorkflowID:  req.WorkflowID,
		AssetID:     req.AssetID,
		InputParams: req.InputParams,
	}

	task, err := h.svc.Create(c.Request().Context(), createReq)
	if err != nil {
		return err
	}

	return c.JSON(201, dto.TaskToResponseWithRelations(task, h.cfg.MinIO.Endpoint, h.cfg.MinIO.BucketName, h.cfg.MinIO.UseSSL))
}

func (h *taskHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid task id",
		})
	}

	withRelations := c.QueryParam("with_relations") == "true"
	if withRelations {
		task, err := h.svc.GetWithRelations(c.Request().Context(), id)
		if err != nil {
			return err
		}
		return c.JSON(200, dto.TaskToResponseWithRelations(task, h.cfg.MinIO.Endpoint, h.cfg.MinIO.BucketName, h.cfg.MinIO.UseSSL))
	}

	task, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.TaskToResponse(task))
}

func (h *taskHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid task id",
		})
	}

	var req dto.TaskUpdateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	updateReq := &app.UpdateTaskRequest{
		Progress:    req.Progress,
		CurrentNode: req.CurrentNode,
		Error:       req.Error,
	}

	if req.Status != nil {
		s := domain.TaskStatus(*req.Status)
		updateReq.Status = &s
	}

	task, err := h.svc.Update(c.Request().Context(), id, updateReq)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.TaskToResponse(task))
}

func (h *taskHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid task id",
		})
	}

	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		return err
	}

	return c.NoContent(204)
}

func (h *taskHandler) Start(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid task id",
		})
	}

	task, err := h.svc.Start(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.TaskToResponse(task))
}

func (h *taskHandler) Complete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid task id",
		})
	}

	task, err := h.svc.Complete(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.TaskToResponse(task))
}

func (h *taskHandler) Fail(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid task id",
		})
	}

	var req struct {
		Error string `json:"error"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	task, err := h.svc.Fail(c.Request().Context(), id, req.Error)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.TaskToResponse(task))
}

func (h *taskHandler) Cancel(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid task id",
		})
	}

	task, err := h.svc.Cancel(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.TaskToResponse(task))
}

func (h *taskHandler) Stats(c echo.Context) error {
	var workflowID *uuid.UUID
	if wfIDStr := c.QueryParam("workflow_id"); wfIDStr != "" {
		id, err := uuid.Parse(wfIDStr)
		if err != nil {
			return c.JSON(400, dto.ErrorResponse{
				Error:   "Bad Request",
				Message: "invalid workflow_id",
			})
		}
		workflowID = &id
	}

	stats, err := h.svc.GetStats(c.Request().Context(), workflowID)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.TaskStatsToResponse(stats))
}
