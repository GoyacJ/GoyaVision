package handler

import (
	"net/http"

	"goyavision/internal/api/dto"
	"goyavision/internal/app"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterMenu(g *echo.Group, d Deps) {
	svc := app.NewMenuService(d.Repo)
	permSvc := app.NewPermissionService(d.Repo)
	h := menuHandler{svc: svc, permSvc: permSvc}

	g.GET("/menus", h.List)
	g.GET("/menus/tree", h.ListTree)
	g.POST("/menus", h.Create)
	g.GET("/menus/:id", h.Get)
	g.PUT("/menus/:id", h.Update)
	g.DELETE("/menus/:id", h.Delete)

	g.GET("/permissions", h.ListPermissions)
}

type menuHandler struct {
	svc     *app.MenuService
	permSvc *app.PermissionService
}

func (h *menuHandler) List(c echo.Context) error {
	var query dto.MenuListQuery
	if err := c.Bind(&query); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid query parameters",
		})
	}

	menus, err := h.svc.List(c.Request().Context(), query.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.MenusToResponse(menus))
}

func (h *menuHandler) ListTree(c echo.Context) error {
	var query dto.MenuListQuery
	if err := c.Bind(&query); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid query parameters",
		})
	}

	menus, err := h.svc.ListTree(c.Request().Context(), query.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.MenusToResponse(menus))
}

func (h *menuHandler) Create(c echo.Context) error {
	var req dto.MenuCreateRequest
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

	if req.Type < 1 || req.Type > 3 {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "type must be 1 (directory), 2 (menu), or 3 (button)",
		})
	}

	menu, err := h.svc.Create(c.Request().Context(), &app.CreateMenuRequest{
		ParentID:   req.ParentID,
		Code:       req.Code,
		Name:       req.Name,
		Type:       req.Type,
		Path:       req.Path,
		Icon:       req.Icon,
		Component:  req.Component,
		Permission: req.Permission,
		Sort:       req.Sort,
		Visible:    req.Visible,
		Status:     req.Status,
	})
	if err != nil {
		if err == app.ErrMenuCodeExists {
			return c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "Conflict",
				Message: "menu code already exists",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, dto.MenuToResponse(menu))
}

func (h *menuHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid menu id",
		})
	}

	menu, err := h.svc.Get(c.Request().Context(), id)
	if err != nil {
		if err == app.ErrMenuNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "menu not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.MenuToResponse(menu))
}

func (h *menuHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid menu id",
		})
	}

	var req dto.MenuUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	menu, err := h.svc.Update(c.Request().Context(), id, &app.UpdateMenuRequest{
		ParentID:   req.ParentID,
		Name:       req.Name,
		Type:       req.Type,
		Path:       req.Path,
		Icon:       req.Icon,
		Component:  req.Component,
		Permission: req.Permission,
		Sort:       req.Sort,
		Visible:    req.Visible,
		Status:     req.Status,
	})
	if err != nil {
		if err == app.ErrMenuNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "menu not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.MenuToResponse(menu))
}

func (h *menuHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid menu id",
		})
	}

	if err := h.svc.Delete(c.Request().Context(), id); err != nil {
		if err == app.ErrMenuNotFound {
			return c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "Not Found",
				Message: "menu not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.NoContent(http.StatusNoContent)
}

func (h *menuHandler) ListPermissions(c echo.Context) error {
	permissions, err := h.permSvc.List(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.JSON(http.StatusOK, dto.PermissionsToResponse(permissions))
}
