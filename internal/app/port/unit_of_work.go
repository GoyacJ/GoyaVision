package port

import (
	"context"

	"goyavision/internal/domain/identity"
	"goyavision/internal/domain/media"
	"goyavision/internal/domain/operator"
	"goyavision/internal/domain/storage"
	"goyavision/internal/domain/workflow"
)

// UnitOfWork 工作单元接口，提供事务边界和仓储访问
//
// 职责：
//  1. 管理数据库事务生命周期（开始、提交、回滚）
//  2. 提供对所有聚合仓储的访问
//  3. 确保一个业务用例内的所有数据变更原子化
//
// 使用示例：
//
//	err := uow.Do(ctx, func(ctx context.Context, repos *Repositories) error {
//	    source, err := repos.Sources.Get(ctx, sourceID)
//	    if err != nil {
//	        return err
//	    }
//	    source.Enable()
//	    return repos.Sources.Update(ctx, source)
//	})
type UnitOfWork interface {
	// Do 在事务内执行 fn
	// 如果 fn 返回 error，事务自动回滚
	// 如果 fn 返回 nil，事务自动提交
	Do(ctx context.Context, fn func(ctx context.Context, repos *Repositories) error) error
}

// Repositories 聚合所有仓储的访问器
type Repositories struct {
	Sources     media.SourceRepository
	Assets      media.AssetRepository
	Operators   operator.Repository
	Workflows   workflow.Repository
	Tasks       workflow.TaskRepository
	Artifacts   workflow.ArtifactRepository
	Users       identity.UserRepository
	Roles       identity.RoleRepository
	Permissions identity.PermissionRepository
	Menus       identity.MenuRepository
	Files       storage.FileRepository
}
