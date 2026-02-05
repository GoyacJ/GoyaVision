package mediamtx

import (
	"context"

	"goyavision/internal/adapter/mediamtx"
	"goyavision/internal/app/port"
)

// Gateway MediaMTX 网关实现（适配现有 Client）
type Gateway struct {
	client *mediamtx.Client
}

// NewGateway 创建 MediaMTX 网关
func NewGateway(baseURL string) port.MediaGateway {
	return &Gateway{
		client: mediamtx.NewClient(baseURL),
	}
}

// PatchPath 更新流路径配置
func (g *Gateway) PatchPath(ctx context.Context, pathName, source string) error {
	cfg := &mediamtx.PathConfig{
		Source: source,
	}
	return g.client.PatchPath(ctx, pathName, cfg)
}

// AddPath 添加流路径
func (g *Gateway) AddPath(ctx context.Context, pathName, source string) error {
	cfg := &mediamtx.PathConfig{
		Source: source,
	}
	return g.client.AddPath(ctx, pathName, cfg)
}

// DeletePath 删除流路径
func (g *Gateway) DeletePath(ctx context.Context, pathName string) error {
	return g.client.DeletePath(ctx, pathName)
}

// GetPathStatus 获取路径状态
func (g *Gateway) GetPathStatus(ctx context.Context, pathName string) (*port.PathStatus, error) {
	path, err := g.client.GetPath(ctx, pathName)
	if err != nil {
		return nil, err
	}

	source := ""
	if path.Source != nil {
		source = path.Source.ID
	}
	status := &port.PathStatus{
		Name:          path.Name,
		Source:        source,
		Ready:         path.Ready,
		NumReaders:    len(path.Readers),
		BytesReceived: path.BytesReceived,
	}
	return status, nil
}

// ListPaths 列出所有路径
func (g *Gateway) ListPaths(ctx context.Context) ([]port.PathInfo, error) {
	paths, err := g.client.ListPaths(ctx, 1, 1000)
	if err != nil {
		return nil, err
	}

	result := make([]port.PathInfo, 0, len(paths.Items))
	for _, path := range paths.Items {
		source := ""
		if path.Source != nil {
			source = path.Source.ID
		}
		result = append(result, port.PathInfo{
			Name:   path.Name,
			Source: source,
		})
	}
	return result, nil
}

// StartRecord 开始录制
func (g *Gateway) StartRecord(ctx context.Context, pathName string) error {
	return g.client.EnableRecording(ctx, pathName, "", "fmp4", "")
}

// StopRecord 停止录制
func (g *Gateway) StopRecord(ctx context.Context, pathName string) error {
	return g.client.DisableRecording(ctx, pathName)
}

// GetRecordStatus 获取录制状态
func (g *Gateway) GetRecordStatus(ctx context.Context, pathName string) (*port.RecordStatus, error) {
	rec, err := g.client.GetRecordings(ctx, pathName)
	if err != nil {
		return nil, err
	}

	recording := len(rec.Segments) > 0
	var startTime *int64
	if recording && len(rec.Segments) > 0 {
		t := rec.Segments[0].Start.Unix()
		startTime = &t
	}

	status := &port.RecordStatus{
		PathName:  pathName,
		Recording: recording,
		StartTime: startTime,
	}
	return status, nil
}

// Ping 健康检查
func (g *Gateway) Ping(ctx context.Context) error {
	_, err := g.client.GetInfo(ctx)
	return err
}
