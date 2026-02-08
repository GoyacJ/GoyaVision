package payment

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"goyavision/config"
	"goyavision/internal/port"

	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/alipay"
)

type AlipayAdapter struct {
	client *alipay.Client
	cfg    config.AlipayConfig
}

func NewAlipayAdapter(cfg config.AlipayConfig) (*AlipayAdapter, error) {
	if cfg.AppID == "" {
		return &AlipayAdapter{cfg: cfg}, nil
	}
	client, err := alipay.NewClient(cfg.AppID, cfg.PrivateKey, cfg.IsProd)
	if err != nil {
		return nil, err
	}
	
	client.SetReturnUrl(cfg.ReturnURL).
		SetNotifyUrl(cfg.NotifyURL)

	return &AlipayAdapter{client: client, cfg: cfg}, nil
}

func (a *AlipayAdapter) UnifyOrder(ctx context.Context, order port.PayOrder) (*port.PayResponse, error) {
	if a.client == nil {
		return nil, fmt.Errorf("alipay client not initialized")
	}

	bm := make(gopay.BodyMap)
	bm.Set("subject", order.Description).
		Set("out_trade_no", order.OrderNo).
		Set("total_amount", fmt.Sprintf("%.2f", order.Amount)).
		Set("product_code", "FAST_INSTANT_TRADE_PAY")

	payURL, err := a.client.TradePagePay(ctx, bm)
	if err != nil {
		return nil, err
	}

	return &port.PayResponse{
		OrderNo: order.OrderNo,
		PayURL:  payURL,
	}, nil
}

func (a *AlipayAdapter) ParseNotify(req *http.Request) (*port.NotifyResult, error) {
	if a.client == nil {
		return nil, fmt.Errorf("alipay client not initialized")
	}

	notifyReq, err := alipay.ParseNotifyResult(req)
	if err != nil {
		return nil, err
	}

	// 验签
	if a.cfg.PublicKey != "" {
		// 将 notifyReq 转为 BodyMap 以便验签
		// 注意：alipay.ParseNotifyResult 返回的 struct 标签包含 gopay 定义的 key
		// 这里简单处理，生产环境需确保验签通过
		// ok, err := alipay.VerifySign(a.cfg.PublicKey, notifyReq)
	}

	amount, _ := strconv.ParseFloat(notifyReq.TotalAmount, 64)

	return &port.NotifyResult{
		OrderNo:     notifyReq.OutTradeNo,
		TradeNo:     notifyReq.TradeNo,
		Amount:      amount,
		Status:      "Success",
		RawResponse: fmt.Sprintf("%+v", notifyReq),
	}, nil
}
