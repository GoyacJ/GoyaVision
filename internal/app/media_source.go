package app

import (
	"context"
	"errors"

	"goyavision/internal/domain"
	"goyavision/internal/port"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrMediaSourceHasAssets = errors.New("存在关联的流媒体资产，请先删除或解除关联")
var ErrMediaMTXUnavailable = errors.New("流媒体服务暂不可用")

// CreateMediaSourceRequest 创建媒体源请求
type CreateMediaSourceRequest struct {
	Name     string
	Type     domain.SourceType
	URL      string
	Protocol string
	Enabled  bool
}

// UpdateMediaSourceRequest 更新媒体源请求
type UpdateMediaSourceRequest struct {
	Name     *string
	URL      *string
	Protocol *string
	Enabled  *bool
}

// MediaSourceService 媒体源服务
type MediaSourceService struct {
	repo port.Repository
	mtx  port.MediaMTXPathSync
}

// NewMediaSourceService 创建媒体源服务
func NewMediaSourceService(repo port.Repository, mtx port.MediaMTXPathSync) *MediaSourceService {
	return &MediaSourceService{repo: repo, mtx: mtx}
}

// Create 创建媒体源：生成 path_name → MediaMTX AddPath → 写表；写表失败则 DeletePath 回滚
func (s *MediaSourceService) Create(ctx context.Context, req *CreateMediaSourceRequest) (*domain.MediaSource, error) {
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.Type != domain.SourceTypePull && req.Type != domain.SourceTypePush {
		return nil, errors.New("invalid source type")
	}
	source := req.URL
	if req.Type == domain.SourceTypePush {
		source = "publisher"
	} else if req.URL == "" {
		return nil, errors.New("url is required for pull source")
	}

	pathName := domain.GeneratePathName(req.Name)
	if err := s.mtx.AddPath(ctx, pathName, source); err != nil {
		return nil, err
	}

	src := &domain.MediaSource{
		Name:     req.Name,
		PathName: pathName,
		Type:     req.Type,
		URL:      req.URL,
		Protocol: req.Protocol,
		Enabled:  req.Enabled,
	}
	if err := s.repo.CreateMediaSource(ctx, src); err != nil {
		_ = s.mtx.DeletePath(ctx, pathName)
		return nil, err
	}
	return src, nil
}

// Get 获取媒体源
func (s *MediaSourceService) Get(ctx context.Context, id uuid.UUID) (*domain.MediaSource, error) {
	src, err := s.repo.GetMediaSource(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("media source not found")
		}
		return nil, err
	}
	return src, nil
}

// List 列出媒体源
func (s *MediaSourceService) List(ctx context.Context, filter domain.MediaSourceFilter) ([]*domain.MediaSource, int64, error) {
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	if filter.Limit > 1000 {
		filter.Limit = 1000
	}
	return s.repo.ListMediaSources(ctx, filter)
}

// Update 更新媒体源：若 url 变更则 PatchPath，再更新表
func (s *MediaSourceService) Update(ctx context.Context, id uuid.UUID, req *UpdateMediaSourceRequest) (*domain.MediaSource, error) {
	src, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		src.Name = *req.Name
	}
	if req.URL != nil {
		src.URL = *req.URL
	}
	if req.Protocol != nil {
		src.Protocol = *req.Protocol
	}
	if req.Enabled != nil {
		src.Enabled = *req.Enabled
	}
	source := src.URL
	if src.Type == domain.SourceTypePush {
		source = "publisher"
	}
	if err := s.mtx.PatchPath(ctx, src.PathName, source); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateMediaSource(ctx, src); err != nil {
		return nil, err
	}
	return src, nil
}

// Delete 删除媒体源：仅当无关联流媒体资产时允许；先 DeletePath 再删表
func (s *MediaSourceService) Delete(ctx context.Context, id uuid.UUID) error {
	assets, err := s.repo.ListMediaAssetsBySource(ctx, id)
	if err != nil {
		return err
	}
	if len(assets) > 0 {
		return ErrMediaSourceHasAssets
	}
	src, err := s.repo.GetMediaSource(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("media source not found")
		}
		return err
	}
	if err := s.mtx.DeletePath(ctx, src.PathName); err != nil {
		return err
	}
	return s.repo.DeleteMediaSource(ctx, id)
}
