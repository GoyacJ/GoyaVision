package mapper

import (
	"encoding/json"

	"goyavision/internal/domain/media"
	"goyavision/internal/infra/persistence/model"

	"gorm.io/datatypes"
)

func SourceToModel(s *media.Source) *model.MediaSourceModel {
	m := &model.MediaSourceModel{
		ID:             s.ID,
		TenantID:       s.TenantID,
		OwnerID:        s.OwnerID,
		Visibility:     int(s.Visibility),
		Name:           s.Name,
		PathName:      s.PathName,
		Type:          string(s.Type),
		URL:           s.URL,
		Protocol:      s.Protocol,
		Enabled:       s.Enabled,
		RecordEnabled: s.RecordEnabled,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
	}
	if s.VisibleRoleIDs != nil {
		data, _ := json.Marshal(s.VisibleRoleIDs)
		m.VisibleRoleIDs = datatypes.JSON(data)
	}
	return m
}

func SourceToDomain(m *model.MediaSourceModel) *media.Source {
	s := &media.Source{
		ID:             m.ID,
		TenantID:       m.TenantID,
		OwnerID:        m.OwnerID,
		Visibility:     media.Visibility(m.Visibility),
		Name:           m.Name,
		PathName:      m.PathName,
		Type:          media.SourceType(m.Type),
		URL:           m.URL,
		Protocol:      m.Protocol,
		Enabled:       m.Enabled,
		RecordEnabled: m.RecordEnabled,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
	if m.VisibleRoleIDs != nil {
		_ = json.Unmarshal(m.VisibleRoleIDs, &s.VisibleRoleIDs)
	}
	return s
}

func AssetToModel(a *media.Asset) *model.MediaAssetModel {
	m := &model.MediaAssetModel{
		ID:         a.ID,
		TenantID:   a.TenantID,
		OwnerID:    a.OwnerID,
		Visibility: int(a.Visibility),
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
	if a.VisibleRoleIDs != nil {
		data, _ := json.Marshal(a.VisibleRoleIDs)
		m.VisibleRoleIDs = datatypes.JSON(data)
	}
	return m
}

func AssetToDomain(m *model.MediaAssetModel) *media.Asset {
	a := &media.Asset{
		ID:         m.ID,
		TenantID:   m.TenantID,
		OwnerID:    m.OwnerID,
		Visibility: media.Visibility(m.Visibility),
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
	if m.VisibleRoleIDs != nil {
		_ = json.Unmarshal(m.VisibleRoleIDs, &a.VisibleRoleIDs)
	}
	return a
}
