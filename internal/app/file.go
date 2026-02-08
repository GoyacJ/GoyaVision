package app

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"goyavision/internal/domain/storage"
	appport "goyavision/internal/app/port"
	"goyavision/internal/port"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateFileRequest 创建文件请求
type CreateFileRequest struct {
	Name         string
	OriginalName string
	Path         string
	Size         int64
	MimeType     string
	Type         storage.FileType
	Extension    string
	Hash         string
	UploaderID   *uuid.UUID
	Metadata     map[string]interface{}
}

// UpdateFileRequest 更新文件请求
type UpdateFileRequest struct {
	Name     *string
	Status   *storage.FileStatus
	Metadata map[string]interface{}
}

// ListFilesRequest 列出文件请求
type ListFilesRequest struct {
	Type       *storage.FileType
	Status     *storage.FileStatus
	UploaderID *uuid.UUID
	Search     string
	From       *time.Time
	To         *time.Time
	Limit      int
	Offset     int
}

type FileService struct {
	repo        port.Repository
	fileStorage appport.FileStorage
}

func NewFileService(repo port.Repository, fileStorage appport.FileStorage) *FileService {
	return &FileService{
		repo:        repo,
		fileStorage: fileStorage,
	}
}

// UploadFile 上传文件
func (s *FileService) UploadFile(ctx context.Context, reader io.Reader, filename string, size int64, uploaderID *uuid.UUID) (*storage.File, error) {
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

	_, err = s.fileStorage.Upload(ctx, objectName, &buf, int64(buf.Len()), mimeType)
	if err != nil {
		return nil, fmt.Errorf("failed to upload file to storage: %w", err)
	}

	// 创建文件记录
	file := &storage.File{
		Name:         strings.TrimSuffix(filename, ext),
		OriginalName: filename,
		Path:         objectName,
		Size:         size,
		MimeType:     mimeType,
		Type:         fileType,
		Extension:    strings.TrimPrefix(ext, "."),
		Status:       storage.FileStatusCompleted,
		Hash:         fileHash,
		UploaderID:   uploaderID,
	}

	if err := s.repo.CreateFile(ctx, file); err != nil {
		_ = s.fileStorage.Delete(ctx, objectName)
		return nil, err
	}

	return file, nil
}

// GetFile 获取文件
func (s *FileService) GetFile(ctx context.Context, id uuid.UUID) (*storage.File, error) {
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
func (s *FileService) ListFiles(ctx context.Context, req *ListFilesRequest) ([]*storage.File, int64, error) {
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 1000 {
		req.Limit = 1000
	}

	filter := storage.FileFilter{
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
func (s *FileService) UpdateFile(ctx context.Context, id uuid.UUID, req *UpdateFileRequest) (*storage.File, error) {
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
		file.Metadata = req.Metadata
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

	if err := s.fileStorage.Delete(ctx, file.Path); err != nil {
		return fmt.Errorf("failed to delete file from storage: %w", err)
	}

	// 删除数据库记录
	return s.repo.DeleteFile(ctx, id)
}

// detectFileType 根据 MIME 类型和扩展名检测文件类型
func (s *FileService) detectFileType(mimeType, ext string) storage.FileType {
	mimeLower := strings.ToLower(mimeType)
	extLower := strings.ToLower(ext)

	if strings.HasPrefix(mimeLower, "image/") {
		return storage.FileTypeImage
	}
	if strings.HasPrefix(mimeLower, "video/") {
		return storage.FileTypeVideo
	}
	if strings.HasPrefix(mimeLower, "audio/") {
		return storage.FileTypeAudio
	}
	if strings.Contains(mimeLower, "pdf") || strings.Contains(mimeLower, "document") ||
		strings.Contains(mimeLower, "text") || strings.Contains(mimeLower, "word") ||
		strings.Contains(mimeLower, "excel") || strings.Contains(mimeLower, "powerpoint") {
		return storage.FileTypeDocument
	}
	if extLower == ".zip" || extLower == ".rar" || extLower == ".7z" ||
		extLower == ".tar" || extLower == ".gz" {
		return storage.FileTypeArchive
	}

	return storage.FileTypeOther
}
