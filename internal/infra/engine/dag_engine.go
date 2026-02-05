package engine

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"goyavision/internal/app/port"
	"goyavision/internal/domain/operator"
	"goyavision/internal/domain/workflow"

	"github.com/google/uuid"
)

var _ workflow.Engine = (*DAGWorkflowEngine)(nil)

// DAGWorkflowEngine implements parallel DAG execution with topological sorting
type DAGWorkflowEngine struct {
	uow      port.UnitOfWork
	executor workflow.OperatorExecutor
	tasks    map[uuid.UUID]*taskExecution
	mu       sync.RWMutex
}

type taskExecution struct {
	ctx         context.Context
	cancel      context.CancelFunc
	progress    int
	currentNode string
	nodeResults map[string]*operator.Output
	mu          sync.RWMutex
}

// NewDAGWorkflowEngine creates a new DAG workflow engine
func NewDAGWorkflowEngine(uow port.UnitOfWork, executor workflow.OperatorExecutor) *DAGWorkflowEngine {
	return &DAGWorkflowEngine{
		uow:      uow,
		executor: executor,
		tasks:    make(map[uuid.UUID]*taskExecution),
	}
}

// Execute executes a workflow using DAG topology with parallel execution
func (e *DAGWorkflowEngine) Execute(ctx context.Context, wf *workflow.Workflow, task *workflow.Task) error {
	if len(wf.Nodes) == 0 {
		return errors.New("workflow has no nodes")
	}

	// Build execution layers from DAG
	layers, err := e.buildExecutionLayers(wf.Nodes, wf.Edges)
	if err != nil {
		return fmt.Errorf("failed to build execution layers: %w", err)
	}

	// Create node map for quick lookup
	nodeMap := make(map[string]*workflow.Node)
	for i := range wf.Nodes {
		nodeMap[wf.Nodes[i].NodeKey] = &wf.Nodes[i]
	}

	// Setup execution context
	execCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	exec := &taskExecution{
		ctx:         execCtx,
		cancel:      cancel,
		progress:    0,
		nodeResults: make(map[string]*operator.Output),
	}

	e.mu.Lock()
	e.tasks[task.ID] = exec
	e.mu.Unlock()

	defer func() {
		e.mu.Lock()
		delete(e.tasks, task.ID)
		e.mu.Unlock()
	}()

	// Update task status to running
	err = e.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		now := time.Now()
		task.Status = workflow.TaskStatusRunning
		task.StartedAt = &now
		task.Progress = 0
		return repos.Tasks.Update(ctx, task)
	})
	if err != nil {
		return fmt.Errorf("failed to update task status: %w", err)
	}

	// Execute layers sequentially, nodes within layer in parallel
	totalLayers := len(layers)
	for i, layer := range layers {
		select {
		case <-execCtx.Done():
			e.updateTaskStatus(ctx, task, workflow.TaskStatusCancelled, "execution cancelled")
			return execCtx.Err()
		default:
		}

		if err := e.executeLayer(execCtx, layer, nodeMap, task, exec); err != nil {
			e.updateTaskStatus(ctx, task, workflow.TaskStatusFailed, err.Error())
			return fmt.Errorf("layer %d execution failed: %w", i+1, err)
		}

		// Update progress
		progress := int((float64(i+1) / float64(totalLayers)) * 100)
		exec.mu.Lock()
		exec.progress = progress
		exec.mu.Unlock()

		e.updateTaskProgress(ctx, task, progress)
	}

	// Update task to success
	return e.updateTaskStatus(ctx, task, workflow.TaskStatusSuccess, "")
}

// Cancel cancels a running workflow execution
func (e *DAGWorkflowEngine) Cancel(ctx context.Context, taskID uuid.UUID) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	exec, ok := e.tasks[taskID]
	if !ok {
		return errors.New("task is not running")
	}

	exec.cancel()
	return nil
}

