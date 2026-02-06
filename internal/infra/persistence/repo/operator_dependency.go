package repo

import (
	"context"
	"fmt"

	"goyavision/internal/domain/operator"
	"goyavision/internal/infra/persistence/mapper"
	"goyavision/internal/infra/persistence/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OperatorDependencyRepo struct {
	db *gorm.DB
}

func NewOperatorDependencyRepo(db *gorm.DB) *OperatorDependencyRepo {
	return &OperatorDependencyRepo{db: db}
}

func (r *OperatorDependencyRepo) Create(ctx context.Context, dep *operator.OperatorDependency) error {
	if dep.ID == uuid.Nil {
		dep.ID = uuid.New()
	}
	m := mapper.OperatorDependencyToModel(dep)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *OperatorDependencyRepo) ListByOperator(ctx context.Context, operatorID uuid.UUID) ([]*operator.OperatorDependency, error) {
	var models []*model.OperatorDependencyModel
	if err := r.db.WithContext(ctx).Where("operator_id = ?", operatorID).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*operator.OperatorDependency, len(models))
	for i, m := range models {
		result[i] = mapper.OperatorDependencyToDomain(m)
	}
	return result, nil
}

func (r *OperatorDependencyRepo) DeleteByOperator(ctx context.Context, operatorID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("operator_id = ?", operatorID).Delete(&model.OperatorDependencyModel{}).Error
}

func (r *OperatorDependencyRepo) CheckDependenciesSatisfied(ctx context.Context, operatorID uuid.UUID) (bool, []string, error) {
	deps, err := r.ListByOperator(ctx, operatorID)
	if err != nil {
		return false, nil, err
	}

	if len(deps) == 0 {
		return true, nil, nil
	}

	var messages []string
	for _, dep := range deps {
		var cnt int64
		err := r.db.WithContext(ctx).
			Model(&model.OperatorModel{}).
			Where("id = ? AND status IN ?", dep.DependsOnID, []string{string(operator.StatusPublished)}).
			Count(&cnt).Error
		if err != nil {
			return false, nil, err
		}
		if cnt == 0 && !dep.IsOptional {
			messages = append(messages, fmt.Sprintf("依赖算子 %s 未发布或不存在", dep.DependsOnID.String()))
		}
	}

	return len(messages) == 0, messages, nil
}
