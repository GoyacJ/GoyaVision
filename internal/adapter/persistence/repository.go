package persistence

import (
	"context"
	"errors"
	"time"

	"goyavision/internal/domain"
	"goyavision/internal/domain/ai_model"
	"goyavision/internal/domain/identity"
	"goyavision/internal/domain/media"
	"goyavision/internal/domain/operator"
	"goyavision/internal/domain/storage"
	"goyavision/internal/domain/system"
	"goyavision/internal/domain/workflow"
	"goyavision/internal/infra/persistence/model"
	"goyavision/internal/infra/persistence/repo"
	"goyavision/internal/port"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrDBNotConfigured = errors.New("database not configured")

var _ port.Repository = (*repository)(nil)

type repository struct {
	db *gorm.DB

	users       *repo.UserRepo
	roles       *repo.RoleRepo
	permissions *repo.PermissionRepo
	menus       *repo.MenuRepo

	assets  *repo.MediaAssetRepo
	sources *repo.MediaSourceRepo

	operators  *repo.OperatorRepo
	algorithms *repo.AlgorithmRepo

	workflows     *repo.WorkflowRepo
	tasks         *repo.TaskRepo
	artifacts     *repo.ArtifactRepo
	contexts      *repo.ContextRepo
	agentSessions *repo.AgentSessionRepo
	runEvents     *repo.RunEventRepo
	toolPolicies  *repo.ToolPolicyRepo

	files          *repo.FileRepo
	aiModels       *repo.AIModelRepo
	userIdentities *repo.UserIdentityRepo
	systemConfigs  *repo.SystemConfigRepo
	userAssets     *repo.UserAssetRepo
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{
		db:             db,
		users:          repo.NewUserRepo(db),
		roles:          repo.NewRoleRepo(db),
		permissions:    repo.NewPermissionRepo(db),
		menus:          repo.NewMenuRepo(db),
		assets:         repo.NewMediaAssetRepo(db),
		sources:        repo.NewMediaSourceRepo(db),
		operators:      repo.NewOperatorRepo(db),
		algorithms:     repo.NewAlgorithmRepo(db),
		workflows:      repo.NewWorkflowRepo(db),
		tasks:          repo.NewTaskRepo(db),
		artifacts:      repo.NewArtifactRepo(db),
		contexts:       repo.NewContextRepo(db),
		agentSessions:  repo.NewAgentSessionRepo(db),
		runEvents:      repo.NewRunEventRepo(db),
		toolPolicies:   repo.NewToolPolicyRepo(db),
		files:          repo.NewFileRepo(db),
		aiModels:       repo.NewAIModelRepo(db),
		userIdentities: repo.NewUserIdentityRepo(db),
		systemConfigs:  repo.NewSystemConfigRepo(db),
		userAssets:     repo.NewUserAssetRepo(db),
	}
}

func (r *repository) checkDB() error {
	if r.db == nil {
		return ErrDBNotConfigured
	}
	return nil
}

func AutoMigrate(db *gorm.DB) error {
	if db == nil {
		return ErrDBNotConfigured
	}
	return db.AutoMigrate(
		&model.TenantModel{},
		&model.UserModel{},
		&model.RoleModel{},
		&model.PermissionModel{},
		&model.MenuModel{},
		&model.MediaAssetModel{},
		&model.MediaSourceModel{},
		&model.OperatorModel{},
		&model.OperatorVersionModel{},
		&model.OperatorTemplateModel{},
		&model.OperatorDependencyModel{},
		&model.AlgorithmModel{},
		&model.AlgorithmVersionModel{},
		&model.AlgorithmImplementationModel{},
		&model.AlgorithmEvaluationModel{},
		&model.WorkflowModel{},
		&model.WorkflowRevisionModel{},
		&model.WorkflowNodeModel{},
		&model.WorkflowEdgeModel{},
		&model.TaskModel{},
		&model.TaskContextStateModel{},
		&model.TaskContextPatchModel{},
		&model.TaskContextSnapshotModel{},
		&model.AgentSessionModel{},
		&model.RunEventModel{},
		&model.ToolPolicyModel{},
		&model.ArtifactModel{},
		&model.FileModel{},
		&model.AIModelModel{},
		&model.UserIdentityModel{},
		&model.SystemConfigModel{},
		&model.UserBalance{},
		&model.UserSubscription{},
		&model.TransactionRecord{},
		&model.PointRecord{},
		&model.UsageStat{},
	)
}

// User methods
func (r *repository) CreateUser(ctx context.Context, u *identity.User) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.users.Create(ctx, u)
}

