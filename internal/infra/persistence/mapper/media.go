package mapper

import (
	"encoding/json"

	"goyavision/internal/domain/media"
	"goyavision/internal/infra/persistence/model"

	"gorm.io/datatypes"
)

func SourceToModel(s *media.Source) *model.MediaSourceModel {
	return &model.MediaSourceModel{
		ID:            s.ID,
		Name:          s.Name,
		PathName:      s.PathName,
		Type:          string(s.Type),
		URL:           s.URL,
		Protocol:      s.Protocol,
		Enabled:       s.Enabled,
		RecordEnabled: s.RecordEnabled,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
}

func SourceToDomain(m *model.MediaSourceModel) *media.Source {
	return &media.Source{
		ID:            m.ID,
		Name:          m.Name,
		PathName:      m.PathName,
		Type:          media.SourceType(m.Type),
		URL:           m.URL,
		Protocol:      m.Protocol,
		Enabled:       m.Enabled,
		RecordEnabled: m.RecordEnabled,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func AssetToModel(a *media.Asset) *model.MediaAssetModel {
	m := &model.MediaAssetModel{
		ID:         a.ID,
		Type:       string(a.Type),
		SourceType: string(a.SourceType),
		SourceID:   a.SourceID,
		ParentID:   a.ParentID,
		Name:       a.Name,
		Path:       a.Path,
		Duration:   a.Duration,
		Size:       a.Size,
		Format:     a.Format,
		Status:     string(a.Status),
		CreatedAt:  a.CreatedAt,
		UpdatedAt:  a.UpdatedAt,
	}
	if a.Metadata != nil {
		data, _ := json.Marshal(a.Metadata)
		m.Metadata = datatypes.JSON(data)
	}
	if a.Tags != nil {
		data, _ := json.Marshal(a.Tags)
		m.Tags = datatypes.JSON(data)
	}
	return m
}

func AssetToDomain(m *model.MediaAssetModel) *media.Asset {
	a := &media.Asset{
		ID:         m.ID,
		Type:       media.AssetType(m.Type),
		SourceType: media.AssetSourceType(m.SourceType),
		SourceID:   m.SourceID,
		ParentID:   m.ParentID,
		Name:       m.Name,
		Path:       m.Path,
		Duration:   m.Duration,
		Size:       m.Size,
		Format:     m.Format,
		Status:     media.AssetStatus(m.Status),
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
	if m.Metadata != nil {
		_ = json.Unmarshal(m.Metadata, &a.Metadata)
	}
	if m.Tags != nil {
		_ = json.Unmarshal(m.Tags, &a.Tags)
	}
	return a
}
