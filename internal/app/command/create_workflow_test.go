package command

import (
	"context"
	"testing"

	"goyavision/internal/app/dto"
	"goyavision/internal/app/port"
	"goyavision/internal/domain/operator"
	"goyavision/internal/domain/workflow"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// --- Mocks ---

type MockUnitOfWork struct {
	mock.Mock
	Repos *port.Repositories
}

func (m *MockUnitOfWork) Do(ctx context.Context, fn func(ctx context.Context, repos *port.Repositories) error) error {
	args := m.Called(ctx, fn)
	if args.Error(0) != nil {
		return args.Error(0)
	}
	// Execute the function with the mock repositories
	return fn(ctx, m.Repos)
}

type MockWorkflowRepo struct {
	mock.Mock
}

func (m *MockWorkflowRepo) Create(ctx context.Context, w *workflow.Workflow) error {
	args := m.Called(ctx, w)
	// Assign ID to workflow to simulate DB creation
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return args.Error(0)
}
func (m *MockWorkflowRepo) Get(ctx context.Context, id uuid.UUID) (*workflow.Workflow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*workflow.Workflow), args.Error(1)
}
func (m *MockWorkflowRepo) GetByCode(ctx context.Context, code string) (*workflow.Workflow, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*workflow.Workflow), args.Error(1)
}
func (m *MockWorkflowRepo) GetWithNodes(ctx context.Context, id uuid.UUID) (*workflow.Workflow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*workflow.Workflow), args.Error(1)
}
func (m *MockWorkflowRepo) List(ctx context.Context, filter workflow.Filter) ([]*workflow.Workflow, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*workflow.Workflow), args.Get(1).(int64), args.Error(2)
}
func (m *MockWorkflowRepo) Update(ctx context.Context, w *workflow.Workflow) error {
	return m.Called(ctx, w).Error(0)
}
func (m *MockWorkflowRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockWorkflowRepo) ListEnabled(ctx context.Context) ([]*workflow.Workflow, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*workflow.Workflow), args.Error(1)
}
func (m *MockWorkflowRepo) CreateNode(ctx context.Context, n *workflow.Node) error {
	return m.Called(ctx, n).Error(0)
}
func (m *MockWorkflowRepo) ListNodes(ctx context.Context, workflowID uuid.UUID) ([]*workflow.Node, error) {
	args := m.Called(ctx, workflowID)
	return args.Get(0).([]*workflow.Node), args.Error(1)
}
func (m *MockWorkflowRepo) DeleteNodes(ctx context.Context, workflowID uuid.UUID) error {
	return m.Called(ctx, workflowID).Error(0)
}
func (m *MockWorkflowRepo) CreateEdge(ctx context.Context, e *workflow.Edge) error {
	return m.Called(ctx, e).Error(0)
}
func (m *MockWorkflowRepo) ListEdges(ctx context.Context, workflowID uuid.UUID) ([]*workflow.Edge, error) {
	args := m.Called(ctx, workflowID)
	return args.Get(0).([]*workflow.Edge), args.Error(1)
}
func (m *MockWorkflowRepo) DeleteEdges(ctx context.Context, workflowID uuid.UUID) error {
	return m.Called(ctx, workflowID).Error(0)
}

type MockOperatorRepo struct {
	mock.Mock
}

func (m *MockOperatorRepo) Create(ctx context.Context, o *operator.Operator) error {
	return m.Called(ctx, o).Error(0)
}
func (m *MockOperatorRepo) Get(ctx context.Context, id uuid.UUID) (*operator.Operator, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*operator.Operator), args.Error(1)
}
func (m *MockOperatorRepo) GetWithActiveVersion(ctx context.Context, id uuid.UUID) (*operator.Operator, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*operator.Operator), args.Error(1)
}
func (m *MockOperatorRepo) GetByCode(ctx context.Context, code string) (*operator.Operator, error) {
	args := m.Called(ctx, code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*operator.Operator), args.Error(1)
}
func (m *MockOperatorRepo) List(ctx context.Context, filter operator.Filter) ([]*operator.Operator, int64, error) {
	args := m.Called(ctx, filter)
	return args.Get(0).([]*operator.Operator), args.Get(1).(int64), args.Error(2)
}
func (m *MockOperatorRepo) Update(ctx context.Context, o *operator.Operator) error {
	return m.Called(ctx, o).Error(0)
}
func (m *MockOperatorRepo) Delete(ctx context.Context, id uuid.UUID) error {
	return m.Called(ctx, id).Error(0)
}
func (m *MockOperatorRepo) ListEnabled(ctx context.Context) ([]*operator.Operator, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*operator.Operator), args.Error(1)
}
func (m *MockOperatorRepo) ListPublished(ctx context.Context) ([]*operator.Operator, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*operator.Operator), args.Error(1)
}
func (m *MockOperatorRepo) ListByCategory(ctx context.Context, category operator.Category) ([]*operator.Operator, error) {
	args := m.Called(ctx, category)
	return args.Get(0).([]*operator.Operator), args.Error(1)
}