func (r *repository) GetUser(ctx context.Context, id uuid.UUID) (*identity.User, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.users.Get(ctx, id)
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (*identity.User, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.users.GetByUsername(ctx, username)
}

func (r *repository) GetUserWithRoles(ctx context.Context, id uuid.UUID) (*identity.User, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.users.GetWithRoles(ctx, id)
}

func (r *repository) ListUsers(ctx context.Context, status *int, limit, offset int) ([]*identity.User, int64, error) {
	if err := r.checkDB(); err != nil {
		return nil, 0, err
	}
	return r.users.List(ctx, status, limit, offset)
}

func (r *repository) UpdateUser(ctx context.Context, u *identity.User) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.users.Update(ctx, u)
}

func (r *repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.users.Delete(ctx, id)
}

func (r *repository) SetUserRoles(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.users.SetUserRoles(ctx, userID, roleIDs)
}

// Role methods
func (r *repository) CreateRole(ctx context.Context, role *identity.Role) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.roles.Create(ctx, role)
}

func (r *repository) GetRole(ctx context.Context, id uuid.UUID) (*identity.Role, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.roles.Get(ctx, id)
}

func (r *repository) GetRoleByCode(ctx context.Context, code string) (*identity.Role, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.roles.GetByCode(ctx, code)
}

func (r *repository) GetRoleWithPermissions(ctx context.Context, id uuid.UUID) (*identity.Role, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.roles.GetWithPermissions(ctx, id)
}

func (r *repository) ListRoles(ctx context.Context, status *int) ([]*identity.Role, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.roles.List(ctx, status)
}

func (r *repository) UpdateRole(ctx context.Context, role *identity.Role) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.roles.Update(ctx, role)
}

func (r *repository) DeleteRole(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.roles.Delete(ctx, id)
}

func (r *repository) SetRolePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.roles.SetPermissions(ctx, roleID, permissionIDs)
}

func (r *repository) SetRoleMenus(ctx context.Context, roleID uuid.UUID, menuIDs []uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.roles.SetMenus(ctx, roleID, menuIDs)
}

func (r *repository) GetDefaultRoles(ctx context.Context) ([]*identity.Role, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.roles.GetDefaultRoles(ctx)
}

// UserIdentity methods
func (r *repository) CreateUserIdentity(ctx context.Context, i *identity.UserIdentity) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.userIdentities.Create(ctx, i)
}

func (r *repository) GetUserIdentity(ctx context.Context, id uuid.UUID) (*identity.UserIdentity, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.userIdentities.Get(ctx, id)
}

func (r *repository) GetUserIdentityByIdentifier(ctx context.Context, identityType identity.IdentityType, identifier string) (*identity.UserIdentity, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.userIdentities.GetByIdentifier(ctx, identityType, identifier)
}

func (r *repository) ListUserIdentities(ctx context.Context, userID uuid.UUID) ([]*identity.UserIdentity, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.userIdentities.ListByUserID(ctx, userID)
}

func (r *repository) DeleteUserIdentity(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.userIdentities.Delete(ctx, id)
}

// Permission methods
func (r *repository) CreatePermission(ctx context.Context, p *identity.Permission) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.permissions.Create(ctx, p)
}

func (r *repository) GetPermission(ctx context.Context, id uuid.UUID) (*identity.Permission, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.permissions.Get(ctx, id)
}

func (r *repository) GetPermissionByCode(ctx context.Context, code string) (*identity.Permission, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.permissions.GetByCode(ctx, code)
}

func (r *repository) ListPermissions(ctx context.Context) ([]*identity.Permission, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.permissions.List(ctx)
}

func (r *repository) UpdatePermission(ctx context.Context, p *identity.Permission) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.permissions.Update(ctx, p)
}

func (r *repository) DeletePermission(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.permissions.Delete(ctx, id)
}

func (r *repository) GetPermissionsByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]*identity.Permission, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.permissions.GetByRoleIDs(ctx, roleIDs)
}

// Menu methods
func (r *repository) CreateMenu(ctx context.Context, m *identity.Menu) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.menus.Create(ctx, m)
}

func (r *repository) GetMenu(ctx context.Context, id uuid.UUID) (*identity.Menu, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.menus.Get(ctx, id)
}

