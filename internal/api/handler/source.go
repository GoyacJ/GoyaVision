package handler

import (
	"net/http"
	"strings"

	"goyavision/internal/api/dto"
	authmiddleware "goyavision/internal/api/middleware"
	appdto "goyavision/internal/app/dto"
	"goyavision/internal/domain/media"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterSource(g *echo.Group, h *Handlers) {
	handler := &sourceHandler{h: h}
	g.GET("/sources", handler.List)
	g.POST("/sources", handler.Create)
	g.GET("/sources/:id", handler.Get)
	g.PUT("/sources/:id", handler.Update)
	g.DELETE("/sources/:id", handler.Delete)
	g.GET("/sources/:id/preview", handler.GetPreview)
}

type sourceHandler struct {
	h *Handlers
}

func (h *sourceHandler) List(c echo.Context) error {
	var query dto.SourceListQuery
	if err := c.Bind(&query); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameters")
	}

	q := appdto.ListSourcesQuery{
		Pagination: appdto.Pagination{
			Limit:  query.Limit,
			Offset: query.Offset,
		},
	}
	if query.Type != nil {
		t := media.SourceType(*query.Type)
		q.Type = &t
	}

	result, err := h.h.ListSources.Handle(c.Request().Context(), q)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.SourceListResponse{
		Items: dto.SourcesToResponse(result.Items),
		Total: result.Total,
	})
}

func (h *sourceHandler) Create(c echo.Context) error {
	var req dto.SourceCreateReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	visibility := media.VisibilityPrivate
	if req.Visibility != nil {
		visibility = media.Visibility(*req.Visibility)
	}

	visibleRoleIDs := req.VisibleRoleIDs
	if visibility == media.VisibilityRole && len(visibleRoleIDs) == 0 {
		if ids, ok := c.Get(authmiddleware.ContextKeyRoleIDs).([]uuid.UUID); ok {
			for _, id := range ids {
				visibleRoleIDs = append(visibleRoleIDs, id.String())
			}
		}
	}

	cmd := appdto.CreateSourceCommand{
		Name:           req.Name,
		Type:           media.SourceType(req.Type),
		URL:            req.URL,
		Protocol:       req.Protocol,
		Enabled:        req.Enabled,
		Visibility:     visibility,
		VisibleRoleIDs: visibleRoleIDs,
	}

	source, err := h.h.CreateSource.Handle(c.Request().Context(), cmd)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, dto.SourceToResponse(source))
}

func (h *sourceHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid source id")
	}

	source, err := h.h.GetSource.Handle(c.Request().Context(), appdto.GetSourceQuery{ID: id})
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.SourceToResponse(source))
}

func (h *sourceHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid source id")
	}

	var req dto.SourceUpdateReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	cmd := appdto.UpdateSourceCommand{
		ID:             id,
		Name:           req.Name,
		URL:            req.URL,
		Protocol:       req.Protocol,
		Enabled:        req.Enabled,
		VisibleRoleIDs: req.VisibleRoleIDs,
	}

	if req.Visibility != nil {
		v := media.Visibility(*req.Visibility)
		cmd.Visibility = &v
	}

	source, err := h.h.UpdateSource.Handle(c.Request().Context(), cmd)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.SourceToResponse(source))
}

func (h *sourceHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid source id")
	}

	err = h.h.DeleteSource.Handle(c.Request().Context(), appdto.DeleteSourceCommand{ID: id})
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *sourceHandler) GetPreview(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid source id")
	}

	source, err := h.h.GetSource.Handle(c.Request().Context(), appdto.GetSourceQuery{ID: id})
	if err != nil {
		return err
	}

	src := source
	m := h.h.Cfg.MediaMTX
	baseHLS := strings.TrimSuffix(m.HLSAddress, "/")
	baseRTSP := strings.TrimSuffix(m.RTSPAddress, "/")
	baseRTMP := strings.TrimSuffix(m.RTMPAddress, "/")
	baseWebRTC := strings.TrimSuffix(m.WebRTCAddress, "/")
	pathName := src.PathName

	resp := dto.SourcePreviewResponse{
		PathName:  pathName,
		HLSURL:    baseHLS + "/" + pathName + "/index.m3u8",
		RTSPURL:   baseRTSP + "/" + pathName,
		RTMPURL:   baseRTMP + "/" + pathName,
		WebRTCURL: baseWebRTC + "/" + pathName + "/whep",
	}
	if src.Type == media.SourceTypePush {
		resp.PushURL = baseRTMP + "/" + pathName
	}

	return c.JSON(http.StatusOK, resp)
}
