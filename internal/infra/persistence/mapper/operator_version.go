package mapper

import (
	"encoding/json"

	"goyavision/internal/domain/operator"
	"goyavision/internal/infra/persistence/model"

	"gorm.io/datatypes"
)

func OperatorVersionToModel(v *operator.OperatorVersion) *model.OperatorVersionModel {
	m := &model.OperatorVersionModel{
		ID:         v.ID,
		OperatorID: v.OperatorID,
		Version:    v.Version,
		ExecMode:   string(v.ExecMode),
		Changelog:  v.Changelog,
		Status:     string(v.Status),
		CreatedAt:  v.CreatedAt,
		UpdatedAt:  v.UpdatedAt,
	}
	if v.ExecConfig != nil {
		data, _ := json.Marshal(v.ExecConfig)
		m.ExecConfig = datatypes.JSON(data)
	}
	if v.InputSchema != nil {
		data, _ := json.Marshal(v.InputSchema)
		m.InputSchema = datatypes.JSON(data)
	}
	if v.OutputSpec != nil {
		data, _ := json.Marshal(v.OutputSpec)
		m.OutputSpec = datatypes.JSON(data)
	}
	if v.Config != nil {
		data, _ := json.Marshal(v.Config)
		m.Config = datatypes.JSON(data)
	}
	return m
}

func OperatorVersionToDomain(m *model.OperatorVersionModel) *operator.OperatorVersion {
	v := &operator.OperatorVersion{
		ID:         m.ID,
		OperatorID: m.OperatorID,
		Version:    m.Version,
		ExecMode:   operator.ExecMode(m.ExecMode),
		Changelog:  m.Changelog,
		Status:     operator.VersionStatus(m.Status),
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
	if m.ExecConfig != nil {
		var cfg operator.ExecConfig
		if err := json.Unmarshal(m.ExecConfig, &cfg); err == nil {
			v.ExecConfig = &cfg
		}
	}
	if m.InputSchema != nil {
		_ = json.Unmarshal(m.InputSchema, &v.InputSchema)
	}
	if m.OutputSpec != nil {
		_ = json.Unmarshal(m.OutputSpec, &v.OutputSpec)
	}
	if m.Config != nil {
		_ = json.Unmarshal(m.Config, &v.Config)
	}
	return v
}
