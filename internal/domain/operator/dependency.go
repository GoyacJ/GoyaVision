package operator

import (
	"time"

	"github.com/google/uuid"
)

type OperatorDependency struct {
	ID          uuid.UUID
	OperatorID  uuid.UUID
	DependsOnID uuid.UUID
	MinVersion  string
	IsOptional  bool
	CreatedAt   time.Time
}
