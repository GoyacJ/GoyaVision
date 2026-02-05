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

type RefreshTokenCommand struct {
	RefreshToken string
}

type ChangePasswordCommand struct {
	UserID      uuid.UUID
	OldPassword string
	NewPassword string
}

// Operator Commands

type CreateOperatorCommand struct {
	Code        string
	Name        string
	Description string
	Category    operator.Category
	Type        operator.Type
	Version     string
	Endpoint    string
	Method      string
	InputSchema map[string]interface{}
	OutputSpec  map[string]interface{}
	Config      map[string]interface{}
	Status      operator.Status
	IsBuiltin   bool
	Tags        []string
}

type UpdateOperatorCommand struct {
	ID          uuid.UUID
	Name        *string
	Description *string
	Endpoint    *string
	Method      *string
	InputSchema map[string]interface{}
	OutputSpec  map[string]interface{}
	Config      map[string]interface{}
	Status      *operator.Status
	Tags        []string
}

type DeleteOperatorCommand struct {
	ID uuid.UUID
}

type EnableOperatorCommand struct {
	ID      uuid.UUID
	Enabled bool
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
