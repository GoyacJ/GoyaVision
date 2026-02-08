package handler

import (
	"time"

	"goyavision/config"
	"goyavision/internal/api/dto"
	"goyavision/internal/api/middleware"
	"goyavision/internal/app"
	appport "goyavision/internal/app/port"
	"goyavision/internal/domain/storage"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterFile(g *echo.Group, h *Handlers) {
	svc := app.NewFileService(h.Repo, h.FileStorage)
	fh := fileHandler{
		svc:    svc,
		cfg:    h.Cfg,
		urlCfg: h.StorageURLConfig,
	}
	g.POST("/files", fh.Upload)
	g.GET("/files", fh.List)
	g.GET("/files/:id", fh.Get)
	g.PUT("/files/:id", fh.Update)
	g.DELETE("/files/:id", fh.Delete)
	g.GET("/files/:id/download", fh.Download)
}

type fileHandler struct {
	svc    *app.FileService
	cfg    *config.Config
	urlCfg appport.StorageURLConfig
}

// Upload 上传文件
func (h *fileHandler) Upload(c echo.Context) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "no file provided",
		})
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

	uploadedFile, err := h.svc.UploadFile(c.Request().Context(), src, file.Filename, file.Size, uploaderID)
	if err != nil {
		return c.JSON(500, dto.ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
	}

	return c.JSON(201, dto.FileToResponse(uploadedFile, h.urlCfg.Endpoint, h.urlCfg.BucketName, h.urlCfg.UseSSL))
}

// List 列出文件
func (h *fileHandler) List(c echo.Context) error {
	var query dto.FileListQuery
	if err := c.Bind(&query); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid query parameters",
		})
	}

	req := &app.ListFilesRequest{
		Limit:  query.Limit,
		Offset: query.Offset,
		Search: query.Search,
	}

	if query.Type != nil {
		t := storage.FileType(*query.Type)
		req.Type = &t
	}
	if query.Status != nil {
		s := storage.FileStatus(*query.Status)
		req.Status = &s
	}
	if query.UploaderID != nil {
		req.UploaderID = query.UploaderID
	}
	if query.From != nil {
		t := time.Unix(*query.From, 0)
		req.From = &t
	}
	if query.To != nil {
		t := time.Unix(*query.To, 0)
		req.To = &t
	}

	files, total, err := h.svc.ListFiles(c.Request().Context(), req)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.FileListResponse{
		Items: dto.FilesToResponse(files, h.urlCfg.Endpoint, h.urlCfg.BucketName, h.urlCfg.UseSSL),
		Total: total,
	})
}

// Get 获取文件
func (h *fileHandler) Get(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid file id",
		})
	}

	file, err := h.svc.GetFile(c.Request().Context(), id)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.FileToResponse(file, h.urlCfg.Endpoint, h.urlCfg.BucketName, h.urlCfg.UseSSL))
}

// Update 更新文件
func (h *fileHandler) Update(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid file id",
		})
	}

	var req dto.FileUpdateReq
	if err := c.Bind(&req); err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid request body",
		})
	}

	updateReq := &app.UpdateFileRequest{
		Name:     req.Name,
		Metadata: req.Metadata,
	}

	if req.Status != nil {
		s := storage.FileStatus(*req.Status)
		updateReq.Status = &s
	}

	file, err := h.svc.UpdateFile(c.Request().Context(), id, updateReq)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.FileToResponse(file, h.urlCfg.Endpoint, h.urlCfg.BucketName, h.urlCfg.UseSSL))
}

// Delete 删除文件
func (h *fileHandler) Delete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid file id",
		})
	}

	if err := h.svc.DeleteFile(c.Request().Context(), id); err != nil {
		return err
	}

	return c.NoContent(204)
}

// Download 下载文件（返回重定向到文件 URL）
func (h *fileHandler) Download(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(400, dto.ErrorResponse{
			Error:   "Bad Request",
			Message: "invalid file id",
		})
	}

	file, err := h.svc.GetFile(c.Request().Context(), id)
	if err != nil {
		return err
	}

	var url string
	if h.urlCfg.PublicBase != "" {
		if h.urlCfg.BucketName != "" {
			url = h.urlCfg.PublicBase + "/" + h.urlCfg.BucketName + "/" + file.Path
		} else {
			url = h.urlCfg.PublicBase + "/" + file.Path
		}
	} else {
		protocol := "http"
		if h.urlCfg.UseSSL {
			protocol = "https"
		}
		url = protocol + "://" + h.urlCfg.Endpoint + "/" + h.urlCfg.BucketName + "/" + file.Path
	}
	return c.Redirect(302, url)
}
