package repo

import (
	"context"

	"goyavision/internal/domain/media"
	"goyavision/internal/infra/persistence/mapper"
	"goyavision/internal/infra/persistence/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaSourceRepo struct {
	db *gorm.DB
}

func NewMediaSourceRepo(db *gorm.DB) *MediaSourceRepo {
	return &MediaSourceRepo{db: db}
}

func (r *MediaSourceRepo) Create(ctx context.Context, s *media.Source) error {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	m := mapper.SourceToModel(s)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *MediaSourceRepo) Get(ctx context.Context, id uuid.UUID) (*media.Source, error) {
	var m model.MediaSourceModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.SourceToDomain(&m), nil
}

func (r *MediaSourceRepo) GetByPathName(ctx context.Context, pathName string) (*media.Source, error) {
	var m model.MediaSourceModel
	if err := r.db.WithContext(ctx).Where("path_name = ?", pathName).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.SourceToDomain(&m), nil
}

func (r *MediaSourceRepo) List(ctx context.Context, filter media.SourceFilter) ([]*media.Source, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.MediaSourceModel{})
	if filter.Type != nil {
		q = q.Where("type = ?", string(*filter.Type))
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 1000 {
		filter.Limit = 1000
	}
	var models []*model.MediaSourceModel
	if err := q.Limit(filter.Limit).Offset(filter.Offset).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}
	result := make([]*media.Source, len(models))
	for i, m := range models {
		result[i] = mapper.SourceToDomain(m)
	}
	return result, total, nil
}

func (r *MediaSourceRepo) Update(ctx context.Context, s *media.Source) error {
	m := mapper.SourceToModel(s)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *MediaSourceRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.MediaSourceModel{}).Error
}
