package engine

import (
	"context"
	"errors"
	"testing"
	"time"

	"goyavision/internal/app/port"
	"goyavision/internal/domain/operator"
	"goyavision/internal/domain/workflow"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock implementations
type MockUnitOfWork struct {
	mock.Mock
}

func (m *MockUnitOfWork) Do(ctx context.Context, fn func(ctx context.Context, repos *port.Repositories) error) error {
	args := m.Called(ctx, fn)
	if fn != nil {
		fn(ctx, &port.Repositories{})
	}
	return args.Error(0)
}

type MockOperatorExecutor struct {
	mock.Mock
}

func (m *MockOperatorExecutor) Execute(ctx context.Context, op *operator.Operator, input *operator.Input) (*operator.Output, error) {
	args := m.Called(ctx, op, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*operator.Output), args.Error(1)
}

// Test topological sort
func TestTopologicalSort(t *testing.T) {
	tests := []struct {
		name    string
		nodes   []workflow.Node
		edges   []workflow.Edge
		wantLen int
		wantErr bool
	}{
		{
			name: "linear graph",
			nodes: []workflow.Node{
				{NodeKey: "A"},
				{NodeKey: "B"},
				{NodeKey: "C"},
			},
			edges: []workflow.Edge{
				{SourceKey: "A", TargetKey: "B"},
				{SourceKey: "B", TargetKey: "C"},
			},
			wantLen: 3,
			wantErr: false,
		},
		{
			name: "parallel branches",
			nodes: []workflow.Node{
				{NodeKey: "A"},
				{NodeKey: "B"},
				{NodeKey: "C"},
				{NodeKey: "D"},
			},
			edges: []workflow.Edge{
				{SourceKey: "A", TargetKey: "B"},
				{SourceKey: "A", TargetKey: "C"},
				{SourceKey: "B", TargetKey: "D"},
				{SourceKey: "C", TargetKey: "D"},
			},
			wantLen: 4,
			wantErr: false,
		},
		{
			name: "cycle detection",
			nodes: []workflow.Node{
				{NodeKey: "A"},
				{NodeKey: "B"},
				{NodeKey: "C"},
			},
			edges: []workflow.Edge{
				{SourceKey: "A", TargetKey: "B"},
				{SourceKey: "B", TargetKey: "C"},
				{SourceKey: "C", TargetKey: "A"},
			},
			wantLen: 0,
			wantErr: true,
		},
		{
			name: "disconnected nodes",
			nodes: []workflow.Node{
				{NodeKey: "A"},
				{NodeKey: "B"},
				{NodeKey: "C"},
			},
			edges:   []workflow.Edge{},
			wantLen: 3,
			wantErr: false,
		},
		{
			name: "complex DAG",
			nodes: []workflow.Node{
				{NodeKey: "start"},
				{NodeKey: "parallel1"},
				{NodeKey: "parallel2"},
				{NodeKey: "parallel3"},
				{NodeKey: "join"},
				{NodeKey: "end"},
			},
			edges: []workflow.Edge{
				{SourceKey: "start", TargetKey: "parallel1"},
				{SourceKey: "start", TargetKey: "parallel2"},
				{SourceKey: "start", TargetKey: "parallel3"},
				{SourceKey: "parallel1", TargetKey: "join"},
				{SourceKey: "parallel2", TargetKey: "join"},
				{SourceKey: "parallel3", TargetKey: "join"},
				{SourceKey: "join", TargetKey: "end"},
			},
			wantLen: 6,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := &DAGWorkflowEngine{}
			got, err := engine.topologicalSort(tt.nodes, tt.edges)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cycle")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantLen, len(got))

				// Verify all nodes are present
				nodeSet := make(map[string]bool)
				for _, nodeKey := range got {
					nodeSet[nodeKey] = true
				}
				for _, node := range tt.nodes {
					assert.True(t, nodeSet[node.NodeKey], "node %s not in sorted result", node.NodeKey)
				}
			}
		})
	}
}

