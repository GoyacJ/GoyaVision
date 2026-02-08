package dto

import (
	"time"

	"goyavision/internal/domain"

	"github.com/google/uuid"
)

type GetUserAssetSummaryQuery struct {
	UserID uuid.UUID
}

type UserAssetSummaryResult struct {
	Points       int64
	Balance      float64
	Subscription string
	MemberLevel  string
}

type ListUserTransactionsQuery struct {
	UserID uuid.UUID
	Limit  int
	Offset int
}

type ListUserTransactionsResult struct {
	Items []*domain.TransactionRecord
	Total int64
}

type ListUserPointRecordsQuery struct {
	UserID uuid.UUID
	Limit  int
	Offset int
}

type ListUserPointRecordsResult struct {
	Items []*domain.PointRecord
	Total int64
}

type GetUserUsageStatsQuery struct {
	UserID uuid.UUID
	Start  time.Time
	End    time.Time
}

type UsageStats struct {
	OperatorCalls int64
	AIModelCalls  int64
	TokenUsage    int64
}

type RechargeCommand struct {
	UserID      uuid.UUID
	Amount      float64
	Channel     string
	Description string
	ClientIP    string
}

type RechargeResult struct {
	OrderNo string
	PayURL  string
	QRCode  string
}

type PaymentNotifyCommand struct {
	Channel string
	Payload interface{} // 原始请求对象
}

type CheckInCommand struct {
	UserID uuid.UUID
}

type SubscribeCommand struct {
	UserID   uuid.UUID
	PlanName string
}
