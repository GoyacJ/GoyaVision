package repo

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"goyavision/internal/domain/workflow"
	"goyavision/internal/infra/persistence/mapper"
	"goyavision/internal/infra/persistence/model"
	"goyavision/internal/infra/persistence/scope"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const contextSnapshotInterval int64 = 20

type ContextRepo struct {
	db *gorm.DB
}

func NewContextRepo(db *gorm.DB) *ContextRepo {
	return &ContextRepo{db: db}
}

func (r *ContextRepo) InitializeState(ctx context.Context, state *workflow.TaskContextState) error {
	if state.TaskID == uuid.Nil {
		return errors.New("task_id is required")
	}
	if state.Version <= 0 {
		state.Version = 1
	}
	if state.Data == nil {
		state.Data = map[string]interface{}{}
	}
	if err := r.ensureTaskVisible(ctx, state.TaskID); err != nil {
		return err
	}
	m := mapper.TaskContextStateToModel(state)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *ContextRepo) GetState(ctx context.Context, taskID uuid.UUID) (*workflow.TaskContextState, error) {
	var m model.TaskContextStateModel
	if err := r.db.WithContext(ctx).
		Table("task_context_state AS s").
		Select("s.*").
		Joins("JOIN tasks t ON t.id = s.task_id").
		Scopes(scope.ScopeTenantOnly(ctx)).
		Where("s.task_id = ?", taskID).
		First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.TaskContextStateToDomain(&m), nil
}

func (r *ContextRepo) ApplyPatch(ctx context.Context, patch *workflow.TaskContextPatch) error {
	if patch.TaskID == uuid.Nil {
		return errors.New("task_id is required")
	}
	if patch.ID == uuid.Nil {
		patch.ID = uuid.New()
	}

	db := r.db.WithContext(ctx)
	var state model.TaskContextStateModel
	if err := db.
		Table("task_context_state AS s").
		Select("s.*").
		Joins("JOIN tasks t ON t.id = s.task_id").
		Scopes(scope.ScopeTenantOnly(ctx)).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("s.task_id = ?", patch.TaskID).
		First(&state).Error; err != nil {
		return err
	}
	if state.Version != patch.BeforeVersion {
		return workflow.ErrContextVersionConflict
	}

	stateData := map[string]interface{}{}
	if len(state.Data) > 0 {
		_ = json.Unmarshal(state.Data, &stateData)
	}
	applyDiff(stateData, patch.Diff)

	patch.AfterVersion = patch.BeforeVersion + 1
	patchModel := mapper.TaskContextPatchToModel(patch)
	if err := db.Create(patchModel).Error; err != nil {
		return err
	}

	raw, err := json.Marshal(stateData)
	if err != nil {
		return err
	}

	result := db.Model(&model.TaskContextStateModel{}).
		Where("task_id = ? AND version = ?", patch.TaskID, patch.BeforeVersion).
		Updates(map[string]interface{}{
			"version":    patch.AfterVersion,
			"data":       datatypes.JSON(raw),
			"updated_at": time.Now(),
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return workflow.ErrContextVersionConflict
	}

	if patch.AfterVersion%contextSnapshotInterval == 0 {
		snapshot := &workflow.TaskContextSnapshot{
			ID:      uuid.New(),
			TaskID:  patch.TaskID,
			Version: patch.AfterVersion,
			Data:    stateData,
			Trigger: "periodic",
		}
		if err := r.CreateSnapshot(ctx, snapshot); err != nil {
			return err
		}
	}

	return nil
}

func (r *ContextRepo) CreateSnapshot(ctx context.Context, snapshot *workflow.TaskContextSnapshot) error {
	if snapshot.ID == uuid.Nil {
		snapshot.ID = uuid.New()
	}
	if snapshot.Trigger == "" {
		snapshot.Trigger = "manual"
	}
	if err := r.ensureTaskVisible(ctx, snapshot.TaskID); err != nil {
		return err
	}
	m := mapper.TaskContextSnapshotToModel(snapshot)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *ContextRepo) ListPatches(ctx context.Context, taskID uuid.UUID, limit, offset int) ([]*workflow.TaskContextPatch, int64, error) {
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}
	if offset < 0 {
		offset = 0
	}

	q := r.db.WithContext(ctx).
		Table("task_context_patches AS p").
		Select("p.*").
		Joins("JOIN tasks t ON t.id = p.task_id").
		Scopes(scope.ScopeTenantOnly(ctx)).
		Where("p.task_id = ?", taskID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []*model.TaskContextPatchModel
	if err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	out := make([]*workflow.TaskContextPatch, len(models))
	for i := range models {
		out[i] = mapper.TaskContextPatchToDomain(models[i])
	}
	return out, total, nil
}

func (r *ContextRepo) ensureTaskVisible(ctx context.Context, taskID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.TaskModel{}).
		Scopes(scope.ScopeTenantOnly(ctx)).
		Where("id = ?", taskID).
		Take(&model.TaskModel{}).Error
}

func applyDiff(target map[string]interface{}, diff workflow.ContextDiff) {
	for _, path := range diff.Unset {
		unsetByPath(target, path)
	}
	for path, value := range diff.Set {
		setByPath(target, path, value)
	}
}

func setByPath(target map[string]interface{}, path string, value interface{}) {
	path = strings.TrimSpace(path)
	if path == "" {
		return
	}
	parts := strings.Split(path, ".")
	cur := target
	for i := 0; i < len(parts)-1; i++ {
		key := parts[i]
		next, ok := cur[key].(map[string]interface{})
		if !ok {
			next = map[string]interface{}{}
			cur[key] = next
		}
		cur = next
	}
	cur[parts[len(parts)-1]] = value
}

func unsetByPath(target map[string]interface{}, path string) {
	path = strings.TrimSpace(path)
	if path == "" {
		return
	}
	parts := strings.Split(path, ".")
	cur := target
	for i := 0; i < len(parts)-1; i++ {
		next, ok := cur[parts[i]].(map[string]interface{})
		if !ok {
			return
		}
		cur = next
	}
	delete(cur, parts[len(parts)-1])
}
