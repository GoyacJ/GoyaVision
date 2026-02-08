package dto

import (
	"time"

	"goyavision/internal/domain/identity"
	"goyavision/internal/domain/media"
	"goyavision/internal/domain/operator"
	"goyavision/internal/domain/storage"
	"goyavision/internal/domain/workflow"

	"github.com/google/uuid"
)

// PagedResult 分页结果
type PagedResult[T any] struct {
	Items  []T
	Total  int64
	Limit  int
	Offset int
}

// Media Source Results

type SourceResult struct {
	Source *media.Source
}

type SourceListResult struct {
	Sources []*media.Source
	Total   int64
}

// Media Asset Results

type AssetResult struct {
	Asset *media.Asset
}

type AssetListResult struct {
	Assets []*media.Asset
	Total  int64
}

type AssetTagsResult struct {
	Tags []string
}

// Auth Results

type LoginResult struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	User         *UserInfo
}

type UserInfo struct {
	ID          uuid.UUID
	Username    string
	Nickname    string
	Email       string
	Phone       string
	Avatar      string
	Roles       []string
	Permissions []string
	Menus       []*identity.Menu
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Operator Results

type OperatorResult struct {
	Operator *operator.Operator
}

type OperatorListResult struct {
	Operators []*operator.Operator
	Total     int64
}

// Workflow Results

type WorkflowResult struct {
	Workflow *workflow.Workflow
}

type WorkflowListResult struct {
	Workflows []*workflow.Workflow
	Total     int64
}

// Task Results

type TaskResult struct {
	Task *workflow.Task
}

type TaskListResult struct {
	Tasks []*workflow.Task
	Total int64
}

type TaskStatsResult struct {
	Stats *workflow.TaskStats
}

// User Management Results

type UserResult struct {
	User *identity.User
}

type UserListResult struct {
	Users []*identity.User
	Total int64
}

type RoleResult struct {
	Role *identity.Role
}

type RoleListResult struct {
	Roles []*identity.Role
}

type PermissionResult struct {
	Permission *identity.Permission
}

type PermissionListResult struct {
	Permissions []*identity.Permission
}

type MenuResult struct {
	Menu *identity.Menu
}

type MenuListResult struct {
	Menus []*identity.Menu
}

// File Results

type FileResult struct {
	File *storage.File
}

type FileListResult struct {
	Files []*storage.File
	Total int64
}
