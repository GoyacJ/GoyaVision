package port

import "context"

// MediaGateway MediaMTX 网关接口
//
// 职责：
//  1. 管理 MediaMTX 流媒体服务器的路径（Path）
//  2. 查询流的状态和统计信息
//  3. 控制录制功能
//
// 实现：
//   - infra/mediamtx/gateway.go (HTTP 客户端实现)
type MediaGateway interface {
	// AddPath 添加流路径
	AddPath(ctx context.Context, pathName, source string) error

	// PatchPath 更新流路径配置
	PatchPath(ctx context.Context, pathName, source string) error

	// DeletePath 删除流路径
	DeletePath(ctx context.Context, pathName string) error

	// GetPathStatus 获取路径状态
	GetPathStatus(ctx context.Context, pathName string) (*PathStatus, error)

	// ListPaths 列出所有路径
	ListPaths(ctx context.Context) ([]PathInfo, error)

	// StartRecord 开始录制
	StartRecord(ctx context.Context, pathName string) error

	// StopRecord 停止录制
	StopRecord(ctx context.Context, pathName string) error

	// GetRecordStatus 获取录制状态
	GetRecordStatus(ctx context.Context, pathName string) (*RecordStatus, error)

	// Ping 健康检查
	Ping(ctx context.Context) error
}

// PathStatus 路径状态
type PathStatus struct {
	Name          string
	Source        string
	Ready         bool
	NumReaders    int
	BytesReceived uint64
}

// PathInfo 路径信息
type PathInfo struct {
	Name   string
	Source string
}

// RecordStatus 录制状态
type RecordStatus struct {
	PathName  string
	Recording bool
	StartTime *int64 // Unix timestamp
}
