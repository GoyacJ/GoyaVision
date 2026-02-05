package dto

import (
	"time"

	"github.com/google/uuid"
	"goyavision/internal/domain/media"
	"goyavision/internal/domain/operator"
	"goyavision/internal/domain/workflow"
)

// Pagination 分页参数
type Pagination struct {
	Limit  int
	Offset int
}

func (p *Pagination) Normalize() {
	if p.Limit <= 0 {
		p.Limit = 20
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
	if p.Offset < 0 {
		p.Offset = 0
	}
}

// Media Source Queries

type GetSourceQuery struct {
	ID uuid.UUID
}

type ListSourcesQuery struct {
	Type       *media.SourceType
	Enabled    *bool
	Pagination Pagination
}

type GetSourceByPathQuery struct {
	PathName string
}

// Media Asset Queries

type GetAssetQuery struct {
	ID uuid.UUID
}

type ListAssetsQuery struct {
	Type       *media.AssetType
	SourceType *media.AssetSourceType
	SourceID   *uuid.UUID
	ParentID   *uuid.UUID
	Status     *media.AssetStatus
	Tags       []string
	From       *time.Time
	To         *time.Time
	Pagination Pagination
}

type ListAssetsBySourceQuery struct {
	SourceID uuid.UUID
}

type ListAssetChildrenQuery struct {
	ParentID uuid.UUID
}

type GetAssetTagsQuery struct {
}

// Auth Queries

type GetProfileQuery struct {
	UserID uuid.UUID
}

// Operator Queries

type GetOperatorQuery struct {
	ID uuid.UUID
}

type GetOperatorByCodeQuery struct {
	Code string
}

type ListOperatorsQuery struct {
	Category   *operator.Category
	Type       *operator.Type
	Status     *operator.Status
	IsBuiltin  *bool
	Tags       []string
	Keyword    string
	Pagination Pagination
}

// Workflow Queries

type GetWorkflowQuery struct {
	ID uuid.UUID
}

type GetWorkflowWithNodesQuery struct {
	ID uuid.UUID
}

type GetWorkflowByCodeQuery struct {
	Code string
}

type ListWorkflowsQuery struct {
	Status      *workflow.Status
	TriggerType *workflow.TriggerType
	Tags        []string
	Keyword     string
	Pagination  Pagination
}

// Task Queries

type GetTaskQuery struct {
	ID uuid.UUID
}

type GetTaskWithRelationsQuery struct {
	ID uuid.UUID
}

type ListTasksQuery struct {
	WorkflowID *uuid.UUID
	AssetID    *uuid.UUID
	Status     *workflow.TaskStatus
	From       *time.Time
	To         *time.Time
	Pagination Pagination
}

type GetTaskStatsQuery struct {
	WorkflowID *uuid.UUID
}

type ListRunningTasksQuery struct {
}

// User Management Queries

type GetUserQuery struct {
	ID uuid.UUID
}

type GetUserByUsernameQuery struct {
	Username string
}

type ListUsersQuery struct {
	Status     *int
	Pagination Pagination
}

type GetRoleQuery struct {
	ID uuid.UUID
}

type GetRoleByCodeQuery struct {
	Code string
}

type ListRolesQuery struct {
	Status *int
}

type GetPermissionQuery struct {
	ID uuid.UUID
}

type GetPermissionByCodeQuery struct {
	Code string
}

type ListPermissionsQuery struct {
}

type GetMenuQuery struct {
	ID uuid.UUID
}

type GetMenuByCodeQuery struct {
	Code string
}

type ListMenusQuery struct {
	Status *int
}

// File Queries

type GetFileQuery struct {
	ID uuid.UUID
}

type ListFilesQuery struct {
	Category   string
	Pagination Pagination
}
