package handler

import (
	"net/http"
	"strconv"
	"time"

	"goyavision/internal/api/dto"
	"goyavision/internal/api/middleware"
	"goyavision/internal/api/response"
	appdto "goyavision/internal/app/dto"
	"goyavision/internal/port"

	"github.com/labstack/echo/v4"
)

// GetUserAssetSummary 获取用户资产概览
func (h *Handlers) GetUserAssetSummary(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	result, err := h.GetUserAssetSummaryHandler.Handle(c.Request().Context(), appdto.GetUserAssetSummaryQuery{UserID: userID})
	if err != nil {
		return err
	}

	return response.OK(c, dto.UserAssetSummary{
		Points:       result.Points,
		Balance:      result.Balance,
		Subscription: result.Subscription,
		MemberLevel:  result.MemberLevel,
	})
}

// GetUserTransactions 获取交易记录
func (h *Handlers) GetUserTransactions(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if limit <= 0 {
		limit = 10
	}

	result, err := h.ListUserTransactionsHandler.Handle(c.Request().Context(), appdto.ListUserTransactionsQuery{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return err
	}

	records := make([]dto.TransactionRecord, len(result.Items))
	for i, item := range result.Items {
		records[i] = dto.TransactionRecord{
			ID:        item.ID,
			Type:      item.Type,
			Method:    item.Method,
			Amount:    item.Amount,
			Status:    item.Status,
			CreatedAt: item.CreatedAt,
		}
	}

	return response.OK(c, records)
}

// GetUserPoints 获取积分记录
func (h *Handlers) GetUserPoints(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	if limit <= 0 {
		limit = 10
	}

	result, err := h.ListUserPointRecordsHandler.Handle(c.Request().Context(), appdto.ListUserPointRecordsQuery{
		UserID: userID,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return err
	}

	records := make([]dto.PointRecord, len(result.Items))
	for i, item := range result.Items {
		records[i] = dto.PointRecord{
			ID:        int64(item.ID),
			Type:      item.Type,
			Change:    item.Change,
			Balance:   item.Balance,
			CreatedAt: item.CreatedAt,
		}
	}

	return response.OK(c, records)
}

// GetUserUsageStats 获取使用统计
func (h *Handlers) GetUserUsageStats(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	// 默认获取最近 30 天
	end := time.Now()
	start := end.AddDate(0, 0, -30)

	result, err := h.GetUserUsageStatsHandler.Handle(c.Request().Context(), appdto.GetUserUsageStatsQuery{
		UserID: userID,
		Start:  start,
		End:    end,
	})
	if err != nil {
		return err
	}

	// 聚合统计
	var totalStats dto.UsageStats
	for _, s := range result {
		totalStats.OperatorCalls += s.OperatorCalls
		totalStats.AIModelCalls += s.AIModelCalls
		totalStats.TokenUsage += s.TokenUsage
	}

	return response.OK(c, totalStats)
}

// Recharge 发起充值
func (h *Handlers) Recharge(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	var req struct {
		Amount  float64 `json:"amount"`
		Channel string  `json:"channel"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	result, err := h.RechargeHandler.Handle(c.Request().Context(), appdto.RechargeCommand{
		UserID:      userID,
		Amount:      req.Amount,
		Channel:     req.Channel,
		Description: "Balance Recharge",
		ClientIP:    c.RealIP(),
	})
	if err != nil {
		return err
	}

	return response.OK(c, result)
}

// CheckIn 签到
func (h *Handlers) CheckIn(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	err := h.CheckInHandler.Handle(c.Request().Context(), appdto.CheckInCommand{UserID: userID})
	if err != nil {
		return err
	}

	return response.OK(c, map[string]string{"message": "Check-in successful"})
}

// Subscribe 订阅
func (h *Handlers) Subscribe(c echo.Context) error {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user not authenticated")
	}

	var req struct {
		PlanName string `json:"plan_name"`
	}
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid request")
	}

	err := h.SubscribeHandler.Handle(c.Request().Context(), appdto.SubscribeCommand{
		UserID:   userID,
		PlanName: req.PlanName,
	})
	if err != nil {
		return err
	}

	return response.OK(c, map[string]string{"message": "Subscription updated"})
}

// HandlePaymentNotify 处理支付回调
func (h *Handlers) HandlePaymentNotify(c echo.Context) error {
	channel := c.Param("channel")
	// 支付网关需要原始 Request
	_, err := h.RechargeHandler.HandleNotify(c.Request().Context(), port.PaymentChannel(channel), c.Request())
	if err != nil {
		return err
	}

	// 这里的 Recharge.HandleNotify 应该在 RechargeHandler 中增加或者作为独立逻辑
	// 为了简化，我直接在这里处理逻辑，或者完善 RechargeHandler
	
	// 根据 result.Status 更新订单和余额
	// ... (在正式生产环境，这部分逻辑应放在 Service/Handler 中)
	
	return c.String(http.StatusOK, "success")
}

// RegisterUserAssetRoutes 注册用户资产路由
func RegisterUserAssetRoutes(g *echo.Group, h *Handlers) {
	asset := g.Group("/user/assets")
	{
		asset.GET("/summary", h.GetUserAssetSummary)
		asset.GET("/transactions", h.GetUserTransactions)
		asset.GET("/points", h.GetUserPoints)
		asset.GET("/usage", h.GetUserUsageStats)
		asset.POST("/recharge", h.Recharge)
		asset.POST("/checkin", h.CheckIn)
		asset.POST("/subscribe", h.Subscribe)
	}

	// 支付回调路由 (不需要 JWT)
	payment := g.Group("/payment/notify")
	{
		payment.POST("/:channel", h.HandlePaymentNotify)
	}
}
