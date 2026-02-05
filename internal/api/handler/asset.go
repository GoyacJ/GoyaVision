package handler

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	appdto "goyavision/internal/app/dto"
	"goyavision/internal/api/dto"
	"goyavision/internal/domain/media"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func inferProtocol(rawURL string) string {
	rawURL = strings.TrimSpace(rawURL)
	if rawURL == "" {
		return ""
	}
	u, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}
	scheme := strings.ToLower(u.Scheme)
	switch scheme {
	case "rtsp", "rtmp", "rtmps", "hls", "https", "http", "webrtc", "srt":
		return scheme
	}
	if scheme != "" {
		return scheme
	}
	return ""
}

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
		Items: dto.AssetsToResponse(result.Items, h.h.Cfg.MinIO.Endpoint, h.h.Cfg.MinIO.BucketName, h.h.Cfg.MinIO.UseSSL),
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

	cmd := appdto.CreateAssetCommand{
		Type:       media.AssetType(req.Type),
		SourceType: media.AssetSourceType(req.SourceType),
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

	if req.Type == string(media.AssetTypeStream) && req.SourceType == string(media.AssetSourceLive) {
		if req.StreamURL != "" {
			srcCmd := appdto.CreateSourceCommand{
				Name:     req.Name,
				Type:     media.SourceTypePull,
				URL:      req.StreamURL,
				Protocol: inferProtocol(req.StreamURL),
				Enabled:  true,
			}
			source, err := h.h.CreateSource.Handle(c.Request().Context(), srcCmd)
			if err != nil {
				return err
			}
			cmd.SourceID = &source.ID
			cmd.Path = source.PathName
		} else if req.SourceID != nil && req.Path == "" {
			source, err := h.h.GetSource.Handle(c.Request().Context(), appdto.GetSourceQuery{ID: *req.SourceID})
			if err != nil {
				return err
			}
			cmd.Path = source.PathName
		}
	}

	asset, err := h.h.CreateAsset.Handle(c.Request().Context(), cmd)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, dto.AssetToResponse(asset, h.h.Cfg.MinIO.Endpoint, h.h.Cfg.MinIO.BucketName, h.h.Cfg.MinIO.UseSSL))
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

	return c.JSON(http.StatusOK, dto.AssetToResponse(asset, h.h.Cfg.MinIO.Endpoint, h.h.Cfg.MinIO.BucketName, h.h.Cfg.MinIO.UseSSL))
}

func (h *assetHandler) Update(c echo.Context) error {
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
		ID:       id,
		Name:     req.Name,
		Metadata: req.Metadata,
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

	return c.JSON(http.StatusOK, dto.AssetToResponse(asset, h.h.Cfg.MinIO.Endpoint, h.h.Cfg.MinIO.BucketName, h.h.Cfg.MinIO.UseSSL))
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

	return c.JSON(http.StatusOK, dto.AssetsToResponse(result.Assets, h.h.Cfg.MinIO.Endpoint, h.h.Cfg.MinIO.BucketName, h.h.Cfg.MinIO.UseSSL))
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
