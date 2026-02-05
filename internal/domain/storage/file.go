package storage

import (
	"time"

	"github.com/google/uuid"
)

type FileType string

const (
	FileTypeImage    FileType = "image"
	FileTypeVideo    FileType = "video"
	FileTypeAudio    FileType = "audio"
	FileTypeDocument FileType = "document"
	FileTypeArchive  FileType = "archive"
	FileTypeOther    FileType = "other"
)

type FileStatus string

const (
	FileStatusUploading FileStatus = "uploading"
	FileStatusCompleted FileStatus = "completed"
	FileStatusFailed    FileStatus = "failed"
	FileStatusDeleted   FileStatus = "deleted"
)

type File struct {
	ID           uuid.UUID
	Name         string
	OriginalName string
	Path         string
	Size         int64
	MimeType     string
	Type         FileType
	Extension    string
	Status       FileStatus
	Hash         string
	UploaderID   *uuid.UUID
	Metadata     map[string]interface{}
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}

func (f *File) IsImage() bool {
	return f.Type == FileTypeImage
}

func (f *File) IsVideo() bool {
	return f.Type == FileTypeVideo
}

func (f *File) IsAudio() bool {
	return f.Type == FileTypeAudio
}

func (f *File) IsDocument() bool {
	return f.Type == FileTypeDocument
}

func (f *File) IsCompleted() bool {
	return f.Status == FileStatusCompleted
}

func (f *File) IsFailed() bool {
	return f.Status == FileStatusFailed
}

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