// GetProgress returns the current execution progress
func (e *DAGWorkflowEngine) GetProgress(ctx context.Context, taskID uuid.UUID) (int, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	exec, ok := e.tasks[taskID]
	if !ok {
		// Task not running, fetch from database
		var progress int
		err := e.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
			task, err := repos.Tasks.Get(ctx, taskID)
			if err != nil {
				return err
			}
			progress = task.Progress
			return nil
		})
		return progress, err
	}

	exec.mu.RLock()
	defer exec.mu.RUnlock()
	return exec.progress, nil
}

// topologicalSort returns nodes in execution order using Kahn's algorithm
// Returns error if cycle detected
func (e *DAGWorkflowEngine) topologicalSort(nodes []workflow.Node, edges []workflow.Edge) ([]string, error) {
	// Build adjacency list and in-degree map
	graph := make(map[string][]string)
	inDegree := make(map[string]int)

	// Initialize all nodes with in-degree 0
	for _, node := range nodes {
		inDegree[node.NodeKey] = 0
		graph[node.NodeKey] = []string{}
	}

	// Build graph from edges
	for _, edge := range edges {
		graph[edge.SourceKey] = append(graph[edge.SourceKey], edge.TargetKey)
		inDegree[edge.TargetKey]++
	}

	// Find all nodes with in-degree 0 (starting nodes)
	queue := []string{}
	for nodeKey, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, nodeKey)
		}
	}

	// Kahn's algorithm
	sorted := []string{}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		sorted = append(sorted, current)

		for _, neighbor := range graph[current] {
			inDegree[neighbor]--
			if inDegree[neighbor] == 0 {
				queue = append(queue, neighbor)
			}
		}
	}

	// Check for cycles
	if len(sorted) != len(nodes) {
		return nil, errors.New("workflow contains cycles")
	}

	return sorted, nil
}

// buildExecutionLayers groups nodes that can execute in parallel
// Nodes in the same layer have no dependencies on each other
func (e *DAGWorkflowEngine) buildExecutionLayers(nodes []workflow.Node, edges []workflow.Edge) ([][]string, error) {
	// Build adjacency list and in-degree map
	graph := make(map[string][]string)
	inDegree := make(map[string]int)

	// Initialize all nodes
	for _, node := range nodes {
		inDegree[node.NodeKey] = 0
		graph[node.NodeKey] = []string{}
	}

	// Build graph from edges
	for _, edge := range edges {
		graph[edge.SourceKey] = append(graph[edge.SourceKey], edge.TargetKey)
		inDegree[edge.TargetKey]++
	}

	layers := [][]string{}
	remaining := len(nodes)

	for remaining > 0 {
		// Find all nodes with in-degree 0 in this layer
		layer := []string{}
		for nodeKey, degree := range inDegree {
			if degree == 0 {
				layer = append(layer, nodeKey)
			}
		}

		if len(layer) == 0 {
			return nil, errors.New("workflow contains cycles")
		}

		layers = append(layers, layer)

		// Remove processed nodes and update in-degrees
		for _, nodeKey := range layer {
			delete(inDegree, nodeKey)
			for _, neighbor := range graph[nodeKey] {
				if _, exists := inDegree[neighbor]; exists {
					inDegree[neighbor]--
				}
			}
			remaining--
		}
	}

	return layers, nil
}

