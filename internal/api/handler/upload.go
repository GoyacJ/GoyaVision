package handler

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"goyavision/config"
	"goyavision/internal/api/dto"
	"goyavision/internal/api/middleware"
	"goyavision/internal/app"
	"goyavision/internal/domain"
	"goyavision/pkg/storage"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterUpload(g *echo.Group, d Deps) {
	h := uploadHandler{
		assetSvc:    app.NewMediaAssetService(d.Repo),
		fileSvc:     app.NewFileService(d.Repo, d.MinIOClient),
		minioClient: d.MinIOClient,
		cfg:         d.Cfg,
	}
	g.POST("/upload", h.Upload)
}

type uploadHandler struct {
	assetSvc    *app.MediaAssetService
	fileSvc     *app.FileService
	minioClient *storage.MinIOClient
	cfg         *config.Config
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

	createReq := &app.CreateMediaAssetRequest{
		Type:       domain.AssetType(assetType),
		SourceType: domain.AssetSourceUpload,
		Name:       name,
		Path:       uploadedFile.Path,
		Size:       file.Size,
		Format:     format,
		Status:     domain.AssetStatusReady,
		Tags:       tags,
	}

	asset, err := h.assetSvc.Create(c.Request().Context(), createReq)
	if err != nil {
		_ = h.fileSvc.DeleteFile(c.Request().Context(), uploadedFile.ID)
		return err
	}

	return c.JSON(201, dto.AssetToResponse(asset, h.cfg.MinIO.Endpoint, h.cfg.MinIO.BucketName, h.cfg.MinIO.UseSSL))
}
