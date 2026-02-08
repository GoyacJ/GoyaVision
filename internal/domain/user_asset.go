package domain

import (
	"time"

	"github.com/google/uuid"
)

// UserBalance 用户余额与积分
type UserBalance struct {
	UserID    uuid.UUID  `json:"user_id"`
	Balance   float64    `json:"balance"`
	Points    int64      `json:"points"`
	Level     string     `json:"level"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// UserSubscription 用户订阅
type UserSubscription struct {
	ID        uint       `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	PlanName  string     `json:"plan_name"`
	Status    string     `json:"status"`
	StartDate *time.Time `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// TransactionRecord 交易记录
type TransactionRecord struct {
	ID        string     `json:"id"`
	UserID    uuid.UUID  `json:"user_id"`
	Type      string     `json:"type"`   // Recharge, Subscription, etc.
	Method    string     `json:"method"` // Alipay, Wechat, UnionPay, Balance
	Amount    float64    `json:"amount"`
	Status    string     `json:"status"` // Pending, Success, Failed
	Remark    string     `json:"remark"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

// PointRecord 积分记录
type PointRecord struct {
	ID        uint      `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Type      string    `json:"type"` // Check-in, Purchase, Reward
	Change    int64     `json:"change"`
	Balance   int64     `json:"balance"`
	Remark    string    `json:"remark"`
	CreatedAt time.Time `json:"created_at"`
}

// UsageStats 使用统计
type UsageStats struct {
	ID            uint      `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	OperatorCalls int64     `json:"operator_calls"`
	AIModelCalls  int64     `json:"ai_model_calls"`
	TokenUsage    int64     `json:"token_usage"`
	Date          time.Time `json:"date"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at"`
}
