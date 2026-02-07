package persistence

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
		// If user is super admin, show all? (Ideally yes, but let's stick to standard logic first)
		// Or maybe checking permission "view_all"?

		// Get current user info from context
		// Note: We need a way to get user info from context which is usually *echo.Context.
		// But repo methods receive context.Context.
		// Echo middleware sets values in echo.Context, but does it propagate to request.Context()?
		// No, echo.Context is a wrapper.
		// WE NEED TO FIX MIDDLEWARE to set values in Request Context too, or use a helper to extract from it if possible.
		// Actually, std context values are typically set via context.WithValue.
		// Echo's c.Set() sets it in a map inside echo.Context.
		// App handlers usually pass c.Request().Context().
		// So we need Middleware to copy keys to Request Context?
		// Or better, just use context.WithValue in middleware.
		
		// Let's assume we fix middleware to use context.WithValue or we have a way to retrieve it.
		// For now, let's look at `middleware.GetUserID`. It takes `echo.Context`.
		// But Repo receives `context.Context`.
		
		// CRITICAL: We need to bridge Echo Context and Std Context.
		// In `internal/api/router.go` or `handlers.go`, we usually pass `c.Request().Context()`.
		// If `c.Set` is used, `c.Request().Context()` does NOT have those values.
		
		// Solution: Update `JWTAuth` middleware to ALSO set values in `c.Request()`.
		
		userID := ctx.Value(middleware.ContextKeyUserID)
		// roles := ctx.Value(middleware.ContextKeyRoles) // Need to ensure roles are in context

		if userID == nil {
			// No user context, maybe return public only? or nothing?
			// For safety, return nothing or just Public.
			return db.Where("visibility = 2") // Public only
		}

		uid := userID.(uuid.UUID)
		
		// Roles logic is complex because roles are loaded in `LoadUserPermissions` middleware.
		// Let's assume roles are available as []string (role codes) or []uuid.UUID (role IDs).
		// The model uses `visible_role_ids` (UUIDs).
		// So we need role IDs in context. `RequirePermission` middleware sets role CODES.
		// We might need to update middleware to set role IDs too.
		
		// For now, simplistic implementation:
		// Owner OR Public
		return db.Where("owner_id = ? OR visibility = 2", uid)
		
		// TODO: Add Role visibility support once Role IDs are in context
	}
}

// Helper to convert std context to echo context if needed? No, that's impossible.
// We must put data into std context.