// Test execution layers
func TestBuildExecutionLayers(t *testing.T) {
	tests := []struct {
		name       string
		nodes      []workflow.Node
		edges      []workflow.Edge
		wantLayers int
		wantErr    bool
		validate   func(*testing.T, [][]string)
	}{
		{
			name: "single node",
			nodes: []workflow.Node{
				{NodeKey: "A"},
			},
			edges:      []workflow.Edge{},
			wantLayers: 1,
			wantErr:    false,
			validate: func(t *testing.T, layers [][]string) {
				assert.Equal(t, 1, len(layers[0]))
				assert.Equal(t, "A", layers[0][0])
			},
		},
		{
			name: "linear sequence",
			nodes: []workflow.Node{
				{NodeKey: "A"},
				{NodeKey: "B"},
				{NodeKey: "C"},
			},
			edges: []workflow.Edge{
				{SourceKey: "A", TargetKey: "B"},
				{SourceKey: "B", TargetKey: "C"},
			},
			wantLayers: 3,
			wantErr:    false,
			validate: func(t *testing.T, layers [][]string) {
				assert.Equal(t, 1, len(layers[0]))
				assert.Equal(t, 1, len(layers[1]))
				assert.Equal(t, 1, len(layers[2]))
			},
		},
		{
			name: "diamond pattern",
			nodes: []workflow.Node{
				{NodeKey: "A"},
				{NodeKey: "B"},
				{NodeKey: "C"},
				{NodeKey: "D"},
			},
			edges: []workflow.Edge{
				{SourceKey: "A", TargetKey: "B"},
				{SourceKey: "A", TargetKey: "C"},
				{SourceKey: "B", TargetKey: "D"},
				{SourceKey: "C", TargetKey: "D"},
			},
			wantLayers: 3,
			wantErr:    false,
			validate: func(t *testing.T, layers [][]string) {
				// Layer 1: A
				assert.Equal(t, 1, len(layers[0]))
				assert.Equal(t, "A", layers[0][0])
				// Layer 2: B, C (parallel)
				assert.Equal(t, 2, len(layers[1]))
				assert.Contains(t, layers[1], "B")
				assert.Contains(t, layers[1], "C")
				// Layer 3: D
				assert.Equal(t, 1, len(layers[2]))
				assert.Equal(t, "D", layers[2][0])
			},
		},
		{
			name: "wide parallelism",
			nodes: []workflow.Node{
				{NodeKey: "start"},
				{NodeKey: "p1"},
				{NodeKey: "p2"},
				{NodeKey: "p3"},
				{NodeKey: "p4"},
				{NodeKey: "end"},
			},
			edges: []workflow.Edge{
				{SourceKey: "start", TargetKey: "p1"},
				{SourceKey: "start", TargetKey: "p2"},
				{SourceKey: "start", TargetKey: "p3"},
				{SourceKey: "start", TargetKey: "p4"},
				{SourceKey: "p1", TargetKey: "end"},
				{SourceKey: "p2", TargetKey: "end"},
				{SourceKey: "p3", TargetKey: "end"},
				{SourceKey: "p4", TargetKey: "end"},
			},
			wantLayers: 3,
			wantErr:    false,
			validate: func(t *testing.T, layers [][]string) {
				// Layer 1: start
				assert.Equal(t, 1, len(layers[0]))
				// Layer 2: p1, p2, p3, p4 (parallel)
				assert.Equal(t, 4, len(layers[1]))
				// Layer 3: end
				assert.Equal(t, 1, len(layers[2]))
			},
		},
		{
			name: "cycle detection",
			nodes: []workflow.Node{
				{NodeKey: "A"},
				{NodeKey: "B"},
			},
			edges: []workflow.Edge{
				{SourceKey: "A", TargetKey: "B"},
				{SourceKey: "B", TargetKey: "A"},
			},
			wantLayers: 0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := &DAGWorkflowEngine{}
			layers, err := engine.buildExecutionLayers(tt.nodes, tt.edges)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "cycle")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantLayers, len(layers))

				if tt.validate != nil {
					tt.validate(t, layers)
				}

				// Verify all nodes are present across layers
				nodeSet := make(map[string]bool)
				for _, layer := range layers {
					for _, nodeKey := range layer {
						nodeSet[nodeKey] = true
					}
				}
				for _, node := range tt.nodes {
					assert.True(t, nodeSet[node.NodeKey], "node %s not in layers", node.NodeKey)
				}
			}
		})
	}
}

