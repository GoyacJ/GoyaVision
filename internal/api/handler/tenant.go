package handler

import (
	"net/http"

	"goyavision/internal/api/dto"
	"goyavision/internal/infra/persistence/model"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterTenant(g *echo.Group, h *Handlers) {
	handler := &tenantHandler{h: h}
	g.GET("/tenants", handler.List)
	g.POST("/tenants", handler.Create)
	g.GET("/tenants/:id", handler.Get)
	g.PUT("/tenants/:id", handler.Update)
	g.DELETE("/tenants/:id", handler.Delete)
}

type tenantHandler struct {
	h *Handlers
}

func (h *tenantHandler) List(c echo.Context) error {
	var models []model.TenantModel
	if err := h.h.DB.Find(&models).Error; err != nil {
		return err
	}

	resp := make([]dto.TenantResponse, len(models))
	for i, m := range models {
		resp[i] = dto.TenantResponse{
			ID:        m.ID,
			Name:      m.Name,
			Code:      m.Code,
			Status:    m.Status,
			CreatedAt: m.CreatedAt,
			UpdatedAt: m.UpdatedAt,
		}
	}
	return c.JSON(http.StatusOK, resp)
}

func (h *tenantHandler) Create(c echo.Context) error {
	var req dto.TenantCreateReq
	if err := c.Bind(&req); err != nil {
		return err
	}

	m := model.TenantModel{
		ID:     uuid.New(),
		Name:   req.Name,
		Code:   req.Code,
		Status: req.Status,
	}
	if m.Status == 0 {
		m.Status = 1
	}

	if err := h.h.DB.Create(&m).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, dto.TenantResponse{
		ID:        m.ID,
		Name:      m.Name,
		Code:      m.Code,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	})
}

func (h *tenantHandler) Get(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}

	var m model.TenantModel
	if err := h.h.DB.First(&m, id).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.TenantResponse{
		ID:        m.ID,
		Name:      m.Name,
		Code:      m.Code,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	})
}

func (h *tenantHandler) Update(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}

	var req dto.TenantUpdateReq
	if err := c.Bind(&req); err != nil {
		return err
	}

	var m model.TenantModel
	if err := h.h.DB.First(&m, id).Error; err != nil {
		return err
	}

	if req.Name != nil {
		m.Name = *req.Name
	}
	if req.Status != nil {
		m.Status = *req.Status
	}

	if err := h.h.DB.Save(&m).Error; err != nil {
		return err
	}

	return c.JSON(http.StatusOK, dto.TenantResponse{
		ID:        m.ID,
		Name:      m.Name,
		Code:      m.Code,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	})
}

func (h *tenantHandler) Delete(c echo.Context) error {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		return err
	}

	if err := h.h.DB.Delete(&model.TenantModel{}, id).Error; err != nil {
		return err
	}

	return c.NoContent(http.StatusNoContent)
}
