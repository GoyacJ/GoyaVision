package engine

import (
	"fmt"
	"sync"

	"goyavision/internal/domain/operator"
	"goyavision/internal/port"
)

var _ port.ExecutorRegistry = (*ExecutorRegistry)(nil)

// ExecutorRegistry 执行器注册表
type ExecutorRegistry struct {
	mu        sync.RWMutex
	executors map[operator.ExecMode]port.OperatorExecutor
}

func NewExecutorRegistry() *ExecutorRegistry {
	return &ExecutorRegistry{executors: make(map[operator.ExecMode]port.OperatorExecutor)}
}

func (r *ExecutorRegistry) Register(mode operator.ExecMode, executor port.OperatorExecutor) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.executors[mode] = executor
}

func (r *ExecutorRegistry) Get(mode operator.ExecMode) (port.OperatorExecutor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	executor, ok := r.executors[mode]
	if !ok {
		return nil, fmt.Errorf("executor for mode %s not found", mode)
	}
	return executor, nil
}
