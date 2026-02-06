package mapper

import (
	"encoding/json"

	"goyavision/internal/domain/operator"
	"goyavision/internal/infra/persistence/model"

	"gorm.io/datatypes"
)

func OperatorTemplateToModel(t *operator.OperatorTemplate) *model.OperatorTemplateModel {
	m := &model.OperatorTemplateModel{
		ID:          t.ID,
		Code:        t.Code,
		Name:        t.Name,
		Description: t.Description,
		Category:    string(t.Category),
		Type:        string(t.Type),
		ExecMode:    string(t.ExecMode),
		Author:      t.Author,
		IconURL:     t.IconURL,
		Downloads:   t.Downloads,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
	}
	if t.ExecConfig != nil {
		data, _ := json.Marshal(t.ExecConfig)
		m.ExecConfig = datatypes.JSON(data)
	}
	if t.InputSchema != nil {
		data, _ := json.Marshal(t.InputSchema)
		m.InputSchema = datatypes.JSON(data)
	}
	if t.OutputSpec != nil {
		data, _ := json.Marshal(t.OutputSpec)
		m.OutputSpec = datatypes.JSON(data)
	}
	if t.Config != nil {
		data, _ := json.Marshal(t.Config)
		m.Config = datatypes.JSON(data)
	}
	if t.Tags != nil {
		data, _ := json.Marshal(t.Tags)
		m.Tags = datatypes.JSON(data)
	}
	return m
}

func OperatorTemplateToDomain(m *model.OperatorTemplateModel) *operator.OperatorTemplate {
	t := &operator.OperatorTemplate{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		Category:    operator.Category(m.Category),
		Type:        operator.Type(m.Type),
		ExecMode:    operator.ExecMode(m.ExecMode),
		Author:      m.Author,
		IconURL:     m.IconURL,
		Downloads:   m.Downloads,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
	if m.ExecConfig != nil {
		var cfg operator.ExecConfig
		if err := json.Unmarshal(m.ExecConfig, &cfg); err == nil {
			t.ExecConfig = &cfg
		}
	}
	if m.InputSchema != nil {
		_ = json.Unmarshal(m.InputSchema, &t.InputSchema)
	}
	if m.OutputSpec != nil {
		_ = json.Unmarshal(m.OutputSpec, &t.OutputSpec)
	}
	if m.Config != nil {
		_ = json.Unmarshal(m.Config, &t.Config)
	}
	if m.Tags != nil {
		_ = json.Unmarshal(m.Tags, &t.Tags)
	}
	return t
}