// Test prepareNodeInput
func TestPrepareNodeInput(t *testing.T) {
	engine := &DAGWorkflowEngine{}

	assetID := uuid.New()
	task := &workflow.Task{
		AssetID: &assetID,
		InputParams: map[string]interface{}{
			"task_param1": "value1",
			"task_param2": 42,
		},
	}

	node := &workflow.Node{
		NodeKey: "test_node",
		Config: &workflow.NodeConfig{
			Params: map[string]interface{}{
				"node_param1": "override",
				"node_param2": "new_value",
			},
		},
	}

	exec := &taskExecution{
		nodeResults: map[string]*operator.Output{
			"upstream_node": {
				OutputAssets: []operator.OutputAsset{
					{Type: "video", Path: "/path/to/video.mp4"},
				},
				Results: []operator.Result{
					{Type: "detection", Data: map[string]interface{}{"objects": 5}},
				},
			},
		},
	}

	input := engine.prepareNodeInput(task, node, exec)

	// Check asset ID
	assert.Equal(t, assetID, input.AssetID)

	// Check task params are present
	assert.Equal(t, "value1", input.Params["task_param1"])
	assert.Equal(t, 42, input.Params["task_param2"])

	// Check node params override
	assert.Equal(t, "override", input.Params["node_param1"])
	assert.Equal(t, "new_value", input.Params["node_param2"])

	// Check upstream node outputs are available
	assert.NotNil(t, input.Params["upstream_node_output"])
	assert.NotNil(t, input.Params["upstream_node_assets"])
	assert.NotNil(t, input.Params["upstream_node_results"])
}

// Test execute with mock dependencies
func TestExecute_Success(t *testing.T) {
	mockUOW := new(MockUnitOfWork)
	mockExecutor := new(MockOperatorExecutor)

	engine := NewDAGWorkflowEngine(mockUOW, mockExecutor)

	operatorID := uuid.New()
	wf := &workflow.Workflow{
		ID: uuid.New(),
		Nodes: []workflow.Node{
			{
				ID:         uuid.New(),
				NodeKey:    "node1",
				OperatorID: &operatorID,
			},
		},
		Edges: []workflow.Edge{},
	}

	task := &workflow.Task{
		ID:         uuid.New(),
		WorkflowID: wf.ID,
		Status:     workflow.TaskStatusPending,
	}

	mockUOW.On("Do", mock.Anything, mock.Anything).Return(nil)
	mockExecutor.On("Execute", mock.Anything, mock.Anything, mock.Anything).Return(
		&operator.Output{
			Results: []operator.Result{
				{Type: "test", Data: map[string]interface{}{"result": "success"}},
			},
		},
		nil,
	)

	ctx := context.Background()
	err := engine.Execute(ctx, wf, task)

	assert.NoError(t, err)
	mockUOW.AssertExpectations(t)
	mockExecutor.AssertExpectations(t)
}

// Test execute with parallel nodes
func TestExecute_ParallelNodes(t *testing.T) {
	mockUOW := new(MockUnitOfWork)
	mockExecutor := new(MockOperatorExecutor)

	engine := NewDAGWorkflowEngine(mockUOW, mockExecutor)

	op1ID := uuid.New()
	op2ID := uuid.New()
	op3ID := uuid.New()

	wf := &workflow.Workflow{
		ID: uuid.New(),
		Nodes: []workflow.Node{
			{ID: uuid.New(), NodeKey: "start", OperatorID: &op1ID},
			{ID: uuid.New(), NodeKey: "parallel1", OperatorID: &op2ID},
			{ID: uuid.New(), NodeKey: "parallel2", OperatorID: &op3ID},
		},
		Edges: []workflow.Edge{
			{SourceKey: "start", TargetKey: "parallel1"},
			{SourceKey: "start", TargetKey: "parallel2"},
		},
	}

	task := &workflow.Task{
		ID:         uuid.New(),
		WorkflowID: wf.ID,
		Status:     workflow.TaskStatusPending,
	}

	mockUOW.On("Do", mock.Anything, mock.Anything).Return(nil)
	mockExecutor.On("Execute", mock.Anything, mock.Anything, mock.Anything).Return(
		&operator.Output{},
		nil,
	)

	ctx := context.Background()
	err := engine.Execute(ctx, wf, task)

	assert.NoError(t, err)
	// Should execute 3 times (once per node)
	assert.Equal(t, 3, len(mockExecutor.Calls))
}

// Test execute with cycle detection
func TestExecute_CycleDetection(t *testing.T) {
	mockUOW := new(MockUnitOfWork)
	mockExecutor := new(MockOperatorExecutor)

	engine := NewDAGWorkflowEngine(mockUOW, mockExecutor)

	wf := &workflow.Workflow{
		ID: uuid.New(),
		Nodes: []workflow.Node{
			{ID: uuid.New(), NodeKey: "A"},
			{ID: uuid.New(), NodeKey: "B"},
			{ID: uuid.New(), NodeKey: "C"},
		},
		Edges: []workflow.Edge{
			{SourceKey: "A", TargetKey: "B"},
			{SourceKey: "B", TargetKey: "C"},
			{SourceKey: "C", TargetKey: "A"},
		},
	}

	task := &workflow.Task{
		ID:         uuid.New(),
		WorkflowID: wf.ID,
		Status:     workflow.TaskStatusPending,
	}

	ctx := context.Background()
	err := engine.Execute(ctx, wf, task)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "cycle")
}

