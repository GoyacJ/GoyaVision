package command

import (
	"context"
	"errors"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/operator"
	"goyavision/pkg/apperr"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type InstallTemplateHandler struct {
	uow port.UnitOfWork
}

func NewInstallTemplateHandler(uow port.UnitOfWork) *InstallTemplateHandler {
	return &InstallTemplateHandler{uow: uow}
}

func (h *InstallTemplateHandler) Handle(ctx context.Context, cmd dto.InstallTemplateCommand) (*operator.Operator, error) {
	if cmd.TemplateID == uuid.Nil {
		return nil, apperr.InvalidInput("template_id is required")
	}
	if cmd.OperatorCode == "" || cmd.OperatorName == "" {
		return nil, apperr.InvalidInput("operator_code and operator_name are required")
	}

	var result *operator.Operator
	err := h.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		tpl, err := repos.OperatorTemplates.Get(ctx, cmd.TemplateID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperr.NotFound("operator template", cmd.TemplateID.String())
			}
			return apperr.Wrap(err, apperr.CodeDBError, "failed to get template")
		}

		if _, err := repos.Operators.GetByCode(ctx, cmd.OperatorCode); err == nil {
			return apperr.Conflict("operator code already exists")
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to check operator code uniqueness")
		}

		tags := cmd.Tags
		if len(tags) == 0 {
			tags = tpl.Tags
		}

		op := &operator.Operator{
			ID:          uuid.New(),
			Code:        cmd.OperatorCode,
			Name:        cmd.OperatorName,
			Description: tpl.Description,
			Category:    tpl.Category,
			Type:        tpl.Type,
			Origin:      operator.OriginMarketplace,
			Status:      operator.StatusDraft,
			Tags:        tags,
		}
		if err := repos.Operators.Create(ctx, op); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to create operator from template")
		}

		version := &operator.OperatorVersion{
			ID:          uuid.New(),
			OperatorID:  op.ID,
			Version:     "1.0.0",
			ExecMode:    tpl.ExecMode,
			ExecConfig:  tpl.ExecConfig,
			InputSchema: tpl.InputSchema,
			OutputSpec:  tpl.OutputSpec,
			Config:      tpl.Config,
			Status:      operator.VersionStatusActive,
		}
		if err := repos.OperatorVersions.Create(ctx, version); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to create initial version from template")
		}

		op.ActiveVersionID = &version.ID
		op.ActiveVersion = version
		syncOperatorCompatFieldsFromVersion(op, version)
		if err := repos.Operators.Update(ctx, op); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to bind active version")
		}

		if err := repos.OperatorTemplates.IncrementDownloads(ctx, tpl.ID); err != nil {
			return apperr.Wrap(err, apperr.CodeDBError, "failed to increase template downloads")
		}

		result = op
		return nil
	})

	return result, err
}
