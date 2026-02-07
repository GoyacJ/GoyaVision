package dto

import (
	"github.com/google/uuid"
	"goyavision/internal/domain/identity"
	"goyavision/internal/domain/media"
	"goyavision/internal/domain/operator"
	"goyavision/internal/domain/workflow"
)

// Media Source Commands

type CreateSourceCommand struct {
	Name     string
	Type     media.SourceType
	URL      string
	Protocol string
	Enabled  bool
}

type UpdateSourceCommand struct {
	ID       uuid.UUID
	Name     *string
	URL      *string
	Protocol *string
	Enabled  *bool
}

type DeleteSourceCommand struct {
	ID uuid.UUID
}

// Media Asset Commands

type CreateAssetCommand struct {
	Type       media.AssetType
	SourceType media.AssetSourceType
	SourceID   *uuid.UUID
	ParentID   *uuid.UUID
	Name       string
	Path       string
	Duration   *float64
	Size       int64
	Format     string
	Metadata   map[string]interface{}
	Status     media.AssetStatus
	Tags       []string
}

type UpdateAssetCommand struct {
	ID       uuid.UUID
	Name     *string
	Status   *media.AssetStatus
	Metadata map[string]interface{}
	Tags     *[]string
}

type DeleteAssetCommand struct {
	ID uuid.UUID
}

// Auth Commands

type LoginCommand struct {
	Username string
	Password string
}

type RegisterCommand struct {
	Username string
	Password string
	Nickname string
	Email    string
	Phone    string
}

type RefreshTokenCommand struct {
	RefreshToken string
}

type ChangePasswordCommand struct {
	UserID      uuid.UUID
	OldPassword string
	NewPassword string
}

type LoginOAuthCommand struct {
	Provider string
	Code     string
	State    string
}

type BindIdentityCommand struct {
	UserID     uuid.UUID
	Provider   string
	Identifier string
	Credential string
	Meta       map[string]interface{}
}

// Operator Commands

type CreateOperatorCommand struct {
	Code        string
	Name        string
	Description string
	Category    operator.Category
	Type        operator.Type
	Origin      operator.Origin
	ExecMode    operator.ExecMode
	ExecConfig  *operator.ExecConfig
	Status      operator.Status
	Tags        []string
}

type UpdateOperatorCommand struct {
	ID          uuid.UUID
	Name        *string
	Description *string
	Category    *operator.Category
	Tags        []string
}

type DeleteOperatorCommand struct {
	ID uuid.UUID
}

type PublishOperatorCommand struct {
	ID uuid.UUID
}

type DeprecateOperatorCommand struct {
	ID uuid.UUID
}

type TestOperatorCommand struct {
	ID      uuid.UUID
	AssetID *uuid.UUID
	Params  map[string]interface{}
}

type CreateOperatorVersionCommand struct {
	OperatorID  uuid.UUID
	Version     string
	ExecMode    operator.ExecMode
	ExecConfig  *operator.ExecConfig
	InputSchema map[string]interface{}
	OutputSpec  map[string]interface{}
	Config      map[string]interface{}
	Changelog   string
	Status      operator.VersionStatus
}

type ActivateVersionCommand struct {
	OperatorID uuid.UUID
	VersionID  uuid.UUID
}

type RollbackVersionCommand struct {
	OperatorID uuid.UUID
	VersionID  uuid.UUID
}

type ArchiveVersionCommand struct {
	OperatorID uuid.UUID
	VersionID  uuid.UUID
}

type TestOperatorResult struct {
	Success     bool                   `json:"success"`
	Message     string                 `json:"message"`
	Diagnostics map[string]interface{} `json:"diagnostics,omitempty"`
}

type InstallMCPOperatorCommand struct {
	ServerID     string
	ToolName     string
	OperatorCode string
	OperatorName string
	Category     *operator.Category
	Type         *operator.Type
	TimeoutSec   int
	Tags         []string
}

type SyncMCPTemplatesCommand struct {
	ServerID string
}

type SyncMCPTemplatesResult struct {
	ServerID string `json:"server_id"`
	Total    int    `json:"total"`
	Created  int    `json:"created"`
	Updated  int    `json:"updated"`
}

type InstallTemplateCommand struct {
	TemplateID   uuid.UUID
	OperatorCode string
	OperatorName string
	Tags         []string
}

