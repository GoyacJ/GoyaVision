package mapper

import (
	"encoding/json"

	"goyavision/internal/domain/storage"
	"goyavision/internal/infra/persistence/model"

	"gorm.io/datatypes"
)

func FileToModel(f *storage.File) *model.FileModel {
	m := &model.FileModel{
		ID:           f.ID,
		Name:         f.Name,
		OriginalName: f.OriginalName,
		Path:         f.Path,
		Size:         f.Size,
		MimeType:     f.MimeType,
		Type:         string(f.Type),
		Extension:    f.Extension,
		Status:       string(f.Status),
		Hash:         f.Hash,
		UploaderID:   f.UploaderID,
		CreatedAt:    f.CreatedAt,
		UpdatedAt:    f.UpdatedAt,
		DeletedAt:    f.DeletedAt,
	}
	if f.Metadata != nil {
		data, _ := json.Marshal(f.Metadata)
		m.Metadata = datatypes.JSON(data)
	}
	return m
}

func FileToDomain(m *model.FileModel) *storage.File {
	f := &storage.File{
		ID:           m.ID,
		Name:         m.Name,
		OriginalName: m.OriginalName,
		Path:         m.Path,
		Size:         m.Size,
		MimeType:     m.MimeType,
		Type:         storage.FileType(m.Type),
		Extension:    m.Extension,
		Status:       storage.FileStatus(m.Status),
		Hash:         m.Hash,
		UploaderID:   m.UploaderID,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		DeletedAt:    m.DeletedAt,
	}
	if m.Metadata != nil {
		_ = json.Unmarshal(m.Metadata, &f.Metadata)
	}
	return f
}
