package repo

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"goyavision/internal/domain/operator"
	"goyavision/internal/infra/persistence/mapper"
	"goyavision/internal/infra/persistence/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OperatorDependencyRepo struct {
	db *gorm.DB
}

func NewOperatorDependencyRepo(db *gorm.DB) *OperatorDependencyRepo {
	return &OperatorDependencyRepo{db: db}
}

func (r *OperatorDependencyRepo) Create(ctx context.Context, dep *operator.OperatorDependency) error {
	if dep.ID == uuid.Nil {
		dep.ID = uuid.New()
	}
	m := mapper.OperatorDependencyToModel(dep)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *OperatorDependencyRepo) ListByOperator(ctx context.Context, operatorID uuid.UUID) ([]*operator.OperatorDependency, error) {
	var models []*model.OperatorDependencyModel
	if err := r.db.WithContext(ctx).Where("operator_id = ?", operatorID).Order("created_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	result := make([]*operator.OperatorDependency, len(models))
	for i, m := range models {
		result[i] = mapper.OperatorDependencyToDomain(m)
	}
	return result, nil
}

func (r *OperatorDependencyRepo) DeleteByOperator(ctx context.Context, operatorID uuid.UUID) error {
	return r.db.WithContext(ctx).Where("operator_id = ?", operatorID).Delete(&model.OperatorDependencyModel{}).Error
}

func (r *OperatorDependencyRepo) CheckDependenciesSatisfied(ctx context.Context, operatorID uuid.UUID) (bool, []string, error) {
	deps, err := r.ListByOperator(ctx, operatorID)
	if err != nil {
		return false, nil, err
	}

	if len(deps) == 0 {
		return true, nil, nil
	}

	var messages []string
	for _, dep := range deps {
		var depOp model.OperatorModel
		err := r.db.WithContext(ctx).
			Select("operators.id", "operators.status", "operators.active_version_id").
			Where("operators.id = ?", dep.DependsOnID).
			First(&depOp).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				if !dep.IsOptional {
					messages = append(messages, fmt.Sprintf("依赖算子 %s 不存在", dep.DependsOnID.String()))
				}
				continue
			}
			return false, nil, err
		}

		if depOp.Status != string(operator.StatusPublished) {
			if !dep.IsOptional {
				messages = append(messages, fmt.Sprintf("依赖算子 %s 未发布", dep.DependsOnID.String()))
			}
			continue
		}

		if dep.MinVersion != "" {
			if depOp.ActiveVersionID == nil {
				if !dep.IsOptional {
					messages = append(messages, fmt.Sprintf("依赖算子 %s 缺少激活版本", dep.DependsOnID.String()))
				}
				continue
			}

			var activeVersion model.OperatorVersionModel
			if err := r.db.WithContext(ctx).
				Select("id", "version").
				Where("id = ?", *depOp.ActiveVersionID).
				First(&activeVersion).Error; err != nil {
				if !dep.IsOptional {
					if err == gorm.ErrRecordNotFound {
						messages = append(messages, fmt.Sprintf("依赖算子 %s 激活版本不存在", dep.DependsOnID.String()))
					} else {
						return false, nil, err
					}
				}
				continue
			}

			ok, err := isVersionGE(activeVersion.Version, dep.MinVersion)
			if err != nil {
				if !dep.IsOptional {
					messages = append(messages, fmt.Sprintf("依赖算子 %s 版本解析失败: %v", dep.DependsOnID.String(), err))
				}
				continue
			}
			if !ok && !dep.IsOptional {
				messages = append(messages, fmt.Sprintf("依赖算子 %s 版本不足，当前 %s，要求 >= %s", dep.DependsOnID.String(), activeVersion.Version, dep.MinVersion))
			}
		}
	}

	return len(messages) == 0, messages, nil
}

func isVersionGE(current string, required string) (bool, error) {
	cv, err := parseVersion(current)
	if err != nil {
		return false, fmt.Errorf("current version %q invalid: %w", current, err)
	}
	rv, err := parseVersion(required)
	if err != nil {
		return false, fmt.Errorf("required version %q invalid: %w", required, err)
	}

	maxLen := len(cv)
	if len(rv) > maxLen {
		maxLen = len(rv)
	}
	for i := 0; i < maxLen; i++ {
		c := 0
		if i < len(cv) {
			c = cv[i]
		}
		r := 0
		if i < len(rv) {
			r = rv[i]
		}
		if c > r {
			return true, nil
		}
		if c < r {
			return false, nil
		}
	}
	return true, nil
}

func parseVersion(raw string) ([]int, error) {
	v := strings.TrimSpace(raw)
	v = strings.TrimPrefix(v, "v")
	if v == "" {
		return nil, fmt.Errorf("empty version")
	}

	if idx := strings.Index(v, "+"); idx >= 0 {
		v = v[:idx]
	}
	if idx := strings.Index(v, "-"); idx >= 0 {
		v = v[:idx]
	}

	parts := strings.Split(v, ".")
	res := make([]int, len(parts))
	for i := range parts {
		if parts[i] == "" {
			return nil, fmt.Errorf("invalid segment in version %q", raw)
		}
		n, err := strconv.Atoi(parts[i])
		if err != nil {
			return nil, fmt.Errorf("invalid segment %q in version %q", parts[i], raw)
		}
		if n < 0 {
			return nil, fmt.Errorf("negative segment in version %q", raw)
		}
		res[i] = n
	}
	return res, nil
}
