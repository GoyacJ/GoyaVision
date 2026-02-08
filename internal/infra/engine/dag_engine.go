package engine

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"goyavision/internal/app/port"
	"goyavision/internal/domain/agent"
	"goyavision/internal/domain/algorithm"
	"goyavision/internal/domain/operator"
	"goyavision/internal/domain/workflow"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

var _ workflow.Engine = (*DAGWorkflowEngine)(nil)

// DAGWorkflowEngine implements parallel DAG execution with topological sorting
type DAGWorkflowEngine struct {
	uow             port.UnitOfWork
	executor        workflow.OperatorExecutor
	schemaValidator port.SchemaValidator
	tasks           map[uuid.UUID]*taskExecution
	mu              sync.RWMutex
}

type taskExecution struct {
	ctx            context.Context
	cancel         context.CancelFunc
	progress       int
	currentNode    string
	contextVersion int64
	nodeResults    map[string]*operator.Output
	nodeExecutions map[string]*workflow.NodeExecution
	mu             sync.RWMutex
}

// NewDAGWorkflowEngine creates a new DAG workflow engine
func NewDAGWorkflowEngine(uow port.UnitOfWork, executor workflow.OperatorExecutor, validators ...port.SchemaValidator) *DAGWorkflowEngine {
	var schemaValidator port.SchemaValidator
	if len(validators) > 0 {
		schemaValidator = validators[0]
	}

	return &DAGWorkflowEngine{
		uow:             uow,
		executor:        executor,
		schemaValidator: schemaValidator,
		tasks:           make(map[uuid.UUID]*taskExecution),
	}
}

