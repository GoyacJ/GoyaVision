package mapper

import (
	"encoding/json"

	"goyavision/internal/domain/operator"
	"goyavision/internal/infra/persistence/model"

	"gorm.io/datatypes"
)

func OperatorToModel(o *operator.Operator) *model.OperatorModel {
	m := &model.OperatorModel{
		ID:             o.ID,
		TenantID:       o.TenantID,
		OwnerID:        o.OwnerID,
		Visibility:     int(o.Visibility),
		Code:           o.Code,
		Name:           o.Name,
		Description: o.Description,
		Category:    string(o.Category),
		Type:        string(o.Type),
		Origin:      string(o.Origin),
		ActiveVersionID: o.ActiveVersionID,
		Status:      string(o.Status),
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
	if o.Tags != nil {
		data, _ := json.Marshal(o.Tags)
		m.Tags = datatypes.JSON(data)
	}
	if o.VisibleRoleIDs != nil {
		data, _ := json.Marshal(o.VisibleRoleIDs)
		m.VisibleRoleIDs = datatypes.JSON(data)
	}
	return m
}

func OperatorToDomain(m *model.OperatorModel) *operator.Operator {
	o := &operator.Operator{
		ID:             m.ID,
		TenantID:       m.TenantID,
		OwnerID:        m.OwnerID,
		Visibility:     operator.Visibility(m.Visibility),
		Code:           m.Code,
		Name:           m.Name,
		Description: m.Description,
		Category:    operator.Category(m.Category),
		Type:        operator.Type(m.Type),
		Origin:      operator.Origin(m.Origin),
		ActiveVersionID: m.ActiveVersionID,
		Status:      operator.Status(m.Status),
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
	if m.Tags != nil {
		_ = json.Unmarshal(m.Tags, &o.Tags)
	}
	if m.VisibleRoleIDs != nil {
		_ = json.Unmarshal(m.VisibleRoleIDs, &o.VisibleRoleIDs)
	}
	if o.Origin == "" {
		o.Origin = operator.OriginCustom
	}
	if m.ActiveVersion != nil {
		o.ActiveVersion = OperatorVersionToDomain(m.ActiveVersion)
	}
	return o
}
