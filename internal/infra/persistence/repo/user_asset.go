package repo

import (
	"context"
	"time"

	"goyavision/internal/domain"
	"goyavision/internal/infra/persistence/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserAssetRepo struct {
	db *gorm.DB
}

func NewUserAssetRepo(db *gorm.DB) *UserAssetRepo {
	return &UserAssetRepo{db: db}
}

func (r *UserAssetRepo) GetUserBalance(ctx context.Context, userID uuid.UUID) (*domain.UserBalance, error) {
	var m model.UserBalance
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			m = model.UserBalance{
				UserID:  userID,
				Balance: 0,
				Points:  0,
				Level:   "Free",
			}
			if err := r.db.WithContext(ctx).Create(&m).Error; err != nil {
				return nil, err
			}
			return r.toDomainBalance(&m), nil
		}
		return nil, err
	}
	return r.toDomainBalance(&m), nil
}

func (r *UserAssetRepo) UpdateUserBalance(ctx context.Context, ub *domain.UserBalance) error {
	m := r.toModelBalance(ub)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *UserAssetRepo) CreateTransactionRecord(ctx context.Context, tr *domain.TransactionRecord) error {
	m := r.toModelTransaction(tr)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *UserAssetRepo) GetTransactionRecord(ctx context.Context, id string) (*domain.TransactionRecord, error) {
	var m model.TransactionRecord
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&m).Error; err != nil {
		return nil, err
	}
	return r.toDomainTransaction(&m), nil
}

func (r *UserAssetRepo) UpdateTransactionRecord(ctx context.Context, tr *domain.TransactionRecord) error {
	m := r.toModelTransaction(tr)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *UserAssetRepo) ListTransactionRecords(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.TransactionRecord, int64, error) {
	var ms []model.TransactionRecord
	var total int64
	db := r.db.WithContext(ctx).Where("user_id = ?", userID)
	db.Model(&model.TransactionRecord{}).Count(&total)
	if err := db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&ms).Error; err != nil {
		return nil, 0, err
	}
	res := make([]*domain.TransactionRecord, len(ms))
	for i, m := range ms {
		res[i] = r.toDomainTransaction(&m)
	}
	return res, total, nil
}

func (r *UserAssetRepo) CreatePointRecord(ctx context.Context, pr *domain.PointRecord) error {
	m := r.toModelPointRecord(pr)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *UserAssetRepo) ListPointRecords(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.PointRecord, int64, error) {
	var ms []model.PointRecord
	var total int64
	db := r.db.WithContext(ctx).Where("user_id = ?", userID)
	db.Model(&model.PointRecord{}).Count(&total)
	if err := db.Order("created_at DESC").Limit(limit).Offset(offset).Find(&ms).Error; err != nil {
		return nil, 0, err
	}
	res := make([]*domain.PointRecord, len(ms))
	for i, m := range ms {
		res[i] = r.toDomainPointRecord(&m)
	}
	return res, total, nil
}

func (r *UserAssetRepo) GetUserSubscription(ctx context.Context, userID uuid.UUID) (*domain.UserSubscription, error) {
	var m model.UserSubscription
	err := r.db.WithContext(ctx).Where("user_id = ? AND status = 'active'", userID).Order("end_date DESC").First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomainSubscription(&m), nil
}

func (r *UserAssetRepo) UpdateUserSubscription(ctx context.Context, us *domain.UserSubscription) error {
	m := r.toModelSubscription(us)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *UserAssetRepo) GetUsageStats(ctx context.Context, userID uuid.UUID, date time.Time) (*domain.UsageStats, error) {
	var m model.UsageStat
	day := date.Truncate(24 * time.Hour)
	err := r.db.WithContext(ctx).Where("user_id = ? AND date = ?", userID, day).First(&m).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return r.toDomainUsageStats(&m), nil
}

func (r *UserAssetRepo) UpdateUsageStats(ctx context.Context, us *domain.UsageStats) error {
	m := r.toModelUsageStats(us)
	return r.db.WithContext(ctx).Save(m).Error
}

