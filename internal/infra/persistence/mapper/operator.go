package mapper

import (
	"encoding/json"

	"goyavision/internal/domain/operator"
	"goyavision/internal/infra/persistence/model"

	"gorm.io/datatypes"
)

func OperatorToModel(o *operator.Operator) *model.OperatorModel {
	m := &model.OperatorModel{
		ID:          o.ID,
		Code:        o.Code,
		Name:        o.Name,
		Description: o.Description,
		Category:    string(o.Category),
		Type:        string(o.Type),
		Origin:      string(o.Origin),
		ActiveVersionID: o.ActiveVersionID,
		Version:     o.Version,
		Endpoint:    o.Endpoint,
		Method:      o.Method,
		Status:      string(o.Status),
		IsBuiltin:   o.IsBuiltin,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
	if o.InputSchema != nil {
		data, _ := json.Marshal(o.InputSchema)
		m.InputSchema = datatypes.JSON(data)
	}
	if o.OutputSpec != nil {
		data, _ := json.Marshal(o.OutputSpec)
		m.OutputSpec = datatypes.JSON(data)
	}
	if o.Config != nil {
		data, _ := json.Marshal(o.Config)
		m.Config = datatypes.JSON(data)
	}
	if o.Tags != nil {
		data, _ := json.Marshal(o.Tags)
		m.Tags = datatypes.JSON(data)
	}
	return m
}

func OperatorToDomain(m *model.OperatorModel) *operator.Operator {
	o := &operator.Operator{
		ID:          m.ID,
		Code:        m.Code,
		Name:        m.Name,
		Description: m.Description,
		Category:    operator.Category(m.Category),
		Type:        operator.Type(m.Type),
		Origin:      operator.Origin(m.Origin),
		ActiveVersionID: m.ActiveVersionID,
		Version:     m.Version,
		Endpoint:    m.Endpoint,
		Method:      m.Method,
		Status:      operator.Status(m.Status),
		IsBuiltin:   m.IsBuiltin,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
	if m.InputSchema != nil {
		_ = json.Unmarshal(m.InputSchema, &o.InputSchema)
	}
	if m.OutputSpec != nil {
		_ = json.Unmarshal(m.OutputSpec, &o.OutputSpec)
	}
	if m.Config != nil {
		_ = json.Unmarshal(m.Config, &o.Config)
	}
	if m.Tags != nil {
		_ = json.Unmarshal(m.Tags, &o.Tags)
	}
	if o.Origin == "" {
		if o.IsBuiltin {
			o.Origin = operator.OriginBuiltin
		} else {
			o.Origin = operator.OriginCustom
		}
	}
	if m.ActiveVersion != nil {
		o.ActiveVersion = OperatorVersionToDomain(m.ActiveVersion)
	}
	return o
}
