package repo

import (
	"context"
	"encoding/json"

	"goyavision/internal/domain/system"
	"goyavision/internal/infra/persistence/model"

	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SystemConfigRepo struct {
	db *gorm.DB
}

func NewSystemConfigRepo(db *gorm.DB) *SystemConfigRepo {
	return &SystemConfigRepo{db: db}
}

func (r *SystemConfigRepo) Get(ctx context.Context, key string) (*system.SystemConfig, error) {
	var m model.SystemConfigModel
	if err := r.db.WithContext(ctx).Where("key = ?", key).First(&m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Return nil if not found
		}
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *SystemConfigRepo) List(ctx context.Context) ([]*system.SystemConfig, error) {
	var models []*model.SystemConfigModel
	if err := r.db.WithContext(ctx).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*system.SystemConfig, len(models))
	for i, m := range models {
		result[i] = toDomain(m)
	}
	return result, nil
}

func (r *SystemConfigRepo) Save(ctx context.Context, config *system.SystemConfig) error {
	m := toModel(config)
	return r.db.WithContext(ctx).Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "description", "updated_at"}),
	}).Create(m).Error
}

func (r *SystemConfigRepo) Delete(ctx context.Context, key string) error {
	return r.db.WithContext(ctx).Where("key = ?", key).Delete(&model.SystemConfigModel{}).Error
}

func toDomain(m *model.SystemConfigModel) *system.SystemConfig {
	return &system.SystemConfig{
		Key:         m.Key,
		Value:       json.RawMessage(m.Value),
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
		UpdatedAt:   m.UpdatedAt,
	}
}

func toModel(c *system.SystemConfig) *model.SystemConfigModel {
	return &model.SystemConfigModel{
		Key:         c.Key,
		Value:       datatypes.JSON(c.Value),
		Description: c.Description,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
