package operator

import (
	"time"

	"github.com/google/uuid"
)

type OperatorTemplate struct {
	ID          uuid.UUID
	Code        string
	Name        string
	Description string
	Category    Category
	Type        Type
	ExecMode    ExecMode
	ExecConfig  *ExecConfig
	InputSchema map[string]interface{}
	OutputSpec  map[string]interface{}
	Config      map[string]interface{}
	Author      string
	Tags        []string
	IconURL     string
	Downloads   int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