// Execute executes a workflow using DAG topology with parallel execution
func (e *DAGWorkflowEngine) Execute(ctx context.Context, wf *workflow.Workflow, task *workflow.Task) error {
	execWF, err := e.resolveWorkflowForTask(ctx, wf, task)
	if err != nil {
		return fmt.Errorf("failed to resolve workflow revision: %w", err)
	}

	if len(execWF.Nodes) == 0 {
		return errors.New("workflow has no nodes")
	}

	// Build execution layers from DAG
	layers, err := e.buildExecutionLayers(execWF.Nodes, execWF.Edges)
	if err != nil {
		return fmt.Errorf("failed to build execution layers: %w", err)
	}

	// Create node map for quick lookup
	nodeMap := make(map[string]*workflow.Node)
	for i := range execWF.Nodes {
		nodeMap[execWF.Nodes[i].NodeKey] = &execWF.Nodes[i]
	}

	// Setup execution context
	execCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	exec := &taskExecution{
		ctx:            execCtx,
		cancel:         cancel,
		progress:       0,
		nodeResults:    make(map[string]*operator.Output),
		nodeExecutions: make(map[string]*workflow.NodeExecution),
	}

	// Initialize node executions
	for _, node := range execWF.Nodes {
		exec.nodeExecutions[node.NodeKey] = &workflow.NodeExecution{
			NodeKey: node.NodeKey,
			Status:  workflow.NodeExecPending,
		}
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

	if err := e.initializeTaskContextState(ctx, execWF, task, exec); err != nil {
		e.updateTaskStatus(ctx, task, workflow.TaskStatusFailed, err.Error())
		return fmt.Errorf("failed to initialize task context: %w", err)
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

		if err := e.executeLayer(execCtx, layer, nodeMap, execWF, execWF.Edges, task, exec); err != nil {
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
	wf *workflow.Workflow,
	edges []workflow.Edge,
	task *workflow.Task,
	exec *taskExecution,
) error {
	if len(layer) == 0 {
		return nil
	}

	// Check conditions for all nodes in layer first
	nodesToExecute := []string{}
	for _, nodeKey := range layer {
		node := nodeMap[nodeKey]
		if e.shouldExecuteNode(node, edges, exec) {
			nodesToExecute = append(nodesToExecute, nodeKey)
		} else {
			// Mark as skipped
			exec.mu.Lock()
			if execNode, ok := exec.nodeExecutions[nodeKey]; ok {
				execNode.Status = workflow.NodeExecSkipped
			}
			exec.mu.Unlock()
		}
	}

	// Sync skipped status
	if len(nodesToExecute) < len(layer) {
		if err := e.syncTaskNodeExecutions(ctx, task, exec); err != nil {
			return err
		}
	}

	if len(nodesToExecute) == 0 {
		return nil
	}

	// For single node, execute directly
	if len(nodesToExecute) == 1 {
		return e.executeNode(ctx, nodeMap[nodesToExecute[0]], wf, task, exec)
	}

	// Parallel execution for multiple nodes
	var wg sync.WaitGroup
	errChan := make(chan error, len(nodesToExecute))

	for _, nodeKey := range nodesToExecute {
		wg.Add(1)
		go func(nk string) {
			defer wg.Done()
			if err := e.executeNode(ctx, nodeMap[nk], wf, task, exec); err != nil {
				errChan <- fmt.Errorf("node %s: %w", nk, err)
			}
		}(nodeKey)
	}

	wg.Wait()
	close(errChan)

	var errs []error
	for err := range errChan {
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}
	return nil
}

// executeNode executes a single node (operator)
func (e *DAGWorkflowEngine) executeNode(
	ctx context.Context,
	node *workflow.Node,
	wf *workflow.Workflow,
	task *workflow.Task,
	exec *taskExecution,
) error {
	// Update node execution status to running
	exec.mu.Lock()
	exec.currentNode = node.NodeKey
	if execNode, ok := exec.nodeExecutions[node.NodeKey]; ok {
		now := time.Now()
		execNode.Status = workflow.NodeExecRunning
		execNode.StartedAt = &now
	}
	exec.mu.Unlock()

	// Sync task state
	if err := e.syncTaskNodeExecutions(ctx, task, exec); err != nil {
		return err
	}
	e.emitRunEvent(ctx, &agent.RunEvent{
		TaskID:    task.ID,
		EventType: agent.EventTypeNodeStarted,
		Source:    "dag_engine",
		NodeKey:   node.NodeKey,
		Payload: map[string]interface{}{
			"status": "running",
		},
	})

	contextState, err := e.getTaskContextState(ctx, task.ID)
	if err != nil {
		return e.failNode(ctx, task, exec, node.NodeKey, fmt.Errorf("failed to get task context: %w", err))
	}

	versionToExecute, err := e.resolveNodeOperatorVersion(ctx, node)
	if err != nil {
		return e.failNode(ctx, task, exec, node.NodeKey, err)
	}

	// Skip if no executable operator
	if versionToExecute == nil {
		diff, err := e.buildNodeContextDiff(wf, node, nil, nil)
		if err != nil {
			return e.failNode(ctx, task, exec, node.NodeKey, err)
		}
		if err := e.applyNodeContextPatch(ctx, task, exec, node.NodeKey, diff); err != nil {
			return e.failNode(ctx, task, exec, node.NodeKey, fmt.Errorf("failed to persist node context diff: %w", err))
		}

		// Mark as success immediately
		exec.mu.Lock()
		if execNode, ok := exec.nodeExecutions[node.NodeKey]; ok {
			now := time.Now()
			execNode.Status = workflow.NodeExecSuccess
			execNode.CompletedAt = &now
		}
		exec.mu.Unlock()
		return e.syncTaskNodeExecutions(ctx, task, exec)
	}
	if err := e.enforceNodeToolPolicy(ctx, node, versionToExecute); err != nil {
		return e.failNode(ctx, task, exec, node.NodeKey, err)
	}

	// Prepare input (merge task input + node config + previous outputs)
	input := e.prepareNodeInput(task, node, exec, contextState)
	if err := e.validateNodeInput(ctx, versionToExecute, input); err != nil {
		return e.failNode(ctx, task, exec, node.NodeKey, err)
	}

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
		output, lastErr = e.executor.Execute(nodeCtx, versionToExecute, input)
		if lastErr == nil {
			break
		}
		if attempt < retryCount-1 {
			// Wait before retry (exponential backoff)
			time.Sleep(time.Duration(1<<uint(attempt)) * time.Second)
		}
	}

	if lastErr != nil {
		return e.failNode(ctx, task, exec, node.NodeKey, fmt.Errorf("node %s failed after %d attempts: %w", node.NodeKey, retryCount, lastErr))
	}

	if err := e.validateNodeOutput(nodeCtx, versionToExecute, output); err != nil {
		return e.failNode(ctx, task, exec, node.NodeKey, err)
	}

	// Store output for downstream nodes
	exec.mu.Lock()
	exec.nodeResults[node.NodeKey] = output
	exec.mu.Unlock()

	// Save artifacts
	var artifactIDs []uuid.UUID
	if output != nil {
		artifactIDs, err = e.saveArtifactsWithIDs(ctx, task.ID, node.NodeKey, output)
		if err != nil {
			return e.failNode(ctx, task, exec, node.NodeKey, fmt.Errorf("failed to save artifacts: %w", err))
		}
	}

	diff, err := e.buildNodeContextDiff(wf, node, output, artifactIDs)
	if err != nil {
		return e.failNode(ctx, task, exec, node.NodeKey, err)
	}
	if err := e.applyNodeContextPatch(ctx, task, exec, node.NodeKey, diff); err != nil {
		return e.failNode(ctx, task, exec, node.NodeKey, fmt.Errorf("failed to persist node context diff: %w", err))
	}

	// Mark as success
	exec.mu.Lock()
	if execNode, ok := exec.nodeExecutions[node.NodeKey]; ok {
		now := time.Now()
		execNode.Status = workflow.NodeExecSuccess
		execNode.CompletedAt = &now
		execNode.ArtifactIDs = artifactIDs
	}
	exec.mu.Unlock()
	e.emitRunEvent(ctx, &agent.RunEvent{
		TaskID:    task.ID,
		EventType: agent.EventTypeNodeSucceeded,
		Source:    "dag_engine",
		NodeKey:   node.NodeKey,
		Payload: map[string]interface{}{
			"artifact_ids":    artifactIDs,
			"context_version": task.ContextVersion,
		},
	})

	return e.syncTaskNodeExecutions(ctx, task, exec)
}

func (e *DAGWorkflowEngine) failNode(ctx context.Context, task *workflow.Task, exec *taskExecution, nodeKey string, err error) error {
	exec.mu.Lock()
	if execNode, ok := exec.nodeExecutions[nodeKey]; ok {
		now := time.Now()
		execNode.Status = workflow.NodeExecFailed
		execNode.Error = err.Error()
		execNode.CompletedAt = &now
	}
	exec.mu.Unlock()
	e.emitRunEvent(ctx, &agent.RunEvent{
		TaskID:    task.ID,
		EventType: agent.EventTypeNodeFailed,
		Source:    "dag_engine",
		NodeKey:   nodeKey,
		Payload: map[string]interface{}{
			"error": err.Error(),
		},
	})
	_ = e.syncTaskNodeExecutions(ctx, task, exec)
	return err
}

func (e *DAGWorkflowEngine) shouldExecuteNode(
	node *workflow.Node,
	edges []workflow.Edge,
	exec *taskExecution,
) bool {
	// Find all incoming edges
	incoming := []workflow.Edge{}
	for _, edge := range edges {
		if edge.TargetKey == node.NodeKey {
			incoming = append(incoming, edge)
		}
	}

	if len(incoming) == 0 {
		return true // No dependencies
	}

	// Check conditions
	for _, edge := range incoming {
		upstreamNodeKey := edge.SourceKey
		exec.mu.RLock()
		upstreamExec, ok := exec.nodeExecutions[upstreamNodeKey]
		exec.mu.RUnlock()

		if !ok || upstreamExec == nil {
			return false // Upstream not executed (should not happen in DAG if sorted correctly)
		}

		// Default condition: always execute if upstream success
		conditionType := "always"
		if edge.Condition != nil && edge.Condition.Type != "" {
			conditionType = edge.Condition.Type
		}

		switch conditionType {
		case "always":
			if upstreamExec.Status != workflow.NodeExecSuccess && upstreamExec.Status != workflow.NodeExecFailed {
				return false // Upstream skipped or not finished
			}
		case "on_success":
			if upstreamExec.Status != workflow.NodeExecSuccess {
				return false
			}
		case "on_failure":
			if upstreamExec.Status != workflow.NodeExecFailed {
				return false
			}
		}
	}

	return true
}

func (e *DAGWorkflowEngine) syncTaskNodeExecutions(ctx context.Context, task *workflow.Task, exec *taskExecution) error {
	exec.mu.RLock()
	executions := make([]workflow.NodeExecution, 0, len(exec.nodeExecutions))
	for _, v := range exec.nodeExecutions {
		executions = append(executions, *v)
	}
	currentNode := exec.currentNode
	contextVersion := exec.contextVersion
	exec.mu.RUnlock()

	return e.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		task.NodeExecutions = executions
		task.CurrentNode = currentNode
		task.ContextVersion = contextVersion
		return repos.Tasks.Update(ctx, task)
	})
}

