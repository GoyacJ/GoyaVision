package repo

import (
	"context"
	"time"

	"goyavision/internal/domain/agent"
	"goyavision/internal/infra/persistence/mapper"
	"goyavision/internal/infra/persistence/model"
	"goyavision/internal/infra/persistence/scope"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AgentSessionRepo struct {
	db *gorm.DB
}

func NewAgentSessionRepo(db *gorm.DB) *AgentSessionRepo {
	return &AgentSessionRepo{db: db}
}

func (r *AgentSessionRepo) Create(ctx context.Context, session *agent.Session) error {
	if session.ID == uuid.Nil {
		session.ID = uuid.New()
	}
	if session.StartedAt.IsZero() {
		session.StartedAt = time.Now().UTC()
	}
	if err := r.ensureTaskVisible(ctx, session.TaskID); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(mapper.AgentSessionToModel(session)).Error
}

func (r *AgentSessionRepo) Get(ctx context.Context, id uuid.UUID) (*agent.Session, error) {
	var m model.AgentSessionModel
	if err := r.db.WithContext(ctx).
		Table("agent_sessions AS s").
		Select("s.*").
		Joins("JOIN tasks t ON t.id = s.task_id").
		Scopes(scope.ScopeTenantOnly(ctx)).
		Where("s.id = ?", id).
		First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.AgentSessionToDomain(&m), nil
}

func (r *AgentSessionRepo) Update(ctx context.Context, session *agent.Session) error {
	visibleTasks := r.db.WithContext(ctx).
		Model(&model.TaskModel{}).
		Scopes(scope.ScopeTenantOnly(ctx)).
		Select("id")
	return r.db.WithContext(ctx).
		Where("id = ? AND task_id IN (?)", session.ID, visibleTasks).
		Updates(mapper.AgentSessionToModel(session)).Error
}

func (r *AgentSessionRepo) List(ctx context.Context, filter agent.SessionFilter) ([]*agent.Session, int64, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 1000 {
		limit = 1000
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	base := r.db.WithContext(ctx).
		Table("agent_sessions AS s").
		Joins("JOIN tasks t ON t.id = s.task_id").
		Scopes(scope.ScopeTenantOnly(ctx))

	if filter.TaskID != nil {
		base = base.Where("s.task_id = ?", *filter.TaskID)
	}
	if filter.Status != nil {
		base = base.Where("s.status = ?", string(*filter.Status))
	}

	var total int64
	if err := base.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []*model.AgentSessionModel
	if err := base.Session(&gorm.Session{}).
		Select("s.*").
		Order("s.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}
	out := make([]*agent.Session, len(models))
	for i := range models {
		out[i] = mapper.AgentSessionToDomain(models[i])
	}
	return out, total, nil
}

type RunEventRepo struct {
	db *gorm.DB
}

func NewRunEventRepo(db *gorm.DB) *RunEventRepo {
	return &RunEventRepo{db: db}
}

func (r *RunEventRepo) Create(ctx context.Context, event *agent.RunEvent) error {
	if event.ID == uuid.Nil {
		event.ID = uuid.New()
	}
	if err := r.ensureTaskVisible(ctx, event.TaskID); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Create(mapper.RunEventToModel(event)).Error
}

func (r *RunEventRepo) List(ctx context.Context, filter agent.EventFilter) ([]*agent.RunEvent, int64, error) {
	limit := filter.Limit
	if limit <= 0 {
		limit = 100
	}
	if limit > 1000 {
		limit = 1000
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	base := r.db.WithContext(ctx).
		Table("run_events AS e").
		Joins("JOIN tasks t ON t.id = e.task_id").
		Scopes(scope.ScopeTenantOnly(ctx))
	if filter.TaskID != nil {
		base = base.Where("e.task_id = ?", *filter.TaskID)
	}
	if filter.SessionID != nil {
		base = base.Where("e.session_id = ?", *filter.SessionID)
	}
	if filter.Source != "" {
		base = base.Where("e.source = ?", filter.Source)
	}
	if filter.NodeKey != "" {
		base = base.Where("e.node_key = ?", filter.NodeKey)
	}

	var total int64
	if err := base.Session(&gorm.Session{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var models []*model.RunEventModel
	if err := base.Session(&gorm.Session{}).
		Select("e.*").
		Order("e.created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&models).Error; err != nil {
		return nil, 0, err
	}
	out := make([]*agent.RunEvent, len(models))
	for i := range models {
		out[i] = mapper.RunEventToDomain(models[i])
	}
	return out, total, nil
}

type ToolPolicyRepo struct {
	db *gorm.DB
}

func NewToolPolicyRepo(db *gorm.DB) *ToolPolicyRepo {
	return &ToolPolicyRepo{db: db}
}

func (r *ToolPolicyRepo) Upsert(ctx context.Context, policy *agent.ToolPolicy) error {
	if policy.ID == uuid.Nil {
		policy.ID = uuid.New()
	}
	m := mapper.ToolPolicyToModel(policy)
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "tool_name"}},
			DoUpdates: clause.AssignmentColumns([]string{"risk_level", "permissions", "data_access", "determinism", "limits", "enabled", "updated_at"}),
		}).
		Create(m).Error
}

func (r *ToolPolicyRepo) GetByToolName(ctx context.Context, toolName string) (*agent.ToolPolicy, error) {
	var m model.ToolPolicyModel
	if err := r.db.WithContext(ctx).Where("tool_name = ?", toolName).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.ToolPolicyToDomain(&m), nil
}

func (r *ToolPolicyRepo) ListEnabled(ctx context.Context) ([]*agent.ToolPolicy, error) {
	var models []*model.ToolPolicyModel
	if err := r.db.WithContext(ctx).Where("enabled = ?", true).Order("tool_name ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]*agent.ToolPolicy, len(models))
	for i := range models {
		out[i] = mapper.ToolPolicyToDomain(models[i])
	}
	return out, nil
}

func (r *AgentSessionRepo) ensureTaskVisible(ctx context.Context, taskID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.TaskModel{}).
		Scopes(scope.ScopeTenantOnly(ctx)).
		Where("id = ?", taskID).
		Take(&model.TaskModel{}).Error
}

func (r *RunEventRepo) ensureTaskVisible(ctx context.Context, taskID uuid.UUID) error {
	return r.db.WithContext(ctx).
		Model(&model.TaskModel{}).
		Scopes(scope.ScopeTenantOnly(ctx)).
		Where("id = ?", taskID).
		Take(&model.TaskModel{}).Error
}