// executeLayer executes all nodes in a layer concurrently
func (e *DAGWorkflowEngine) executeLayer(
	ctx context.Context,
	layer []string,
	nodeMap map[string]*workflow.Node,
	task *workflow.Task,
	exec *taskExecution,
) error {
	if len(layer) == 0 {
		return nil
	}

	// For single node, execute directly
	if len(layer) == 1 {
		return e.executeNode(ctx, nodeMap[layer[0]], task, exec)
	}

	// Parallel execution for multiple nodes
	var wg sync.WaitGroup
	errChan := make(chan error, len(layer))

	for _, nodeKey := range layer {
		wg.Add(1)
		go func(nk string) {
			defer wg.Done()
			if err := e.executeNode(ctx, nodeMap[nk], task, exec); err != nil {
				errChan <- fmt.Errorf("node %s: %w", nk, err)
			}
		}(nodeKey)
	}

	wg.Wait()
	close(errChan)

	// Check for errors
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

// executeNode executes a single node (operator)
func (e *DAGWorkflowEngine) executeNode(
	ctx context.Context,
	node *workflow.Node,
	task *workflow.Task,
	exec *taskExecution,
) error {
	// Update current node
	exec.mu.Lock()
	exec.currentNode = node.NodeKey
	exec.mu.Unlock()

	// Update task current node in database
	e.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		task.CurrentNode = node.NodeKey
		return repos.Tasks.Update(ctx, task)
	})

	// Skip if no operator
	if node.OperatorID == nil {
		return nil
	}

	// Get operator using UnitOfWork
	var op *operator.Operator
	err := e.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		var err error
		op, err = repos.Operators.Get(ctx, *node.OperatorID)
		return err
	})
	if err != nil {
		return fmt.Errorf("failed to get operator: %w", err)
	}

	// Prepare input (merge task input + node config + previous outputs)
	input := e.prepareNodeInput(task, node, exec)

	// Apply timeout if configured
	nodeCtx := ctx
	if node.Config != nil && node.Config.TimeoutSeconds > 0 {
		var cancel context.CancelFunc
		nodeCtx, cancel = context.WithTimeout(ctx, time.Duration(node.Config.TimeoutSeconds)*time.Second)
		defer cancel()
	}

	// Execute operator with retry logic
	var output *operator.Output
	retryCount := 1
	if node.Config != nil && node.Config.RetryCount > 0 {
		retryCount = node.Config.RetryCount + 1
	}

	var lastErr error
	for attempt := 0; attempt < retryCount; attempt++ {
		output, lastErr = e.executor.Execute(nodeCtx, op, input)
		if lastErr == nil {
			break
		}
		if attempt < retryCount-1 {
			// Wait before retry (exponential backoff)
			time.Sleep(time.Duration(1<<uint(attempt)) * time.Second)
		}
	}

	if lastErr != nil {
		return fmt.Errorf("node %s failed after %d attempts: %w", node.NodeKey, retryCount, lastErr)
	}

	// Store output for downstream nodes
	exec.mu.Lock()
	exec.nodeResults[node.NodeKey] = output
	exec.mu.Unlock()

	// Save artifacts
	if err := e.saveArtifacts(ctx, task.ID, node.NodeKey, output); err != nil {
		return fmt.Errorf("failed to save artifacts: %w", err)
	}

	return nil
}

// prepareNodeInput prepares input for node execution
func (e *DAGWorkflowEngine) prepareNodeInput(
	task *workflow.Task,
	node *workflow.Node,
	exec *taskExecution,
) *operator.Input {
	input := &operator.Input{
		Params: make(map[string]interface{}),
	}

	// Add task asset if present
	if task.AssetID != nil {
		input.AssetID = *task.AssetID
	}

	// Add task input params
	if task.InputParams != nil {
		for k, v := range task.InputParams {
			input.Params[k] = v
		}
	}

	// Add node config params (override task params)
	if node.Config != nil && node.Config.Params != nil {
		for k, v := range node.Config.Params {
			input.Params[k] = v
		}
	}

	// Add outputs from upstream nodes (for data flow)
	exec.mu.RLock()
	for nodeKey, output := range exec.nodeResults {
		// Store full output under node key
		input.Params[nodeKey+"_output"] = output

		// Also flatten output assets for easier access
		if len(output.OutputAssets) > 0 {
			input.Params[nodeKey+"_assets"] = output.OutputAssets
		}
		if len(output.Results) > 0 {
			input.Params[nodeKey+"_results"] = output.Results
		}
		if len(output.Timeline) > 0 {
			input.Params[nodeKey+"_timeline"] = output.Timeline
		}
	}
	exec.mu.RUnlock()

	return input
}