func (e *DAGWorkflowEngine) validateNodeInput(ctx context.Context, version *operator.OperatorVersion, input *operator.Input) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if e.schemaValidator == nil || version == nil || len(version.InputSchema) == 0 {
		return nil
	}

	inputPayload := make(map[string]interface{})
	if input != nil {
		for k, v := range input.Params {
			inputPayload[k] = v
		}
		if input.AssetID != uuid.Nil {
			inputPayload["asset_id"] = input.AssetID.String()
		}
	}

	if err := e.schemaValidator.ValidateInput(ctx, version.InputSchema, inputPayload); err != nil {
		return fmt.Errorf("runtime input schema validation failed: %w", err)
	}

	return nil
}

func (e *DAGWorkflowEngine) validateNodeOutput(ctx context.Context, version *operator.OperatorVersion, output *operator.Output) error {
	if ctx == nil {
		ctx = context.Background()
	}
	if e.schemaValidator == nil || version == nil || len(version.OutputSpec) == 0 {
		return nil
	}

	if output == nil {
		output = &operator.Output{}
	}

	outputPayload := make(map[string]interface{})
	b, err := json.Marshal(output)
	if err != nil {
		return fmt.Errorf("marshal runtime output failed: %w", err)
	}
	if err := json.Unmarshal(b, &outputPayload); err != nil {
		return fmt.Errorf("unmarshal runtime output failed: %w", err)
	}

	if err := e.schemaValidator.ValidateOutput(ctx, version.OutputSpec, outputPayload); err != nil {
		return fmt.Errorf("runtime output schema validation failed: %w", err)
	}

	return nil
}

