package query

import (
	"context"

	"goyavision/internal/app/dto"
	appport "goyavision/internal/app/port"
)

type GetUserAssetSummaryHandler struct {
	uow appport.UnitOfWork
}

func NewGetUserAssetSummaryHandler(uow appport.UnitOfWork) *GetUserAssetSummaryHandler {
	return &GetUserAssetSummaryHandler{uow: uow}
}

func (h *GetUserAssetSummaryHandler) Handle(ctx context.Context, q dto.GetUserAssetSummaryQuery) (*dto.UserAssetSummaryResult, error) {
	var res *dto.UserAssetSummaryResult
	err := h.uow.Do(ctx, func(ctx context.Context, repos *appport.Repositories) error {
		balance, err := repos.UserAssets.GetUserBalance(ctx, q.UserID)
		if err != nil {
			return err
		}

		sub, err := repos.UserAssets.GetUserSubscription(ctx, q.UserID)
		if err != nil {
			return err
		}

		plan := "Free"
		if sub != nil {
			plan = sub.PlanName
		}

		res = &dto.UserAssetSummaryResult{
			Points:       balance.Points,
			Balance:      balance.Balance,
			Subscription: plan,
			MemberLevel:  balance.Level,
		}
		return nil
	})
	return res, err
}

type ListUserTransactionsHandler struct {
	uow appport.UnitOfWork
}

func NewListUserTransactionsHandler(uow appport.UnitOfWork) *ListUserTransactionsHandler {
	return &ListUserTransactionsHandler{uow: uow}
}

func (h *ListUserTransactionsHandler) Handle(ctx context.Context, q dto.ListUserTransactionsQuery) (*dto.ListUserTransactionsResult, error) {
	var res *dto.ListUserTransactionsResult
	err := h.uow.Do(ctx, func(ctx context.Context, repos *appport.Repositories) error {
		items, total, err := repos.UserAssets.ListTransactionRecords(ctx, q.UserID, q.Limit, q.Offset)
		if err != nil {
			return err
		}
		res = &dto.ListUserTransactionsResult{
			Items: items,
			Total: total,
		}
		return nil
	})
	return res, err
}

type ListUserPointRecordsHandler struct {
	uow appport.UnitOfWork
}

func NewListUserPointRecordsHandler(uow appport.UnitOfWork) *ListUserPointRecordsHandler {
	return &ListUserPointRecordsHandler{uow: uow}
}

func (h *ListUserPointRecordsHandler) Handle(ctx context.Context, q dto.ListUserPointRecordsQuery) (*dto.ListUserPointRecordsResult, error) {
	var res *dto.ListUserPointRecordsResult
	err := h.uow.Do(ctx, func(ctx context.Context, repos *appport.Repositories) error {
		items, total, err := repos.UserAssets.ListPointRecords(ctx, q.UserID, q.Limit, q.Offset)
		if err != nil {
			return err
		}
		res = &dto.ListUserPointRecordsResult{
			Items: items,
			Total: total,
		}
		return nil
	})
	return res, err
}

type GetUserUsageStatsHandler struct {
	uow appport.UnitOfWork
}

func NewGetUserUsageStatsHandler(uow appport.UnitOfWork) *GetUserUsageStatsHandler {
	return &GetUserUsageStatsHandler{uow: uow}
}

func (h *GetUserUsageStatsHandler) Handle(ctx context.Context, q dto.GetUserUsageStatsQuery) ([]*dto.UsageStats, error) {
	var res []*dto.UsageStats
	err := h.uow.Do(ctx, func(ctx context.Context, repos *appport.Repositories) error {
		stats, err := repos.UserAssets.ListUsageStats(ctx, q.UserID, q.Start, q.End)
		if err != nil {
			return err
		}
		
		// 转换 Domain 到 DTO
		for _, s := range stats {
			res = append(res, &dto.UsageStats{
				OperatorCalls: s.OperatorCalls,
				AIModelCalls:  s.AIModelCalls,
				TokenUsage:    s.TokenUsage,
			})
		}
		return nil
	})
	return res, err
}
