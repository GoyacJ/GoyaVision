package repo

import (
	"context"

	"goyavision/internal/domain/identity"
	"goyavision/internal/infra/persistence/mapper"
	"goyavision/internal/infra/persistence/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, u *identity.User) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	m := mapper.UserToModel(u)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *UserRepo) Get(ctx context.Context, id uuid.UUID) (*identity.User, error) {
	var m model.UserModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.UserToDomain(&m), nil
}

func (r *UserRepo) GetByUsername(ctx context.Context, username string) (*identity.User, error) {
	var m model.UserModel
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.UserToDomain(&m), nil
}

func (r *UserRepo) GetWithRoles(ctx context.Context, id uuid.UUID) (*identity.User, error) {
	var m model.UserModel
	if err := r.db.WithContext(ctx).Preload("Roles").Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.UserToDomain(&m), nil
}

func (r *UserRepo) List(ctx context.Context, status *int, limit, offset int) ([]*identity.User, int64, error) {
	q := r.db.WithContext(ctx).Model(&model.UserModel{})
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var models []*model.UserModel
	if err := q.Preload("Roles").Limit(limit).Offset(offset).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, 0, err
	}
	result := make([]*identity.User, len(models))
	for i, m := range models {
		result[i] = mapper.UserToDomain(m)
	}
	return result, total, nil
}

func (r *UserRepo) Update(ctx context.Context, u *identity.User) error {
	m := mapper.UserToModel(u)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *UserRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Exec("DELETE FROM user_roles WHERE user_model_id = ?", id).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.UserModel{}).Error
}

func (r *UserRepo) SetUserRoles(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error {
	var user model.UserModel
	if err := r.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return err
	}
	var roles []model.RoleModel
	if len(roleIDs) > 0 {
		if err := r.db.WithContext(ctx).Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
			return err
		}
	}
	return r.db.WithContext(ctx).Model(&user).Association("Roles").Replace(roles)
}

type RoleRepo struct {
	db *gorm.DB
}

func NewRoleRepo(db *gorm.DB) *RoleRepo {
	return &RoleRepo{db: db}
}

func (r *RoleRepo) Create(ctx context.Context, role *identity.Role) error {
	if role.ID == uuid.Nil {
		role.ID = uuid.New()
	}
	m := mapper.RoleToModel(role)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *RoleRepo) Get(ctx context.Context, id uuid.UUID) (*identity.Role, error) {
	var m model.RoleModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.RoleToDomain(&m), nil
}

func (r *RoleRepo) GetByCode(ctx context.Context, code string) (*identity.Role, error) {
	var m model.RoleModel
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.RoleToDomain(&m), nil
}

func (r *RoleRepo) GetWithPermissions(ctx context.Context, id uuid.UUID) (*identity.Role, error) {
	var m model.RoleModel
	if err := r.db.WithContext(ctx).Preload("Permissions").Preload("Menus").Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.RoleToDomain(&m), nil
}

func (r *RoleRepo) List(ctx context.Context, status *int) ([]*identity.Role, error) {
	q := r.db.WithContext(ctx)
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	var models []*model.RoleModel
	if err := q.Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*identity.Role, len(models))
	for i, m := range models {
		result[i] = mapper.RoleToDomain(m)
	}
	return result, nil
}

func (r *RoleRepo) Update(ctx context.Context, role *identity.Role) error {
	m := mapper.RoleToModel(role)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *RoleRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Exec("DELETE FROM role_permissions WHERE role_model_id = ?", id).Error; err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Exec("DELETE FROM role_menus WHERE role_model_id = ?", id).Error; err != nil {
		return err
	}
	if err := r.db.WithContext(ctx).Exec("DELETE FROM user_roles WHERE role_model_id = ?", id).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.RoleModel{}).Error
}

func (r *RoleRepo) SetPermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	var role model.RoleModel
	if err := r.db.WithContext(ctx).First(&role, roleID).Error; err != nil {
		return err
	}
	var permissions []model.PermissionModel
	if len(permissionIDs) > 0 {
		if err := r.db.WithContext(ctx).Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
			return err
		}
	}
	return r.db.WithContext(ctx).Model(&role).Association("Permissions").Replace(permissions)
}

func (r *RoleRepo) SetMenus(ctx context.Context, roleID uuid.UUID, menuIDs []uuid.UUID) error {
	var role model.RoleModel
	if err := r.db.WithContext(ctx).First(&role, roleID).Error; err != nil {
		return err
	}
	var menus []model.MenuModel
	if len(menuIDs) > 0 {
		if err := r.db.WithContext(ctx).Where("id IN ?", menuIDs).Find(&menus).Error; err != nil {
			return err
		}
	}
	return r.db.WithContext(ctx).Model(&role).Association("Menus").Replace(menus)
}

