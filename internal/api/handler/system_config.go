package handler

import (
	"encoding/json"
	"net/http"

	"goyavision/internal/api/dto"
	authmiddleware "goyavision/internal/api/middleware"
	"goyavision/internal/domain/system"
	"goyavision/internal/port"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type SystemConfigHandler struct {
	repo port.Repository
}

func NewSystemConfigHandler(repo port.Repository) *SystemConfigHandler {
	return &SystemConfigHandler{repo: repo}
}

// RegisterSystemConfig registers system config routes
func RegisterSystemConfig(public *echo.Group, protected *echo.Group, h *Handlers) {
	handler := NewSystemConfigHandler(h.Repo)
	
	// Public Routes
	public.GET("/public/config", handler.GetPublicConfig)
	
	// Protected Routes
	protected.PUT("/system/config", handler.UpdateConfig)
}

// GetPublicConfigResponse response structure
type GetPublicConfigResponse struct {
	HomePath    string              `json:"home_path"`
	PublicMenus []*dto.MenuResponse `json:"public_menus"`
}

// GetPublicConfig gets public system configuration
func (h *SystemConfigHandler) GetPublicConfig(c echo.Context) error {
	ctx := c.Request().Context()

	// Get Home Path
	homePathConfig, err := h.repo.GetSystemConfig(ctx, system.ConfigKeyHomePath)
	if err != nil {
		return err
	}
	
	homePath := "/assets" // default
	if homePathConfig != nil {
		var val string
		if err := json.Unmarshal(homePathConfig.Value, &val); err == nil && val != "" {
			homePath = val
		}
	}

	// Get Public Menus
	publicMenusConfig, err := h.repo.GetSystemConfig(ctx, system.ConfigKeyPublicMenus)
	if err != nil {
		return err
	}

	var publicMenuIDs []string
	if publicMenusConfig != nil {
		_ = json.Unmarshal(publicMenusConfig.Value, &publicMenuIDs)
	}

	var publicMenus []*dto.MenuResponse
	if len(publicMenuIDs) > 0 {
		var ids []uuid.UUID
		for _, idStr := range publicMenuIDs {
			if id, err := uuid.Parse(idStr); err == nil {
				ids = append(ids, id)
			}
		}

		if len(ids) > 0 {
			allMenus, err := h.repo.ListMenus(ctx, nil)
			if err != nil {
				return err
			}
			
			idMap := make(map[uuid.UUID]bool)
			for _, id := range ids {
				idMap[id] = true
			}
			
			// Use map to reconstruct tree for selected nodes
			// We only include nodes that are in ids list.
			// If a node's parent is NOT in list, it becomes a root.
			
			// 1. Convert to DTOs and put in map
			dtoMap := make(map[uuid.UUID]*dto.MenuResponse)
			for _, m := range allMenus {
				if idMap[m.ID] {
					dtoMap[m.ID] = dto.MenuToResponse(m)
					// Clear children from DTO as MenuToResponse might have copied them if domain model had them
					dtoMap[m.ID].Children = nil
				}
			}
			
			// 2. Build tree
			var roots []*dto.MenuResponse
			for _, m := range allMenus {
				if !idMap[m.ID] {
					continue
				}
				
				node := dtoMap[m.ID]
				if m.ParentID != nil && idMap[*m.ParentID] {
					// Parent is also public, attach to parent
					if parent, ok := dtoMap[*m.ParentID]; ok {
						parent.Children = append(parent.Children, node)
					} else {
						// Should not happen if idMap logic is correct
						roots = append(roots, node)
					}
				} else {
					// Parent is not public or root
					roots = append(roots, node)
				}
			}
			
			publicMenus = roots
		}
	}

	return c.JSON(http.StatusOK, GetPublicConfigResponse{
		HomePath:    homePath,
		PublicMenus: publicMenus,
	})
}

// UpdateConfig updates system configuration
func (h *SystemConfigHandler) UpdateConfig(c echo.Context) error {
	if !authmiddleware.HasPermission(c, "system:config:update") {
		return c.JSON(http.StatusForbidden, dto.ErrorResponse{
			Error:   "Forbidden",
			Message: "无修改权限",
		})
	}

	var req map[string]interface{}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request body")
	}

	ctx := c.Request().Context()

	for k, v := range req {
		// Validate keys
		if k != system.ConfigKeyHomePath && k != system.ConfigKeyPublicMenus {
			continue 
		}

		valBytes, err := json.Marshal(v)
		if err != nil {
			continue
		}

		config := &system.SystemConfig{
			Key:   k,
			Value: json.RawMessage(valBytes),
		}
		if err := h.repo.SaveSystemConfig(ctx, config); err != nil {
			return err
		}
	}

	return c.NoContent(http.StatusOK)
}