type MockSchemaValidator struct {
	mock.Mock
}

func (m *MockSchemaValidator) ValidateInput(ctx context.Context, schema map[string]interface{}, input map[string]interface{}) error {
	return m.Called(ctx, schema, input).Error(0)
}
func (m *MockSchemaValidator) ValidateOutput(ctx context.Context, schema map[string]interface{}, output map[string]interface{}) error {
	return m.Called(ctx, schema, output).Error(0)
}
func (m *MockSchemaValidator) IsValidJSONSchema(ctx context.Context, schema map[string]interface{}) error {
	return m.Called(ctx, schema).Error(0)
}
func (m *MockSchemaValidator) ValidateConnection(ctx context.Context, upstreamOutputSpec map[string]interface{}, downstreamInputSchema map[string]interface{}) error {
	return m.Called(ctx, upstreamOutputSpec, downstreamInputSchema).Error(0)
}

// --- Tests ---

func TestCreateWorkflow_Persistence(t *testing.T) {
	// Setup Mocks
	mockUOW := new(MockUnitOfWork)
	mockWFRepo := new(MockWorkflowRepo)
	mockOpRepo := new(MockOperatorRepo)
	mockValidator := new(MockSchemaValidator)

	mockUOW.Repos = &port.Repositories{
		Workflows: mockWFRepo,
		Operators: mockOpRepo,
	}

	handler := NewCreateWorkflowHandler(mockUOW, mockValidator)

	// Test Data
	opID := uuid.New()
	nodeKey := "node-1"
	targetKey := "node-2"
	
	cmd := dto.CreateWorkflowCommand{
		Code:        "test-wf",
		Name:        "Test Workflow",
		TriggerType: workflow.TriggerTypeManual,
		Nodes: []dto.WorkflowNodeInput{
			{
				NodeKey:    nodeKey,
				NodeType:   "operator",
				OperatorID: &opID,
				Config: map[string]interface{}{
					"retry_count":     float64(3), // JSON unmarshal typically produces float64 for numbers
					"timeout_seconds": float64(60),
					"params": map[string]interface{}{
						"key": "value",
					},
				},
				Position: map[string]interface{}{
					"x": float64(100),
					"y": float64(200),
				},
			},
			{
				NodeKey:  targetKey,
				NodeType: "end",
			},
		},
		Edges: []dto.WorkflowEdgeInput{
			{
				SourceKey: nodeKey,
				TargetKey: targetKey,
				Condition: map[string]interface{}{
					"type": "on_success",
				},
			},
		},
	}

	// Expectations
	mockUOW.On("Do", mock.Anything, mock.Anything).Return(nil)
	
	// Workflow check existing
	mockWFRepo.On("GetByCode", mock.Anything, "test-wf").Return(nil, assert.AnError).Maybe() // Return error to simulate not found
	
	// Operator check
	mockOpRepo.On("Get", mock.Anything, opID).Return(&operator.Operator{ID: opID, ActiveVersion: &operator.OperatorVersion{}}, nil)
	mockOpRepo.On("GetWithActiveVersion", mock.Anything, opID).Return(&operator.Operator{ID: opID, ActiveVersion: &operator.OperatorVersion{}}, nil)

	// Create Workflow
	mockWFRepo.On("Create", mock.Anything, mock.AnythingOfType("*workflow.Workflow")).Return(nil)

	// Create Nodes - Validate Config and Position persistence
	mockWFRepo.On("CreateNode", mock.Anything, mock.MatchedBy(func(n *workflow.Node) bool {
		if n.NodeKey == nodeKey {
			// Verify Config
			if n.Config == nil || n.Config.RetryCount != 3 || n.Config.TimeoutSeconds != 60 {
				return false
			}
			if n.Config.Params["key"] != "value" {
				return false
			}
			// Verify Position
			if n.Position == nil || n.Position.X != 100 || n.Position.Y != 200 {
				return false
			}
			return true
		}
		return true // Other nodes
	})).Return(nil)

	// Create Edge - Validate Condition persistence
	mockWFRepo.On("CreateEdge", mock.Anything, mock.MatchedBy(func(e *workflow.Edge) bool {
		if e.SourceKey == nodeKey && e.TargetKey == targetKey {
			// Verify Condition
			if e.Condition == nil || e.Condition.Type != "on_success" {
				return false
			}
			return true
		}
		return false
	})).Return(nil)

	// GetWithNodes for return
	mockWFRepo.On("GetWithNodes", mock.Anything, mock.Anything).Return(&workflow.Workflow{}, nil)

	// Execute
	_, err := handler.Handle(context.Background(), cmd)

	// Assert
	assert.NoError(t, err)
	mockWFRepo.AssertExpectations(t)
}
