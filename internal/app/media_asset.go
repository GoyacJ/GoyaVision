package app

import (
	"context"
	"errors"
	"time"

	"goyavision/internal/domain"
	"goyavision/internal/port"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateMediaAssetRequest 创建媒体资产请求
type CreateMediaAssetRequest struct {
	Type       domain.AssetType       `json:"type"`
	SourceType domain.AssetSourceType `json:"source_type"`
	SourceID   *uuid.UUID             `json:"source_id,omitempty"`
	ParentID   *uuid.UUID             `json:"parent_id,omitempty"`
	Name       string                 `json:"name"`
	Path       string                 `json:"path"`
	Duration   *float64               `json:"duration,omitempty"`
	Size       int64                  `json:"size"`
	Format     string                 `json:"format,omitempty"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Status     domain.AssetStatus     `json:"status,omitempty"`
	Tags       []string               `json:"tags,omitempty"`
}

// UpdateMediaAssetRequest 更新媒体资产请求
type UpdateMediaAssetRequest struct {
	Name     *string                `json:"name,omitempty"`
	Status   *domain.AssetStatus    `json:"status,omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
	Tags     []string               `json:"tags,omitempty"`
}

// ListMediaAssetsRequest 列出媒体资产请求
type ListMediaAssetsRequest struct {
	Type       *domain.AssetType
	SourceType *domain.AssetSourceType
	SourceID   *uuid.UUID
	ParentID   *uuid.UUID
	Status     *domain.AssetStatus
	Tags       []string
	From       *time.Time
	To         *time.Time
	Limit      int
	Offset     int
}

type MediaAssetService struct {
	repo port.Repository
}

func NewMediaAssetService(repo port.Repository) *MediaAssetService {
	return &MediaAssetService{
		repo: repo,
	}
}

// Create 创建媒体资产
func (s *MediaAssetService) Create(ctx context.Context, req *CreateMediaAssetRequest) (*domain.MediaAsset, error) {
	if req.Type == "" {
		return nil, errors.New("type is required")
	}
	if req.SourceType == "" {
		return nil, errors.New("source_type is required")
	}
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.Path == "" {
		return nil, errors.New("path is required")
	}

	if req.Type != domain.AssetTypeVideo && req.Type != domain.AssetTypeImage && req.Type != domain.AssetTypeAudio {
		return nil, errors.New("invalid asset type")
	}

	if req.SourceType != domain.AssetSourceLive &&
		req.SourceType != domain.AssetSourceVOD &&
		req.SourceType != domain.AssetSourceUpload &&
		req.SourceType != domain.AssetSourceGenerated {
		return nil, errors.New("invalid source type")
	}

	if req.ParentID != nil {
		if _, err := s.repo.GetMediaAsset(ctx, *req.ParentID); err != nil {
			return nil, errors.New("parent asset not found")
		}
	}

	status := domain.AssetStatusPending
	if req.Status != "" {
		status = req.Status
	}

	asset := &domain.MediaAsset{
		Type:       req.Type,
		SourceType: req.SourceType,
		SourceID:   req.SourceID,
		ParentID:   req.ParentID,
		Name:       req.Name,
		Path:       req.Path,
		Duration:   req.Duration,
		Size:       req.Size,
		Format:     req.Format,
		Status:     status,
	}

	if err := s.repo.CreateMediaAsset(ctx, asset); err != nil {
		return nil, err
	}

	return asset, nil
}

// Get 获取媒体资产
func (s *MediaAssetService) Get(ctx context.Context, id uuid.UUID) (*domain.MediaAsset, error) {
	asset, err := s.repo.GetMediaAsset(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("asset not found")
		}
		return nil, err
	}
	return asset, nil
}

// List 列出媒体资产
func (s *MediaAssetService) List(ctx context.Context, req *ListMediaAssetsRequest) ([]*domain.MediaAsset, int64, error) {
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 1000 {
		req.Limit = 1000
	}

	filter := domain.MediaAssetFilter{
		Type:       req.Type,
		SourceType: req.SourceType,
		SourceID:   req.SourceID,
		ParentID:   req.ParentID,
		Status:     req.Status,
		Tags:       req.Tags,
		From:       req.From,
		To:         req.To,
		Limit:      req.Limit,
		Offset:     req.Offset,
	}

	return s.repo.ListMediaAssets(ctx, filter)
}

// Update 更新媒体资产
func (s *MediaAssetService) Update(ctx context.Context, id uuid.UUID, req *UpdateMediaAssetRequest) (*domain.MediaAsset, error) {
	asset, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		asset.Name = *req.Name
	}
	if req.Status != nil {
		asset.Status = *req.Status
	}

	if err := s.repo.UpdateMediaAsset(ctx, asset); err != nil {
		return nil, err
	}

	return asset, nil
}

// Delete 删除媒体资产
func (s *MediaAssetService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}

	children, err := s.repo.ListMediaAssetsByParent(ctx, id)
	if err != nil {
		return err
	}
	if len(children) > 0 {
		return errors.New("cannot delete asset with children")
	}

	return s.repo.DeleteMediaAsset(ctx, id)
}

// ListBySource 列出指定媒体源的所有资产
func (s *MediaAssetService) ListBySource(ctx context.Context, sourceID uuid.UUID) ([]*domain.MediaAsset, error) {
	return s.repo.ListMediaAssetsBySource(ctx, sourceID)
}

// ListChildren 列出子资产（派生资产）
func (s *MediaAssetService) ListChildren(ctx context.Context, parentID uuid.UUID) ([]*domain.MediaAsset, error) {
	if _, err := s.Get(ctx, parentID); err != nil {
		return nil, err
	}
	return s.repo.ListMediaAssetsByParent(ctx, parentID)
}
