package persistence

import (
	"context"
	"errors"

	"goyavision/internal/domain"
	"goyavision/internal/port"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrDBNotConfigured = errors.New("database not configured")

var _ port.Repository = (*repository)(nil)

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) checkDB() error {
	if r.db == nil {
		return ErrDBNotConfigured
	}
	return nil
}

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&domain.Stream{},
		&domain.Algorithm{},
		&domain.AlgorithmBinding{},
		&domain.RecordSession{},
		&domain.InferenceResult{},
		&domain.User{},
		&domain.Role{},
		&domain.Permission{},
		&domain.Menu{},
		&domain.MediaAsset{},
		&domain.Operator{},
		&domain.Workflow{},
		&domain.WorkflowNode{},
		&domain.WorkflowEdge{},
	); err != nil {
		return err
	}

	if err := addIndexesAndConstraints(db); err != nil {
		return err
	}

	return nil
}

func addIndexesAndConstraints(db *gorm.DB) error {
	if db.Migrator().HasIndex(&domain.RecordSession{}, "idx_record_sessions_stream_running") {
		return nil
	}

	if err := db.Exec(`
		CREATE UNIQUE INDEX IF NOT EXISTS idx_record_sessions_stream_running 
		ON record_sessions (stream_id) 
		WHERE status = 'running'
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_inference_results_stream_ts 
		ON inference_results (stream_id, ts DESC)
	`).Error; err != nil {
		return err
	}

	if err := db.Exec(`
		CREATE INDEX IF NOT EXISTS idx_inference_results_binding_ts 
		ON inference_results (algorithm_binding_id, ts DESC)
	`).Error; err != nil {
		return err
	}

	return nil
}

func ensureID(id *uuid.UUID) {
	if *id == uuid.Nil {
		*id = uuid.New()
	}
}

func (r *repository) CreateStream(ctx context.Context, s *domain.Stream) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	ensureID(&s.ID)
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *repository) GetStream(ctx context.Context, id uuid.UUID) (*domain.Stream, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var s domain.Stream
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&s).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *repository) ListStreams(ctx context.Context, enabled *bool) ([]*domain.Stream, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var list []*domain.Stream
	q := r.db.WithContext(ctx)
	if enabled != nil {
		q = q.Where("enabled = ?", *enabled)
	}
	if err := q.Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) UpdateStream(ctx context.Context, s *domain.Stream) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *repository) DeleteStream(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Stream{}).Error
}

func (r *repository) CreateAlgorithm(ctx context.Context, a *domain.Algorithm) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	ensureID(&a.ID)
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *repository) GetAlgorithm(ctx context.Context, id uuid.UUID) (*domain.Algorithm, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var a domain.Algorithm
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *repository) ListAlgorithms(ctx context.Context) ([]*domain.Algorithm, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var list []*domain.Algorithm
	if err := r.db.WithContext(ctx).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) UpdateAlgorithm(ctx context.Context, a *domain.Algorithm) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(a).Error
}

func (r *repository) DeleteAlgorithm(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Algorithm{}).Error
}

func (r *repository) CreateAlgorithmBinding(ctx context.Context, b *domain.AlgorithmBinding) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	ensureID(&b.ID)
	return r.db.WithContext(ctx).Create(b).Error
}

func (r *repository) GetAlgorithmBinding(ctx context.Context, id uuid.UUID) (*domain.AlgorithmBinding, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var b domain.AlgorithmBinding
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&b).Error; err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *repository) ListAlgorithmBindingsByStream(ctx context.Context, streamID uuid.UUID) ([]*domain.AlgorithmBinding, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var list []*domain.AlgorithmBinding
	if err := r.db.WithContext(ctx).Where("stream_id = ?", streamID).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) UpdateAlgorithmBinding(ctx context.Context, b *domain.AlgorithmBinding) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(b).Error
}

func (r *repository) DeleteAlgorithmBinding(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.AlgorithmBinding{}).Error
}

func (r *repository) CreateRecordSession(ctx context.Context, rec *domain.RecordSession) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	ensureID(&rec.ID)
	return r.db.WithContext(ctx).Create(rec).Error
}

