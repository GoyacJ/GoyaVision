package persistence

import (
	"context"
	"errors"

	"goyavision/internal/domain"
	"goyavision/internal/port"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrDBNotConfigured = errors.New("database not configured")

var _ port.Repository = (*repository)(nil)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) checkDB() error {
	if r.db == nil {
		return ErrDBNotConfigured
	}
	return nil
}

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&domain.Stream{},
		&domain.Algorithm{},
		&domain.AlgorithmBinding{},
		&domain.RecordSession{},
		&domain.InferenceResult{},
	); err != nil {
		return err
	}

	if err := addIndexesAndConstraints(db); err != nil {
		return err
	}

	return nil
}

func addIndexesAndConstraints(db *gorm.DB) error {
	if db.Migrator().HasIndex(&domain.RecordSession{}, "idx_record_sessions_stream_running") {
		return nil
	}

	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_record_sessions_stream_running 
		ON record_sessions (stream_id) 
		WHERE status = 'running'
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_inference_results_stream_ts 
		ON inference_results (stream_id, ts DESC)
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_inference_results_binding_ts 
		ON inference_results (algorithm_binding_id, ts DESC)
	`).Error; err != nil {
		return err
	}

	return nil
}

func ensureID(id *uuid.UUID) {
	if *id == uuid.Nil {
		*id = uuid.New()
	}
}

func (r *repository) CreateStream(ctx context.Context, s *domain.Stream) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	ensureID(&s.ID)
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *repository) GetStream(ctx context.Context, id uuid.UUID) (*domain.Stream, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var s domain.Stream
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) ListStreams(ctx context.Context, enabled *bool) ([]*domain.Stream, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var list []*domain.Stream
	q := r.db.WithContext(ctx)
	if enabled != nil {
		q = q.Where("enabled = ?", *enabled)
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) UpdateStream(ctx context.Context, s *domain.Stream) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *repository) DeleteStream(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Stream{}).Error
}

func (r *repository) CreateAlgorithm(ctx context.Context, a *domain.Algorithm) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	ensureID(&a.ID)
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *repository) GetAlgorithm(ctx context.Context, id uuid.UUID) (*domain.Algorithm, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var a domain.Algorithm
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *repository) ListAlgorithms(ctx context.Context) ([]*domain.Algorithm, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var list []*domain.Algorithm
	if err := r.db.WithContext(ctx).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) UpdateAlgorithm(ctx context.Context, a *domain.Algorithm) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(a).Error
}

func (r *repository) DeleteAlgorithm(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Algorithm{}).Error
}

func (r *repository) CreateAlgorithmBinding(ctx context.Context, b *domain.AlgorithmBinding) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	ensureID(&b.ID)
	return r.db.WithContext(ctx).Create(b).Error
}

func (r *repository) GetAlgorithmBinding(ctx context.Context, id uuid.UUID) (*domain.AlgorithmBinding, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var b domain.AlgorithmBinding
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&b).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *repository) ListAlgorithmBindingsByStream(ctx context.Context, streamID uuid.UUID) ([]*domain.AlgorithmBinding, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var list []*domain.AlgorithmBinding
	if err := r.db.WithContext(ctx).Where("stream_id = ?", streamID).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) UpdateAlgorithmBinding(ctx context.Context, b *domain.AlgorithmBinding) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(b).Error
}

func (r *repository) DeleteAlgorithmBinding(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.AlgorithmBinding{}).Error
}

func (r *repository) CreateRecordSession(ctx context.Context, rec *domain.RecordSession) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	ensureID(&rec.ID)
	return r.db.WithContext(ctx).Create(rec).Error
}

func (r *repository) GetRecordSession(ctx context.Context, id uuid.UUID) (*domain.RecordSession, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var rec domain.RecordSession
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&rec).Error; err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *repository) GetRunningRecordSessionByStream(ctx context.Context, streamID uuid.UUID) (*domain.RecordSession, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var rec domain.RecordSession
	if err := r.db.WithContext(ctx).Where("stream_id = ? AND status = ?", streamID, domain.RecordStatusRunning).First(&rec).Error; err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *repository) ListRecordSessionsByStream(ctx context.Context, streamID uuid.UUID) ([]*domain.RecordSession, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var list []*domain.RecordSession
	if err := r.db.WithContext(ctx).Where("stream_id = ?", streamID).Order("started_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) UpdateRecordSession(ctx context.Context, rec *domain.RecordSession) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(rec).Error
}

func (r *repository) CreateInferenceResult(ctx context.Context, ir *domain.InferenceResult) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	ensureID(&ir.ID)
	return r.db.WithContext(ctx).Create(ir).Error
}

func (r *repository) ListInferenceResults(ctx context.Context, streamID, bindingID *uuid.UUID, from, to *int64, limit, offset int) ([]*domain.InferenceResult, int64, error) {
	if err := r.checkDB(); err != nil {
		return nil, 0, err
	}
	q := r.db.WithContext(ctx).Model(&domain.InferenceResult{})
	if streamID != nil {
		q = q.Where("stream_id = ?", *streamID)
	}
	if bindingID != nil {
		q = q.Where("algorithm_binding_id = ?", *bindingID)
	}
	if from != nil {
		q = q.Where("EXTRACT(EPOCH FROM ts) >= ?", *from)
	}
	if to != nil {
		q = q.Where("EXTRACT(EPOCH FROM ts) <= ?", *to)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []*domain.InferenceResult
	if err := q.Limit(limit).Offset(offset).Order("ts DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
