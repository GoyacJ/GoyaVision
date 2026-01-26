# 贡献指南

感谢您对 GoyaVision 项目的关注！我们欢迎所有形式的贡献。

## 如何贡献

### 报告问题

如果您发现了 bug 或有功能建议，请通过 [GitHub Issues](https://github.com/GoyaDo/GoyaVision/issues) 提交。

在提交 issue 前，请：
- 搜索现有 issue，避免重复
- 使用清晰的标题和描述
- 提供复现步骤（如适用）
- 包含环境信息（Go 版本、操作系统等）

### 提交代码

1. **Fork 项目**
   ```bash
   git clone https://github.com/GoyaDo/GoyaVision.git
   cd GoyaVision
   ```

2. **创建功能分支**
   ```bash
   git checkout -b feature/your-feature-name
   # 或
   git checkout -b fix/your-bug-fix
   ```

3. **开发与测试**
   - 遵循项目代码规范（见 `.cursor/rules/goyavision.mdc`）
   - 确保代码通过测试
   - 更新相关文档

4. **提交代码**
   ```bash
   git add .
   git commit -m "feat: 添加新功能描述"
   ```
   
   提交信息格式遵循 [Conventional Commits](https://www.conventionalcommits.org/)：
   - `feat:` 新功能
   - `fix:` Bug 修复
   - `docs:` 文档更新
   - `style:` 代码格式（不影响功能）
   - `refactor:` 重构
   - `test:` 测试相关
   - `chore:` 构建/工具相关

5. **推送并创建 Pull Request**
   ```bash
   git push origin feature/your-feature-name
   ```
   
   然后在 GitHub 上创建 Pull Request。

## 代码规范

### Go 代码风格

- 遵循 `gofmt`/`goimports` 标准格式
- 使用有意义的变量和函数名
- 添加必要的注释，特别是公共 API
- 不写行尾注释
- 错误处理：不吞掉错误，区分业务错误与基础设施错误

### 项目结构

遵循分层架构：
- `domain/`: 领域实体，纯业务逻辑
- `port/`: 接口定义
- `app/`: 应用服务，编排 domain 与 port
- `adapter/`: 基础设施实现
- `api/`: HTTP 层

详细规则见 `.cursor/rules/goyavision.mdc`。

### 测试

- 为新功能添加单元测试
- 确保所有测试通过
- 保持测试覆盖率

## Pull Request 流程

1. **确保 PR 描述清晰**
   - 说明变更的目的
   - 关联相关 issue（如适用）
   - 列出主要变更点

2. **代码审查**
   - 维护者会审查您的 PR
   - 根据反馈进行修改
   - 保持讨论友好和专业

3. **合并**
   - PR 通过审查后会被合并
   - 感谢您的贡献！

## 开发环境设置

### 前置要求

- Go 1.22+
- PostgreSQL
- FFmpeg（PATH 或配置路径）
- 可选：MediaMTX（用于预览功能）

### 本地开发

```bash
# 克隆项目
git clone https://github.com/GoyaDo/GoyaVision.git
cd GoyaVision

# 安装依赖
go mod download

# 配置数据库
createdb goyavision
# 或通过环境变量配置
export GOYAVISION_DB_DSN="host=localhost user=goyavision password=goyavision dbname=goyavision port=5432 sslmode=disable"

# 运行
go run ./cmd/server
```

## 文档贡献

文档改进同样重要：
- 修复文档中的错误
- 改进示例和说明
- 添加使用案例
- 翻译文档（欢迎！）

## 行为准则

请遵循我们的 [行为准则](CODE_OF_CONDUCT.md)，保持社区友好和包容。

## 获取帮助

- 查看 [文档](README.md)
- 搜索 [现有 Issues](https://github.com/GoyaDo/GoyaVision/issues)
- 创建新的 Issue 提问

再次感谢您的贡献！