func (r *repository) GetMenuByCode(ctx context.Context, code string) (*identity.Menu, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.menus.GetByCode(ctx, code)
}

func (r *repository) ListMenus(ctx context.Context, status *int) ([]*identity.Menu, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.menus.List(ctx, status)
}

func (r *repository) UpdateMenu(ctx context.Context, m *identity.Menu) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.menus.Update(ctx, m)
}

func (r *repository) DeleteMenu(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.menus.Delete(ctx, id)
}

func (r *repository) GetMenusByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]*identity.Menu, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.menus.GetByRoleIDs(ctx, roleIDs)
}

// MediaAsset methods
func (r *repository) CreateMediaAsset(ctx context.Context, a *media.Asset) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.assets.Create(ctx, a)
}

func (r *repository) GetMediaAsset(ctx context.Context, id uuid.UUID) (*media.Asset, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.assets.Get(ctx, id)
}

func (r *repository) ListMediaAssets(ctx context.Context, filter media.AssetFilter) ([]*media.Asset, int64, error) {
	if err := r.checkDB(); err != nil {
		return nil, 0, err
	}
	return r.assets.List(ctx, filter)
}

func (r *repository) UpdateMediaAsset(ctx context.Context, a *media.Asset) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.assets.Update(ctx, a)
}

func (r *repository) DeleteMediaAsset(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.assets.Delete(ctx, id)
}

func (r *repository) ListMediaAssetsBySource(ctx context.Context, sourceID uuid.UUID) ([]*media.Asset, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.assets.ListBySource(ctx, sourceID)
}

func (r *repository) ListMediaAssetsByParent(ctx context.Context, parentID uuid.UUID) ([]*media.Asset, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.assets.ListByParent(ctx, parentID)
}

func (r *repository) GetAllAssetTags(ctx context.Context) ([]string, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.assets.GetAllTags(ctx)
}

// MediaSource methods
func (r *repository) CreateMediaSource(ctx context.Context, s *media.Source) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.sources.Create(ctx, s)
}

func (r *repository) GetMediaSource(ctx context.Context, id uuid.UUID) (*media.Source, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.sources.Get(ctx, id)
}

func (r *repository) GetMediaSourceByPathName(ctx context.Context, pathName string) (*media.Source, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.sources.GetByPathName(ctx, pathName)
}

func (r *repository) ListMediaSources(ctx context.Context, filter media.SourceFilter) ([]*media.Source, int64, error) {
	if err := r.checkDB(); err != nil {
		return nil, 0, err
	}
	return r.sources.List(ctx, filter)
}

func (r *repository) UpdateMediaSource(ctx context.Context, s *media.Source) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.sources.Update(ctx, s)
}

func (r *repository) DeleteMediaSource(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.sources.Delete(ctx, id)
}

// Operator methods
func (r *repository) CreateOperator(ctx context.Context, o *operator.Operator) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.operators.Create(ctx, o)
}

func (r *repository) GetOperator(ctx context.Context, id uuid.UUID) (*operator.Operator, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.operators.Get(ctx, id)
}

func (r *repository) GetOperatorByCode(ctx context.Context, code string) (*operator.Operator, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.operators.GetByCode(ctx, code)
}

func (r *repository) ListOperators(ctx context.Context, filter operator.Filter) ([]*operator.Operator, int64, error) {
	if err := r.checkDB(); err != nil {
		return nil, 0, err
	}
	return r.operators.List(ctx, filter)
}

func (r *repository) UpdateOperator(ctx context.Context, o *operator.Operator) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.operators.Update(ctx, o)
}

func (r *repository) DeleteOperator(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.operators.Delete(ctx, id)
}

func (r *repository) ListEnabledOperators(ctx context.Context) ([]*operator.Operator, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.operators.ListPublished(ctx)
}

func (r *repository) ListOperatorsByCategory(ctx context.Context, category operator.Category) ([]*operator.Operator, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.operators.ListByCategory(ctx, category)
}

// Workflow methods
func (r *repository) CreateWorkflow(ctx context.Context, w *workflow.Workflow) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.workflows.Create(ctx, w)
}

func (r *repository) GetWorkflow(ctx context.Context, id uuid.UUID) (*workflow.Workflow, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.workflows.Get(ctx, id)
}

func (r *repository) GetWorkflowByCode(ctx context.Context, code string) (*workflow.Workflow, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.workflows.GetByCode(ctx, code)
}

