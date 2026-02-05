# dev-start

开始开发前的完整检查清单。在开始任何新功能或 Bug 修复前执行。

## 执行步骤

### 1. 按顺序阅读以下文档

按照以下顺序阅读文档，了解项目当前状态和设计：

1. **开发进度** (`docs/development-progress.md`)
   - 查看功能状态（✅/🚧/⏸️），避免重复工作
   - 确认当前迭代优先级

2. **变更日志** (`CHANGELOG.md`)
   - 重点查看 **[未发布]** 章节
   - 了解最近的改动和影响

3. **需求文档** (`docs/requirements.md`)
   - 确认功能需求和验收标准
   - 理解业务规则

4. **架构文档** (`docs/architecture.md`)
   - 理解分层架构（Domain/Port/App/Adapter/API）
   - 遵循依赖规则

5. **API 文档** (`docs/api.md`)
   - 查看已有端点，避免重复
   - 了解请求/响应格式约定

### 2. 核心规则

遵循以下架构约束：

**分层依赖规则（严格遵守）**：
- Domain → 无外部依赖
- Port → 可依赖 Domain
- App → 可依赖 Domain + Port（**禁止**依赖 Adapter）
- Adapter → 实现 Port 接口，可依赖 Domain
- API → 可依赖 App + Port + Domain（**禁止**直接依赖 Adapter）

**算子标准协议**：
- 输入：`{"asset_id": "uuid", "params": {}}`
- 输出：`{"output_assets": [], "results": [], "timeline": [], "diagnostics": {}}`

**代码风格**：
- Go: 使用 `gofmt`/`goimports`，无行尾注释
- 前端: TypeScript strict，Vue 3 Composition API
- 错误处理: 使用 `internal/api/errors.go` 定义的错误类型

### 3. 快速参考

**默认凭证**: admin / admin123
**服务端口**: 8080 (API), 5432 (DB), 8554 (RTSP), 8888 (HLS)
**构建命令**: `make build`, `make build-web`, `make build-all`
**配置文件**: `configs/config.<env>.yaml`，环境变量前缀 `GOYAVISION_*`

## 开始开发

现在可以开始实现功能了。记住：
- 实现前先读取相关文件
- 遵循分层架构
- 使用 DTO，不直接暴露 Domain 实体
- 完成后使用 `/dev-done` 更新文档
