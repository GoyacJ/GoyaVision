package port

import (
	"context"
	"net/http"
)

// PaymentChannel 支付渠道
type PaymentChannel string

const (
	Alipay   PaymentChannel = "alipay"
	Wechat   PaymentChannel = "wechat"
	UnionPay PaymentChannel = "unionpay"
)

// PayOrder 支付订单信息
type PayOrder struct {
	OrderNo     string
	Amount      float64
	Description string
	ClientIP    string
}

// PayResponse 支付响应
type PayResponse struct {
	OrderNo    string
	PayURL     string // 支付链接或表单 HTML
	QRCode     string // 二维码链接 (微信 Native)
	RawData    string // 原始响应
}

// PaymentGateway 支付网关接口
type PaymentGateway interface {
	// UnifyOrder 统一下单
	UnifyOrder(ctx context.Context, channel PaymentChannel, order PayOrder) (*PayResponse, error)
	// ParseNotify 解析并验证异步通知
	ParseNotify(channel PaymentChannel, req *http.Request) (*NotifyResult, error)
}

// NotifyResult 异步通知结果
type NotifyResult struct {
	OrderNo     string
	TradeNo     string
	Amount      float64
	Status      string // Success, Failed
	RawResponse string
}