// saveArtifacts saves operator output as artifacts
func (e *DAGWorkflowEngine) saveArtifacts(
	ctx context.Context,
	taskID uuid.UUID,
	nodeKey string,
	output *operator.Output,
) error {
	if output == nil {
		return nil
	}

	return e.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		// Save output assets
		for _, asset := range output.OutputAssets {
			data := &workflow.ArtifactData{
				AssetInfo: &workflow.AssetInfo{
					Type:     string(asset.Type),
					Path:     asset.Path,
					Format:   asset.Format,
					Metadata: asset.Metadata,
				},
				Metadata: map[string]interface{}{
					"node_key": nodeKey,
				},
			}

			artifact := &workflow.Artifact{
				ID:     uuid.New(),
				TaskID: taskID,
				Type:   workflow.ArtifactTypeAsset,
				Data:   data,
			}

			if err := repos.Artifacts.Create(ctx, artifact); err != nil {
				return fmt.Errorf("failed to create asset artifact: %w", err)
			}
		}

		// Save analysis results
		if len(output.Results) > 0 {
			resultData := make([]map[string]interface{}, len(output.Results))
			for i, result := range output.Results {
				resultData[i] = map[string]interface{}{
					"type":       result.Type,
					"data":       result.Data,
					"confidence": result.Confidence,
				}
			}

			data := &workflow.ArtifactData{
				Results: resultData,
				Metadata: map[string]interface{}{
					"node_key": nodeKey,
				},
			}

			artifact := &workflow.Artifact{
				ID:     uuid.New(),
				TaskID: taskID,
				Type:   workflow.ArtifactTypeResult,
				Data:   data,
			}

			if err := repos.Artifacts.Create(ctx, artifact); err != nil {
				return fmt.Errorf("failed to create result artifact: %w", err)
			}
		}

		// Save timeline events
		if len(output.Timeline) > 0 {
			timelineData := make([]workflow.TimelineSegment, len(output.Timeline))
			for i, event := range output.Timeline {
				timelineData[i] = workflow.TimelineSegment{
					Start:      event.Start,
					End:        event.End,
					EventType:  event.EventType,
					Confidence: event.Confidence,
					Data:       event.Data,
				}
			}

			data := &workflow.ArtifactData{
				Timeline: timelineData,
				Metadata: map[string]interface{}{
					"node_key": nodeKey,
				},
			}

			artifact := &workflow.Artifact{
				ID:     uuid.New(),
				TaskID: taskID,
				Type:   workflow.ArtifactTypeTimeline,
				Data:   data,
			}

			if err := repos.Artifacts.Create(ctx, artifact); err != nil {
				return fmt.Errorf("failed to create timeline artifact: %w", err)
			}
		}

		// Save diagnostics if present
		if len(output.Diagnostics) > 0 {
			data := &workflow.ArtifactData{
				Diagnostics: output.Diagnostics,
				Metadata: map[string]interface{}{
					"node_key": nodeKey,
				},
			}

			artifact := &workflow.Artifact{
				ID:     uuid.New(),
				TaskID: taskID,
				Type:   workflow.ArtifactTypeReport,
				Data:   data,
			}

			if err := repos.Artifacts.Create(ctx, artifact); err != nil {
				return fmt.Errorf("failed to create diagnostics artifact: %w", err)
			}
		}

		return nil
	})
}

// updateTaskProgress updates task progress
func (e *DAGWorkflowEngine) updateTaskProgress(ctx context.Context, task *workflow.Task, progress int) error {
	return e.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		task.Progress = progress
		return repos.Tasks.Update(ctx, task)
	})
}

// updateTaskStatus updates task status and completion
func (e *DAGWorkflowEngine) updateTaskStatus(
	ctx context.Context,
	task *workflow.Task,
	status workflow.TaskStatus,
	errorMsg string,
) error {
	return e.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		task.Status = status
		task.Progress = 100
		now := time.Now()
		task.CompletedAt = &now
		if errorMsg != "" {
			task.Error = errorMsg
		}
		return repos.Tasks.Update(ctx, task)
	})
}
