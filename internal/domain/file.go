package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

// FileType 文件类型
type FileType string

const (
	FileTypeImage FileType = "image"
	FileTypeVideo FileType = "video"
	FileTypeAudio FileType = "audio"
	FileTypeDocument FileType = "document"
	FileTypeArchive FileType = "archive"
	FileTypeOther FileType = "other"
)

// FileStatus 文件状态
type FileStatus string

const (
	FileStatusUploading FileStatus = "uploading"
	FileStatusCompleted FileStatus = "completed"
	FileStatusFailed FileStatus = "failed"
	FileStatusDeleted FileStatus = "deleted"
)

// File 文件实体
type File struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"type:varchar(255);not null;index:idx_files_name"`
	OriginalName string  `gorm:"type:varchar(255);not null"`
	Path        string    `gorm:"type:varchar(1024);not null;index:idx_files_path"`
	Size        int64     `gorm:"not null;default:0;index:idx_files_size"`
	MimeType    string    `gorm:"type:varchar(100);not null"`
	Type        FileType  `gorm:"type:varchar(20);not null;index:idx_files_type"`
	Extension   string    `gorm:"type:varchar(20)"`
	Status      FileStatus `gorm:"type:varchar(20);not null;default:'completed';index:idx_files_status"`
	Hash        string    `gorm:"type:varchar(64);index:idx_files_hash"`
	UploaderID  *uuid.UUID `gorm:"type:uuid;index:idx_files_uploader"`
	Metadata    datatypes.JSON `gorm:"type:jsonb"`
	CreatedAt   time.Time `gorm:"autoCreateTime;index:idx_files_created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `gorm:"index:idx_files_deleted_at"`
}

func (File) TableName() string { return "files" }

// IsImage 判断是否为图片
func (f *File) IsImage() bool {
	return f.Type == FileTypeImage
}

// IsVideo 判断是否为视频
func (f *File) IsVideo() bool {
	return f.Type == FileTypeVideo
}

// IsAudio 判断是否为音频
func (f *File) IsAudio() bool {
	return f.Type == FileTypeAudio
}

// IsDocument 判断是否为文档
func (f *File) IsDocument() bool {
	return f.Type == FileTypeDocument
}

// IsCompleted 判断是否上传完成
func (f *File) IsCompleted() bool {
	return f.Status == FileStatusCompleted
}

// IsFailed 判断是否上传失败
func (f *File) IsFailed() bool {
	return f.Status == FileStatusFailed
}

// FileFilter 文件过滤器
type FileFilter struct {
	Type       *FileType
	Status     *FileStatus
	UploaderID *uuid.UUID
	Search     string
	From       *time.Time
	To         *time.Time
	Limit      int
	Offset     int
}
