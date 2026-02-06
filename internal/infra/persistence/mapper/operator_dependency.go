package mapper

import (
	"goyavision/internal/domain/operator"
	"goyavision/internal/infra/persistence/model"
)

func OperatorDependencyToModel(d *operator.OperatorDependency) *model.OperatorDependencyModel {
	return &model.OperatorDependencyModel{
		ID:          d.ID,
		OperatorID:  d.OperatorID,
		DependsOnID: d.DependsOnID,
		MinVersion:  d.MinVersion,
		IsOptional:  d.IsOptional,
		CreatedAt:   d.CreatedAt,
	}
}

func OperatorDependencyToDomain(m *model.OperatorDependencyModel) *operator.OperatorDependency {
	return &operator.OperatorDependency{
		ID:          m.ID,
		OperatorID:  m.OperatorID,
		DependsOnID: m.DependsOnID,
		MinVersion:  m.MinVersion,
		IsOptional:  m.IsOptional,
		CreatedAt:   m.CreatedAt,
	}
}
