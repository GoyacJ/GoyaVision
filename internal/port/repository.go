package port

import (
	"context"
	"time"

	"goyavision/internal/domain"
	"goyavision/internal/domain/ai_model"
	"goyavision/internal/domain/identity"
	"goyavision/internal/domain/media"
	"goyavision/internal/domain/operator"
	"goyavision/internal/domain/storage"
	"goyavision/internal/domain/system"
	"goyavision/internal/domain/workflow"

	"github.com/google/uuid"
)

type Repository interface {
	// User
	CreateUser(ctx context.Context, u *identity.User) error
	GetUser(ctx context.Context, id uuid.UUID) (*identity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*identity.User, error)
	GetUserWithRoles(ctx context.Context, id uuid.UUID) (*identity.User, error)
	ListUsers(ctx context.Context, status *int, limit, offset int) ([]*identity.User, int64, error)
	UpdateUser(ctx context.Context, u *identity.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
	SetUserRoles(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error

	// Role
	CreateRole(ctx context.Context, r *identity.Role) error
	GetRole(ctx context.Context, id uuid.UUID) (*identity.Role, error)
	GetRoleByCode(ctx context.Context, code string) (*identity.Role, error)
	GetRoleWithPermissions(ctx context.Context, id uuid.UUID) (*identity.Role, error)
	ListRoles(ctx context.Context, status *int) ([]*identity.Role, error)
	UpdateRole(ctx context.Context, r *identity.Role) error
	DeleteRole(ctx context.Context, id uuid.UUID) error
	SetRolePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error
	SetRoleMenus(ctx context.Context, roleID uuid.UUID, menuIDs []uuid.UUID) error
	GetDefaultRoles(ctx context.Context) ([]*identity.Role, error)

	// UserIdentity
	CreateUserIdentity(ctx context.Context, i *identity.UserIdentity) error
	GetUserIdentity(ctx context.Context, id uuid.UUID) (*identity.UserIdentity, error)
	GetUserIdentityByIdentifier(ctx context.Context, identityType identity.IdentityType, identifier string) (*identity.UserIdentity, error)
	ListUserIdentities(ctx context.Context, userID uuid.UUID) ([]*identity.UserIdentity, error)
	DeleteUserIdentity(ctx context.Context, id uuid.UUID) error

	// Permission
	CreatePermission(ctx context.Context, p *identity.Permission) error
	GetPermission(ctx context.Context, id uuid.UUID) (*identity.Permission, error)
	GetPermissionByCode(ctx context.Context, code string) (*identity.Permission, error)
	ListPermissions(ctx context.Context) ([]*identity.Permission, error)
	UpdatePermission(ctx context.Context, p *identity.Permission) error
	DeletePermission(ctx context.Context, id uuid.UUID) error
	GetPermissionsByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]*identity.Permission, error)

	// Menu
	CreateMenu(ctx context.Context, m *identity.Menu) error
	GetMenu(ctx context.Context, id uuid.UUID) (*identity.Menu, error)
	GetMenuByCode(ctx context.Context, code string) (*identity.Menu, error)
	ListMenus(ctx context.Context, status *int) ([]*identity.Menu, error)
	UpdateMenu(ctx context.Context, m *identity.Menu) error
	DeleteMenu(ctx context.Context, id uuid.UUID) error
	GetMenusByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]*identity.Menu, error)

	// MediaAsset
	CreateMediaAsset(ctx context.Context, a *media.Asset) error
	GetMediaAsset(ctx context.Context, id uuid.UUID) (*media.Asset, error)
	ListMediaAssets(ctx context.Context, filter media.AssetFilter) ([]*media.Asset, int64, error)
	UpdateMediaAsset(ctx context.Context, a *media.Asset) error
	DeleteMediaAsset(ctx context.Context, id uuid.UUID) error
	ListMediaAssetsBySource(ctx context.Context, sourceID uuid.UUID) ([]*media.Asset, error)
	ListMediaAssetsByParent(ctx context.Context, parentID uuid.UUID) ([]*media.Asset, error)
	GetAllAssetTags(ctx context.Context) ([]string, error) // 获取所有标签

	// MediaSource
	CreateMediaSource(ctx context.Context, s *media.Source) error
	GetMediaSource(ctx context.Context, id uuid.UUID) (*media.Source, error)
	GetMediaSourceByPathName(ctx context.Context, pathName string) (*media.Source, error)
	ListMediaSources(ctx context.Context, filter media.SourceFilter) ([]*media.Source, int64, error)
	UpdateMediaSource(ctx context.Context, s *media.Source) error
	DeleteMediaSource(ctx context.Context, id uuid.UUID) error

	// Operator
	CreateOperator(ctx context.Context, o *operator.Operator) error
	GetOperator(ctx context.Context, id uuid.UUID) (*operator.Operator, error)
	GetOperatorByCode(ctx context.Context, code string) (*operator.Operator, error)
	ListOperators(ctx context.Context, filter operator.Filter) ([]*operator.Operator, int64, error)
	UpdateOperator(ctx context.Context, o *operator.Operator) error
	DeleteOperator(ctx context.Context, id uuid.UUID) error
	ListEnabledOperators(ctx context.Context) ([]*operator.Operator, error)
	ListOperatorsByCategory(ctx context.Context, category operator.Category) ([]*operator.Operator, error)

	// Workflow
	CreateWorkflow(ctx context.Context, w *workflow.Workflow) error
	GetWorkflow(ctx context.Context, id uuid.UUID) (*workflow.Workflow, error)
	GetWorkflowByCode(ctx context.Context, code string) (*workflow.Workflow, error)
	GetWorkflowWithNodes(ctx context.Context, id uuid.UUID) (*workflow.Workflow, error)
	ListWorkflows(ctx context.Context, filter workflow.Filter) ([]*workflow.Workflow, int64, error)
	UpdateWorkflow(ctx context.Context, w *workflow.Workflow) error
	DeleteWorkflow(ctx context.Context, id uuid.UUID) error
	ListEnabledWorkflows(ctx context.Context) ([]*workflow.Workflow, error)

	// WorkflowNode
	CreateWorkflowNode(ctx context.Context, n *workflow.Node) error
	ListWorkflowNodes(ctx context.Context, workflowID uuid.UUID) ([]*workflow.Node, error)
	DeleteWorkflowNodes(ctx context.Context, workflowID uuid.UUID) error

	// WorkflowEdge
	CreateWorkflowEdge(ctx context.Context, e *workflow.Edge) error
	ListWorkflowEdges(ctx context.Context, workflowID uuid.UUID) ([]*workflow.Edge, error)
	DeleteWorkflowEdges(ctx context.Context, workflowID uuid.UUID) error

	// Task
	CreateTask(ctx context.Context, t *workflow.Task) error
	GetTask(ctx context.Context, id uuid.UUID) (*workflow.Task, error)
	GetTaskWithRelations(ctx context.Context, id uuid.UUID) (*workflow.Task, error)
	ListTasks(ctx context.Context, filter workflow.TaskFilter) ([]*workflow.Task, int64, error)
	UpdateTask(ctx context.Context, t *workflow.Task) error
	DeleteTask(ctx context.Context, id uuid.UUID) error
	GetTaskStats(ctx context.Context, workflowID *uuid.UUID) (*workflow.TaskStats, error)
	ListRunningTasks(ctx context.Context) ([]*workflow.Task, error)

	// Artifact
	CreateArtifact(ctx context.Context, a *workflow.Artifact) error
	GetArtifact(ctx context.Context, id uuid.UUID) (*workflow.Artifact, error)
	ListArtifacts(ctx context.Context, filter workflow.ArtifactFilter) ([]*workflow.Artifact, int64, error)
	DeleteArtifact(ctx context.Context, id uuid.UUID) error
	ListArtifactsByTask(ctx context.Context, taskID uuid.UUID) ([]*workflow.Artifact, error)
	ListArtifactsByType(ctx context.Context, taskID uuid.UUID, artifactType workflow.ArtifactType) ([]*workflow.Artifact, error)

	// File
	CreateFile(ctx context.Context, f *storage.File) error
	GetFile(ctx context.Context, id uuid.UUID) (*storage.File, error)
	GetFileByPath(ctx context.Context, path string) (*storage.File, error)
	ListFiles(ctx context.Context, filter storage.FileFilter) ([]*storage.File, int64, error)
	UpdateFile(ctx context.Context, f *storage.File) error
	DeleteFile(ctx context.Context, id uuid.UUID) error
	GetFileByHash(ctx context.Context, hash string) (*storage.File, error)

	// AIModel
	CreateAIModel(ctx context.Context, d *ai_model.AIModel) error
	GetAIModel(ctx context.Context, id uuid.UUID) (*ai_model.AIModel, error)
	UpdateAIModel(ctx context.Context, d *ai_model.AIModel) error
	DeleteAIModel(ctx context.Context, id uuid.UUID) error
	ListAIModels(ctx context.Context, filter ai_model.Filter) ([]*ai_model.AIModel, int64, error)

	// SystemConfig
	GetSystemConfig(ctx context.Context, key string) (*system.SystemConfig, error)
	ListSystemConfigs(ctx context.Context) ([]*system.SystemConfig, error)
	SaveSystemConfig(ctx context.Context, config *system.SystemConfig) error
	DeleteSystemConfig(ctx context.Context, key string) error

	UserAssetRepository
}

type UserAssetRepository interface {
	GetUserBalance(ctx context.Context, userID uuid.UUID) (*domain.UserBalance, error)
	UpdateUserBalance(ctx context.Context, ub *domain.UserBalance) error
	CreateTransactionRecord(ctx context.Context, tr *domain.TransactionRecord) error
	GetTransactionRecord(ctx context.Context, id string) (*domain.TransactionRecord, error)
	UpdateTransactionRecord(ctx context.Context, tr *domain.TransactionRecord) error
	ListTransactionRecords(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.TransactionRecord, int64, error)
	CreatePointRecord(ctx context.Context, pr *domain.PointRecord) error
	ListPointRecords(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.PointRecord, int64, error)
	GetUserSubscription(ctx context.Context, userID uuid.UUID) (*domain.UserSubscription, error)
	UpdateUserSubscription(ctx context.Context, us *domain.UserSubscription) error
	GetUsageStats(ctx context.Context, userID uuid.UUID, date time.Time) (*domain.UsageStats, error)
	UpdateUsageStats(ctx context.Context, us *domain.UsageStats) error
	ListUsageStats(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]*domain.UsageStats, error)
}