func (r *repository) GetRecordSession(ctx context.Context, id uuid.UUID) (*domain.RecordSession, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var rec domain.RecordSession
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&rec).Error; err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *repository) GetRunningRecordSessionByStream(ctx context.Context, streamID uuid.UUID) (*domain.RecordSession, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var rec domain.RecordSession
	if err := r.db.WithContext(ctx).Where("stream_id = ? AND status = ?", streamID, domain.RecordStatusRunning).First(&rec).Error; err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *repository) ListRecordSessionsByStream(ctx context.Context, streamID uuid.UUID) ([]*domain.RecordSession, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var list []*domain.RecordSession
	if err := r.db.WithContext(ctx).Where("stream_id = ?", streamID).Order("started_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) UpdateRecordSession(ctx context.Context, rec *domain.RecordSession) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(rec).Error
}

func (r *repository) CreateInferenceResult(ctx context.Context, ir *domain.InferenceResult) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	ensureID(&ir.ID)
	return r.db.WithContext(ctx).Create(ir).Error
}

func (r *repository) ListInferenceResults(ctx context.Context, streamID, bindingID *uuid.UUID, from, to *int64, limit, offset int) ([]*domain.InferenceResult, int64, error) {
	if err := r.checkDB(); err != nil {
		return nil, 0, err
	}
	q := r.db.WithContext(ctx).Model(&domain.InferenceResult{})
	if streamID != nil {
		q = q.Where("stream_id = ?", *streamID)
	}
	if bindingID != nil {
		q = q.Where("algorithm_binding_id = ?", *bindingID)
	}
	if from != nil {
		q = q.Where("EXTRACT(EPOCH FROM ts) >= ?", *from)
	}
	if to != nil {
		q = q.Where("EXTRACT(EPOCH FROM ts) <= ?", *to)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []*domain.InferenceResult
	if err := q.Limit(limit).Offset(offset).Order("ts DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *repository) ListEnabledAlgorithmBindings(ctx context.Context) ([]*domain.AlgorithmBinding, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var list []*domain.AlgorithmBinding
	if err := r.db.WithContext(ctx).Where("enabled = ?", true).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// User methods

func (r *repository) CreateUser(ctx context.Context, u *domain.User) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	ensureID(&u.ID)
	return r.db.WithContext(ctx).Create(u).Error
}

func (r *repository) GetUser(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var u domain.User
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var u domain.User
	if err := r.db.WithContext(ctx).Where("username = ?", username).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) GetUserWithRoles(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var u domain.User
	if err := r.db.WithContext(ctx).Preload("Roles").Where("id = ?", id).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *repository) ListUsers(ctx context.Context, status *int, limit, offset int) ([]*domain.User, int64, error) {
	if err := r.checkDB(); err != nil {
		return nil, 0, err
	}
	q := r.db.WithContext(ctx).Model(&domain.User{})
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var list []*domain.User
	if err := q.Limit(limit).Offset(offset).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, total, nil
}

func (r *repository) UpdateUser(ctx context.Context, u *domain.User) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(u).Error
}

func (r *repository) DeleteUser(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.User{}).Error
}

func (r *repository) SetUserRoles(ctx context.Context, userID uuid.UUID, roleIDs []uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	var user domain.User
	if err := r.db.WithContext(ctx).First(&user, userID).Error; err != nil {
		return err
	}
	var roles []domain.Role
	if len(roleIDs) > 0 {
		if err := r.db.WithContext(ctx).Where("id IN ?", roleIDs).Find(&roles).Error; err != nil {
			return err
		}
	}
	return r.db.WithContext(ctx).Model(&user).Association("Roles").Replace(roles)
}

// Role methods

