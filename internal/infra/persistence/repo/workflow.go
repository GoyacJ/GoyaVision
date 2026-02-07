package repo

import (
	"context"
	"encoding/json"

	"goyavision/internal/domain/workflow"
	"goyavision/internal/infra/persistence/mapper"
	"goyavision/internal/infra/persistence/model"
	"goyavision/internal/infra/persistence/scope"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WorkflowRepo struct {
	db *gorm.DB
}

func NewWorkflowRepo(db *gorm.DB) *WorkflowRepo {
	return &WorkflowRepo{db: db}
}

func (r *WorkflowRepo) Create(ctx context.Context, w *workflow.Workflow) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	tenantID, userID := scope.GetContextInfo(ctx)
	m := mapper.WorkflowToModel(w)
	m.TenantID = tenantID
	m.OwnerID = userID

	return r.db.WithContext(ctx).Create(m).Error
}

func (r *WorkflowRepo) Get(ctx context.Context, id uuid.UUID) (*workflow.Workflow, error) {
	var m model.WorkflowModel
	if err := r.db.WithContext(ctx).Scopes(scope.ScopeTenant(ctx), scope.ScopeVisibility(ctx)).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.WorkflowToDomain(&m), nil
}

func (r *WorkflowRepo) GetByCode(ctx context.Context, code string) (*workflow.Workflow, error) {
	var m model.WorkflowModel
	if err := r.db.WithContext(ctx).Scopes(scope.ScopeTenant(ctx), scope.ScopeVisibility(ctx)).Where("code = ?", code).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.WorkflowToDomain(&m), nil
}

func (r *WorkflowRepo) GetWithNodes(ctx context.Context, id uuid.UUID) (*workflow.Workflow, error) {
	var m model.WorkflowModel
	if err := r.db.WithContext(ctx).
		Scopes(scope.ScopeTenant(ctx), scope.ScopeVisibility(ctx)).
		Preload("Nodes.Operator").
		Preload("Edges").
		Where("id = ?", id).
		First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.WorkflowToDomain(&m), nil
}

func (r *WorkflowRepo) List(ctx context.Context, filter workflow.Filter) ([]*workflow.Workflow, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.WorkflowModel{}).Scopes(scope.ScopeTenant(ctx), scope.ScopeVisibility(ctx))

	if filter.Status != nil {
		q = q.Where("status = ?", string(*filter.Status))
	}
	if filter.TriggerType != nil {
		q = q.Where("trigger_type = ?", string(*filter.TriggerType))
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

	var models []*model.WorkflowModel
	if err := q.Limit(filter.Limit).Offset(filter.Offset).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}

	result := make([]*workflow.Workflow, len(models))
	for i, m := range models {
		result[i] = mapper.WorkflowToDomain(m)
	}
	return result, total, nil
}

func (r *WorkflowRepo) Update(ctx context.Context, w *workflow.Workflow) error {
	m := mapper.WorkflowToModel(w)
	return r.db.WithContext(ctx).Scopes(scope.ScopeTenant(ctx)).Where("id = ?", w.ID).Updates(m).Error
}

func (r *WorkflowRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Scopes(scope.ScopeTenant(ctx)).Where("id = ?", id).Delete(&model.WorkflowModel{}).Error
}

func (r *WorkflowRepo) ListEnabled(ctx context.Context) ([]*workflow.Workflow, error) {
	var models []*model.WorkflowModel
	if err := r.db.WithContext(ctx).
		Scopes(scope.ScopeTenant(ctx), scope.ScopeVisibility(ctx)).
		Preload("Nodes.Operator").
		Preload("Edges").
		Where("status = ?", "enabled").
		Order("created_at DESC").
		Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*workflow.Workflow, len(models))
	for i, m := range models {
		result[i] = mapper.WorkflowToDomain(m)
	}
	return result, nil
}

func (r *WorkflowRepo) CreateNode(ctx context.Context, n *workflow.Node) error {
	if n.ID == uuid.Nil {
		n.ID = uuid.New()
	}
	m := mapper.NodeToModel(n)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *WorkflowRepo) ListNodes(ctx context.Context, workflowID uuid.UUID) ([]*workflow.Node, error) {
	var models []*model.WorkflowNodeModel
	if err := r.db.WithContext(ctx).Preload("Operator").Where("workflow_id = ?", workflowID).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*workflow.Node, len(models))
	for i, m := range models {
		result[i] = mapper.NodeToDomain(m)
	}
	return result, nil
}

func (r *WorkflowRepo) DeleteNodes(ctx context.Context, workflowID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("workflow_id = ?", workflowID).Delete(&model.WorkflowNodeModel{}).Error
}

func (r *WorkflowRepo) CreateEdge(ctx context.Context, e *workflow.Edge) error {
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	m := mapper.EdgeToModel(e)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *WorkflowRepo) ListEdges(ctx context.Context, workflowID uuid.UUID) ([]*workflow.Edge, error) {
	var models []*model.WorkflowEdgeModel
	if err := r.db.WithContext(ctx).Where("workflow_id = ?", workflowID).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*workflow.Edge, len(models))
	for i, m := range models {
		result[i] = mapper.EdgeToDomain(m)
	}
	return result, nil
}

func (r *WorkflowRepo) DeleteEdges(ctx context.Context, workflowID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("workflow_id = ?", workflowID).Delete(&model.WorkflowEdgeModel{}).Error
}