// prepareNodeInput prepares input for node execution
func (e *DAGWorkflowEngine) prepareNodeInput(
	task *workflow.Task,
	node *workflow.Node,
	exec *taskExecution,
	contextState *workflow.TaskContextState,
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

	// Attach full context for operators that can consume unified context.
	if contextState != nil && contextState.Data != nil {
		input.Params["context"] = contextState.Data
	}

	// Resolve mapped inputs from context first.
	if node.Config != nil && len(node.Config.InputMapping) > 0 && contextState != nil {
		for paramKey, contextPath := range node.Config.InputMapping {
			value, ok := resolveAnyByPath(contextState.Data, contextPath)
			if ok {
				input.Params[paramKey] = value
			}
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

// saveArtifactsWithIDs saves operator output as artifacts and returns their IDs
func (e *DAGWorkflowEngine) saveArtifactsWithIDs(
	ctx context.Context,
	taskID uuid.UUID,
	nodeKey string,
	output *operator.Output,
) ([]uuid.UUID, error) {
	if output == nil {
		return nil, nil
	}

	var artifactIDs []uuid.UUID

	err := e.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
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
			artifactIDs = append(artifactIDs, artifact.ID)
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
			artifactIDs = append(artifactIDs, artifact.ID)
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
			artifactIDs = append(artifactIDs, artifact.ID)
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
			artifactIDs = append(artifactIDs, artifact.ID)
		}

		return nil
	})

	return artifactIDs, err
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

func (e *DAGWorkflowEngine) initializeTaskContextState(
	ctx context.Context,
	wf *workflow.Workflow,
	task *workflow.Task,
	exec *taskExecution,
) error {
	return e.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if repos.Contexts == nil {
			return errors.New("context repository is not configured")
		}

		state := &workflow.TaskContextState{
			TaskID:  task.ID,
			Version: 1,
			Data:    buildInitialTaskContextData(wf, task),
		}
		if err := repos.Contexts.InitializeState(ctx, state); err != nil {
			return err
		}

		task.ContextVersion = state.Version
		exec.mu.Lock()
		exec.contextVersion = state.Version
		exec.mu.Unlock()
		return repos.Tasks.Update(ctx, task)
	})
}

func buildInitialTaskContextData(wf *workflow.Workflow, task *workflow.Task) map[string]interface{} {
	vars := map[string]interface{}{}
	if wf != nil && wf.ContextSpec != nil && wf.ContextSpec.Vars != nil {
		for key, spec := range wf.ContextSpec.Vars {
			if spec.Default != nil {
				vars[key] = spec.Default
			}
		}
	}
	for key, value := range task.InputParams {
		vars[key] = value
	}

	meta := map[string]interface{}{
		"task_id":        task.ID.String(),
		"workflow_id":    task.WorkflowID.String(),
		"initialized_at": time.Now().UTC().Format(time.RFC3339),
	}
	if wf != nil {
		meta["workflow_code"] = wf.Code
	}
	if task.AssetID != nil {
		meta["asset_id"] = task.AssetID.String()
	}

	return map[string]interface{}{
		"meta":      meta,
		"vars":      vars,
		"shared":    map[string]interface{}{},
		"nodes":     map[string]interface{}{},
		"artifacts": map[string]interface{}{},
	}
}

func (e *DAGWorkflowEngine) getTaskContextState(ctx context.Context, taskID uuid.UUID) (*workflow.TaskContextState, error) {
	var state *workflow.TaskContextState
	err := e.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if repos.Contexts == nil {
			return errors.New("context repository is not configured")
		}
		var err error
		state, err = repos.Contexts.GetState(ctx, taskID)
		return err
	})
	return state, err
}

