package port

import (
	"context"

	"goyavision/internal/domain"

	"github.com/google/uuid"
)

// OperatorExecutor 算子执行器接口
type OperatorExecutor interface {
	// Execute 执行算子
	Execute(ctx context.Context, operator *domain.Operator, input *domain.OperatorInput) (*domain.OperatorOutput, error)
}

// WorkflowEngine 工作流引擎接口
type WorkflowEngine interface {
	// Execute 执行工作流
	Execute(ctx context.Context, workflow *domain.Workflow, task *domain.Task) error
	
	// Cancel 取消工作流执行
	Cancel(ctx context.Context, taskID uuid.UUID) error
	
	// GetProgress 获取工作流执行进度
	GetProgress(ctx context.Context, taskID uuid.UUID) (int, error)
}