func (r *repository) GetWorkflowWithNodes(ctx context.Context, id uuid.UUID) (*workflow.Workflow, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.workflows.GetWithNodes(ctx, id)
}

func (r *repository) ListWorkflows(ctx context.Context, filter workflow.Filter) ([]*workflow.Workflow, int64, error) {
	if err := r.checkDB(); err != nil {
		return nil, 0, err
	}
	return r.workflows.List(ctx, filter)
}

func (r *repository) UpdateWorkflow(ctx context.Context, w *workflow.Workflow) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.workflows.Update(ctx, w)
}

func (r *repository) DeleteWorkflow(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.workflows.Delete(ctx, id)
}

func (r *repository) ListEnabledWorkflows(ctx context.Context) ([]*workflow.Workflow, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.workflows.ListEnabled(ctx)
}

// WorkflowNode methods
func (r *repository) CreateWorkflowNode(ctx context.Context, n *workflow.Node) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.workflows.CreateNode(ctx, n)
}

func (r *repository) ListWorkflowNodes(ctx context.Context, workflowID uuid.UUID) ([]*workflow.Node, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.workflows.ListNodes(ctx, workflowID)
}

func (r *repository) DeleteWorkflowNodes(ctx context.Context, workflowID uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.workflows.DeleteNodes(ctx, workflowID)
}

// WorkflowEdge methods
func (r *repository) CreateWorkflowEdge(ctx context.Context, e *workflow.Edge) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.workflows.CreateEdge(ctx, e)
}

func (r *repository) ListWorkflowEdges(ctx context.Context, workflowID uuid.UUID) ([]*workflow.Edge, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.workflows.ListEdges(ctx, workflowID)
}

func (r *repository) DeleteWorkflowEdges(ctx context.Context, workflowID uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.workflows.DeleteEdges(ctx, workflowID)
}

// Task methods
func (r *repository) CreateTask(ctx context.Context, t *workflow.Task) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.tasks.Create(ctx, t)
}

func (r *repository) GetTask(ctx context.Context, id uuid.UUID) (*workflow.Task, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.tasks.Get(ctx, id)
}

func (r *repository) GetTaskWithRelations(ctx context.Context, id uuid.UUID) (*workflow.Task, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.tasks.GetWithRelations(ctx, id)
}

func (r *repository) ListTasks(ctx context.Context, filter workflow.TaskFilter) ([]*workflow.Task, int64, error) {
	if err := r.checkDB(); err != nil {
		return nil, 0, err
	}
	return r.tasks.List(ctx, filter)
}

func (r *repository) UpdateTask(ctx context.Context, t *workflow.Task) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.tasks.Update(ctx, t)
}

func (r *repository) DeleteTask(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.tasks.Delete(ctx, id)
}

func (r *repository) GetTaskStats(ctx context.Context, workflowID *uuid.UUID) (*workflow.TaskStats, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.tasks.GetStats(ctx, workflowID)
}

func (r *repository) ListRunningTasks(ctx context.Context) ([]*workflow.Task, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.tasks.ListRunning(ctx)
}

// Artifact methods
func (r *repository) CreateArtifact(ctx context.Context, a *workflow.Artifact) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.artifacts.Create(ctx, a)
}

func (r *repository) GetArtifact(ctx context.Context, id uuid.UUID) (*workflow.Artifact, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.artifacts.Get(ctx, id)
}

func (r *repository) ListArtifacts(ctx context.Context, filter workflow.ArtifactFilter) ([]*workflow.Artifact, int64, error) {
	if err := r.checkDB(); err != nil {
		return nil, 0, err
	}
	return r.artifacts.List(ctx, filter)
}

func (r *repository) DeleteArtifact(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.artifacts.Delete(ctx, id)
}

func (r *repository) ListArtifactsByTask(ctx context.Context, taskID uuid.UUID) ([]*workflow.Artifact, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.artifacts.ListByTask(ctx, taskID)
}

func (r *repository) ListArtifactsByType(ctx context.Context, taskID uuid.UUID, artifactType workflow.ArtifactType) ([]*workflow.Artifact, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.artifacts.ListByType(ctx, taskID, artifactType)
}

// File methods
func (r *repository) CreateFile(ctx context.Context, f *storage.File) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.files.Create(ctx, f)
}

func (r *repository) GetFile(ctx context.Context, id uuid.UUID) (*storage.File, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.files.Get(ctx, id)
}

