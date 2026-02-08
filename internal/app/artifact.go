package app

import (
	"context"
	"errors"
	"time"

	"goyavision/internal/domain/workflow"
	"goyavision/internal/port"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateArtifactRequest 创建产物请求
type CreateArtifactRequest struct {
	TaskID  uuid.UUID              `json:"task_id"`
	Type    workflow.ArtifactType  `json:"type"`
	AssetID *uuid.UUID             `json:"asset_id,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

// ListArtifactsRequest 列出产物请求
type ListArtifactsRequest struct {
	TaskID  *uuid.UUID
	NodeKey *string
	Type    *workflow.ArtifactType
	AssetID *uuid.UUID
	From    *time.Time
	To      *time.Time
	Limit   int
	Offset  int
}

type ArtifactService struct {
	repo port.Repository
}

func NewArtifactService(repo port.Repository) *ArtifactService {
	return &ArtifactService{
		repo: repo,
	}
}

// Create 创建产物
func (s *ArtifactService) Create(ctx context.Context, req *CreateArtifactRequest) (*workflow.Artifact, error) {
	if req.TaskID == uuid.Nil {
		return nil, errors.New("task_id is required")
	}
	if req.Type == "" {
		return nil, errors.New("type is required")
	}

	if req.Type != workflow.ArtifactTypeAsset &&
		req.Type != workflow.ArtifactTypeResult &&
		req.Type != workflow.ArtifactTypeTimeline &&
		req.Type != workflow.ArtifactTypeReport {
		return nil, errors.New("invalid artifact type")
	}

	if _, err := s.repo.GetTask(ctx, req.TaskID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	if req.AssetID != nil {
		if _, err := s.repo.GetMediaAsset(ctx, *req.AssetID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("asset not found")
			}
			return nil, err
		}
	}

	artifact := &workflow.Artifact{
		TaskID:  req.TaskID,
		Type:    req.Type,
		AssetID: req.AssetID,
	}

	if err := s.repo.CreateArtifact(ctx, artifact); err != nil {
		return nil, err
	}

	return s.repo.GetArtifact(ctx, artifact.ID)
}

// Get 获取产物
func (s *ArtifactService) Get(ctx context.Context, id uuid.UUID) (*workflow.Artifact, error) {
	artifact, err := s.repo.GetArtifact(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("artifact not found")
		}
		return nil, err
	}
	return artifact, nil
}

// List 列出产物
func (s *ArtifactService) List(ctx context.Context, req *ListArtifactsRequest) ([]*workflow.Artifact, int64, error) {
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 1000 {
		req.Limit = 1000
	}

	filter := workflow.ArtifactFilter{
		TaskID:  req.TaskID,
		NodeKey: req.NodeKey,
		Type:    req.Type,
		AssetID: req.AssetID,
		From:    req.From,
		To:      req.To,
		Limit:   req.Limit,
		Offset:  req.Offset,
	}

	return s.repo.ListArtifacts(ctx, filter)
}

// Delete 删除产物
func (s *ArtifactService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}

	return s.repo.DeleteArtifact(ctx, id)
}

// ListByTask 列出指定任务的所有产物
func (s *ArtifactService) ListByTask(ctx context.Context, taskID uuid.UUID) ([]*workflow.Artifact, error) {
	if _, err := s.repo.GetTask(ctx, taskID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	return s.repo.ListArtifactsByTask(ctx, taskID)
}

// ListByType 列出指定任务的指定类型产物
func (s *ArtifactService) ListByType(ctx context.Context, taskID uuid.UUID, artifactType workflow.ArtifactType) ([]*workflow.Artifact, error) {
	if _, err := s.repo.GetTask(ctx, taskID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("task not found")
		}
		return nil, err
	}

	return s.repo.ListArtifactsByType(ctx, taskID, artifactType)
}
