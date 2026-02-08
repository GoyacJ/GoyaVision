package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserBalance 用户余额与积分
type UserBalance struct {
	UserID    uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Balance   float64        `gorm:"type:decimal(16,2);default:0"`
	Points    int64          `gorm:"default:0"`
	Level     string         `gorm:"size:32;default:'Free'"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// UserSubscription 用户订阅
type UserSubscription struct {
	ID        uint           `gorm:"primaryKey"`
	UserID    uuid.UUID      `gorm:"type:uuid;index"`
	PlanName  string         `gorm:"size:64"` // Free, Pro, Enterprise
	Status    string         `gorm:"size:32"` // active, expired, cancelled
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// TransactionRecord 交易记录
type TransactionRecord struct {
	ID        string         `gorm:"primaryKey;size:64"`
	UserID    uuid.UUID      `gorm:"type:uuid;index"`
	Type      string         `gorm:"size:32"` // recharge, payment, refund
	Method    string         `gorm:"size:32"` // alipay, wechat, unionpay, balance
	Amount    float64        `gorm:"type:decimal(16,2)"`
	Status    string         `gorm:"size:32"` // pending, success, failed
	Remark    string         `gorm:"type:text"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// PointRecord 积分记录
type PointRecord struct {
	ID        uint           `gorm:"primaryKey"`
	UserID    uuid.UUID      `gorm:"type:uuid;index"`
	Type      string         `gorm:"size:32"` // checkin, task, consumption
	Change    int64
	Balance   int64
	Remark    string         `gorm:"type:text"`
	CreatedAt time.Time
}

// UsageStat 使用统计
type UsageStat struct {
	ID            uint           `gorm:"primaryKey"`
	UserID        uuid.UUID      `gorm:"type:uuid;index"`
	OperatorCalls int64          `gorm:"default:0"`
	AIModelCalls  int64          `gorm:"default:0"`
	TokenUsage    int64          `gorm:"default:0"`
	Date          time.Time      `gorm:"type:date;index"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}
