package repo

import (
	"context"
	"encoding/json"

	"goyavision/internal/domain/operator"
	"goyavision/internal/infra/persistence/mapper"
	"goyavision/internal/infra/persistence/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OperatorRepo struct {
	db *gorm.DB
}

func NewOperatorRepo(db *gorm.DB) *OperatorRepo {
	return &OperatorRepo{db: db}
}

func (r *OperatorRepo) Create(ctx context.Context, o *operator.Operator) error {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	m := mapper.OperatorToModel(o)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *OperatorRepo) Get(ctx context.Context, id uuid.UUID) (*operator.Operator, error) {
	var m model.OperatorModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.OperatorToDomain(&m), nil
}

func (r *OperatorRepo) GetWithActiveVersion(ctx context.Context, id uuid.UUID) (*operator.Operator, error) {
	var m model.OperatorModel
	if err := r.db.WithContext(ctx).
		Preload("ActiveVersion").
		Where("id = ?", id).
		First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.OperatorToDomain(&m), nil
}

func (r *OperatorRepo) GetByCode(ctx context.Context, code string) (*operator.Operator, error) {
	var m model.OperatorModel
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.OperatorToDomain(&m), nil
}

func (r *OperatorRepo) List(ctx context.Context, filter operator.Filter) ([]*operator.Operator, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.OperatorModel{})

	if filter.Category != nil {
		q = q.Where("category = ?", string(*filter.Category))
	}
	if filter.Type != nil {
		q = q.Where("type = ?", string(*filter.Type))
	}
	if filter.Status != nil {
		q = q.Where("status = ?", string(*filter.Status))
	}
	if filter.Origin != nil {
		q = q.Where("origin = ?", string(*filter.Origin))
	}
	if filter.IsBuiltin != nil {
		q = q.Where("is_builtin = ?", *filter.IsBuiltin)
	}
	if filter.ExecMode != nil {
		q = q.Joins("LEFT JOIN operator_versions ov ON ov.id = operators.active_version_id").
			Where("ov.exec_mode = ?", string(*filter.ExecMode))
	}
	if len(filter.Tags) > 0 {
		tagsJSON, _ := json.Marshal(filter.Tags)
		q = q.Where("tags @> ?::jsonb", string(tagsJSON))
	}
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		q = q.Where("name ILIKE ? OR description ILIKE ? OR code ILIKE ?", keyword, keyword, keyword)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []*model.OperatorModel
	if err := q.Limit(filter.Limit).Offset(filter.Offset).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*operator.Operator, len(models))
	for i, m := range models {
		result[i] = mapper.OperatorToDomain(m)
	}
	return result, total, nil
}

func (r *OperatorRepo) Update(ctx context.Context, o *operator.Operator) error {
	m := mapper.OperatorToModel(o)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *OperatorRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.OperatorModel{}).Error
}

func (r *OperatorRepo) ListPublished(ctx context.Context) ([]*operator.Operator, error) {
	var models []*model.OperatorModel
	if err := r.db.WithContext(ctx).Where("status = ?", string(operator.StatusPublished)).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*operator.Operator, len(models))
	for i, m := range models {
		result[i] = mapper.OperatorToDomain(m)
	}
	return result, nil
}

func (r *OperatorRepo) ListByCategory(ctx context.Context, category operator.Category) ([]*operator.Operator, error) {
	var models []*model.OperatorModel
	if err := r.db.WithContext(ctx).Where("category = ?", string(category)).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*operator.Operator, len(models))
	for i, m := range models {
		result[i] = mapper.OperatorToDomain(m)
	}
	return result, nil
}