func (r *repository) CreateRole(ctx context.Context, role *domain.Role) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	ensureID(&role.ID)
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *repository) GetRole(ctx context.Context, id uuid.UUID) (*domain.Role, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var role domain.Role
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *repository) GetRoleByCode(ctx context.Context, code string) (*domain.Role, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var role domain.Role
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *repository) GetRoleWithPermissions(ctx context.Context, id uuid.UUID) (*domain.Role, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var role domain.Role
	if err := r.db.WithContext(ctx).Preload("Permissions").Preload("Menus").Where("id = ?", id).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *repository) ListRoles(ctx context.Context, status *int) ([]*domain.Role, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	q := r.db.WithContext(ctx)
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	var list []*domain.Role
	if err := q.Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) UpdateRole(ctx context.Context, role *domain.Role) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(role).Error
}

func (r *repository) DeleteRole(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Role{}).Error
}

func (r *repository) SetRolePermissions(ctx context.Context, roleID uuid.UUID, permissionIDs []uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	var role domain.Role
	if err := r.db.WithContext(ctx).First(&role, roleID).Error; err != nil {
		return err
	}
	var permissions []domain.Permission
	if len(permissionIDs) > 0 {
		if err := r.db.WithContext(ctx).Where("id IN ?", permissionIDs).Find(&permissions).Error; err != nil {
			return err
		}
	}
	return r.db.WithContext(ctx).Model(&role).Association("Permissions").Replace(permissions)
}

func (r *repository) SetRoleMenus(ctx context.Context, roleID uuid.UUID, menuIDs []uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	var role domain.Role
	if err := r.db.WithContext(ctx).First(&role, roleID).Error; err != nil {
		return err
	}
	var menus []domain.Menu
	if len(menuIDs) > 0 {
		if err := r.db.WithContext(ctx).Where("id IN ?", menuIDs).Find(&menus).Error; err != nil {
			return err
		}
	}
	return r.db.WithContext(ctx).Model(&role).Association("Menus").Replace(menus)
}

// Permission methods

func (r *repository) CreatePermission(ctx context.Context, p *domain.Permission) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	ensureID(&p.ID)
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *repository) GetPermission(ctx context.Context, id uuid.UUID) (*domain.Permission, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var p domain.Permission
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) GetPermissionByCode(ctx context.Context, code string) (*domain.Permission, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var p domain.Permission
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *repository) ListPermissions(ctx context.Context) ([]*domain.Permission, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var list []*domain.Permission
	if err := r.db.WithContext(ctx).Order("code").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) UpdatePermission(ctx context.Context, p *domain.Permission) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(p).Error
}

func (r *repository) DeletePermission(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Permission{}).Error
}

func (r *repository) GetPermissionsByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]*domain.Permission, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	if len(roleIDs) == 0 {
		return []*domain.Permission{}, nil
	}
	var permissions []*domain.Permission
	err := r.db.WithContext(ctx).
		Distinct().
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id IN ?", roleIDs).
		Find(&permissions).Error
	if err != nil {
		return nil, err
	}
	return permissions, nil
}

// Menu methods

func (r *repository) CreateMenu(ctx context.Context, m *domain.Menu) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	ensureID(&m.ID)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *repository) GetMenu(ctx context.Context, id uuid.UUID) (*domain.Menu, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var m domain.Menu
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *repository) GetMenuByCode(ctx context.Context, code string) (*domain.Menu, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var m domain.Menu
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *repository) ListMenus(ctx context.Context, status *int) ([]*domain.Menu, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	q := r.db.WithContext(ctx)
	if status != nil {
		q = q.Where("status = ?", *status)
	}
	var list []*domain.Menu
	if err := q.Order("sort, created_at").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) UpdateMenu(ctx context.Context, m *domain.Menu) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *repository) DeleteMenu(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Menu{}).Error
}

func (r *repository) GetMenusByRoleIDs(ctx context.Context, roleIDs []uuid.UUID) ([]*domain.Menu, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	if len(roleIDs) == 0 {
		return []*domain.Menu{}, nil
	}
	var menus []*domain.Menu
	err := r.db.WithContext(ctx).
		Distinct().
		Joins("JOIN role_menus ON role_menus.menu_id = menus.id").
		Where("role_menus.role_id IN ?", roleIDs).
		Order("menus.sort, menus.created_at").
		Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}

// MediaAsset methods

