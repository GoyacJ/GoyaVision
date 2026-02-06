package repo

import (
	"context"

	"goyavision/internal/domain/operator"
	"goyavision/internal/infra/persistence/mapper"
	"goyavision/internal/infra/persistence/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OperatorVersionRepo struct {
	db *gorm.DB
}

func NewOperatorVersionRepo(db *gorm.DB) *OperatorVersionRepo {
	return &OperatorVersionRepo{db: db}
}

func (r *OperatorVersionRepo) Create(ctx context.Context, version *operator.OperatorVersion) error {
	if version.ID == uuid.Nil {
		version.ID = uuid.New()
	}
	m := mapper.OperatorVersionToModel(version)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *OperatorVersionRepo) Get(ctx context.Context, id uuid.UUID) (*operator.OperatorVersion, error) {
	var m model.OperatorVersionModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.OperatorVersionToDomain(&m), nil
}

func (r *OperatorVersionRepo) ListByOperator(ctx context.Context, operatorID uuid.UUID) ([]*operator.OperatorVersion, error) {
	var models []*model.OperatorVersionModel
	if err := r.db.WithContext(ctx).Where("operator_id = ?", operatorID).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*operator.OperatorVersion, len(models))
	for i, m := range models {
		result[i] = mapper.OperatorVersionToDomain(m)
	}
	return result, nil
}

func (r *OperatorVersionRepo) GetByOperatorAndVersion(ctx context.Context, operatorID uuid.UUID, version string) (*operator.OperatorVersion, error) {
	var m model.OperatorVersionModel
	if err := r.db.WithContext(ctx).Where("operator_id = ? AND version = ?", operatorID, version).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.OperatorVersionToDomain(&m), nil
}

func (r *OperatorVersionRepo) Update(ctx context.Context, version *operator.OperatorVersion) error {
	m := mapper.OperatorVersionToModel(version)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *OperatorVersionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.OperatorVersionModel{}).Error
}
