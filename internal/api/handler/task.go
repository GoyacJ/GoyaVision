package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"goyavision/internal/api/dto"
	authmiddleware "goyavision/internal/api/middleware"
	appdto "goyavision/internal/app/dto"
	"goyavision/internal/domain/media"
	"goyavision/internal/domain/workflow"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterTaskRoutes(public *echo.Group, protected *echo.Group, h *Handlers) {
	handler := &taskHandler{h: h}
	// Public
	public.GET("/tasks", handler.List)
	public.GET("/tasks/stats", handler.Stats)
	public.GET("/tasks/:id", handler.Get)
	public.GET("/tasks/:id/progress/stream", handler.ProgressStream)

	// Protected
	protected.POST("/tasks", handler.Create)
	protected.PUT("/tasks/:id", handler.Update)
	protected.DELETE("/tasks/:id", handler.Delete)
	protected.POST("/tasks/:id/start", handler.Start)
	protected.POST("/tasks/:id/complete", handler.Complete)
	protected.POST("/tasks/:id/fail", handler.Fail)
	protected.POST("/tasks/:id/cancel", handler.Cancel)
}

type taskHandler struct {
	h *Handlers
}

func (h *taskHandler) List(c echo.Context) error {
	var query dto.TaskListQuery
	if err := c.Bind(&query); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameters")
	}

	q := appdto.ListTasksQuery{
		WorkflowID: query.WorkflowID,
		AssetID:    query.AssetID,
		Pagination: appdto.Pagination{
			Limit:  query.Limit,
			Offset: query.Offset,
		},
	}

	if query.Status != nil {
		s := workflow.TaskStatus(*query.Status)
		q.Status = &s
	}

	if query.From != nil {
		t := time.Unix(*query.From, 0)
		q.From = &t
	}

	if query.To != nil {
		t := time.Unix(*query.To, 0)
		q.To = &t
	}

	// 权限过滤：非超级管理员只能查看自己触发的任务
	userID, ok := authmiddleware.GetUserID(c)
	roles := c.Get(authmiddleware.ContextKeyRoles)
	isSuperAdmin := false
	if roleList, ok := roles.([]string); ok {
		for _, r := range roleList {
			if r == "super_admin" {
				isSuperAdmin = true
				break
			}
		}
	}

	if !isSuperAdmin && ok {
		q.TriggeredByUserID = &userID
	}

	result, err := h.h.ListTasks.Handle(c.Request().Context(), q)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.TaskListResponse{
		Items: dto.TasksToResponse(result.Items),
		Total: result.Total,
	})
}

func (h *taskHandler) Create(c echo.Context) error {
	var req dto.TaskCreateReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	cmd := appdto.CreateTaskCommand{
		WorkflowID:  req.WorkflowID,
		AssetID:     req.AssetID,
		InputParams: req.InputParams,
	}

	task, err := h.h.CreateTask.Handle(c.Request().Context(), cmd)
	if err != nil {
		return err
	}

	var wf *workflow.Workflow
	var asset *media.Asset
	if task.WorkflowID != uuid.Nil {
		wf, _ = h.h.GetWorkflow.Handle(c.Request().Context(), appdto.GetWorkflowQuery{ID: task.WorkflowID})
	}
	if task.AssetID != nil {
		asset, _ = h.h.GetAsset.Handle(c.Request().Context(), appdto.GetAssetQuery{ID: *task.AssetID})
	}

	return c.JSON(http.StatusCreated, dto.TaskToResponseWithRelations(task, wf, asset, h.h.StorageURLConfig.Endpoint, h.h.StorageURLConfig.BucketName, h.h.StorageURLConfig.PublicBase, h.h.StorageURLConfig.UseSSL))
}