func (r *repository) GetFileByPath(ctx context.Context, path string) (*storage.File, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.files.GetByPath(ctx, path)
}

func (r *repository) ListFiles(ctx context.Context, filter storage.FileFilter) ([]*storage.File, int64, error) {
	if err := r.checkDB(); err != nil {
		return nil, 0, err
	}
	return r.files.List(ctx, filter)
}

func (r *repository) UpdateFile(ctx context.Context, f *storage.File) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.files.Update(ctx, f)
}

func (r *repository) DeleteFile(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.files.Delete(ctx, id)
}

func (r *repository) GetFileByHash(ctx context.Context, hash string) (*storage.File, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.files.GetByHash(ctx, hash)
}

// AIModel methods
func (r *repository) CreateAIModel(ctx context.Context, d *ai_model.AIModel) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.aiModels.Create(ctx, d)
}

func (r *repository) GetAIModel(ctx context.Context, id uuid.UUID) (*ai_model.AIModel, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.aiModels.Get(ctx, id)
}

func (r *repository) UpdateAIModel(ctx context.Context, d *ai_model.AIModel) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.aiModels.Update(ctx, d)
}

func (r *repository) DeleteAIModel(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.aiModels.Delete(ctx, id)
}

func (r *repository) ListAIModels(ctx context.Context, filter ai_model.Filter) ([]*ai_model.AIModel, int64, error) {
	if err := r.checkDB(); err != nil {
		return nil, 0, err
	}
	return r.aiModels.List(ctx, filter)
}

// SystemConfig methods
func (r *repository) GetSystemConfig(ctx context.Context, key string) (*system.SystemConfig, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.systemConfigs.Get(ctx, key)
}

func (r *repository) ListSystemConfigs(ctx context.Context) ([]*system.SystemConfig, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.systemConfigs.List(ctx)
}

func (r *repository) SaveSystemConfig(ctx context.Context, config *system.SystemConfig) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.systemConfigs.Save(ctx, config)
}

func (r *repository) DeleteSystemConfig(ctx context.Context, key string) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.systemConfigs.Delete(ctx, key)
}

// UserAsset methods
func (r *repository) GetUserBalance(ctx context.Context, userID uuid.UUID) (*domain.UserBalance, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.userAssets.GetUserBalance(ctx, userID)
}

func (r *repository) UpdateUserBalance(ctx context.Context, ub *domain.UserBalance) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.userAssets.UpdateUserBalance(ctx, ub)
}

func (r *repository) CreateTransactionRecord(ctx context.Context, tr *domain.TransactionRecord) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.userAssets.CreateTransactionRecord(ctx, tr)
}

func (r *repository) GetTransactionRecord(ctx context.Context, id string) (*domain.TransactionRecord, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.userAssets.GetTransactionRecord(ctx, id)
}

func (r *repository) UpdateTransactionRecord(ctx context.Context, tr *domain.TransactionRecord) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.userAssets.UpdateTransactionRecord(ctx, tr)
}

func (r *repository) ListTransactionRecords(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.TransactionRecord, int64, error) {
	if err := r.checkDB(); err != nil {
		return nil, 0, err
	}
	return r.userAssets.ListTransactionRecords(ctx, userID, limit, offset)
}

func (r *repository) CreatePointRecord(ctx context.Context, pr *domain.PointRecord) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.userAssets.CreatePointRecord(ctx, pr)
}

func (r *repository) ListPointRecords(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.PointRecord, int64, error) {
	if err := r.checkDB(); err != nil {
		return nil, 0, err
	}
	return r.userAssets.ListPointRecords(ctx, userID, limit, offset)
}

func (r *repository) GetUserSubscription(ctx context.Context, userID uuid.UUID) (*domain.UserSubscription, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.userAssets.GetUserSubscription(ctx, userID)
}

func (r *repository) UpdateUserSubscription(ctx context.Context, us *domain.UserSubscription) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.userAssets.UpdateUserSubscription(ctx, us)
}

func (r *repository) GetUsageStats(ctx context.Context, userID uuid.UUID, date time.Time) (*domain.UsageStats, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.userAssets.GetUsageStats(ctx, userID, date)
}

func (r *repository) UpdateUsageStats(ctx context.Context, us *domain.UsageStats) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.userAssets.UpdateUsageStats(ctx, us)
}

func (r *repository) ListUsageStats(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]*domain.UsageStats, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	return r.userAssets.ListUsageStats(ctx, userID, start, end)
}
