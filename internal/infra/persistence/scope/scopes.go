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
			return db.Where("tenant_id = ?", tenantID)
		}
		// If no tenant in context (e.g. system task), skip filter or handle accordingly.
		// For now, if no tenant, we don't filter (dangerous? maybe. But necessary for background jobs).
		return db
	}
}

// ScopeVisibility adds visibility filter based on context user
func ScopeVisibility(ctx context.Context) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		userID := ctx.Value(middleware.ContextKeyUserID)

		if userID == nil {
			// No user context, return Public only.
			return db.Where("visibility = 2") // Public only
		}

		uid := userID.(uuid.UUID)
		
		// Get Role IDs from context
		var roleIDs []string
		if val := ctx.Value(middleware.ContextKeyRoleIDs); val != nil {
			if ids, ok := val.([]uuid.UUID); ok {
				for _, id := range ids {
					roleIDs = append(roleIDs, id.String())
				}
			}
		}

		// Owner OR Public OR (Role AND Role Match)
		// Note: We cast jsonb to text array for containment check, or use jsonb_exists_any if we format roles as string array
		// But simpler is: EXISTS (SELECT 1 FROM jsonb_array_elements_text(visible_role_ids) WHERE value IN (?))
		
		if len(roleIDs) > 0 {
			return db.Where("owner_id = ? OR visibility = 2 OR (visibility = 1 AND EXISTS (SELECT 1 FROM jsonb_array_elements_text(visible_role_ids) WHERE value IN ?))", uid, roleIDs)
		}

		return db.Where("owner_id = ? OR visibility = 2", uid)
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