func (h *taskHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid task id")
	}

	withRelations := c.QueryParam("with_relations") == "true"
	if withRelations {
		task, err := h.h.GetTaskWithRelations.Handle(c.Request().Context(), appdto.GetTaskWithRelationsQuery{ID: id})
		if err != nil {
			return err
		}

		var wf *workflow.Workflow
		var asset *media.Asset
		if task.WorkflowID != uuid.Nil {
			wf, _ = h.h.GetWorkflow.Handle(c.Request().Context(), appdto.GetWorkflowQuery{ID: task.WorkflowID})
		}
		if task.AssetID != nil {
			asset, _ = h.h.GetAsset.Handle(c.Request().Context(), appdto.GetAssetQuery{ID: *task.AssetID})
		}

		return c.JSON(http.StatusOK, dto.TaskToResponseWithRelations(task, wf, asset, h.h.StorageURLConfig.Endpoint, h.h.StorageURLConfig.BucketName, h.h.StorageURLConfig.PublicBase, h.h.StorageURLConfig.UseSSL))
	}

	task, err := h.h.GetTask.Handle(c.Request().Context(), appdto.GetTaskQuery{ID: id})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.TaskToResponse(task))
}

func (h *taskHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid task id")
	}

	var req dto.TaskUpdateReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	cmd := appdto.UpdateTaskCommand{
		ID:          id,
		Progress:    req.Progress,
		CurrentNode: req.CurrentNode,
		Error:       req.Error,
	}

	if req.Status != nil {
		s := workflow.TaskStatus(*req.Status)
		cmd.Status = &s
	}

	task, err := h.h.UpdateTask.Handle(c.Request().Context(), cmd)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.TaskToResponse(task))
}

func (h *taskHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid task id")
	}

	err = h.h.DeleteTask.Handle(c.Request().Context(), appdto.DeleteTaskCommand{ID: id})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *taskHandler) Start(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid task id")
	}

	task, err := h.h.StartTask.Handle(c.Request().Context(), appdto.StartTaskCommand{ID: id})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.TaskToResponse(task))
}

func (h *taskHandler) Complete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid task id")
	}

	task, err := h.h.CompleteTask.Handle(c.Request().Context(), appdto.CompleteTaskCommand{ID: id})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.TaskToResponse(task))
}

func (h *taskHandler) Fail(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid task id")
	}

	var req struct {
		Error string `json:"error"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	task, err := h.h.FailTask.Handle(c.Request().Context(), appdto.FailTaskCommand{ID: id, ErrorMsg: req.Error})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.TaskToResponse(task))
}

func (h *taskHandler) Cancel(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid task id")
	}

	task, err := h.h.CancelTask.Handle(c.Request().Context(), appdto.CancelTaskCommand{ID: id})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.TaskToResponse(task))
}

func (h *taskHandler) Stats(c echo.Context) error {
	var workflowID *uuid.UUID
	if wfIDStr := c.QueryParam("workflow_id"); wfIDStr != "" {
		id, err := uuid.Parse(wfIDStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid workflow_id")
		}
		workflowID = &id
	}

	stats, err := h.h.GetTaskStats.Handle(c.Request().Context(), appdto.GetTaskStatsQuery{WorkflowID: workflowID})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.TaskStatsToResponse(stats))
}

// ProgressStream SSE 实时推送任务执行进度
func (h *taskHandler) ProgressStream(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid task id")
	}

	c.Response().Header().Set("Content-Type", "text/event-stream")
	c.Response().Header().Set("Cache-Control", "no-cache")
	c.Response().Header().Set("Connection", "keep-alive")
	c.Response().Header().Set("X-Accel-Buffering", "no")
	c.Response().WriteHeader(http.StatusOK)

	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.Request().Context().Done():
			return nil
		case <-ticker.C:
			task, err := h.h.GetTask.Handle(c.Request().Context(), appdto.GetTaskQuery{ID: id})
			if err != nil {
				return nil
			}

			event := map[string]interface{}{
				"status":          string(task.Status),
				"progress":        task.Progress,
				"current_node":    task.CurrentNode,
				"node_executions": dto.TaskToResponse(task).NodeExecutions,
				"error":           task.Error,
			}
			data, _ := json.Marshal(event)
			fmt.Fprintf(c.Response(), "data: %s\n\n", data)
			c.Response().Flush()

			if task.IsCompleted() {
				return nil
			}
		}
	}
}
