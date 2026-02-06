package engine

import (
	"context"
	"fmt"

	"goyavision/internal/domain/operator"
	"goyavision/internal/port"
)

// RoutingOperatorExecutor 基于执行模式路由到具体执行器
// 用于适配仅依赖 Execute(ctx, version, input) 的工作流引擎。
type RoutingOperatorExecutor struct {
	registry port.ExecutorRegistry
}

func NewRoutingOperatorExecutor(registry port.ExecutorRegistry) *RoutingOperatorExecutor {
	return &RoutingOperatorExecutor{registry: registry}
}

func (e *RoutingOperatorExecutor) Execute(ctx context.Context, version *operator.OperatorVersion, input *operator.Input) (*operator.Output, error) {
	if version == nil {
		return nil, fmt.Errorf("operator version is nil")
	}
	executor, err := e.registry.Get(version.ExecMode)
	if err != nil {
		return nil, err
	}
	return executor.Execute(ctx, version, input)
}
