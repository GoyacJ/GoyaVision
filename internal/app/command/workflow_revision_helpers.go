package command

import (
	"context"
	"fmt"
	"time"

	"goyavision/internal/app/port"
	"goyavision/internal/domain/workflow"
	"goyavision/pkg/apperr"

	"github.com/google/uuid"
)

func persistAndActivateWorkflowRevision(
	ctx context.Context,
	repos *port.Repositories,
	wf *workflow.Workflow,
	revision int64,
) (*workflow.WorkflowRevision, error) {
	if wf == nil {
		return nil, apperr.InvalidInput("workflow is required")
	}
	if revision <= 0 {
		return nil, apperr.InvalidInput("revision must be greater than 0")
	}

	rev := &workflow.WorkflowRevision{
		ID:         uuid.New(),
		WorkflowID: wf.ID,
		Revision:   revision,
		Status:     workflow.RevisionStatusActive,
		Definition: workflow.BuildDefinitionFromWorkflow(wf),
		CreatedAt:  time.Now().UTC(),
		UpdatedAt:  time.Now().UTC(),
	}

	if err := repos.Workflows.CreateRevision(ctx, rev); err != nil {
		return nil, apperr.Wrap(err, apperr.CodeDBError, "failed to create workflow revision")
	}
	if err := repos.Workflows.ActivateRevision(ctx, wf.ID, rev.ID); err != nil {
		return nil, apperr.Wrap(err, apperr.CodeDBError, "failed to activate workflow revision")
	}

	wf.CurrentRevisionID = &rev.ID
	wf.CurrentRevision = rev.Revision
	return rev, nil
}

func nextWorkflowRevision(current int64) int64 {
	if current < 0 {
		return 1
	}
	return current + 1
}

func ensureWorkflowHasActiveRevision(wf *workflow.Workflow) error {
	if wf == nil {
		return apperr.InvalidInput("workflow is required")
	}
	if wf.CurrentRevisionID == nil || wf.CurrentRevision == 0 {
		return apperr.InvalidInput(fmt.Sprintf("workflow %s has no active revision", wf.ID.String()))
	}
	return nil
}
