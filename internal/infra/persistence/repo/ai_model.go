package repo

import (
	"context"

	"goyavision/internal/domain/ai_model"
	"goyavision/internal/infra/persistence/mapper"
	"goyavision/internal/infra/persistence/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AIModelRepo struct {
	db *gorm.DB
}

func NewAIModelRepo(db *gorm.DB) *AIModelRepo {
	return &AIModelRepo{db: db}
}

func (r *AIModelRepo) Create(ctx context.Context, d *ai_model.AIModel) error {
	if d.ID == uuid.Nil {
		d.ID = uuid.New()
	}
	m := mapper.AIModelToModel(d)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *AIModelRepo) Get(ctx context.Context, id uuid.UUID) (*ai_model.AIModel, error) {
	var m model.AIModelModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.AIModelToDomain(&m), nil
}

func (r *AIModelRepo) Update(ctx context.Context, d *ai_model.AIModel) error {
	m := mapper.AIModelToModel(d)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *AIModelRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.AIModelModel{}).Error
}

func (r *AIModelRepo) List(ctx context.Context, filter ai_model.Filter) ([]*ai_model.AIModel, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.AIModelModel{})

	if filter.Keyword != "" {
		k := "%" + filter.Keyword + "%"
		q = q.Where("name ILIKE ?", k)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []*model.AIModelModel
	if err := q.Limit(filter.Limit).Offset(filter.Offset).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*ai_model.AIModel, len(models))
	for i, m := range models {
		result[i] = mapper.AIModelToDomain(m)
	}
	return result, total, nil
}
