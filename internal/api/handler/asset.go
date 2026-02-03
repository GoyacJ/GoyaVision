package handler

import (
	"strings"
	"time"

	"goyavision/internal/api/dto"
	"goyavision/internal/app"
	"goyavision/internal/domain"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterAsset(g *echo.Group, d Deps) {
	svc := app.NewMediaAssetService(d.Repo)
	h := assetHandler{svc: svc}
	g.GET("/assets", h.List)
	g.POST("/assets", h.Create)
	g.GET("/assets/:id", h.Get)
	g.PUT("/assets/:id", h.Update)
	g.DELETE("/assets/:id", h.Delete)
	g.GET("/assets/:id/children", h.ListChildren)
	g.GET("/assets/tags", h.GetAllTags) // 获取所有标签
}

type assetHandler struct {
	svc *app.MediaAssetService
}

func (h *assetHandler) List(c echo.Context) error {
	var query dto.AssetListQuery
	if err := c.Bind(&query); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid query parameters",
		})
	}

	req := &app.ListMediaAssetsRequest{
		Limit:  query.Limit,
		Offset: query.Offset,
	}

	if query.Type != nil {
		t := domain.AssetType(*query.Type)
		req.Type = &t
	}

	if query.SourceType != nil {
		st := domain.AssetSourceType(*query.SourceType)
		req.SourceType = &st
	}

	if query.SourceID != nil {
		req.SourceID = query.SourceID
	}

	if query.ParentID != nil {
		req.ParentID = query.ParentID
	}

	if query.Status != nil {
		s := domain.AssetStatus(*query.Status)
		req.Status = &s
	}

	if query.Tags != nil && *query.Tags != "" {
		req.Tags = strings.Split(*query.Tags, ",")
	}

	if query.From != nil {
		t := time.Unix(*query.From, 0)
		req.From = &t
	}

	if query.To != nil {
		t := time.Unix(*query.To, 0)
		req.To = &t
	}

	assets, total, err := h.svc.List(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.AssetListResponse{
		Items: dto.AssetsToResponse(assets),
		Total: total,
	})
}

func (h *assetHandler) Create(c echo.Context) error {
	var req dto.AssetCreateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	status := domain.AssetStatusPending
	if req.Status != "" {
		status = domain.AssetStatus(req.Status)
	}

	createReq := &app.CreateMediaAssetRequest{
		Type:       domain.AssetType(req.Type),
		SourceType: domain.AssetSourceType(req.SourceType),
		SourceID:   req.SourceID,
		ParentID:   req.ParentID,
		Name:       req.Name,
		Path:       req.Path,
		Duration:   req.Duration,
		Size:       req.Size,
		Format:     req.Format,
		Metadata:   req.Metadata,
		Status:     status,
		Tags:       req.Tags,
	}

	asset, err := h.svc.Create(c.Request().Context(), createReq)
	if err != nil {
		return err
	}

	return c.JSON(201, dto.AssetToResponse(asset))
}

func (h *assetHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid asset id",
		})
	}

	asset, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.AssetToResponse(asset))
}

func (h *assetHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid asset id",
		})
	}

	var req dto.AssetUpdateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	updateReq := &app.UpdateMediaAssetRequest{
		Name:     req.Name,
		Metadata: req.Metadata,
		Tags:     req.Tags,
	}

	if req.Status != nil {
		s := domain.AssetStatus(*req.Status)
		updateReq.Status = &s
	}

	asset, err := h.svc.Update(c.Request().Context(), id, updateReq)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.AssetToResponse(asset))
}

func (h *assetHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid asset id",
		})
	}

	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		return err
	}

	return c.NoContent(204)
}

func (h *assetHandler) ListChildren(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid asset id",
		})
	}

	children, err := h.svc.ListChildren(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.AssetsToResponse(children))
}

func (h *assetHandler) GetAllTags(c echo.Context) error {
	tags, err := h.svc.GetAllTags(c.Request().Context())
	if err != nil {
		return err
	}

	return c.JSON(200, map[string]interface{}{
		"tags": tags,
	})
}
