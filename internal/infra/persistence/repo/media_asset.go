package repo

import (
	"context"
	"encoding/json"

	"goyavision/internal/domain/media"
	"goyavision/internal/infra/persistence/mapper"
	"goyavision/internal/infra/persistence/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MediaAssetRepo struct {
	db *gorm.DB
}

func NewMediaAssetRepo(db *gorm.DB) *MediaAssetRepo {
	return &MediaAssetRepo{db: db}
}

func (r *MediaAssetRepo) Create(ctx context.Context, a *media.Asset) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	m := mapper.AssetToModel(a)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *MediaAssetRepo) Get(ctx context.Context, id uuid.UUID) (*media.Asset, error) {
	var m model.MediaAssetModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.AssetToDomain(&m), nil
}

func (r *MediaAssetRepo) List(ctx context.Context, filter media.AssetFilter) ([]*media.Asset, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.MediaAssetModel{})

	if filter.Type != nil {
		q = q.Where("type = ?", string(*filter.Type))
	}
	if filter.SourceType != nil {
		q = q.Where("source_type = ?", string(*filter.SourceType))
	}
	if filter.SourceID != nil {
		q = q.Where("source_id = ?", *filter.SourceID)
	}
	if filter.ParentID != nil {
		q = q.Where("parent_id = ?", *filter.ParentID)
	}
	if filter.Status != nil {
		q = q.Where("status = ?", string(*filter.Status))
	}
	if len(filter.Tags) > 0 {
		tagsJSON, _ := json.Marshal(filter.Tags)
		q = q.Where("tags @> ?::jsonb", string(tagsJSON))
	}
	if filter.From != nil {
		q = q.Where("created_at >= ?", *filter.From)
	}
	if filter.To != nil {
		q = q.Where("created_at <= ?", *filter.To)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []*model.MediaAssetModel
	if err := q.Limit(filter.Limit).Offset(filter.Offset).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*media.Asset, len(models))
	for i, m := range models {
		result[i] = mapper.AssetToDomain(m)
	}
	return result, total, nil
}

func (r *MediaAssetRepo) Update(ctx context.Context, a *media.Asset) error {
	m := mapper.AssetToModel(a)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *MediaAssetRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.MediaAssetModel{}).Error
}

func (r *MediaAssetRepo) ListBySource(ctx context.Context, sourceID uuid.UUID) ([]*media.Asset, error) {
	var models []*model.MediaAssetModel
	if err := r.db.WithContext(ctx).Where("source_id = ?", sourceID).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*media.Asset, len(models))
	for i, m := range models {
		result[i] = mapper.AssetToDomain(m)
	}
	return result, nil
}

func (r *MediaAssetRepo) ListByParent(ctx context.Context, parentID uuid.UUID) ([]*media.Asset, error) {
	var models []*model.MediaAssetModel
	if err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*media.Asset, len(models))
	for i, m := range models {
		result[i] = mapper.AssetToDomain(m)
	}
	return result, nil
}

func (r *MediaAssetRepo) GetAllTags(ctx context.Context) ([]string, error) {
	var models []*model.MediaAssetModel
	if err := r.db.WithContext(ctx).Select("tags").Where("tags IS NOT NULL AND tags != '[]'").Find(&models).Error; err != nil {
		return nil, err
	}

	tagSet := make(map[string]bool)
	for _, m := range models {
		if m.Tags == nil {
			continue
		}
		var tags []string
		if err := json.Unmarshal(m.Tags, &tags); err == nil {
			for _, tag := range tags {
				if tag != "" {
					tagSet[tag] = true
				}
			}
		}
	}

	result := make([]string, 0, len(tagSet))
	for tag := range tagSet {
		result = append(result, tag)
	}
	return result, nil
}