func (r *UserAssetRepo) ListUsageStats(ctx context.Context, userID uuid.UUID, start, end time.Time) ([]*domain.UsageStats, error) {
	var ms []model.UsageStat
	if err := r.db.WithContext(ctx).Where("user_id = ? AND date >= ? AND date <= ?", userID, start, end).Order("date ASC").Find(&ms).Error; err != nil {
		return nil, err
	}
	res := make([]*domain.UsageStats, len(ms))
	for i, m := range ms {
		res[i] = r.toDomainUsageStats(&m)
	}
	return res, nil
}

func (r *UserAssetRepo) toDomainBalance(m *model.UserBalance) *domain.UserBalance {
	return &domain.UserBalance{
		UserID:    m.UserID,
		Balance:   m.Balance,
		Points:    m.Points,
		Level:     m.Level,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (r *UserAssetRepo) toModelBalance(d *domain.UserBalance) *model.UserBalance {
	return &model.UserBalance{
		UserID:    d.UserID,
		Balance:   d.Balance,
		Points:    d.Points,
		Level:     d.Level,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func (r *UserAssetRepo) toDomainTransaction(m *model.TransactionRecord) *domain.TransactionRecord {
	return &domain.TransactionRecord{
		ID:        m.ID,
		UserID:    m.UserID,
		Type:      m.Type,
		Method:    m.Method,
		Amount:    m.Amount,
		Status:    m.Status,
		Remark:    m.Remark,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (r *UserAssetRepo) toModelTransaction(d *domain.TransactionRecord) *model.TransactionRecord {
	return &model.TransactionRecord{
		ID:        d.ID,
		UserID:    d.UserID,
		Type:      d.Type,
		Method:    d.Method,
		Amount:    d.Amount,
		Status:    d.Status,
		Remark:    d.Remark,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func (r *UserAssetRepo) toDomainPointRecord(m *model.PointRecord) *domain.PointRecord {
	return &domain.PointRecord{
		ID:        m.ID,
		UserID:    m.UserID,
		Type:      m.Type,
		Change:    m.Change,
		Balance:   m.Balance,
		Remark:    m.Remark,
		CreatedAt: m.CreatedAt,
	}
}

func (r *UserAssetRepo) toModelPointRecord(d *domain.PointRecord) *model.PointRecord {
	return &model.PointRecord{
		ID:        d.ID,
		UserID:    d.UserID,
		Type:      d.Type,
		Change:    d.Change,
		Balance:   d.Balance,
		Remark:    d.Remark,
		CreatedAt: d.CreatedAt,
	}
}

func (r *UserAssetRepo) toDomainSubscription(m *model.UserSubscription) *domain.UserSubscription {
	return &domain.UserSubscription{
		ID:        m.ID,
		UserID:    m.UserID,
		PlanName:  m.PlanName,
		Status:    m.Status,
		StartDate: &m.StartDate,
		EndDate:   &m.EndDate,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func (r *UserAssetRepo) toModelSubscription(d *domain.UserSubscription) *model.UserSubscription {
	m := &model.UserSubscription{
		ID:        d.ID,
		UserID:    d.UserID,
		PlanName:  d.PlanName,
		Status:    d.Status,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
	if d.StartDate != nil {
		m.StartDate = *d.StartDate
	}
	if d.EndDate != nil {
		m.EndDate = *d.EndDate
	}
	return m
}

func (r *UserAssetRepo) toDomainUsageStats(m *model.UsageStat) *domain.UsageStats {
	return &domain.UsageStats{
		ID:            m.ID,
		UserID:        m.UserID,
		OperatorCalls: m.OperatorCalls,
		AIModelCalls:  m.AIModelCalls,
		TokenUsage:    m.TokenUsage,
		Date:          m.Date,
		CreatedAt:     m.CreatedAt,
		UpdatedAt:     m.UpdatedAt,
	}
}

func (r *UserAssetRepo) toModelUsageStats(d *domain.UsageStats) *model.UsageStat {
	return &model.UsageStat{
		ID:            d.ID,
		UserID:        d.UserID,
		OperatorCalls: d.OperatorCalls,
		AIModelCalls:  d.AIModelCalls,
		TokenUsage:    d.TokenUsage,
		Date:          d.Date,
		CreatedAt:     d.CreatedAt,
		UpdatedAt:     d.UpdatedAt,
	}
}
