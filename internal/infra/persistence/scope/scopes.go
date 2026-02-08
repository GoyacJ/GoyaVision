package scope

import (
	"context"

	"goyavision/internal/api/middleware"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// ScopeTenant adds tenant_id filter based on context
func ScopeTenant(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		tenantID, ok := middleware.GetTenantID(ctx)
		if ok && tenantID != uuid.Nil {
			// 允许本租户数据 OR 全局公开数据 (visibility = 2)
			// 使用这种形式确保 GORM 正确处理括号
			return db.Where(db.Where("tenant_id = ?", tenantID).Or("visibility = ?", 2))
		}
		return db
	}
}

// ScopeTenantOnly adds tenant_id filter based on context, without visibility check
func ScopeTenantOnly(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		tenantID, ok := middleware.GetTenantID(ctx)
		if ok && tenantID != uuid.Nil {
			return db.Where("tenant_id = ?", tenantID)
		}
		return db
	}
}

// ScopeVisibility adds visibility filter based on context user
func ScopeVisibility(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		userID := ctx.Value(middleware.ContextKeyUserID)

		if userID == nil {
			return db.Where("visibility = ?", 2)
		}

		uid := userID.(uuid.UUID)
		
		var roleIDs []string
		if val := ctx.Value(middleware.ContextKeyRoleIDs); val != nil {
			if ids, ok := val.([]uuid.UUID); ok {
				for _, id := range ids {
					roleIDs = append(roleIDs, id.String())
				}
			}
		}

		// (Owner OR Public OR RoleMatch)
		q := db.Where("owner_id = ?", uid).Or("visibility = ?", 2)
		
		if len(roleIDs) > 0 {
			// 修复子查询别名问题：EXISTS (SELECT 1 FROM jsonb_array_elements_text(visible_role_ids) as r WHERE r IN ?)
			q = q.Or("visibility = 1 AND EXISTS (SELECT 1 FROM jsonb_array_elements_text(visible_role_ids) as r WHERE r IN ?)", roleIDs)
		}

		return db.Where(q)
	}
}

// GetContextInfo retrieves tenantID and userID from context
func GetContextInfo(ctx context.Context) (tenantID uuid.UUID, userID uuid.UUID) {
	tenantID, _ = middleware.GetTenantID(ctx)
	if uid := ctx.Value(middleware.ContextKeyUserID); uid != nil {
		userID = uid.(uuid.UUID)
	}
	return
}