// Test cancel
func TestCancel(t *testing.T) {
	mockUOW := new(MockUnitOfWork)
	mockExecutor := new(MockOperatorExecutor)

	engine := NewDAGWorkflowEngine(mockUOW, mockExecutor)

	taskID := uuid.New()

	// Setup a running task
	ctx, cancel := context.WithCancel(context.Background())
	engine.tasks[taskID] = &taskExecution{
		ctx:         ctx,
		cancel:      cancel,
		nodeResults: make(map[string]*operator.Output),
	}

	// Cancel the task
	err := engine.Cancel(context.Background(), taskID)
	assert.NoError(t, err)

	// Verify context is cancelled
	select {
	case <-ctx.Done():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Fatal("context not cancelled")
	}

	// Try to cancel non-existent task
	err = engine.Cancel(context.Background(), uuid.New())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not running")
}

// Test get progress
func TestGetProgress(t *testing.T) {
	mockUOW := new(MockUnitOfWork)
	mockExecutor := new(MockOperatorExecutor)

	engine := NewDAGWorkflowEngine(mockUOW, mockExecutor)

	taskID := uuid.New()

	// Test running task
	engine.tasks[taskID] = &taskExecution{
		progress:    50,
		nodeResults: make(map[string]*operator.Output),
	}

	progress, err := engine.GetProgress(context.Background(), taskID)
	assert.NoError(t, err)
	assert.Equal(t, 50, progress)

	// Test non-running task (fetch from database)
	nonRunningTaskID := uuid.New()
	mockUOW.On("Do", mock.Anything, mock.Anything).Return(nil).Once()

	progress, err = engine.GetProgress(context.Background(), nonRunningTaskID)
	assert.NoError(t, err)
}

// Test node execution with retry
func TestExecuteNode_WithRetry(t *testing.T) {
	mockUOW := new(MockUnitOfWork)
	mockExecutor := new(MockOperatorExecutor)

	engine := NewDAGWorkflowEngine(mockUOW, mockExecutor)

	operatorID := uuid.New()
	node := &workflow.Node{
		NodeKey:    "test_node",
		OperatorID: &operatorID,
		Config: &workflow.NodeConfig{
			RetryCount: 2,
		},
	}

	task := &workflow.Task{
		ID: uuid.New(),
	}

	exec := &taskExecution{
		nodeResults: make(map[string]*operator.Output),
	}

	mockUOW.On("Do", mock.Anything, mock.Anything).Return(nil)

	// First two calls fail, third succeeds
	mockExecutor.On("Execute", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("temporary error")).Twice()
	mockExecutor.On("Execute", mock.Anything, mock.Anything, mock.Anything).
		Return(&operator.Output{}, nil).Once()

	ctx := context.Background()
	err := engine.executeNode(ctx, node, task, exec)

	assert.NoError(t, err)
	// Should have called execute 3 times (1 original + 2 retries)
	assert.Equal(t, 3, len(mockExecutor.Calls))
}

// Test node execution with timeout
func TestExecuteNode_WithTimeout(t *testing.T) {
	mockUOW := new(MockUnitOfWork)
	mockExecutor := new(MockOperatorExecutor)

	engine := NewDAGWorkflowEngine(mockUOW, mockExecutor)

	operatorID := uuid.New()
	node := &workflow.Node{
		NodeKey:    "test_node",
		OperatorID: &operatorID,
		Config: &workflow.NodeConfig{
			TimeoutSeconds: 1,
		},
	}

	task := &workflow.Task{
		ID: uuid.New(),
	}

	exec := &taskExecution{
		nodeResults: make(map[string]*operator.Output),
	}

	mockUOW.On("Do", mock.Anything, mock.Anything).Return(nil)

	// Simulate slow execution
	mockExecutor.On("Execute", mock.Anything, mock.Anything, mock.Anything).
		Run(func(args mock.Arguments) {
			ctx := args.Get(0).(context.Context)
			select {
			case <-time.After(2 * time.Second):
			case <-ctx.Done():
				// Timeout occurred as expected
			}
		}).
		Return(nil, context.DeadlineExceeded)

	ctx := context.Background()
	err := engine.executeNode(ctx, node, task, exec)

	assert.Error(t, err)
}
