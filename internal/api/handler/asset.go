package handler

import (
	"net/http"
	"strings"
	"time"

	"goyavision/internal/api/dto"
	authmiddleware "goyavision/internal/api/middleware"
	appdto "goyavision/internal/app/dto"
	"goyavision/internal/domain/media"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterAsset(g *echo.Group, h *Handlers) {
	handler := &assetHandler{h: h}
	g.GET("/assets", handler.List)
	g.POST("/assets", handler.Create)
	g.GET("/assets/:id", handler.Get)
	g.PUT("/assets/:id", handler.Update)
	g.DELETE("/assets/:id", handler.Delete)
	g.GET("/assets/:id/children", handler.ListChildren)
	g.GET("/assets/tags", handler.GetAllTags)
}

type assetHandler struct {
	h *Handlers
}

func (h *assetHandler) List(c echo.Context) error {
	var query dto.AssetListQuery
	if err := c.Bind(&query); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameters")
	}

	q := appdto.ListAssetsQuery{
		Pagination: appdto.Pagination{
			Limit:  query.Limit,
			Offset: query.Offset,
		},
	}

	if query.Type != nil {
		t := media.AssetType(*query.Type)
		q.Type = &t
	}

	if query.SourceType != nil {
		st := media.AssetSourceType(*query.SourceType)
		q.SourceType = &st
	}

	if query.SourceID != nil {
		q.SourceID = query.SourceID
	}

	if query.ParentID != nil {
		q.ParentID = query.ParentID
	}

	if query.Status != nil {
		s := media.AssetStatus(*query.Status)
		q.Status = &s
	}

	if query.Tags != nil && *query.Tags != "" {
		q.Tags = strings.Split(*query.Tags, ",")
	}

	if query.From != nil {
		t := time.Unix(*query.From, 0)
		q.From = &t
	}

	if query.To != nil {
		t := time.Unix(*query.To, 0)
		q.To = &t
	}

	result, err := h.h.ListAssets.Handle(c.Request().Context(), q)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.AssetListResponse{
		Items: dto.AssetsToResponse(result.Items, h.h.Cfg.MinIO.Endpoint, h.h.Cfg.MinIO.BucketName, h.h.Cfg.MinIO.PublicBase, h.h.Cfg.MinIO.UseSSL),
		Total: result.Total,
	})
}

func (h *assetHandler) Create(c echo.Context) error {
	var req dto.AssetCreateReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	status := media.AssetStatusPending
	if req.Status != "" {
		status = media.AssetStatus(req.Status)
	}

	visibility := media.VisibilityPrivate
	if req.Visibility != nil {
		visibility = media.Visibility(*req.Visibility)
	}

	cmd := appdto.CreateAssetCommand{
		Type:           media.AssetType(req.Type),
		SourceType:     media.AssetSourceType(req.SourceType),
		SourceID:       req.SourceID,
		ParentID:       req.ParentID,
		Name:           req.Name,
		Path:           req.Path,
		Duration:       req.Duration,
		Size:           req.Size,
		Format:         req.Format,
		Metadata:       req.Metadata,
		Status:         status,
		Tags:           req.Tags,
		Visibility:     visibility,
		VisibleRoleIDs: req.VisibleRoleIDs,
	}

	asset, err := h.h.CreateAsset.Handle(c.Request().Context(), cmd)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, dto.AssetToResponse(asset, h.h.Cfg.MinIO.Endpoint, h.h.Cfg.MinIO.BucketName, h.h.Cfg.MinIO.PublicBase, h.h.Cfg.MinIO.UseSSL))
}

func (h *assetHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid asset id")
	}

	asset, err := h.h.GetAsset.Handle(c.Request().Context(), appdto.GetAssetQuery{ID: id})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.AssetToResponse(asset, h.h.Cfg.MinIO.Endpoint, h.h.Cfg.MinIO.BucketName, h.h.Cfg.MinIO.PublicBase, h.h.Cfg.MinIO.UseSSL))
}

func (h *assetHandler) Update(c echo.Context) error {
	if !authmiddleware.HasPermission(c, "asset:update") {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error:   "Forbidden",
			Message: "无编辑权限",
		})
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid asset id")
	}

	var req dto.AssetUpdateReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	cmd := appdto.UpdateAssetCommand{
		ID:             id,
		Name:           req.Name,
		Metadata:       req.Metadata,
		VisibleRoleIDs: req.VisibleRoleIDs,
	}

	if req.Visibility != nil {
		v := media.Visibility(*req.Visibility)
		cmd.Visibility = &v
	}

	if req.Tags != nil {
		tags := req.Tags
		cmd.Tags = &tags
	}

	if req.Status != nil {
		s := media.AssetStatus(*req.Status)
		cmd.Status = &s
	}

	asset, err := h.h.UpdateAsset.Handle(c.Request().Context(), cmd)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.AssetToResponse(asset, h.h.Cfg.MinIO.Endpoint, h.h.Cfg.MinIO.BucketName, h.h.Cfg.MinIO.PublicBase, h.h.Cfg.MinIO.UseSSL))
}

func (h *assetHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid asset id")
	}

	err = h.h.DeleteAsset.Handle(c.Request().Context(), appdto.DeleteAssetCommand{ID: id})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *assetHandler) ListChildren(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid asset id")
	}

	result, err := h.h.ListAssetChildren.Handle(c.Request().Context(), appdto.ListAssetChildrenQuery{ParentID: id})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.AssetsToResponse(result.Assets, h.h.Cfg.MinIO.Endpoint, h.h.Cfg.MinIO.BucketName, h.h.Cfg.MinIO.PublicBase, h.h.Cfg.MinIO.UseSSL))
}

func (h *assetHandler) GetAllTags(c echo.Context) error {
	result, err := h.h.GetAssetTags.Handle(c.Request().Context(), appdto.GetAssetTagsQuery{})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"tags": result.Tags,
	})
}