func (e *DAGWorkflowEngine) applyNodeContextPatch(
	ctx context.Context,
	task *workflow.Task,
	exec *taskExecution,
	nodeKey string,
	diff workflow.ContextDiff,
) error {
	const maxAttempts = 8
	for attempt := 0; attempt < maxAttempts; attempt++ {
		err := e.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
			if repos.Contexts == nil {
				return errors.New("context repository is not configured")
			}
			state, err := repos.Contexts.GetState(ctx, task.ID)
			if err != nil {
				return err
			}

			patch := &workflow.TaskContextPatch{
				TaskID:        task.ID,
				WriterNodeKey: nodeKey,
				BeforeVersion: state.Version,
				Diff:          diff,
			}
			if err := repos.Contexts.ApplyPatch(ctx, patch); err != nil {
				return err
			}

			task.ContextVersion = patch.AfterVersion
			exec.mu.Lock()
			exec.contextVersion = patch.AfterVersion
			exec.mu.Unlock()
			return repos.Tasks.Update(ctx, task)
		})
		if err == nil {
			return nil
		}
		if errors.Is(err, workflow.ErrContextVersionConflict) {
			time.Sleep(time.Duration(1<<uint(attempt)) * 10 * time.Millisecond)
			continue
		}
		return err
	}
	return workflow.ErrContextVersionConflict
}

func (e *DAGWorkflowEngine) buildNodeContextDiff(
	wf *workflow.Workflow,
	node *workflow.Node,
	output *operator.Output,
	artifactIDs []uuid.UUID,
) (workflow.ContextDiff, error) {
	diff := workflow.ContextDiff{
		Set: map[string]interface{}{
			fmt.Sprintf("nodes.%s.status", node.NodeKey): workflow.NodeExecSuccess,
		},
	}

	if len(artifactIDs) > 0 {
		diff.Set[fmt.Sprintf("nodes.%s.artifact_ids", node.NodeKey)] = artifactIDs
	}

	outputPayload := map[string]interface{}{}
	if output != nil {
		var err error
		outputPayload, err = outputToMap(output)
		if err != nil {
			return workflow.ContextDiff{}, fmt.Errorf("failed to serialize node output: %w", err)
		}
		diff.Set[fmt.Sprintf("nodes.%s.output", node.NodeKey)] = outputPayload
	}

	if node.Config != nil && len(node.Config.OutputMapping) > 0 {
		for contextPath, outputSelector := range node.Config.OutputMapping {
			contextPath = strings.TrimSpace(contextPath)
			outputSelector = strings.TrimSpace(outputSelector)
			if contextPath == "" || outputSelector == "" {
				continue
			}
			value, ok := resolveAnyByPath(outputPayload, outputSelector)
			if !ok {
				return workflow.ContextDiff{}, fmt.Errorf("invalid output_mapping: selector %s not found", outputSelector)
			}
			diff.Set[contextPath] = value
		}
	}

	if err := validateContextWritePaths(wf, node.NodeKey, diff); err != nil {
		return workflow.ContextDiff{}, err
	}
	return diff, nil
}

func validateContextWritePaths(wf *workflow.Workflow, nodeKey string, diff workflow.ContextDiff) error {
	isLocalPath := func(path string) bool {
		return strings.HasPrefix(path, "nodes."+nodeKey+".")
	}

	checkPath := func(path string) error {
		if isLocalPath(path) {
			return nil
		}
		if strings.HasPrefix(path, "vars.") {
			return fmt.Errorf("forbidden context write path %s: writing vars.* is not allowed", path)
		}
		if wf == nil || wf.ContextSpec == nil || wf.ContextSpec.SharedKeys == nil {
			return fmt.Errorf("context write path %s is not allowed: shared key is not declared", path)
		}
		spec, ok := wf.ContextSpec.SharedKeys[path]
		if !ok {
			return fmt.Errorf("context write path %s is not declared in context_spec.shared_keys", path)
		}
		if !spec.CAS {
			return fmt.Errorf("shared key %s must enable cas", path)
		}
		if strings.TrimSpace(spec.ConflictPolicy) == "" {
			return fmt.Errorf("shared key %s must declare conflict policy", path)
		}
		return nil
	}

	for path := range diff.Set {
		if err := checkPath(path); err != nil {
			return err
		}
	}
	for _, path := range diff.Unset {
		if err := checkPath(path); err != nil {
			return err
		}
	}
	return nil
}

func outputToMap(output *operator.Output) (map[string]interface{}, error) {
	if output == nil {
		return map[string]interface{}{}, nil
	}
	b, err := json.Marshal(output)
	if err != nil {
		return nil, err
	}
	m := map[string]interface{}{}
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, err
	}
	return m, nil
}

func resolveAnyByPath(root interface{}, path string) (interface{}, bool) {
	path = strings.TrimSpace(path)
	if path == "" {
		return root, true
	}
	current := root
	for _, part := range strings.Split(path, ".") {
		switch typed := current.(type) {
		case map[string]interface{}:
			next, ok := typed[part]
			if !ok {
				return nil, false
			}
			current = next
		case []interface{}:
			index, err := strconv.Atoi(part)
			if err != nil || index < 0 || index >= len(typed) {
				return nil, false
			}
			current = typed[index]
		default:
			return nil, false
		}
	}
	return current, true
}