type DependencyItemInput struct {
	DependsOnID uuid.UUID
	MinVersion  string
	IsOptional  bool
}

type SetOperatorDependenciesCommand struct {
	OperatorID   uuid.UUID
	Dependencies []DependencyItemInput
}

// Workflow Commands

type WorkflowNodeInput struct {
	NodeKey    string
	NodeType   string
	OperatorID *uuid.UUID
	Config     map[string]interface{}
	Position   map[string]interface{}
}

type WorkflowEdgeInput struct {
	SourceKey string
	TargetKey string
	Condition map[string]interface{}
}

type CreateWorkflowCommand struct {
	Code        string
	Name        string
	Description string
	Version     string
	TriggerType workflow.TriggerType
	TriggerConf map[string]interface{}
	Status      workflow.Status
	Tags        []string
	Nodes       []WorkflowNodeInput
	Edges       []WorkflowEdgeInput
}

type UpdateWorkflowCommand struct {
	ID          uuid.UUID
	Name        *string
	Description *string
	TriggerConf map[string]interface{}
	Status      *workflow.Status
	Tags        []string
	Nodes       []WorkflowNodeInput
	Edges       []WorkflowEdgeInput
}

type DeleteWorkflowCommand struct {
	ID uuid.UUID
}

type EnableWorkflowCommand struct {
	ID      uuid.UUID
	Enabled bool
}

// Task Commands

type CreateTaskCommand struct {
	WorkflowID  uuid.UUID
	AssetID     *uuid.UUID
	InputParams map[string]interface{}
}

type UpdateTaskCommand struct {
	ID          uuid.UUID
	Status      *workflow.TaskStatus
	Progress    *int
	CurrentNode *string
	Error       *string
}

type DeleteTaskCommand struct {
	ID uuid.UUID
}

type StartTaskCommand struct {
	ID uuid.UUID
}

type CompleteTaskCommand struct {
	ID uuid.UUID
}

type FailTaskCommand struct {
	ID       uuid.UUID
	ErrorMsg string
}

type CancelTaskCommand struct {
	ID uuid.UUID
}

// User Management Commands

type CreateUserCommand struct {
	Username string
	Password string
	Nickname string
	Email    string
	Phone    string
	Avatar   string
	Status   int
	RoleIDs  []uuid.UUID
}

type UpdateUserCommand struct {
	ID       uuid.UUID
	Nickname *string
	Email    *string
	Phone    *string
	Avatar   *string
	Status   *int
}

type DeleteUserCommand struct {
	ID uuid.UUID
}

type AssignUserRolesCommand struct {
	UserID  uuid.UUID
	RoleIDs []uuid.UUID
}

type CreateRoleCommand struct {
	Name        string
	Code        string
	Description string
	Status      int
}

type UpdateRoleCommand struct {
	ID          uuid.UUID
	Name        *string
	Description *string
	Status      *int
}

type DeleteRoleCommand struct {
	ID uuid.UUID
}

type AssignRolePermissionsCommand struct {
	RoleID        uuid.UUID
	PermissionIDs []uuid.UUID
}

type AssignRoleMenusCommand struct {
	RoleID  uuid.UUID
	MenuIDs []uuid.UUID
}

type CreatePermissionCommand struct {
	Name        string
	Code        string
	Description string
	Category    string
}

type UpdatePermissionCommand struct {
	ID          uuid.UUID
	Name        *string
	Description *string
	Category    *string
}

type DeletePermissionCommand struct {
	ID uuid.UUID
}

type CreateMenuCommand struct {
	ParentID  *uuid.UUID
	Name      string
	Code      string
	Path      string
	Icon      string
	Sort      int
	Type      identity.MenuType
	Component string
	Status    int
}

type UpdateMenuCommand struct {
	ID        uuid.UUID
	ParentID  *uuid.UUID
	Name      *string
	Path      *string
	Icon      *string
	Sort      *int
	Type      *identity.MenuType
	Component *string
	Status    *int
}

type DeleteMenuCommand struct {
	ID uuid.UUID
}

// File Commands

type CreateFileCommand struct {
	Name        string
	Path        string
	Size        int64
	ContentType string
	Category    string
}

type DeleteFileCommand struct {
	ID uuid.UUID
}
