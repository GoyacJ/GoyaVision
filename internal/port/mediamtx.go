package port

import "context"

// MediaMTXPathSync MediaMTX 路径同步接口（创建/更新/删除 path，由 adapter 实现）
type MediaMTXPathSync interface {
	AddPath(ctx context.Context, pathName string, source string) error
	DeletePath(ctx context.Context, pathName string) error
	PatchPath(ctx context.Context, pathName string, source string) error
}
