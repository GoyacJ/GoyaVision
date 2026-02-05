package handler

import (
	"net/http"

	"goyavision/internal/api/dto"
	"goyavision/internal/app"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterRole(g *echo.Group, h *Handlers) {
	svc := app.NewRoleService(h.Repo)
	rh := roleHandler{svc: svc}

	g.GET("/roles", rh.List)
	g.POST("/roles", rh.Create)
	g.GET("/roles/:id", rh.Get)
	g.PUT("/roles/:id", rh.Update)
	g.DELETE("/roles/:id", rh.Delete)
}

type roleHandler struct {
	svc *app.RoleService
}

func (h *roleHandler) List(c echo.Context) error {
	var query dto.RoleListQuery
	if err := c.Bind(&query); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid query parameters",
		})
	}

	roles, err := h.svc.List(c.Request().Context(), query.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.RolesToResponse(roles))
}

func (h *roleHandler) Create(c echo.Context) error {
	var req dto.RoleCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	if req.Code == "" || req.Name == "" {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "code and name are required",
		})
	}

	role, err := h.svc.Create(c.Request().Context(), &app.CreateRoleRequest{
		Code:          req.Code,
		Name:          req.Name,
		Description:   req.Description,
		Status:        req.Status,
		PermissionIDs: req.PermissionIDs,
		MenuIDs:       req.MenuIDs,
	})
	if err != nil {
		if err == app.ErrRoleCodeExists {
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "Conflict",
				Message: "role code already exists",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.RoleToResponse(role))
}

func (h *roleHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid role id",
		})
	}

	role, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		if err == app.ErrRoleNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "role not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.RoleToResponse(role))
}

func (h *roleHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid role id",
		})
	}

	var req dto.RoleUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	role, err := h.svc.Update(c.Request().Context(), id, &app.UpdateRoleRequest{
		Name:          req.Name,
		Description:   req.Description,
		Status:        req.Status,
		PermissionIDs: req.PermissionIDs,
		MenuIDs:       req.MenuIDs,
	})
	if err != nil {
		if err == app.ErrRoleNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "role not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.RoleToResponse(role))
}

func (h *roleHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid role id",
		})
	}

	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		if err == app.ErrRoleNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "role not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}
