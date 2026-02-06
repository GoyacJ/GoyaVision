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

type OperatorTemplateRepo struct {
	db *gorm.DB
}

func NewOperatorTemplateRepo(db *gorm.DB) *OperatorTemplateRepo {
	return &OperatorTemplateRepo{db: db}
}

func (r *OperatorTemplateRepo) Create(ctx context.Context, tpl *operator.OperatorTemplate) error {
	if tpl.ID == uuid.Nil {
		tpl.ID = uuid.New()
	}
	m := mapper.OperatorTemplateToModel(tpl)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *OperatorTemplateRepo) Get(ctx context.Context, id uuid.UUID) (*operator.OperatorTemplate, error) {
	var m model.OperatorTemplateModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.OperatorTemplateToDomain(&m), nil
}

func (r *OperatorTemplateRepo) GetByCode(ctx context.Context, code string) (*operator.OperatorTemplate, error) {
	var m model.OperatorTemplateModel
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.OperatorTemplateToDomain(&m), nil
}

func (r *OperatorTemplateRepo) List(ctx context.Context, filter operator.TemplateFilter) ([]*operator.OperatorTemplate, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.OperatorTemplateModel{})

	if filter.Category != nil {
		q = q.Where("category = ?", string(*filter.Category))
	}
	if filter.Type != nil {
		q = q.Where("type = ?", string(*filter.Type))
	}
	if filter.ExecMode != nil {
		q = q.Where("exec_mode = ?", string(*filter.ExecMode))
	}
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		q = q.Where("name ILIKE ? OR description ILIKE ? OR code ILIKE ?", keyword, keyword, keyword)
	}
	if len(filter.Tags) > 0 {
		tagsJSON, _ := json.Marshal(filter.Tags)
		q = q.Where("tags @> ?::jsonb", string(tagsJSON))
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []*model.OperatorTemplateModel
	if err := q.Limit(filter.Limit).Offset(filter.Offset).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*operator.OperatorTemplate, len(models))
	for i, m := range models {
		result[i] = mapper.OperatorTemplateToDomain(m)
	}
	return result, total, nil
}

func (r *OperatorTemplateRepo) Update(ctx context.Context, tpl *operator.OperatorTemplate) error {
	m := mapper.OperatorTemplateToModel(tpl)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *OperatorTemplateRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.OperatorTemplateModel{}).Error
}

func (r *OperatorTemplateRepo) IncrementDownloads(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Model(&model.OperatorTemplateModel{}).
		Where("id = ?", id).
		UpdateColumn("downloads", gorm.Expr("downloads + 1")).Error
}
