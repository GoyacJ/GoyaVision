package handler

import (
	"time"

	"goyavision/config"
	"goyavision/internal/api/dto"
	"goyavision/internal/api/middleware"
	"goyavision/internal/app"
	"goyavision/internal/domain"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func RegisterFile(g *echo.Group, d Deps) {
	svc := app.NewFileService(d.Repo, d.MinIOClient)
	h := fileHandler{
		svc: svc,
		cfg: d.Cfg,
	}
	g.POST("/files", h.Upload)
	g.GET("/files", h.List)
	g.GET("/files/:id", h.Get)
	g.PUT("/files/:id", h.Update)
	g.DELETE("/files/:id", h.Delete)
	g.GET("/files/:id/download", h.Download)
}

type fileHandler struct {
	svc *app.FileService
	cfg *config.Config
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

	return c.JSON(201, dto.FileToResponse(uploadedFile, h.cfg.MinIO.Endpoint, h.cfg.MinIO.BucketName, h.cfg.MinIO.UseSSL))
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
		t := domain.FileType(*query.Type)
		req.Type = &t
	}
	if query.Status != nil {
		s := domain.FileStatus(*query.Status)
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
		Items: dto.FilesToResponse(files, h.cfg.MinIO.Endpoint, h.cfg.MinIO.BucketName, h.cfg.MinIO.UseSSL),
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

	return c.JSON(200, dto.FileToResponse(file, h.cfg.MinIO.Endpoint, h.cfg.MinIO.BucketName, h.cfg.MinIO.UseSSL))
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
		s := domain.FileStatus(*req.Status)
		updateReq.Status = &s
	}

	file, err := h.svc.UpdateFile(c.Request().Context(), id, updateReq)
	if err != nil {
		return err
	}

	return c.JSON(200, dto.FileToResponse(file, h.cfg.MinIO.Endpoint, h.cfg.MinIO.BucketName, h.cfg.MinIO.UseSSL))
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

	// 生成文件 URL
	protocol := "http"
	if h.cfg.MinIO.UseSSL {
		protocol = "https"
	}
	url := protocol + "://" + h.cfg.MinIO.Endpoint + "/" + h.cfg.MinIO.BucketName + "/" + file.Path

	return c.Redirect(302, url)
}
