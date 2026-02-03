package mediamtx

import (
	"context"

	"goyavision/internal/port"
)

var _ port.MediaMTXPathSync = (*PathSync)(nil)

// PathSync 实现 port.MediaMTXPathSync，用于 MediaSource 与 MediaMTX path 同步
type PathSync struct {
	client *Client
}

// NewPathSync 创建 PathSync
func NewPathSync(client *Client) *PathSync {
	return &PathSync{client: client}
}

// AddPath 添加 path，source 为拉流 URL 或 "publisher"（推流）
func (p *PathSync) AddPath(ctx context.Context, pathName string, source string) error {
	cfg := &PathConfig{Source: source}
	return p.client.AddPath(ctx, pathName, cfg)
}

// DeletePath 删除 path
func (p *PathSync) DeletePath(ctx context.Context, pathName string) error {
	return p.client.DeletePath(ctx, pathName)
}

// PatchPath 更新 path 的 source
func (p *PathSync) PatchPath(ctx context.Context, pathName string, source string) error {
	cfg := &PathConfig{Source: source}
	return p.client.PatchPath(ctx, pathName, cfg)
}
