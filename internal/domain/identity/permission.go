package identity

import (
	"time"

	"github.com/google/uuid"
)

// Permission 权限实体（API 资源，纯域模型，无 ORM 依赖）
type Permission struct {
	ID          uuid.UUID
	Code        string // 权限编码（唯一标识）
	Name        string // 权限名称
	Method      string // HTTP 方法（GET, POST, PUT, DELETE, *）
	Path        string // API 路径
	Description string // 权限描述
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewPermission 创建新权限
func NewPermission(code, name, method, path string) *Permission {
	return &Permission{
		ID:     uuid.New(),
		Code:   code,
		Name:   name,
		Method: method,
		Path:   path,
	}
}

// MatchesRequest 检查权限是否匹配请求
func (p *Permission) MatchesRequest(method, path string) bool {
	// 通配符方法匹配所有
	if p.Method == "*" && p.Path == path {
		return true
	}

	// 精确匹配
	return p.Method == method && p.Path == path
}

// IsPublic 是否为公开权限（无需认证）
func (p *Permission) IsPublic() bool {
	// 可以根据 Code 前缀或特定标记判断
	// 例如：public:* 开头的权限为公开权限
	return false // 默认都需要认证
}
