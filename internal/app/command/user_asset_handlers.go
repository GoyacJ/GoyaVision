package command

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"goyavision/internal/app/dto"
	appport "goyavision/internal/app/port"
	"goyavision/internal/domain"
	"goyavision/internal/port"

	"github.com/google/uuid"
)

type RechargeHandler struct {
	uow            appport.UnitOfWork
	paymentGateway port.PaymentGateway
}

func NewRechargeHandler(uow appport.UnitOfWork, pg port.PaymentGateway) *RechargeHandler {
	return &RechargeHandler{uow: uow, paymentGateway: pg}
}

func (h *RechargeHandler) Handle(ctx context.Context, cmd dto.RechargeCommand) (*dto.RechargeResult, error) {
	orderNo := fmt.Sprintf("RECH%d%s", time.Now().Unix(), uuid.New().String()[:8])
	
	err := h.uow.Do(ctx, func(ctx context.Context, repos *appport.Repositories) error {
		// 创建交易记录
		tr := &domain.TransactionRecord{
			ID:        orderNo,
			UserID:    cmd.UserID,
			Type:      "Recharge",
			Method:    cmd.Channel,
			Amount:    cmd.Amount,
			Status:    "Pending",
			Remark:    cmd.Description,
			CreatedAt: time.Now(),
		}
		return repos.UserAssets.CreateTransactionRecord(ctx, tr)
	})
	if err != nil {
		return nil, err
	}

	// 调用支付网关
	resp, err := h.paymentGateway.UnifyOrder(ctx, port.PaymentChannel(cmd.Channel), port.PayOrder{
		OrderNo:     orderNo,
		Amount:      cmd.Amount,
		Description: cmd.Description,
		ClientIP:    cmd.ClientIP,
	})
	if err != nil {
		return nil, err
	}

	return &dto.RechargeResult{
		OrderNo: resp.OrderNo,
		PayURL:  resp.PayURL,
		QRCode:  resp.QRCode,
	}, nil
}

type HandlePaymentNotifyHandler struct {
	uow            appport.UnitOfWork
	paymentGateway port.PaymentGateway
}

func NewHandlePaymentNotifyHandler(uow appport.UnitOfWork, pg port.PaymentGateway) *HandlePaymentNotifyHandler {
	return &HandlePaymentNotifyHandler{uow: uow, paymentGateway: pg}
}

// HandleNotify 处理支付异步回调
func (h *RechargeHandler) HandleNotify(ctx context.Context, channel port.PaymentChannel, req *http.Request) (*port.NotifyResult, error) {
	// 1. 解析回调
	result, err := h.paymentGateway.ParseNotify(channel, req)
	if err != nil {
		return nil, err
	}

	// 2. 业务处理 (更新订单和余额)
	err = h.uow.Do(ctx, func(ctx context.Context, repos *appport.Repositories) error {
		order, err := repos.UserAssets.GetTransactionRecord(ctx, result.OrderNo)
		if err != nil {
			return err
		}

		if order.Status == "Success" {
			return nil // 已处理，幂等
		}

		order.Status = "Success"
		order.UpdatedAt = time.Now()
		if err := repos.UserAssets.UpdateTransactionRecord(ctx, order); err != nil {
			return err
		}

		// 增加用户余额
		balance, err := repos.UserAssets.GetUserBalance(ctx, order.UserID)
		if err != nil {
			return err
		}
		balance.Balance += order.Amount
		balance.UpdatedAt = time.Now()
		return repos.UserAssets.UpdateUserBalance(ctx, balance)
	})

	return result, err
}

type CheckInHandler struct {
	uow appport.UnitOfWork
}

func NewCheckInHandler(uow appport.UnitOfWork) *CheckInHandler {
	return &CheckInHandler{uow: uow}
}

func (h *CheckInHandler) Handle(ctx context.Context, cmd dto.CheckInCommand) error {
	return h.uow.Do(ctx, func(ctx context.Context, repos *appport.Repositories) error {
		points := int64(50)
		
		balance, err := repos.UserAssets.GetUserBalance(ctx, cmd.UserID)
		if err != nil {
			return err
		}

		balance.Points += points
		balance.UpdatedAt = time.Now()
		if err := repos.UserAssets.UpdateUserBalance(ctx, balance); err != nil {
			return err
		}

		pr := &domain.PointRecord{
			UserID:    cmd.UserID,
			Type:      "Check-in",
			Change:    points,
			Balance:   balance.Points,
			Remark:    fmt.Sprintf("Daily Check-in at %s", time.Now().Format("2006-01-02")),
			CreatedAt: time.Now(),
		}
		return repos.UserAssets.CreatePointRecord(ctx, pr)
	})
}

type SubscribeHandler struct {
	uow appport.UnitOfWork
}

func NewSubscribeHandler(uow appport.UnitOfWork) *SubscribeHandler {
	return &SubscribeHandler{uow: uow}
}

func (h *SubscribeHandler) Handle(ctx context.Context, cmd dto.SubscribeCommand) error {
	return h.uow.Do(ctx, func(ctx context.Context, repos *appport.Repositories) error {
		now := time.Now()
		endDate := now.AddDate(0, 1, 0)

		sub, err := repos.UserAssets.GetUserSubscription(ctx, cmd.UserID)
		if err != nil {
			return err
		}

		if sub == nil {
			sub = &domain.UserSubscription{
				UserID:    cmd.UserID,
				PlanName:  cmd.PlanName,
				Status:    "active",
				StartDate: &now,
				EndDate:   &endDate,
				CreatedAt: now,
			}
		} else {
			sub.PlanName = cmd.PlanName
			sub.EndDate = &endDate
			sub.UpdatedAt = now
		}

		return repos.UserAssets.UpdateUserSubscription(ctx, sub)
	})
}