func (e *DAGWorkflowEngine) resolveNodeOperatorVersion(ctx context.Context, node *workflow.Node) (*operator.OperatorVersion, error) {
	if node == nil {
		return nil, nil
	}

	if node.OperatorID != nil {
		var op *operator.Operator
		err := e.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
			var err error
			op, err = repos.Operators.GetWithActiveVersion(ctx, *node.OperatorID)
			return err
		})
		if err != nil {
			return nil, fmt.Errorf("failed to get operator: %w", err)
		}
		if op == nil || op.ActiveVersion == nil {
			return nil, fmt.Errorf("operator has no active version")
		}
		return op.ActiveVersion, nil
	}

	if node.Config == nil || node.Config.AlgorithmRef == nil {
		return nil, nil
	}

	ref := node.Config.AlgorithmRef
	var selected *operator.OperatorVersion
	err := e.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if repos.Algorithms == nil {
			return fmt.Errorf("algorithm repository is not configured")
		}
		if repos.OperatorVersions == nil {
			return fmt.Errorf("operator version repository is not configured")
		}

		var algo *algorithm.Algorithm
		if ref.AlgorithmID != nil {
			var err error
			algo, err = repos.Algorithms.GetWithRelations(ctx, *ref.AlgorithmID)
			if err != nil {
				return err
			}
		} else if strings.TrimSpace(ref.AlgorithmCode) != "" {
			a, err := repos.Algorithms.GetByCode(ctx, ref.AlgorithmCode)
			if err != nil {
				return err
			}
			algo, err = repos.Algorithms.GetWithRelations(ctx, a.ID)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("algorithm_ref requires algorithm_id or algorithm_code")
		}

		version, err := resolveAlgorithmVersion(ctx, repos.Algorithms, algo, ref.Version)
		if err != nil {
			return err
		}
		impl := selectAlgorithmImplementation(version, ref)
		if impl == nil {
			return fmt.Errorf("algorithm %s revision %s has no operator_version implementation", algo.Code, version.Version)
		}

		opVersionID, err := uuid.Parse(strings.TrimSpace(impl.BindingRef))
		if err != nil {
			return fmt.Errorf("invalid implementation binding_ref: %w", err)
		}
		versionRef, err := repos.OperatorVersions.Get(ctx, opVersionID)
		if err != nil {
			return err
		}
		selected = versionRef
		return nil
	})
	if err != nil {
		return nil, err
	}
	if selected == nil {
		return nil, fmt.Errorf("no executable operator version resolved")
	}
	return selected, nil
}

func resolveAlgorithmVersion(
	ctx context.Context,
	repo algorithm.Repository,
	algo *algorithm.Algorithm,
	versionName string,
) (*algorithm.Version, error) {
	if algo == nil {
		return nil, fmt.Errorf("algorithm is not found")
	}
	if strings.TrimSpace(versionName) != "" {
		version, err := repo.GetVersionByName(ctx, algo.ID, strings.TrimSpace(versionName))
		if err != nil {
			return nil, err
		}
		return version, nil
	}

	versions, err := repo.ListVersions(ctx, algo.ID)
	if err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, fmt.Errorf("algorithm has no versions")
	}
	for i := range versions {
		if versions[i].Status == algorithm.VersionStatusPublished {
			return versions[i], nil
		}
	}
	return versions[0], nil
}

func selectAlgorithmImplementation(version *algorithm.Version, ref *workflow.AlgorithmRef) *algorithm.Implementation {
	if version == nil || len(version.Implementations) == 0 {
		return nil
	}

	candidates := make([]algorithm.Implementation, 0, len(version.Implementations))
	for i := range version.Implementations {
		if version.Implementations[i].Type == algorithm.ImplementationOperatorVersion {
			candidates = append(candidates, version.Implementations[i])
		}
	}
	if len(candidates) == 0 {
		return nil
	}

	if ref != nil && strings.TrimSpace(ref.Tier) != "" {
		for i := range candidates {
			if strings.EqualFold(candidates[i].Tier, strings.TrimSpace(ref.Tier)) {
				return &candidates[i]
			}
		}
	}

	policy := string(version.SelectionPolicy)
	if ref != nil && strings.TrimSpace(ref.SelectionPolicy) != "" {
		policy = strings.TrimSpace(ref.SelectionPolicy)
	}

	switch policy {
	case string(algorithm.SelectionPolicyHighQuality):
		best := candidates[0]
		for i := 1; i < len(candidates); i++ {
			if candidates[i].QualityScore > best.QualityScore {
				best = candidates[i]
			}
		}
		return &best
	case string(algorithm.SelectionPolicyLowCost):
		best := candidates[0]
		for i := 1; i < len(candidates); i++ {
			if candidates[i].CostScore < best.CostScore {
				best = candidates[i]
			}
		}
		return &best
	default:
		for i := range candidates {
			if candidates[i].IsDefault {
				return &candidates[i]
			}
		}
		first := candidates[0]
		return &first
	}
}

