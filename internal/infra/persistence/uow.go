package persistence

import (
	"context"

	appport "goyavision/internal/app/port"
	"goyavision/internal/infra/persistence/repo"

	"gorm.io/gorm"
)

type gormUoW struct {
	db *gorm.DB
}

func NewUnitOfWork(db *gorm.DB) appport.UnitOfWork {
	return &gormUoW{db: db}
}

func (u *gormUoW) Do(ctx context.Context, fn func(ctx context.Context, repos *appport.Repositories) error) error {
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		repos := newRepositories(tx)
		return fn(ctx, repos)
	})
}

func newRepositories(db *gorm.DB) *appport.Repositories {
	return &appport.Repositories{
		Sources:     repo.NewMediaSourceRepo(db),
		Assets:      repo.NewMediaAssetRepo(db),
		Operators:   repo.NewOperatorRepo(db),
		OperatorVersions:     repo.NewOperatorVersionRepo(db),
		OperatorTemplates:    repo.NewOperatorTemplateRepo(db),
		OperatorDependencies: repo.NewOperatorDependencyRepo(db),
		Workflows:   repo.NewWorkflowRepo(db),
		Tasks:       repo.NewTaskRepo(db),
		Artifacts:   repo.NewArtifactRepo(db),
		Users:       repo.NewUserRepo(db),
		Roles:       repo.NewRoleRepo(db),
		Permissions: repo.NewPermissionRepo(db),
		Menus:       repo.NewMenuRepo(db),
		Files:          repo.NewFileRepo(db),
		AIModels:       repo.NewAIModelRepo(db),
		UserIdentities: repo.NewUserIdentityRepo(db),
	}
}

func NewRepositories(db *gorm.DB) *appport.Repositories {
	return newRepositories(db)
}
