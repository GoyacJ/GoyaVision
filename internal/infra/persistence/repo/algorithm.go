package repo

import (
	"context"
	"encoding/json"

	"goyavision/internal/domain/algorithm"
	"goyavision/internal/infra/persistence/mapper"
	"goyavision/internal/infra/persistence/model"
	"goyavision/internal/infra/persistence/scope"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AlgorithmRepo struct {
	db *gorm.DB
}

func NewAlgorithmRepo(db *gorm.DB) *AlgorithmRepo {
	return &AlgorithmRepo{db: db}
}

func (r *AlgorithmRepo) Create(ctx context.Context, a *algorithm.Algorithm) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	tenantID, userID := scope.GetContextInfo(ctx)
	m := mapper.AlgorithmToModel(a)
	m.TenantID = tenantID
	m.OwnerID = userID
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *AlgorithmRepo) Get(ctx context.Context, id uuid.UUID) (*algorithm.Algorithm, error) {
	var m model.AlgorithmModel
	if err := r.db.WithContext(ctx).
		Scopes(scope.ScopeTenant(ctx), scope.ScopeVisibility(ctx)).
		Where("id = ?", id).
		First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.AlgorithmToDomain(&m), nil
}

func (r *AlgorithmRepo) GetByCode(ctx context.Context, code string) (*algorithm.Algorithm, error) {
	var m model.AlgorithmModel
	if err := r.db.WithContext(ctx).
		Scopes(scope.ScopeTenant(ctx), scope.ScopeVisibility(ctx)).
		Where("code = ?", code).
		First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.AlgorithmToDomain(&m), nil
}

func (r *AlgorithmRepo) GetWithRelations(ctx context.Context, id uuid.UUID) (*algorithm.Algorithm, error) {
	var m model.AlgorithmModel
	if err := r.db.WithContext(ctx).
		Scopes(scope.ScopeTenant(ctx), scope.ScopeVisibility(ctx)).
		Preload("Versions", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Preload("Versions.Implementations", func(db *gorm.DB) *gorm.DB {
			return db.Order("is_default DESC, created_at ASC")
		}).
		Preload("Versions.Evaluations", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC")
		}).
		Where("id = ?", id).
		First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.AlgorithmToDomain(&m), nil
}

