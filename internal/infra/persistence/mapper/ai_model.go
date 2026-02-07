package mapper

import (
	"encoding/json"

	"goyavision/internal/domain/ai_model"
	"goyavision/internal/infra/persistence/model"

	"gorm.io/datatypes"
)

func AIModelToModel(d *ai_model.AIModel) *model.AIModelModel {
	m := &model.AIModelModel{
		ID:             d.ID,
		TenantID:       d.TenantID,
		OwnerID:        d.OwnerID,
		Visibility:     int(d.Visibility),
		Name:           d.Name,
		Description: d.Description,
		Provider:    string(d.Provider),
		Endpoint:  d.Endpoint,
		APIKey:    d.APIKey,
		ModelName: d.ModelName,
		Status:    string(d.Status),
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
	if d.Config != nil {
		b, _ := json.Marshal(d.Config)
		m.Config = datatypes.JSON(b)
	}
	if d.VisibleRoleIDs != nil {
		data, _ := json.Marshal(d.VisibleRoleIDs)
		m.VisibleRoleIDs = datatypes.JSON(data)
	}
	return m
}

func AIModelToDomain(m *model.AIModelModel) *ai_model.AIModel {
	d := &ai_model.AIModel{
		ID:             m.ID,
		TenantID:       m.TenantID,
		OwnerID:        m.OwnerID,
		Visibility:     ai_model.Visibility(m.Visibility),
		Name:           m.Name,
		Description: m.Description,
		Provider:    ai_model.Provider(m.Provider),
		Endpoint:  m.Endpoint,
		APIKey:    m.APIKey,
		ModelName: m.ModelName,
		Status:    ai_model.Status(m.Status),
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
	if m.Config != nil {
		_ = json.Unmarshal(m.Config, &d.Config)
	}
	if m.VisibleRoleIDs != nil {
		_ = json.Unmarshal(m.VisibleRoleIDs, &d.VisibleRoleIDs)
	}
	return d
}
