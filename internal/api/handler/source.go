package handler

import (
	"net/http"
	"strings"

	"goyavision/config"
	"goyavision/internal/api/dto"
	"goyavision/internal/app"
	"goyavision/internal/domain"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterSource(g *echo.Group, d Deps) {
	if d.MediaSourceService == nil {
		return
	}
	h := &sourceHandler{svc: d.MediaSourceService, cfg: d.Cfg}
	g.GET("/sources", h.List)
	g.POST("/sources", h.Create)
	g.GET("/sources/:id", h.Get)
	g.PUT("/sources/:id", h.Update)
	g.DELETE("/sources/:id", h.Delete)
	g.GET("/sources/:id/preview", h.GetPreview)
}

type sourceHandler struct {
	svc *app.MediaSourceService
	cfg *config.Config
}

func (h *sourceHandler) List(c echo.Context) error {
	var query dto.SourceListQuery
	if err := c.Bind(&query); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid query parameters",
		})
	}
	filter := domain.MediaSourceFilter{
		Limit:  query.Limit,
		Offset: query.Offset,
	}
	if query.Type != nil {
		t := domain.SourceType(*query.Type)
		filter.Type = &t
	}
	list, total, err := h.svc.List(c.Request().Context(), filter)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dto.SourceListResponse{
		Items: dto.SourcesToResponse(list),
		Total: total,
	})
}

func (h *sourceHandler) Create(c echo.Context) error {
	var req dto.SourceCreateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}
	createReq := &app.CreateMediaSourceRequest{
		Name:     req.Name,
		Type:     domain.SourceType(req.Type),
		URL:      req.URL,
		Protocol: req.Protocol,
		Enabled:  req.Enabled,
	}
	src, err := h.svc.Create(c.Request().Context(), createReq)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, dto.SourceToResponse(src))
}

func (h *sourceHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid source id",
		})
	}
	src, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dto.SourceToResponse(src))
}

func (h *sourceHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid source id",
		})
	}
	var req dto.SourceUpdateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}
	updateReq := &app.UpdateMediaSourceRequest{
		Name:     req.Name,
		URL:      req.URL,
		Protocol: req.Protocol,
		Enabled:  req.Enabled,
	}
	src, err := h.svc.Update(c.Request().Context(), id, updateReq)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dto.SourceToResponse(src))
}

func (h *sourceHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid source id",
		})
	}
	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		return err
	}
	return c.NoContent(http.StatusNoContent)
}

func (h *sourceHandler) GetPreview(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid source id",
		})
	}
	src, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		return err
	}
	m := h.cfg.MediaMTX
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
	if src.Type == domain.SourceTypePush {
		resp.PushURL = baseRTMP + "/" + pathName
	}
	return c.JSON(http.StatusOK, resp)
}
