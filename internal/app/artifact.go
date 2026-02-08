package app

import (
	"context"
	"errors"
	"time"

	"goyavision/internal/domain/workflow"
	"goyavision/internal/port"
	"goyavision/pkg/apperr"

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
		return nil, apperr.InvalidInput("task_id is required")
	}
	if req.Type == "" {
		return nil, apperr.InvalidInput("type is required")
	}

	if req.Type != workflow.ArtifactTypeAsset &&
		req.Type != workflow.ArtifactTypeResult &&
		req.Type != workflow.ArtifactTypeTimeline &&
		req.Type != workflow.ArtifactTypeReport {
		return nil, apperr.InvalidInput("invalid artifact type")
	}

	if _, err := s.repo.GetTask(ctx, req.TaskID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("task", req.TaskID)
		}
		return nil, apperr.Internal("get task", err)
	}

	if req.AssetID != nil {
		if _, err := s.repo.GetMediaAsset(ctx, *req.AssetID); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, apperr.NotFound("asset", *req.AssetID)
			}
			return nil, apperr.Internal("get asset", err)
		}
	}

	artifact := &workflow.Artifact{
		TaskID:  req.TaskID,
		Type:    req.Type,
		AssetID: req.AssetID,
	}

	if err := s.repo.CreateArtifact(ctx, artifact); err != nil {
		return nil, apperr.Internal("create artifact", err)
	}

	out, err := s.repo.GetArtifact(ctx, artifact.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("artifact", artifact.ID)
		}
		return nil, apperr.Internal("get artifact", err)
	}
	return out, nil
}

// Get 获取产物
func (s *ArtifactService) Get(ctx context.Context, id uuid.UUID) (*workflow.Artifact, error) {
	artifact, err := s.repo.GetArtifact(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("artifact", id)
		}
		return nil, apperr.Internal("get artifact", err)
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
	list, total, err := s.repo.ListArtifacts(ctx, filter)
	if err != nil {
		return nil, 0, apperr.Internal("list artifacts", err)
	}
	return list, total, nil
}

// Delete 删除产物
func (s *ArtifactService) Delete(ctx context.Context, id uuid.UUID) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	if err := s.repo.DeleteArtifact(ctx, id); err != nil {
		return apperr.Internal("delete artifact", err)
	}
	return nil
}

// ListByTask 列出指定任务的所有产物
func (s *ArtifactService) ListByTask(ctx context.Context, taskID uuid.UUID) ([]*workflow.Artifact, error) {
	if _, err := s.repo.GetTask(ctx, taskID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("task", taskID)
		}
		return nil, apperr.Internal("get task", err)
	}
	list, err := s.repo.ListArtifactsByTask(ctx, taskID)
	if err != nil {
		return nil, apperr.Internal("list artifacts by task", err)
	}
	return list, nil
}

// ListByType 列出指定任务的指定类型产物
func (s *ArtifactService) ListByType(ctx context.Context, taskID uuid.UUID, artifactType workflow.ArtifactType) ([]*workflow.Artifact, error) {
	if _, err := s.repo.GetTask(ctx, taskID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperr.NotFound("task", taskID)
		}
		return nil, apperr.Internal("get task", err)
	}
	list, err := s.repo.ListArtifactsByType(ctx, taskID, artifactType)
	if err != nil {
		return nil, apperr.Internal("list artifacts by type", err)
	}
	return list, nil
}