func (r *repository) CreateMediaAsset(ctx context.Context, a *domain.MediaAsset) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	ensureID(&a.ID)
	return r.db.WithContext(ctx).Create(a).Error
}

func (r *repository) GetMediaAsset(ctx context.Context, id uuid.UUID) (*domain.MediaAsset, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var a domain.MediaAsset
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&a).Error; err != nil {
		return nil, err
	}
	return &a, nil
}

func (r *repository) ListMediaAssets(ctx context.Context, filter domain.MediaAssetFilter) ([]*domain.MediaAsset, int64, error) {
	if err := r.checkDB(); err != nil {
		return nil, 0, err
	}

	q := r.db.WithContext(ctx).Model(&domain.MediaAsset{})

	if filter.Type != nil {
		q = q.Where("type = ?", *filter.Type)
	}
	if filter.SourceType != nil {
		q = q.Where("source_type = ?", *filter.SourceType)
	}
	if filter.SourceID != nil {
		q = q.Where("source_id = ?", *filter.SourceID)
	}
	if filter.ParentID != nil {
		q = q.Where("parent_id = ?", *filter.ParentID)
	}
	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
	}
	if len(filter.Tags) > 0 {
		q = q.Where("tags @> ?", filter.Tags)
	}
	if filter.From != nil {
		q = q.Where("created_at >= ?", *filter.From)
	}
	if filter.To != nil {
		q = q.Where("created_at <= ?", *filter.To)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []*domain.MediaAsset
	if err := q.Limit(filter.Limit).Offset(filter.Offset).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *repository) UpdateMediaAsset(ctx context.Context, a *domain.MediaAsset) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(a).Error
}

func (r *repository) DeleteMediaAsset(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.MediaAsset{}).Error
}

func (r *repository) ListMediaAssetsBySource(ctx context.Context, sourceID uuid.UUID) ([]*domain.MediaAsset, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var list []*domain.MediaAsset
	if err := r.db.WithContext(ctx).Where("source_id = ?", sourceID).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) ListMediaAssetsByParent(ctx context.Context, parentID uuid.UUID) ([]*domain.MediaAsset, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var list []*domain.MediaAsset
	if err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// Operator methods

func (r *repository) CreateOperator(ctx context.Context, o *domain.Operator) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	ensureID(&o.ID)
	return r.db.WithContext(ctx).Create(o).Error
}

func (r *repository) GetOperator(ctx context.Context, id uuid.UUID) (*domain.Operator, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var o domain.Operator
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&o).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *repository) GetOperatorByCode(ctx context.Context, code string) (*domain.Operator, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var o domain.Operator
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&o).Error; err != nil {
		return nil, err
	}
	return &o, nil
}

func (r *repository) ListOperators(ctx context.Context, filter domain.OperatorFilter) ([]*domain.Operator, int64, error) {
	if err := r.checkDB(); err != nil {
		return nil, 0, err
	}

	q := r.db.WithContext(ctx).Model(&domain.Operator{})

	if filter.Category != nil {
		q = q.Where("category = ?", *filter.Category)
	}
	if filter.Type != nil {
		q = q.Where("type = ?", *filter.Type)
	}
	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
	}
	if filter.IsBuiltin != nil {
		q = q.Where("is_builtin = ?", *filter.IsBuiltin)
	}
	if len(filter.Tags) > 0 {
		q = q.Where("tags @> ?", filter.Tags)
	}
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		q = q.Where("name ILIKE ? OR description ILIKE ? OR code ILIKE ?", keyword, keyword, keyword)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []*domain.Operator
	if err := q.Limit(filter.Limit).Offset(filter.Offset).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *repository) UpdateOperator(ctx context.Context, o *domain.Operator) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(o).Error
}

func (r *repository) DeleteOperator(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Operator{}).Error
}

