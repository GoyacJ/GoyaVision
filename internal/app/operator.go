package app

import (
	"context"
	"errors"

	"goyavision/internal/domain"
	"goyavision/internal/port"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// CreateOperatorRequest 创建算子请求
type CreateOperatorRequest struct {
	Code        string                 `json:"code"`
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Category    domain.OperatorCategory `json:"category"`
	Type        domain.OperatorType    `json:"type"`
	Version     string                 `json:"version,omitempty"`
	Endpoint    string                 `json:"endpoint"`
	Method      string                 `json:"method,omitempty"`
	InputSchema map[string]interface{} `json:"input_schema,omitempty"`
	OutputSpec  map[string]interface{} `json:"output_spec,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Status      domain.OperatorStatus  `json:"status,omitempty"`
	IsBuiltin   bool                   `json:"is_builtin,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
}

// UpdateOperatorRequest 更新算子请求
type UpdateOperatorRequest struct {
	Name        *string                `json:"name,omitempty"`
	Description *string                `json:"description,omitempty"`
	Endpoint    *string                `json:"endpoint,omitempty"`
	Method      *string                `json:"method,omitempty"`
	InputSchema map[string]interface{} `json:"input_schema,omitempty"`
	OutputSpec  map[string]interface{} `json:"output_spec,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty"`
	Status      *domain.OperatorStatus `json:"status,omitempty"`
	Tags        []string               `json:"tags,omitempty"`
}

// ListOperatorsRequest 列出算子请求
type ListOperatorsRequest struct {
	Category  *domain.OperatorCategory
	Type      *domain.OperatorType
	Status    *domain.OperatorStatus
	IsBuiltin *bool
	Tags      []string
	Keyword   string
	Limit     int
	Offset    int
}

type OperatorService struct {
	repo port.Repository
}

func NewOperatorService(repo port.Repository) *OperatorService {
	return &OperatorService{
		repo: repo,
	}
}

// Create 创建算子
func (s *OperatorService) Create(ctx context.Context, req *CreateOperatorRequest) (*domain.Operator, error) {
	if req.Code == "" {
		return nil, errors.New("code is required")
	}
	if req.Name == "" {
		return nil, errors.New("name is required")
	}
	if req.Category == "" {
		return nil, errors.New("category is required")
	}
	if req.Type == "" {
		return nil, errors.New("type is required")
	}
	if req.Endpoint == "" {
		return nil, errors.New("endpoint is required")
	}

	if req.Category != domain.OperatorCategoryAnalysis &&
		req.Category != domain.OperatorCategoryProcessing &&
		req.Category != domain.OperatorCategoryGeneration &&
		req.Category != domain.OperatorCategoryUtility {
		return nil, errors.New("invalid category")
	}

	if _, err := s.repo.GetOperatorByCode(ctx, req.Code); err == nil {
		return nil, errors.New("operator code already exists")
	}

	version := "1.0.0"
	if req.Version != "" {
		version = req.Version
	}

	method := "POST"
	if req.Method != "" {
		method = req.Method
	}

	status := domain.OperatorStatusDraft
	if req.Status != "" {
		status = req.Status
	}

	operator := &domain.Operator{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
		Category:    req.Category,
		Type:        req.Type,
		Version:     version,
		Endpoint:    req.Endpoint,
		Method:      method,
		Status:      status,
		IsBuiltin:   req.IsBuiltin,
	}

	if err := s.repo.CreateOperator(ctx, operator); err != nil {
		return nil, err
	}

	return operator, nil
}

// Get 获取算子
func (s *OperatorService) Get(ctx context.Context, id uuid.UUID) (*domain.Operator, error) {
	operator, err := s.repo.GetOperator(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("operator not found")
		}
		return nil, err
	}
	return operator, nil
}

// GetByCode 根据代码获取算子
func (s *OperatorService) GetByCode(ctx context.Context, code string) (*domain.Operator, error) {
	operator, err := s.repo.GetOperatorByCode(ctx, code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("operator not found")
		}
		return nil, err
	}
	return operator, nil
}

// List 列出算子
func (s *OperatorService) List(ctx context.Context, req *ListOperatorsRequest) ([]*domain.Operator, int64, error) {
	if req.Limit <= 0 {
		req.Limit = 20
	}
	if req.Limit > 1000 {
		req.Limit = 1000
	}

	filter := domain.OperatorFilter{
		Category:  req.Category,
		Type:      req.Type,
		Status:    req.Status,
		IsBuiltin: req.IsBuiltin,
		Tags:      req.Tags,
		Keyword:   req.Keyword,
		Limit:     req.Limit,
		Offset:    req.Offset,
	}

	return s.repo.ListOperators(ctx, filter)
}

// Update 更新算子
func (s *OperatorService) Update(ctx context.Context, id uuid.UUID, req *UpdateOperatorRequest) (*domain.Operator, error) {
	operator, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if operator.IsBuiltin {
		return nil, errors.New("cannot update builtin operator")
	}

	if req.Name != nil {
		operator.Name = *req.Name
	}
	if req.Description != nil {
		operator.Description = *req.Description
	}
	if req.Endpoint != nil {
		operator.Endpoint = *req.Endpoint
	}
	if req.Method != nil {
		operator.Method = *req.Method
	}
	if req.Status != nil {
		operator.Status = *req.Status
	}

	if err := s.repo.UpdateOperator(ctx, operator); err != nil {
		return nil, err
	}

	return operator, nil
}

// Delete 删除算子
func (s *OperatorService) Delete(ctx context.Context, id uuid.UUID) error {
	operator, err := s.Get(ctx, id)
	if err != nil {
		return err
	}

	if operator.IsBuiltin {
		return errors.New("cannot delete builtin operator")
	}

	return s.repo.DeleteOperator(ctx, id)
}

// Enable 启用算子
func (s *OperatorService) Enable(ctx context.Context, id uuid.UUID) (*domain.Operator, error) {
	operator, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	operator.Status = domain.OperatorStatusEnabled
	if err := s.repo.UpdateOperator(ctx, operator); err != nil {
		return nil, err
	}

	return operator, nil
}

// Disable 禁用算子
func (s *OperatorService) Disable(ctx context.Context, id uuid.UUID) (*domain.Operator, error) {
	operator, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	operator.Status = domain.OperatorStatusDisabled
	if err := s.repo.UpdateOperator(ctx, operator); err != nil {
		return nil, err
	}

	return operator, nil
}

// ListEnabled 列出所有启用的算子
func (s *OperatorService) ListEnabled(ctx context.Context) ([]*domain.Operator, error) {
	return s.repo.ListEnabledOperators(ctx)
}

// ListByCategory 根据分类列出算子
func (s *OperatorService) ListByCategory(ctx context.Context, category domain.OperatorCategory) ([]*domain.Operator, error) {
	return s.repo.ListOperatorsByCategory(ctx, category)
}
