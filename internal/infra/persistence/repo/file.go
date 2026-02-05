package repo

import (
	"context"

	"goyavision/internal/domain/storage"
	"goyavision/internal/infra/persistence/mapper"
	"goyavision/internal/infra/persistence/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileRepo struct {
	db *gorm.DB
}

func NewFileRepo(db *gorm.DB) *FileRepo {
	return &FileRepo{db: db}
}

func (r *FileRepo) Create(ctx context.Context, f *storage.File) error {
	if f.ID == uuid.Nil {
		f.ID = uuid.New()
	}
	m := mapper.FileToModel(f)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *FileRepo) Get(ctx context.Context, id uuid.UUID) (*storage.File, error) {
	var m model.FileModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.FileToDomain(&m), nil
}

func (r *FileRepo) GetByPath(ctx context.Context, path string) (*storage.File, error) {
	var m model.FileModel
	if err := r.db.WithContext(ctx).Where("path = ?", path).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.FileToDomain(&m), nil
}

func (r *FileRepo) GetByHash(ctx context.Context, hash string) (*storage.File, error) {
	var m model.FileModel
	if err := r.db.WithContext(ctx).Where("hash = ?", hash).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.FileToDomain(&m), nil
}

func (r *FileRepo) List(ctx context.Context, filter storage.FileFilter) ([]*storage.File, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.FileModel{})

	if filter.Type != nil {
		q = q.Where("type = ?", string(*filter.Type))
	}
	if filter.Status != nil {
		q = q.Where("status = ?", string(*filter.Status))
	}
	if filter.UploaderID != nil {
		q = q.Where("uploader_id = ?", *filter.UploaderID)
	}
	if filter.Search != "" {
		q = q.Where("name ILIKE ? OR original_name ILIKE ?", "%"+filter.Search+"%", "%"+filter.Search+"%")
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

	var models []*model.FileModel
	if err := q.Limit(filter.Limit).Offset(filter.Offset).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*storage.File, len(models))
	for i, m := range models {
		result[i] = mapper.FileToDomain(m)
	}
	return result, total, nil
}

func (r *FileRepo) Update(ctx context.Context, f *storage.File) error {
	m := mapper.FileToModel(f)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *FileRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.FileModel{}).Error
}
