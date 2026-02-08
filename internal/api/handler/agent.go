package handler

import (
	"net/http"

	"goyavision/internal/api/dto"
	appdto "goyavision/internal/app/dto"
	"goyavision/internal/domain/agent"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterAgentRoutes(public *echo.Group, protected *echo.Group, h *Handlers) {
	handler := &agentHandler{h: h}

	public.GET("/agent/sessions/:id", handler.GetSession)
	protected.GET("/agent/sessions", handler.ListSessions)
	protected.GET("/agent/sessions/:id/events", handler.ListSessionEvents)
	protected.POST("/agent/sessions", handler.CreateSession)
	protected.POST("/agent/sessions/:id/run", handler.RunSession)
	protected.POST("/agent/sessions/:id/stop", handler.StopSession)
}

type agentHandler struct {
	h *Handlers
}

func (h *agentHandler) CreateSession(c echo.Context) error {
	var req dto.AgentSessionCreateReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	session, err := h.h.CreateAgentSession.Handle(c.Request().Context(), appdto.CreateAgentSessionCommand{
		TaskID: req.TaskID,
		Budget: req.Budget,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, dto.AgentSessionToResponse(session))
}

func (h *agentHandler) GetSession(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid agent session id")
	}

	session, err := h.h.GetAgentSession.Handle(c.Request().Context(), appdto.GetAgentSessionQuery{ID: id})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dto.AgentSessionToResponse(session))
}

func (h *agentHandler) StopSession(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid agent session id")
	}

	var req dto.AgentSessionStopReq
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	status := agent.SessionStatus(req.Status)
	session, err := h.h.StopAgentSession.Handle(c.Request().Context(), appdto.StopAgentSessionCommand{
		SessionID: id,
		Status:    status,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dto.AgentSessionToResponse(session))
}

func (h *agentHandler) RunSession(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid agent session id")
	}

	req := dto.AgentSessionRunReq{}
	if c.Request().ContentLength > 0 {
		if err := c.Bind(&req); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
		}
	}

	session, err := h.h.RunAgentSessionStep.Handle(c.Request().Context(), appdto.RunAgentSessionStepCommand{
		SessionID:  id,
		MaxActions: req.MaxActions,
	})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, dto.AgentSessionToResponse(session))
}

func (h *agentHandler) ListSessions(c echo.Context) error {
	var req dto.AgentSessionListQuery
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameters")
	}

	query := appdto.ListAgentSessionsQuery{
		TaskID: req.TaskID,
		Pagination: appdto.Pagination{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	}
	if req.Status != nil && *req.Status != "" {
		status := agent.SessionStatus(*req.Status)
		query.Status = &status
	}

	result, err := h.h.ListAgentSessions.Handle(c.Request().Context(), query)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.AgentSessionListResponse{
		Items: dto.AgentSessionsToResponse(result.Items),
		Total: result.Total,
	})
}

func (h *agentHandler) ListSessionEvents(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid agent session id")
	}

	var req dto.AgentSessionEventListQuery
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid query parameters")
	}

	query := appdto.ListAgentSessionEventsQuery{
		SessionID: id,
		Pagination: appdto.Pagination{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	}
	if req.Source != nil {
		query.Source = *req.Source
	}
	if req.NodeKey != nil {
		query.NodeKey = *req.NodeKey
	}

	result, err := h.h.ListAgentSessionEvents.Handle(c.Request().Context(), query)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.AgentSessionEventListResponse{
		Items: dto.AgentSessionEventsToResponse(result.Items),
		Total: result.Total,
	})
}