func (r *AlgorithmRepo) List(ctx context.Context, filter algorithm.Filter) ([]*algorithm.Algorithm, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.AlgorithmModel{}).Scopes(scope.ScopeTenant(ctx), scope.ScopeVisibility(ctx))
	if filter.Status != nil {
		q = q.Where("status = ?", string(*filter.Status))
	}
	if filter.Scenario != "" {
		q = q.Where("scenario = ?", filter.Scenario)
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

	var models []*model.AlgorithmModel
	if err := q.Order("created_at DESC").Limit(filter.Limit).Offset(filter.Offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	out := make([]*algorithm.Algorithm, len(models))
	for i := range models {
		out[i] = mapper.AlgorithmToDomain(models[i])
	}
	return out, total, nil
}

func (r *AlgorithmRepo) Update(ctx context.Context, a *algorithm.Algorithm) error {
	m := mapper.AlgorithmToModel(a)
	return r.db.WithContext(ctx).
		Scopes(scope.ScopeTenant(ctx)).
		Where("id = ?", a.ID).
		Omit("Versions").
		Updates(m).Error
}

func (r *AlgorithmRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Scopes(scope.ScopeTenant(ctx)).Where("id = ?", id).Delete(&model.AlgorithmModel{}).Error
}

func (r *AlgorithmRepo) CreateVersion(ctx context.Context, v *algorithm.Version) error {
	if v.ID == uuid.Nil {
		v.ID = uuid.New()
	}
	m := mapper.AlgorithmVersionToModel(v)
	return r.db.WithContext(ctx).Omit("Implementations", "Evaluations").Create(m).Error
}

func (r *AlgorithmRepo) GetVersion(ctx context.Context, id uuid.UUID) (*algorithm.Version, error) {
	var m model.AlgorithmVersionModel
	if err := r.db.WithContext(ctx).
		Table("algorithm_versions AS v").
		Select("v.*").
		Joins("JOIN algorithms a ON a.id = v.algorithm_id").
		Scopes(scope.ScopeTenant(ctx), scope.ScopeVisibility(ctx)).
		Preload("Implementations", func(db *gorm.DB) *gorm.DB { return db.Order("is_default DESC, created_at ASC") }).
		Preload("Evaluations", func(db *gorm.DB) *gorm.DB { return db.Order("created_at DESC") }).
		Where("v.id = ?", id).
		First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.AlgorithmVersionToDomain(&m), nil
}

func (r *AlgorithmRepo) ListVersions(ctx context.Context, algorithmID uuid.UUID) ([]*algorithm.Version, error) {
	var models []*model.AlgorithmVersionModel
	if err := r.db.WithContext(ctx).
		Table("algorithm_versions AS v").
		Select("v.*").
		Joins("JOIN algorithms a ON a.id = v.algorithm_id").
		Scopes(scope.ScopeTenant(ctx), scope.ScopeVisibility(ctx)).
		Preload("Implementations", func(db *gorm.DB) *gorm.DB { return db.Order("is_default DESC, created_at ASC") }).
		Preload("Evaluations", func(db *gorm.DB) *gorm.DB { return db.Order("created_at DESC") }).
		Where("v.algorithm_id = ?", algorithmID).
		Order("v.created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]*algorithm.Version, len(models))
	for i := range models {
		out[i] = mapper.AlgorithmVersionToDomain(models[i])
	}
	return out, nil
}

func (r *AlgorithmRepo) GetVersionByName(ctx context.Context, algorithmID uuid.UUID, version string) (*algorithm.Version, error) {
	var m model.AlgorithmVersionModel
	if err := r.db.WithContext(ctx).
		Table("algorithm_versions AS v").
		Select("v.*").
		Joins("JOIN algorithms a ON a.id = v.algorithm_id").
		Scopes(scope.ScopeTenant(ctx), scope.ScopeVisibility(ctx)).
		Preload("Implementations", func(db *gorm.DB) *gorm.DB { return db.Order("is_default DESC, created_at ASC") }).
		Preload("Evaluations", func(db *gorm.DB) *gorm.DB { return db.Order("created_at DESC") }).
		Where("v.algorithm_id = ? AND v.version = ?", algorithmID, version).
		First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.AlgorithmVersionToDomain(&m), nil
}

func (r *AlgorithmRepo) UpdateVersion(ctx context.Context, v *algorithm.Version) error {
	m := mapper.AlgorithmVersionToModel(v)
	visibleAlgorithms := r.db.WithContext(ctx).
		Model(&model.AlgorithmModel{}).
		Scopes(scope.ScopeTenant(ctx), scope.ScopeVisibility(ctx)).
		Select("id")
	return r.db.WithContext(ctx).
		Where("id = ? AND algorithm_id IN (?)", v.ID, visibleAlgorithms).
		Omit("Implementations", "Evaluations").
		Updates(m).Error
}

func (r *AlgorithmRepo) ReplaceImplementations(ctx context.Context, versionID uuid.UUID, impls []algorithm.Implementation) error {
	if err := r.ensureVersionVisible(ctx, versionID); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("version_id = ?", versionID).Delete(&model.AlgorithmImplementationModel{}).Error; err != nil {
			return err
		}

		var firstID uuid.UUID
		var defaultID *uuid.UUID
		for i := range impls {
			if impls[i].ID == uuid.Nil {
				impls[i].ID = uuid.New()
			}
			impls[i].VersionID = versionID

			if i == 0 {
				firstID = impls[i].ID
			}
			if impls[i].IsDefault {
				id := impls[i].ID
				defaultID = &id
			}

			if err := tx.Create(mapper.AlgorithmImplementationToModel(&impls[i])).Error; err != nil {
				return err
			}
		}

		if defaultID == nil && len(impls) > 0 {
			defaultID = &firstID
			if err := tx.Model(&model.AlgorithmImplementationModel{}).
				Where("id = ?", firstID).
				Update("is_default", true).Error; err != nil {
				return err
			}
		}

		return tx.Model(&model.AlgorithmVersionModel{}).
			Where("id = ?", versionID).
			Update("default_implementation_id", defaultID).Error
	})
}

func (r *AlgorithmRepo) ListImplementations(ctx context.Context, versionID uuid.UUID) ([]*algorithm.Implementation, error) {
	if err := r.ensureVersionVisible(ctx, versionID); err != nil {
		return nil, err
	}
	var models []*model.AlgorithmImplementationModel
	if err := r.db.WithContext(ctx).
		Where("version_id = ?", versionID).
		Order("is_default DESC, created_at ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]*algorithm.Implementation, len(models))
	for i := range models {
		out[i] = mapper.AlgorithmImplementationToDomain(models[i])
	}
	return out, nil
}

func (r *AlgorithmRepo) ReplaceEvaluations(ctx context.Context, versionID uuid.UUID, profiles []algorithm.EvaluationProfile) error {
	if err := r.ensureVersionVisible(ctx, versionID); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("version_id = ?", versionID).Delete(&model.AlgorithmEvaluationModel{}).Error; err != nil {
			return err
		}
		for i := range profiles {
			if profiles[i].ID == uuid.Nil {
				profiles[i].ID = uuid.New()
			}
			profiles[i].VersionID = versionID
			if err := tx.Create(mapper.AlgorithmEvaluationToModel(&profiles[i])).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

func (r *AlgorithmRepo) ListEvaluations(ctx context.Context, versionID uuid.UUID) ([]*algorithm.EvaluationProfile, error) {
	if err := r.ensureVersionVisible(ctx, versionID); err != nil {
		return nil, err
	}
	var models []*model.AlgorithmEvaluationModel
	if err := r.db.WithContext(ctx).
		Where("version_id = ?", versionID).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]*algorithm.EvaluationProfile, len(models))
	for i := range models {
		out[i] = mapper.AlgorithmEvaluationToDomain(models[i])
	}
	return out, nil
}

func (r *AlgorithmRepo) ensureVersionVisible(ctx context.Context, versionID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Table("algorithm_versions AS v").
		Select("v.id").
		Joins("JOIN algorithms a ON a.id = v.algorithm_id").
		Scopes(scope.ScopeTenant(ctx), scope.ScopeVisibility(ctx)).
		Where("v.id = ?", versionID).
		Take(&model.AlgorithmVersionModel{}).Error
}