func (e *DAGWorkflowEngine) enforceNodeToolPolicy(ctx context.Context, node *workflow.Node, version *operator.OperatorVersion) error {
	if version == nil {
		return nil
	}

	return e.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if repos.ToolPolicies == nil || repos.Operators == nil {
			return nil
		}

		op, err := repos.Operators.Get(ctx, version.OperatorID)
		if err != nil {
			return err
		}
		policy, err := repos.ToolPolicies.GetByToolName(ctx, op.Code)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}
		if policy == nil {
			return nil
		}
		if !policy.Enabled {
			return fmt.Errorf("tool policy denied: %s is disabled", op.Code)
		}
		if strings.EqualFold(policy.RiskLevel, "high") {
			approved := false
			if node != nil && node.Config != nil && node.Config.Params != nil {
				if v, ok := node.Config.Params["policy_approved"].(bool); ok {
					approved = v
				}
			}
			if !approved {
				return fmt.Errorf("tool policy denied: high-risk tool %s requires policy_approved=true", op.Code)
			}
		}
		if len(policy.Permissions) > 0 {
			granted := extractStringSlice(nil)
			if node != nil && node.Config != nil && node.Config.Params != nil {
				granted = extractStringSlice(node.Config.Params["granted_permissions"])
			}
			missing := missingItems(policy.Permissions, granted)
			if len(missing) > 0 {
				return fmt.Errorf("tool policy denied: %s missing permissions %s", op.Code, strings.Join(missing, ","))
			}
		}

		requestedAccess := extractNodeDataAccess(node)
		if err := validateDataAccessPolicy(op.Code, policy.DataAccess, requestedAccess); err != nil {
			return fmt.Errorf("tool policy denied: %w", err)
		}
		return nil
	})
}

func extractNodeDataAccess(node *workflow.Node) agent.DataAccess {
	out := agent.DataAccess{}
	if node == nil || node.Config == nil || node.Config.Params == nil {
		return out
	}

	params := node.Config.Params
	if raw, ok := params["data_access"]; ok {
		if m, ok := raw.(map[string]interface{}); ok {
			if v, ok := m["read_scopes"]; ok {
				out.ReadScopes = append(out.ReadScopes, extractStringSlice(v)...)
			}
			if v, ok := m["write_scopes"]; ok {
				out.WriteScopes = append(out.WriteScopes, extractStringSlice(v)...)
			}
			if v, ok := m["network_allowlist"]; ok {
				out.NetworkAllowlist = append(out.NetworkAllowlist, extractStringSlice(v)...)
			}
		}
	}

	if v, ok := params["read_scopes"]; ok {
		out.ReadScopes = append(out.ReadScopes, extractStringSlice(v)...)
	}
	if v, ok := params["write_scopes"]; ok {
		out.WriteScopes = append(out.WriteScopes, extractStringSlice(v)...)
	}
	if v, ok := params["network_allowlist"]; ok {
		out.NetworkAllowlist = append(out.NetworkAllowlist, extractStringSlice(v)...)
	}

	out.ReadScopes = dedupeStrings(out.ReadScopes)
	out.WriteScopes = dedupeStrings(out.WriteScopes)
	out.NetworkAllowlist = dedupeStrings(out.NetworkAllowlist)
	return out
}

func validateDataAccessPolicy(toolName string, allowed agent.DataAccess, requested agent.DataAccess) error {
	if err := validateScopeAccess(toolName, "read_scopes", allowed.ReadScopes, requested.ReadScopes); err != nil {
		return err
	}
	if err := validateScopeAccess(toolName, "write_scopes", allowed.WriteScopes, requested.WriteScopes); err != nil {
		return err
	}
	if err := validateScopeAccess(toolName, "network_allowlist", allowed.NetworkAllowlist, requested.NetworkAllowlist); err != nil {
		return err
	}
	return nil
}

func validateScopeAccess(toolName string, scopeType string, allowed []string, requested []string) error {
	if len(requested) == 0 {
		return nil
	}
	if len(allowed) == 0 {
		return fmt.Errorf("%s requests %s but policy allowlist is empty", toolName, scopeType)
	}

	allowSet := make(map[string]struct{}, len(allowed))
	for _, v := range allowed {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		allowSet[v] = struct{}{}
	}
	for _, req := range requested {
		req = strings.TrimSpace(req)
		if req == "" {
			continue
		}
		if _, ok := allowSet[req]; !ok {
			return fmt.Errorf("%s scope %s is not allowed for %s", scopeType, req, toolName)
		}
	}
	return nil
}