func (r *repository) ListEnabledOperators(ctx context.Context) ([]*domain.Operator, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var list []*domain.Operator
	if err := r.db.WithContext(ctx).Where("status = ?", domain.OperatorStatusEnabled).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) ListOperatorsByCategory(ctx context.Context, category domain.OperatorCategory) ([]*domain.Operator, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var list []*domain.Operator
	if err := r.db.WithContext(ctx).Where("category = ?", category).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// Workflow methods

func (r *repository) CreateWorkflow(ctx context.Context, w *domain.Workflow) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	ensureID(&w.ID)
	return r.db.WithContext(ctx).Create(w).Error
}

func (r *repository) GetWorkflow(ctx context.Context, id uuid.UUID) (*domain.Workflow, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var w domain.Workflow
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&w).Error; err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *repository) GetWorkflowByCode(ctx context.Context, code string) (*domain.Workflow, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var w domain.Workflow
	if err := r.db.WithContext(ctx).Where("code = ?", code).First(&w).Error; err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *repository) GetWorkflowWithNodes(ctx context.Context, id uuid.UUID) (*domain.Workflow, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var w domain.Workflow
	if err := r.db.WithContext(ctx).
		Preload("Nodes.Operator").
		Preload("Edges").
		Where("id = ?", id).
		First(&w).Error; err != nil {
		return nil, err
	}
	return &w, nil
}

func (r *repository) ListWorkflows(ctx context.Context, filter domain.WorkflowFilter) ([]*domain.Workflow, int64, error) {
	if err := r.checkDB(); err != nil {
		return nil, 0, err
	}

	q := r.db.WithContext(ctx).Model(&domain.Workflow{})

	if filter.Status != nil {
		q = q.Where("status = ?", *filter.Status)
	}
	if filter.TriggerType != nil {
		q = q.Where("trigger_type = ?", *filter.TriggerType)
	}
	if len(filter.Tags) > 0 {
		q = q.Where("tags @> ?", filter.Tags)
	}
	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		q = q.Where("name ILIKE ? OR description ILIKE ? OR code ILIKE ?", keyword, keyword, keyword)
	}

	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var list []*domain.Workflow
	if err := q.Limit(filter.Limit).Offset(filter.Offset).Order("created_at DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *repository) UpdateWorkflow(ctx context.Context, w *domain.Workflow) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Save(w).Error
}

func (r *repository) DeleteWorkflow(ctx context.Context, id uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&domain.Workflow{}).Error
}

func (r *repository) ListEnabledWorkflows(ctx context.Context) ([]*domain.Workflow, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var list []*domain.Workflow
	if err := r.db.WithContext(ctx).
		Preload("Nodes.Operator").
		Preload("Edges").
		Where("status = ?", domain.WorkflowStatusEnabled).
		Order("created_at DESC").
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// WorkflowNode methods

func (r *repository) CreateWorkflowNode(ctx context.Context, n *domain.WorkflowNode) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	ensureID(&n.ID)
	return r.db.WithContext(ctx).Create(n).Error
}

func (r *repository) ListWorkflowNodes(ctx context.Context, workflowID uuid.UUID) ([]*domain.WorkflowNode, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var list []*domain.WorkflowNode
	if err := r.db.WithContext(ctx).
		Preload("Operator").
		Where("workflow_id = ?", workflowID).
		Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) DeleteWorkflowNodes(ctx context.Context, workflowID uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("workflow_id = ?", workflowID).Delete(&domain.WorkflowNode{}).Error
}

// WorkflowEdge methods

func (r *repository) CreateWorkflowEdge(ctx context.Context, e *domain.WorkflowEdge) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	ensureID(&e.ID)
	return r.db.WithContext(ctx).Create(e).Error
}

func (r *repository) ListWorkflowEdges(ctx context.Context, workflowID uuid.UUID) ([]*domain.WorkflowEdge, error) {
	if err := r.checkDB(); err != nil {
		return nil, err
	}
	var list []*domain.WorkflowEdge
	if err := r.db.WithContext(ctx).Where("workflow_id = ?", workflowID).Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

func (r *repository) DeleteWorkflowEdges(ctx context.Context, workflowID uuid.UUID) error {
	if err := r.checkDB(); err != nil {
		return err
	}
	return r.db.WithContext(ctx).Where("workflow_id = ?", workflowID).Delete(&domain.WorkflowEdge{}).Error
}
