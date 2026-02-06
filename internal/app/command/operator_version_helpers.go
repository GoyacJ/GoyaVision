package command

import "goyavision/internal/domain/operator"

func syncOperatorCompatFieldsFromVersion(op *operator.Operator, version *operator.OperatorVersion) {
	if op == nil || version == nil {
		return
	}
	// 兼容字段收口：不再在写路径同步旧执行字段，统一以 ActiveVersion 为准。
}
