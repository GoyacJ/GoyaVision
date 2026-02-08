package system

import (
	"context"
	"encoding/json"
	"time"
)

// SystemConfig 系统配置实体
type SystemConfig struct {
	Key         string          `json:"key"`
	Value       json.RawMessage `json:"value"`
	Description string          `json:"description"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

// ConfigKeys 定义系统配置的 Key 常量
const (
	ConfigKeyHomePath    = "system.home_path"
	ConfigKeyPublicMenus = "system.public_menus"
)

// ConfigRepository 系统配置仓储接口
type ConfigRepository interface {
	Get(ctx context.Context, key string) (*SystemConfig, error)
	List(ctx context.Context) ([]*SystemConfig, error)
	Save(ctx context.Context, config *SystemConfig) error
	Delete(ctx context.Context, key string) error
}
