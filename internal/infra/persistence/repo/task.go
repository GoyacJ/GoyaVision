package repo

import (
	"context"

	"goyavision/internal/domain/workflow"
	"goyavision/internal/infra/persistence/mapper"
	"goyavision/internal/infra/persistence/model"
	"goyavision/internal/infra/persistence/scope"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TaskRepo struct {
	db *gorm.DB
}

func NewTaskRepo(db *gorm.DB) *TaskRepo {
	return &TaskRepo{db: db}
}

func (r *TaskRepo) Create(ctx context.Context, t *workflow.Task) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	tenantID, userID := scope.GetContextInfo(ctx)
	m := mapper.TaskToModel(t)
	m.TenantID = tenantID
	m.TriggeredByUserID = &userID

	return r.db.WithContext(ctx).Create(m).Error
}

func (r *TaskRepo) Get(ctx context.Context, id uuid.UUID) (*workflow.Task, error) {
	var m model.TaskModel
	if err := r.db.WithContext(ctx).Scopes(scope.ScopeTenant(ctx)).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.TaskToDomain(&m), nil
}

func (r *TaskRepo) GetWithRelations(ctx context.Context, id uuid.UUID) (*workflow.Task, error) {
	var m model.TaskModel
	if err := r.db.WithContext(ctx).
		Scopes(scope.ScopeTenant(ctx)).
		Preload("Workflow").
		Preload("Asset").
		Preload("Artifacts").
		Where("id = ?", id).
		First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.TaskToDomain(&m), nil
}

func (r *TaskRepo) List(ctx context.Context, filter workflow.TaskFilter) ([]*workflow.Task, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.TaskModel{}).Scopes(scope.ScopeTenant(ctx))

	if filter.WorkflowID != nil {
		q = q.Where("workflow_id = ?", *filter.WorkflowID)
	}
	if filter.AssetID != nil {
		q = q.Where("asset_id = ?", *filter.AssetID)
	}
	if filter.Status != nil {
		q = q.Where("status = ?", string(*filter.Status))
	}
	if filter.TriggeredByUserID != nil {
		q = q.Where("triggered_by_user_id = ?", *filter.TriggeredByUserID)
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

	var models []*model.TaskModel
	if err := q.Limit(filter.Limit).Offset(filter.Offset).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*workflow.Task, len(models))
	for i, m := range models {
		result[i] = mapper.TaskToDomain(m)
	}
	return result, total, nil
}

func (r *TaskRepo) Update(ctx context.Context, t *workflow.Task) error {
	m := mapper.TaskToModel(t)
	return r.db.WithContext(ctx).Scopes(scope.ScopeTenant(ctx)).Where("id = ?", t.ID).Updates(m).Error
}

func (r *TaskRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Scopes(scope.ScopeTenant(ctx)).Where("id = ?", id).Delete(&model.TaskModel{}).Error
}

func (r *TaskRepo) GetStats(ctx context.Context, workflowID *uuid.UUID) (*workflow.TaskStats, error) {
	q := r.db.WithContext(ctx).Model(&model.TaskModel{}).Scopes(scope.ScopeTenant(ctx))
	if workflowID != nil {
		q = q.Where("workflow_id = ?", *workflowID)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, err
	}

	stats := &workflow.TaskStats{Total: total}

	type statusCount struct {
		Status string
		Count  int64
	}
	var counts []statusCount
	sq := r.db.WithContext(ctx).Model(&model.TaskModel{}).Scopes(scope.ScopeTenant(ctx)).
		Select("status, COUNT(*) as count").
		Group("status")
	if workflowID != nil {
		sq = sq.Where("workflow_id = ?", *workflowID)
	}
	if err := sq.Find(&counts).Error; err != nil {
		return nil, err
	}

	for _, sc := range counts {
		switch workflow.TaskStatus(sc.Status) {
		case workflow.TaskStatusPending:
			stats.Pending = sc.Count
		case workflow.TaskStatusRunning:
			stats.Running = sc.Count
		case workflow.TaskStatusSuccess:
			stats.Success = sc.Count
		case workflow.TaskStatusFailed:
			stats.Failed = sc.Count
		case workflow.TaskStatusCancelled:
			stats.Cancelled = sc.Count
		}
	}

	return stats, nil
}

func (r *TaskRepo) ListRunning(ctx context.Context) ([]*workflow.Task, error) {
	var models []*model.TaskModel
	if err := r.db.WithContext(ctx).
		Scopes(scope.ScopeTenant(ctx)).
		Where("status = ?", "running").
		Order("created_at ASC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*workflow.Task, len(models))
	for i, m := range models {
		result[i] = mapper.TaskToDomain(m)
	}
	return result, nil
}

type ArtifactRepo struct {
	db *gorm.DB
}

func NewArtifactRepo(db *gorm.DB) *ArtifactRepo {
	return &ArtifactRepo{db: db}
}

func (r *ArtifactRepo) Create(ctx context.Context, a *workflow.Artifact) error {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	tenantID, _ := scope.GetContextInfo(ctx)
	m := mapper.ArtifactToModel(a)
	m.TenantID = tenantID
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *ArtifactRepo) Get(ctx context.Context, id uuid.UUID) (*workflow.Artifact, error) {
	var m model.ArtifactModel
	if err := r.db.WithContext(ctx).
		Scopes(scope.ScopeTenant(ctx)).
		Preload("Task").
		Preload("Asset").
		Where("id = ?", id).
		First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.ArtifactToDomain(&m), nil
}

func (r *ArtifactRepo) List(ctx context.Context, filter workflow.ArtifactFilter) ([]*workflow.Artifact, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.ArtifactModel{}).Scopes(scope.ScopeTenant(ctx))

	if filter.TaskID != nil {
		q = q.Where("task_id = ?", *filter.TaskID)
	}
	if filter.Type != nil {
		q = q.Where("type = ?", string(*filter.Type))
	}
	if filter.AssetID != nil {
		q = q.Where("asset_id = ?", *filter.AssetID)
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

	var models []*model.ArtifactModel
	if err := q.Limit(filter.Limit).Offset(filter.Offset).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*workflow.Artifact, len(models))
	for i, m := range models {
		result[i] = mapper.ArtifactToDomain(m)
	}
	return result, total, nil
}

func (r *ArtifactRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Scopes(scope.ScopeTenant(ctx)).Where("id = ?", id).Delete(&model.ArtifactModel{}).Error
}

func (r *ArtifactRepo) ListByTask(ctx context.Context, taskID uuid.UUID) ([]*workflow.Artifact, error) {
	var models []*model.ArtifactModel
	if err := r.db.WithContext(ctx).
		Scopes(scope.ScopeTenant(ctx)).
		Preload("Asset").
		Where("task_id = ?", taskID).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*workflow.Artifact, len(models))
	for i, m := range models {
		result[i] = mapper.ArtifactToDomain(m)
	}
	return result, nil
}

func (r *ArtifactRepo) ListByType(ctx context.Context, taskID uuid.UUID, artifactType workflow.ArtifactType) ([]*workflow.Artifact, error) {
	var models []*model.ArtifactModel
	if err := r.db.WithContext(ctx).
		Scopes(scope.ScopeTenant(ctx)).
		Preload("Asset").
		Where("task_id = ? AND type = ?", taskID, string(artifactType)).
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*workflow.Artifact, len(models))
	for i, m := range models {
		result[i] = mapper.ArtifactToDomain(m)
	}
	return result, nil
}
