package payment

import (
	"context"
	"fmt"
	"net/http"

	"goyavision/config"
	"goyavision/internal/port"
)

type GoPayAdapter struct {
	alipay *AlipayAdapter
	wechat *WechatAdapter
}

func NewGoPayAdapter(cfg config.Payment) (*GoPayAdapter, error) {
	alipayAdapter, err := NewAlipayAdapter(cfg.Alipay)
	if err != nil {
		return nil, err
	}
	wechatAdapter, err := NewWechatAdapter(cfg.Wechat)
	if err != nil {
		return nil, err
	}
	return &GoPayAdapter{
		alipay: alipayAdapter,
		wechat: wechatAdapter,
	}, nil
}

func (g *GoPayAdapter) UnifyOrder(ctx context.Context, channel port.PaymentChannel, order port.PayOrder) (*port.PayResponse, error) {
	switch channel {
	case port.Alipay:
		return g.alipay.UnifyOrder(ctx, order)
	case port.Wechat:
		return g.wechat.UnifyOrder(ctx, order)
	case port.UnionPay:
		return nil, fmt.Errorf("unionpay not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported payment channel: %s", channel)
	}
}

func (g *GoPayAdapter) ParseNotify(channel port.PaymentChannel, req *http.Request) (*port.NotifyResult, error) {
	switch channel {
	case port.Alipay:
		return g.alipay.ParseNotify(req)
	case port.Wechat:
		return g.wechat.ParseNotify(req)
	case port.UnionPay:
		return nil, fmt.Errorf("unionpay not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported payment channel: %s", channel)
	}
}
