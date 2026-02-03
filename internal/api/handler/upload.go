package handler

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"goyavision/config"
	"goyavision/internal/api/dto"
	"goyavision/internal/app"
	"goyavision/internal/domain"
	"goyavision/pkg/storage"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterUpload(g *echo.Group, d Deps) {
	h := uploadHandler{
		svc:         app.NewMediaAssetService(d.Repo),
		minioClient: d.MinIOClient,
		cfg:         d.Cfg,
	}
	g.POST("/upload", h.Upload)
}

type uploadHandler struct {
	svc         *app.MediaAssetService
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

	ext := filepath.Ext(file.Filename)
	objectName := fmt.Sprintf("%s/%s%s", assetType, uuid.New().String(), ext)

	contentType := file.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = h.minioClient.Upload(c.Request().Context(), objectName, src, file.Size, contentType)
	if err != nil {
		return c.JSON(500, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: fmt.Sprintf("failed to upload file: %v", err),
		})
	}

	format := ext
	if len(format) > 0 && format[0] == '.' {
		format = format[1:]
	}

	createReq := &app.CreateMediaAssetRequest{
		Type:       domain.AssetType(assetType),
		SourceType: domain.AssetSourceUpload,
		Name:       name,
		Path:       objectName,
		Size:       file.Size,
		Format:     format,
		Status:     domain.AssetStatusReady,
		Tags:       tags,
	}

	asset, err := h.svc.Create(c.Request().Context(), createReq)
	if err != nil {
		_ = h.minioClient.Delete(c.Request().Context(), objectName)
		return err
	}

	return c.JSON(201, dto.AssetToResponse(asset, h.cfg.MinIO.Endpoint, h.cfg.MinIO.BucketName, h.cfg.MinIO.UseSSL))
}
