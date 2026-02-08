package handler

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"strconv"

	"goyavision/config"
	"goyavision/internal/api/dto"
	"goyavision/internal/api/middleware"
	"goyavision/internal/app"
	appdto "goyavision/internal/app/dto"
	appport "goyavision/internal/app/port"
	"goyavision/internal/domain/media"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterUpload(g *echo.Group, h *Handlers) {
	uh := uploadHandler{
		h:       h,
		fileSvc: app.NewFileService(h.Repo, h.FileStorage),
		cfg:     h.Cfg,
		urlCfg:  h.StorageURLConfig,
	}
	g.POST("/upload", uh.Upload)
}

type uploadHandler struct {
	h       *Handlers
	fileSvc *app.FileService
	cfg     *config.Config
	urlCfg  appport.StorageURLConfig
}

func (h *uploadHandler) Upload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "no file provided",
		})
	}

	assetType := c.FormValue("type")
	if assetType == "" {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "type is required",
		})
	}

	name := c.FormValue("name")
	if name == "" {
		name = file.Filename
	}

	// 处理标签
	tagsStr := c.FormValue("tags")
	var tags []string
	if tagsStr != "" {
		if err := json.Unmarshal([]byte(tagsStr), &tags); err != nil {
			tags = []string{} // 解析失败时使用空数组
		}
	}

	// 处理可见性
	visibility := media.VisibilityPrivate
	if vStr := c.FormValue("visibility"); vStr != "" {
		if v, err := strconv.Atoi(vStr); err == nil {
			visibility = media.Visibility(v)
		}
	}

	var visibleRoleIDs []string
	if visibility == media.VisibilityRole {
		// 自动填充当前用户的角色
		if ids, ok := c.Get(middleware.ContextKeyRoleIDs).([]uuid.UUID); ok {
			for _, id := range ids {
				visibleRoleIDs = append(visibleRoleIDs, id.String())
			}
		}
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(500, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: "failed to open file",
		})
	}
	defer src.Close()

	// 获取上传者 ID（从 JWT token 中）
	var uploaderID *uuid.UUID
	if id, ok := middleware.GetUserID(c); ok {
		uploaderID = &id
	}

	// 使用文件服务上传文件
	uploadedFile, err := h.fileSvc.UploadFile(c.Request().Context(), src, file.Filename, file.Size, uploaderID)
	if err != nil {
		return c.JSON(500, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: fmt.Sprintf("failed to upload file: %v", err),
		})
	}

	// 创建媒体资产记录
	ext := filepath.Ext(file.Filename)
	format := ext
	if len(format) > 0 && format[0] == '.' {
		format = format[1:]
	}

	cmd := appdto.CreateAssetCommand{
		Type:           media.AssetType(assetType),
		SourceType:     media.AssetSourceUpload,
		Name:           name,
		Path:           uploadedFile.Path,
		Size:           file.Size,
		Format:         format,
		Status:         media.AssetStatusReady,
		Tags:           tags,
		Visibility:     visibility,
		VisibleRoleIDs: visibleRoleIDs,
	}

	asset, err := h.h.CreateAsset.Handle(c.Request().Context(), cmd)
	if err != nil {
		_ = h.fileSvc.DeleteFile(c.Request().Context(), uploadedFile.ID)
		return err
	}

	return c.JSON(201, dto.AssetToResponse(asset, h.urlCfg.Endpoint, h.urlCfg.BucketName, h.urlCfg.PublicBase, h.urlCfg.UseSSL))
}
