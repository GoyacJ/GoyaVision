package repo

import (
	"context"

	"goyavision/internal/domain/identity"
	"goyavision/internal/infra/persistence/mapper"
	"goyavision/internal/infra/persistence/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserIdentityRepo struct {
	db *gorm.DB
}

func NewUserIdentityRepo(db *gorm.DB) *UserIdentityRepo {
	return &UserIdentityRepo{db: db}
}

func (r *UserIdentityRepo) Create(ctx context.Context, i *identity.UserIdentity) error {
	if i.ID == uuid.Nil {
		i.ID = uuid.New()
	}
	m := mapper.UserIdentityToModel(i)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *UserIdentityRepo) Get(ctx context.Context, id uuid.UUID) (*identity.UserIdentity, error) {
	var m model.UserIdentityModel
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.UserIdentityToDomain(&m), nil
}

func (r *UserIdentityRepo) GetByIdentifier(ctx context.Context, identityType identity.IdentityType, identifier string) (*identity.UserIdentity, error) {
	var m model.UserIdentityModel
	if err := r.db.WithContext(ctx).Where("identity_type = ? AND identifier = ?", string(identityType), identifier).First(&m).Error; err != nil {
		return nil, err
	}
	return mapper.UserIdentityToDomain(&m), nil
}

func (r *UserIdentityRepo) ListByUserID(ctx context.Context, userID uuid.UUID) ([]*identity.UserIdentity, error) {
	var models []*model.UserIdentityModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}
	var result []*identity.UserIdentity
	for _, m := range models {
		result = append(result, mapper.UserIdentityToDomain(m))
	}
	return result, nil
}

func (r *UserIdentityRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return r.db.WithContext(ctx).Delete(&model.UserIdentityModel{}, id).Error
}