func extractStringSlice(raw interface{}) []string {
	switch v := raw.(type) {
	case nil:
		return nil
	case []string:
		out := make([]string, 0, len(v))
		for _, item := range v {
			item = strings.TrimSpace(item)
			if item != "" {
				out = append(out, item)
			}
		}
		return out
	case []interface{}:
		out := make([]string, 0, len(v))
		for _, item := range v {
			switch typed := item.(type) {
			case string:
				typed = strings.TrimSpace(typed)
				if typed != "" {
					out = append(out, typed)
				}
			default:
				text := strings.TrimSpace(fmt.Sprint(typed))
				if text != "" && text != "<nil>" {
					out = append(out, text)
				}
			}
		}
		return out
	case string:
		parts := strings.Split(v, ",")
		out := make([]string, 0, len(parts))
		for _, part := range parts {
			part = strings.TrimSpace(part)
			if part != "" {
				out = append(out, part)
			}
		}
		return out
	default:
		text := strings.TrimSpace(fmt.Sprint(v))
		if text == "" || text == "<nil>" {
			return nil
		}
		return []string{text}
	}
}

func missingItems(required []string, granted []string) []string {
	grantSet := make(map[string]struct{}, len(granted))
	for _, item := range granted {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		grantSet[item] = struct{}{}
	}

	missing := make([]string, 0)
	for _, item := range required {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := grantSet[item]; !ok {
			missing = append(missing, item)
		}
	}
	return missing
}

func dedupeStrings(items []string) []string {
	if len(items) == 0 {
		return nil
	}
	seen := make(map[string]struct{}, len(items))
	out := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		out = append(out, item)
	}
	return out
}

func (e *DAGWorkflowEngine) emitRunEvent(ctx context.Context, event *agent.RunEvent) {
	if event == nil {
		return
	}
	_ = e.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		if repos.RunEvents == nil {
			return nil
		}
		return repos.RunEvents.Create(ctx, event)
	})
}

func (e *DAGWorkflowEngine) resolveWorkflowForTask(ctx context.Context, wf *workflow.Workflow, task *workflow.Task) (*workflow.Workflow, error) {
	if task == nil {
		return wf, nil
	}
	if task.WorkflowRevisionID == nil {
		if wf != nil && wf.CurrentRevisionID != nil {
			task.WorkflowRevisionID = wf.CurrentRevisionID
			task.WorkflowRevision = wf.CurrentRevision
		}
		return wf, nil
	}

	var rev *workflow.WorkflowRevision
	err := e.uow.Do(ctx, func(ctx context.Context, repos *port.Repositories) error {
		var err error
		rev, err = repos.Workflows.GetRevision(ctx, *task.WorkflowRevisionID)
		return err
	})
	if err != nil {
		return nil, err
	}
	if rev == nil {
		return wf, nil
	}

	if task.WorkflowRevision == 0 {
		task.WorkflowRevision = rev.Revision
	}
	return buildWorkflowFromRevisionDefinition(wf, rev), nil
}

func buildWorkflowFromRevisionDefinition(base *workflow.Workflow, rev *workflow.WorkflowRevision) *workflow.Workflow {
	if rev == nil {
		return base
	}
	out := &workflow.Workflow{
		ID:                rev.WorkflowID,
		CurrentRevisionID: &rev.ID,
		CurrentRevision:   rev.Revision,
		Version:           rev.Definition.Version,
		TriggerType:       rev.Definition.TriggerType,
		TriggerConf:       rev.Definition.TriggerConf,
		ContextSpec:       rev.Definition.ContextSpec,
		Nodes:             rev.Definition.Nodes,
		Edges:             rev.Definition.Edges,
	}
	if base != nil {
		out.ID = base.ID
		out.TenantID = base.TenantID
		out.OwnerID = base.OwnerID
		out.Visibility = base.Visibility
		out.VisibleRoleIDs = base.VisibleRoleIDs
		out.Code = base.Code
		out.Name = base.Name
		out.Description = base.Description
		out.Status = base.Status
		out.Tags = base.Tags
		out.CreatedAt = base.CreatedAt
		out.UpdatedAt = base.UpdatedAt
		if out.Version == "" {
			out.Version = base.Version
		}
		if out.TriggerType == "" {
			out.TriggerType = base.TriggerType
		}
		if out.TriggerConf == nil {
			out.TriggerConf = base.TriggerConf
		}
		if out.ContextSpec == nil {
			out.ContextSpec = base.ContextSpec
		}
	}
	return out
}
