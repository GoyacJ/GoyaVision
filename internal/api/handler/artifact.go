package handler

import (
	"time"

	"goyavision/config"
	"goyavision/internal/api/dto"
	"goyavision/internal/app"
	appport "goyavision/internal/app/port"
	"goyavision/internal/domain/workflow"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterArtifact(g *echo.Group, h *Handlers) {
	svc := app.NewArtifactService(h.Repo)
	ah := artifactHandler{
		svc:    svc,
		cfg:    h.Cfg,
		urlCfg: h.StorageURLConfig,
	}
	g.GET("/artifacts", ah.List)
	g.POST("/artifacts", ah.Create)
	g.GET("/artifacts/:id", ah.Get)
	g.DELETE("/artifacts/:id", ah.Delete)
	g.GET("/tasks/:task_id/artifacts", ah.ListByTask)
}

type artifactHandler struct {
	svc    *app.ArtifactService
	cfg    *config.Config
	urlCfg appport.StorageURLConfig
}

func (h *artifactHandler) List(c echo.Context) error {
	var query dto.ArtifactListQuery
	if err := c.Bind(&query); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid query parameters",
		})
	}

	req := &app.ListArtifactsRequest{
		TaskID:  query.TaskID,
		NodeKey: query.NodeKey,
		AssetID: query.AssetID,
		Limit:   query.Limit,
		Offset:  query.Offset,
	}

	if query.Type != nil {
		t := workflow.ArtifactType(*query.Type)
		req.Type = &t
	}

	if query.From != nil {
		t := time.Unix(*query.From, 0)
		req.From = &t
	}

	if query.To != nil {
		t := time.Unix(*query.To, 0)
		req.To = &t
	}

	artifacts, total, err := h.svc.List(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.ArtifactListResponse{
		Items: dto.ArtifactsToResponse(artifacts, h.urlCfg.Endpoint, h.urlCfg.BucketName, h.urlCfg.UseSSL),
		Total: total,
	})
}

func (h *artifactHandler) Create(c echo.Context) error {
	var req dto.ArtifactCreateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	createReq := &app.CreateArtifactRequest{
		TaskID:  req.TaskID,
		Type:    workflow.ArtifactType(req.Type),
		AssetID: req.AssetID,
		Data:    req.Data,
	}

	artifact, err := h.svc.Create(c.Request().Context(), createReq)
	if err != nil {
		return err
	}

	return c.JSON(201, dto.ArtifactToResponse(artifact, h.urlCfg.Endpoint, h.urlCfg.BucketName, h.urlCfg.UseSSL))
}

func (h *artifactHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid artifact id",
		})
	}

	artifact, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.ArtifactToResponse(artifact, h.urlCfg.Endpoint, h.urlCfg.BucketName, h.urlCfg.UseSSL))
}

func (h *artifactHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid artifact id",
		})
	}

	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		return err
	}

	return c.NoContent(204)
}

func (h *artifactHandler) ListByTask(c echo.Context) error {
	taskIDStr := c.Param("task_id")
	taskID, err := uuid.Parse(taskIDStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid task id",
		})
	}

	artifactType := c.QueryParam("type")
	if artifactType != "" {
		artifacts, err := h.svc.ListByType(c.Request().Context(), taskID, workflow.ArtifactType(artifactType))
		if err != nil {
			return err
		}
		return c.JSON(200, dto.ArtifactsToResponse(artifacts, h.urlCfg.Endpoint, h.urlCfg.BucketName, h.urlCfg.UseSSL))
	}

	artifacts, err := h.svc.ListByTask(c.Request().Context(), taskID)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.ArtifactsToResponse(artifacts, h.urlCfg.Endpoint, h.urlCfg.BucketName, h.urlCfg.UseSSL))
}
