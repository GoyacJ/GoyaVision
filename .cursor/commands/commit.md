# commit

创建符合 Conventional Commits 规范的 Git 提交。

## 分支管理规范

在提交前，请确保您在正确的分支上工作：
- **新功能**：`feature/*` (拉自 `develop`)
- **普通修复**：`fix/*` (拉自 `develop`)
- **紧急修复**：`hotfix/*` (拉自 `main`)
- **重构**：`refactor/*` (拉自 `develop`)

## 提交格式

```
<type>(<scope>): <subject>

[optional body]

[optional footer]
```

## Type 类型

- `feat`: 新功能
- `fix`: Bug 修复
- `docs`: 文档变更
- `refactor`: 代码重构（不改变功能）
- `test`: 测试相关
- `chore`: 构建、配置、依赖更新
- `perf`: 性能优化
- `style`: 代码格式（不影响逻辑）

## Scope 范围（可选）

`asset`, `operator`, `workflow`, `task`, `auth`, `api`, `ui`, `source`

## 提交示例

### ✅ 新功能
```bash
git commit -m "feat(asset): 实现媒体资产管理功能

- 添加 MediaAsset 实体和 Repository
- 实现 MediaAssetService（CRUD、搜索、派生追踪）
- 添加 Asset Handler 和 DTO
- 更新 API 文档和开发进度"
```

### ✅ Bug 修复
```bash
git commit -m "fix(workflow): 修复 DAG 验证时的死循环问题

在验证包含自环的 DAG 时会立即返回错误，
避免进入无限循环。

Closes #123"
```

### ✅ 文档更新
```bash
git commit -m "docs: 更新开发进度、CHANGELOG 与 API 文档"
```

### ✅ 重构
```bash
git commit -m "refactor(operator): 重构算子标准协议实现

提取公共接口，简化算子实现代码。"
```

## 提交前自检清单

- [ ] **代码已测试** - 单元测试、集成测试、手动测试已通过
- [ ] **文档已更新** - 开发进度、CHANGELOG、API/架构文档已同步
- [ ] **代码已格式化** - Go: `go fmt ./...` 和 `go vet ./...`
- [ ] **无调试代码** - 移除 console.log、fmt.Println、注释代码块、临时文件
- [ ] **提交信息符合规范** - 有正确的 type 和 scope，Subject 清晰简洁

## 好提交 vs 坏提交

### ❌ 坏提交（不要这样做）
```bash
git commit -m "update"
git commit -m "fix bug"
git commit -m "完成功能"
git commit -m "WIP"
```

### ✅ 好提交（推荐做法）
```bash
git commit -m "feat(asset): 实现媒体资产管理功能"
git commit -m "fix(workflow): 修复 DAG 验证死循环"
git commit -m "docs: 更新 V1.0 架构文档"
git commit -m "refactor(ui): 优化媒体资产列表加载性能"
```

## 执行提交

请按照规范格式创建提交，并确保所有自检项都已完成。
