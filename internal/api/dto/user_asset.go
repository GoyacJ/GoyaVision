package dto

import "time"

// UserAssetSummary 综合资产概览
type UserAssetSummary struct {
	Points       int64   `json:"points"`
	Balance      float64 `json:"balance"`
	Subscription string  `json:"subscription"`
	MemberLevel  string  `json:"member_level"`
}

// PaymentMethod 支付方式
type PaymentMethod struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

// TransactionRecord 交易记录
type TransactionRecord struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Method    string    `json:"method"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// PointRecord 积分记录
type PointRecord struct {
	ID        int64     `json:"id"`
	Type      string    `json:"type"`
	Change    int64     `json:"change"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

// SubscriptionPlan 订阅计划
type SubscriptionPlan struct {
	Name     string   `json:"name"`
	Price    float64  `json:"price"`
	Features []string `json:"features"`
	IsActive bool     `json:"is_active"`
}

// UsageStats 使用统计
type UsageStats struct {
	OperatorCalls int64 `json:"operator_calls"`
	AIModelCalls  int64 `json:"ai_model_calls"`
	TokenUsage    int64 `json:"token_usage"` // 单位: Token
}
