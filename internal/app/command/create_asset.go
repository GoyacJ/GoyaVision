package command

import (
	"context"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/media"
	"goyavision/pkg/apperr"
)

type CreateAssetHandler struct {
	uow port.UnitOfWork
}

func NewCreateAssetHandler(uow port.UnitOfWork) *CreateAssetHandler {
	return &CreateAssetHandler{uow: uow}
}

func (h *CreateAssetHandler) Handle(ctx context.Context, cmd dto.CreateAssetCommand) (*media.Asset, error) {
	if cmd.Type == "" {
		return nil, apperr.InvalidInput("type is required")
	}
	if cmd.SourceType == "" {
		return nil, apperr.InvalidInput("source_type is required")
	}
	if cmd.Name == "" {
		return nil, apperr.InvalidInput("name is required")
	}
	if cmd.Path == "" {
		return nil, apperr.InvalidInput("path is required")
	}

	if cmd.Type != media.AssetTypeVideo && cmd.Type != media.AssetTypeImage &&
		cmd.Type != media.AssetTypeAudio && cmd.Type != media.AssetTypeStream {
		return nil, apperr.InvalidInput("invalid asset type")
	}

	if cmd.SourceType != media.AssetSourceLive && cmd.SourceType != media.AssetSourceVOD &&
		cmd.SourceType != media.AssetSourceUpload && cmd.SourceType != media.AssetSourceGenerated {
		return nil, apperr.InvalidInput("invalid source type")
	}

	var asset *media.Asset
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if cmd.ParentID != nil {
			if _, err := repos.Assets.Get(ctx, *cmd.ParentID); err != nil {
				return apperr.NotFound("parent asset", *cmd.ParentID)
			}
		}

		status := media.AssetStatusPending
		if cmd.Status != "" {
			status = cmd.Status
		}

		asset = &media.Asset{
			Type:       cmd.Type,
			SourceType: cmd.SourceType,
			SourceID:   cmd.SourceID,
			ParentID:   cmd.ParentID,
			Name:       cmd.Name,
			Path:       cmd.Path,
			Duration:   cmd.Duration,
			Size:       cmd.Size,
			Format:     cmd.Format,
			Metadata:   cmd.Metadata,
			Status:     status,
			Tags:       cmd.Tags,
		}

		return repos.Assets.Create(ctx, asset)
	})

	if err != nil {
		return nil, err
	}

	return asset, nil
}
