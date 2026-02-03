package port

import (
	"context"

	"goyavision/internal/domain"

	"github.com/google/uuid"
)

type Repository interface {
	// User
	CreateUser(ctx context.Context, u *domain.User) error
	GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetUserByUsername(ctx context.Context, username string) (*domain.User, error)
	GetUserWithRoles(ctx context.Context, id uuid.UUID) (*domain.User, error)
	ListUsers(ctx context.Context, status *int, limit, offset int) ([]*domain.User, int64, error)
	UpdateUser(ctx context.Context, u *domain.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	SetUserRoles(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error

	// Role
	CreateRole(ctx context.Context, r *domain.Role) error
	GetRole(ctx context.Context, id uuid.UUID) (*domain.Role, error)
	GetRoleByCode(ctx context.Context, code string) (*domain.Role, error)
	GetRoleWithPermissions(ctx context.Context, id uuid.UUID) (*domain.Role, error)
	ListRoles(ctx context.Context, status *int) ([]*domain.Role, error)
	UpdateRole(ctx context.Context, r *domain.Role) error
	DeleteRole(ctx context.Context, id uuid.UUID) error
	SetRolePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
	SetRoleMenus(ctx context.Context, roleID uuid.UUID, menuIDs []uuid.UUID) error

	// Permission
	CreatePermission(ctx context.Context, p *domain.Permission) error
	GetPermission(ctx context.Context, id uuid.UUID) (*domain.Permission, error)
	GetPermissionByCode(ctx context.Context, code string) (*domain.Permission, error)
	ListPermissions(ctx context.Context) ([]*domain.Permission, error)
	UpdatePermission(ctx context.Context, p *domain.Permission) error
	DeletePermission(ctx context.Context, id uuid.UUID) error
	GetPermissionsByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]*domain.Permission, error)

	// Menu
	CreateMenu(ctx context.Context, m *domain.Menu) error
	GetMenu(ctx context.Context, id uuid.UUID) (*domain.Menu, error)
	GetMenuByCode(ctx context.Context, code string) (*domain.Menu, error)
	ListMenus(ctx context.Context, status *int) ([]*domain.Menu, error)
	UpdateMenu(ctx context.Context, m *domain.Menu) error
	DeleteMenu(ctx context.Context, id uuid.UUID) error
	GetMenusByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]*domain.Menu, error)

	// MediaAsset
	CreateMediaAsset(ctx context.Context, a *domain.MediaAsset) error
	GetMediaAsset(ctx context.Context, id uuid.UUID) (*domain.MediaAsset, error)
	ListMediaAssets(ctx context.Context, filter domain.MediaAssetFilter) ([]*domain.MediaAsset, int64, error)
	UpdateMediaAsset(ctx context.Context, a *domain.MediaAsset) error
	DeleteMediaAsset(ctx context.Context, id uuid.UUID) error
	ListMediaAssetsBySource(ctx context.Context, sourceID uuid.UUID) ([]*domain.MediaAsset, error)
	ListMediaAssetsByParent(ctx context.Context, parentID uuid.UUID) ([]*domain.MediaAsset, error)
	GetAllAssetTags(ctx context.Context) ([]string, error) // 获取所有标签

	// Operator
	CreateOperator(ctx context.Context, o *domain.Operator) error
	GetOperator(ctx context.Context, id uuid.UUID) (*domain.Operator, error)
	GetOperatorByCode(ctx context.Context, code string) (*domain.Operator, error)
	ListOperators(ctx context.Context, filter domain.OperatorFilter) ([]*domain.Operator, int64, error)
	UpdateOperator(ctx context.Context, o *domain.Operator) error
	DeleteOperator(ctx context.Context, id uuid.UUID) error
	ListEnabledOperators(ctx context.Context) ([]*domain.Operator, error)
	ListOperatorsByCategory(ctx context.Context, category domain.OperatorCategory) ([]*domain.Operator, error)

	// Workflow
	CreateWorkflow(ctx context.Context, w *domain.Workflow) error
	GetWorkflow(ctx context.Context, id uuid.UUID) (*domain.Workflow, error)
	GetWorkflowByCode(ctx context.Context, code string) (*domain.Workflow, error)
	GetWorkflowWithNodes(ctx context.Context, id uuid.UUID) (*domain.Workflow, error)
	ListWorkflows(ctx context.Context, filter domain.WorkflowFilter) ([]*domain.Workflow, int64, error)
	UpdateWorkflow(ctx context.Context, w *domain.Workflow) error
	DeleteWorkflow(ctx context.Context, id uuid.UUID) error
	ListEnabledWorkflows(ctx context.Context) ([]*domain.Workflow, error)

	// WorkflowNode
	CreateWorkflowNode(ctx context.Context, n *domain.WorkflowNode) error
	ListWorkflowNodes(ctx context.Context, workflowID uuid.UUID) ([]*domain.WorkflowNode, error)
	DeleteWorkflowNodes(ctx context.Context, workflowID uuid.UUID) error

	// WorkflowEdge
	CreateWorkflowEdge(ctx context.Context, e *domain.WorkflowEdge) error
	ListWorkflowEdges(ctx context.Context, workflowID uuid.UUID) ([]*domain.WorkflowEdge, error)
	DeleteWorkflowEdges(ctx context.Context, workflowID uuid.UUID) error

	// Task
	CreateTask(ctx context.Context, t *domain.Task) error
	GetTask(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	GetTaskWithRelations(ctx context.Context, id uuid.UUID) (*domain.Task, error)
	ListTasks(ctx context.Context, filter domain.TaskFilter) ([]*domain.Task, int64, error)
	UpdateTask(ctx context.Context, t *domain.Task) error
	DeleteTask(ctx context.Context, id uuid.UUID) error
	GetTaskStats(ctx context.Context, workflowID *uuid.UUID) (*domain.TaskStats, error)
	ListRunningTasks(ctx context.Context) ([]*domain.Task, error)

	// Artifact
	CreateArtifact(ctx context.Context, a *domain.Artifact) error
	GetArtifact(ctx context.Context, id uuid.UUID) (*domain.Artifact, error)
	ListArtifacts(ctx context.Context, filter domain.ArtifactFilter) ([]*domain.Artifact, int64, error)
	DeleteArtifact(ctx context.Context, id uuid.UUID) error
	ListArtifactsByTask(ctx context.Context, taskID uuid.UUID) ([]*domain.Artifact, error)
	ListArtifactsByType(ctx context.Context, taskID uuid.UUID, artifactType domain.ArtifactType) ([]*domain.Artifact, error)
}
