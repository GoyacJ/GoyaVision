package port

import (
	"context"

	"goyavision/internal/domain/operator"
	"goyavision/internal/domain/workflow"

	"github.com/google/uuid"
)

// OperatorExecutor 算子执行器接口
type OperatorExecutor interface {
	// Execute 执行算子
	Execute(ctx context.Context, version *operator.OperatorVersion, input *operator.Input) (*operator.Output, error)

	// Mode 返回当前执行器支持的执行模式
	Mode() operator.ExecMode

	// HealthCheck 执行器健康检查
	HealthCheck(ctx context.Context, version *operator.OperatorVersion) error
}

// ExecutorRegistry 执行器注册表
type ExecutorRegistry interface {
	Register(mode operator.ExecMode, executor OperatorExecutor)
	Get(mode operator.ExecMode) (OperatorExecutor, error)
}

// WorkflowEngine 工作流引擎接口
type WorkflowEngine interface {
	// Execute 执行工作流
	Execute(ctx context.Context, workflow *workflow.Workflow, task *workflow.Task) error

	// Cancel 取消工作流执行
	Cancel(ctx context.Context, taskID uuid.UUID) error

	// GetProgress 获取工作流执行进度
	GetProgress(ctx context.Context, taskID uuid.UUID) (int, error)
}
