package app

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"goyavision/internal/domain"
	"goyavision/internal/port"
	"goyavision/pkg/storage"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// CreateFileRequest 创建文件请求
type CreateFileRequest struct {
	Name        string
	OriginalName string
	Path        string
	Size        int64
	MimeType    string
	Type        domain.FileType
	Extension   string
	Hash        string
	UploaderID  *uuid.UUID
	Metadata    map[string]interface{}
}

// UpdateFileRequest 更新文件请求
type UpdateFileRequest struct {
	Name     *string
	Status   *domain.FileStatus
	Metadata map[string]interface{}
}

// ListFilesRequest 列出文件请求
type ListFilesRequest struct {
	Type       *domain.FileType
	Status     *domain.FileStatus
	UploaderID *uuid.UUID
	Search     string
	From       *time.Time
	To         *time.Time
	Limit      int
	Offset     int
}

type FileService struct {
	repo       port.Repository
	minioClient *storage.MinIOClient
}

func NewFileService(repo port.Repository, minioClient *storage.MinIOClient) *FileService {
	return &FileService{
		repo:       repo,
		minioClient: minioClient,
	}
}

// UploadFile 上传文件
func (s *FileService) UploadFile(ctx context.Context, reader io.Reader, filename string, size int64, uploaderID *uuid.UUID) (*domain.File, error) {
	// 确定文件类型
	ext := strings.ToLower(filepath.Ext(filename))
	mimeType := mime.TypeByExtension(ext)
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	fileType := s.detectFileType(mimeType, ext)

	// 使用 TeeReader 同时计算哈希和准备上传
	hash := md5.New()
	teeReader := io.TeeReader(reader, hash)

	// 读取数据到缓冲区（用于上传）
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, teeReader); err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// 计算哈希
	hashBytes := hash.Sum(nil)
	fileHash := hex.EncodeToString(hashBytes)

	// 检查文件是否已存在
	existingFile, err := s.repo.GetFileByHash(ctx, fileHash)
	if err == nil && existingFile != nil {
		return existingFile, nil
	}

	// 生成存储路径
	objectName := fmt.Sprintf("files/%s/%s%s", fileType, uuid.New().String(), ext)

	// 上传到 MinIO
	_, err = s.minioClient.Upload(ctx, objectName, &buf, int64(buf.Len()), mimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to storage: %w", err)
	}

	// 创建文件记录
	file := &domain.File{
		Name:        strings.TrimSuffix(filename, ext),
		OriginalName: filename,
		Path:        objectName,
		Size:        size,
		MimeType:    mimeType,
		Type:        fileType,
		Extension:   strings.TrimPrefix(ext, "."),
		Status:      domain.FileStatusCompleted,
		Hash:        fileHash,
		UploaderID:  uploaderID,
	}

	if err := s.repo.CreateFile(ctx, file); err != nil {
		_ = s.minioClient.Delete(ctx, objectName)
		return nil, err
	}

	return file, nil
}

// GetFile 获取文件
func (s *FileService) GetFile(ctx context.Context, id uuid.UUID) (*domain.File, error) {
	file, err := s.repo.GetFile(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("file not found")
		}
		return nil, err
	}
	return file, nil
}

// ListFiles 列出文件
func (s *FileService) ListFiles(ctx context.Context, req *ListFilesRequest) ([]*domain.File, int64, error) {
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 1000 {
		req.Limit = 1000
	}

	filter := domain.FileFilter{
		Type:       req.Type,
		Status:     req.Status,
		UploaderID: req.UploaderID,
		Search:     req.Search,
		From:       req.From,
		To:         req.To,
		Limit:      req.Limit,
		Offset:     req.Offset,
	}

	return s.repo.ListFiles(ctx, filter)
}

// UpdateFile 更新文件
func (s *FileService) UpdateFile(ctx context.Context, id uuid.UUID, req *UpdateFileRequest) (*domain.File, error) {
	file, err := s.GetFile(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		file.Name = *req.Name
	}
	if req.Status != nil {
		file.Status = *req.Status
	}
	if req.Metadata != nil {
		metadataBytes, err := json.Marshal(req.Metadata)
		if err != nil {
			return nil, errors.New("failed to marshal metadata: " + err.Error())
		}
		file.Metadata = datatypes.JSON(metadataBytes)
	}

	if err := s.repo.UpdateFile(ctx, file); err != nil {
		return nil, err
	}

	return file, nil
}

// DeleteFile 删除文件
func (s *FileService) DeleteFile(ctx context.Context, id uuid.UUID) error {
	file, err := s.GetFile(ctx, id)
	if err != nil {
		return err
	}

	// 删除存储中的文件
	if err := s.minioClient.Delete(ctx, file.Path); err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	// 删除数据库记录
	return s.repo.DeleteFile(ctx, id)
}

// detectFileType 根据 MIME 类型和扩展名检测文件类型
func (s *FileService) detectFileType(mimeType, ext string) domain.FileType {
	mimeLower := strings.ToLower(mimeType)
	extLower := strings.ToLower(ext)

	if strings.HasPrefix(mimeLower, "image/") {
		return domain.FileTypeImage
	}
	if strings.HasPrefix(mimeLower, "video/") {
		return domain.FileTypeVideo
	}
	if strings.HasPrefix(mimeLower, "audio/") {
		return domain.FileTypeAudio
	}
	if strings.Contains(mimeLower, "pdf") || strings.Contains(mimeLower, "document") ||
		strings.Contains(mimeLower, "text") || strings.Contains(mimeLower, "word") ||
		strings.Contains(mimeLower, "excel") || strings.Contains(mimeLower, "powerpoint") {
		return domain.FileTypeDocument
	}
	if extLower == ".zip" || extLower == ".rar" || extLower == ".7z" ||
		extLower == ".tar" || extLower == ".gz" {
		return domain.FileTypeArchive
	}

	return domain.FileTypeOther
}