func (r *RoleRepo) GetDefaultRoles(ctx context.Context) ([]*identity.Role, error) {
	var models []*model.RoleModel
	if err := r.db.WithContext(ctx).Where("is_default = ?", true).Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*identity.Role, len(models))
	for i, m := range models {
		result[i] = mapper.RoleToDomain(m)
	}
	return result, nil
}

type PermissionRepo struct {
	db *gorm.DB
}

func NewPermissionRepo(db *gorm.DB) *PermissionRepo {
	return &PermissionRepo{db: db}
}

func (r *PermissionRepo) Create(ctx context.Context, p *identity.Permission) error {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	m := mapper.PermissionToModel(p)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PermissionRepo) Get(ctx context.Context, id uuid.UUID) (*identity.Permission, error) {
	var m model.PermissionModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.PermissionToDomain(&m), nil
}

func (r *PermissionRepo) GetByCode(ctx context.Context, code string) (*identity.Permission, error) {
	var m model.PermissionModel
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.PermissionToDomain(&m), nil
}

func (r *PermissionRepo) List(ctx context.Context) ([]*identity.Permission, error) {
	var models []*model.PermissionModel
	if err := r.db.WithContext(ctx).Order("code").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*identity.Permission, len(models))
	for i, m := range models {
		result[i] = mapper.PermissionToDomain(m)
	}
	return result, nil
}

func (r *PermissionRepo) Update(ctx context.Context, p *identity.Permission) error {
	m := mapper.PermissionToModel(p)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *PermissionRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Exec("DELETE FROM role_permissions WHERE permission_model_id = ?", id).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.PermissionModel{}).Error
}

func (r *PermissionRepo) GetByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]*identity.Permission, error) {
	if len(roleIDs) == 0 {
		return []*identity.Permission{}, nil
	}
	var models []*model.PermissionModel
	err := r.db.WithContext(ctx).
		Distinct().
		Joins("JOIN role_permissions ON role_permissions.permission_model_id = permissions.id").
		Where("role_permissions.role_model_id IN ?", roleIDs).
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	result := make([]*identity.Permission, len(models))
	for i, m := range models {
		result[i] = mapper.PermissionToDomain(m)
	}
	return result, nil
}

type MenuRepo struct {
	db *gorm.DB
}

func NewMenuRepo(db *gorm.DB) *MenuRepo {
	return &MenuRepo{db: db}
}

func (r *MenuRepo) Create(ctx context.Context, menu *identity.Menu) error {
	if menu.ID == uuid.Nil {
		menu.ID = uuid.New()
	}
	m := mapper.MenuToModel(menu)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *MenuRepo) Get(ctx context.Context, id uuid.UUID) (*identity.Menu, error) {
	var m model.MenuModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.MenuToDomain(&m), nil
}

func (r *MenuRepo) GetByCode(ctx context.Context, code string) (*identity.Menu, error) {
	var m model.MenuModel
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.MenuToDomain(&m), nil
}

func (r *MenuRepo) List(ctx context.Context, status *int) ([]*identity.Menu, error) {
	q := r.db.WithContext(ctx)
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	var models []*model.MenuModel
	if err := q.Order("sort, created_at").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*identity.Menu, len(models))
	for i, m := range models {
		result[i] = mapper.MenuToDomain(m)
	}
	return result, nil
}

func (r *MenuRepo) Update(ctx context.Context, menu *identity.Menu) error {
	m := mapper.MenuToModel(menu)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *MenuRepo) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.db.WithContext(ctx).Exec("DELETE FROM role_menus WHERE menu_model_id = ?", id).Error; err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&model.MenuModel{}).Error
}

func (r *MenuRepo) GetByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]*identity.Menu, error) {
	if len(roleIDs) == 0 {
		return []*identity.Menu{}, nil
	}
	var models []*model.MenuModel
	err := r.db.WithContext(ctx).
		Distinct().
		Joins("JOIN role_menus ON role_menus.menu_model_id = menus.id").
		Where("role_menus.role_model_id IN ?", roleIDs).
		Where("menus.status = ?", 1). // 只返回启用的菜单
		Order("menus.sort, menus.created_at").
		Find(&models).Error
	if err != nil {
		return nil, err
	}
	result := make([]*identity.Menu, len(models))
	for i, m := range models {
		result[i] = mapper.MenuToDomain(m)
	}
	return result, nil
}
