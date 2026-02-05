# Cursor 配置更新日志

## 2024-01-XX - 符合官方规范更新

### 更新内容

#### ✅ Skills（技能）

- **修正 frontmatter 字段**：将 `skill` 字段改为 `name`，符合官方规范
  - `development-workflow/SKILL.md`
  - `goyavision-context/SKILL.md`

#### ✅ Hooks（钩子）

- **修正脚本路径**：将相对路径 `hooks/` 改为项目根目录相对路径 `.cursor/hooks/`
  - 更新 `hooks.json` 中的 `stop` hook 配置
  - 确保脚本路径从项目根目录正确解析
- **重新实现 stop hook**：完全符合官方规范
  - ✅ 从 stdin 读取 JSON 输入（包含 status, loop_count 等字段）
  - ✅ 输出 JSON 格式到 stdout（包含 followup_message）
  - ✅ 检查 loop_count 限制（防止无限循环）
  - ✅ 支持 jq 和 fallback 解析（提高兼容性）
  - ✅ 添加 timeout 和 loop_limit 配置
  - ✅ 使用 followup_message 自动触发检查清单提醒

#### ✅ Commands（命令）

- **创建 Cursor Commands**：在 `.cursor/commands/` 目录下创建命令文件
  - `dev-start.md` - 开始开发检查清单
  - `dev-done.md` - 完成开发检查清单
  - `commit.md` - Git 提交规范
  - `context.md` - 项目上下文
  - `api-doc.md` - API 文档更新指南
  - `progress.md` - 开发进度更新指南

#### ✅ Rules（规则）

- **优化 frontmatter**：为 `frontend-components.mdc` 添加 frontmatter
  - 添加 `description` 字段
  - 添加 `globs` 配置，仅在前端文件时应用
  - 设置 `alwaysApply: false`

### 符合的官方规范

1. **Skills 规范**
   - ✅ 使用 `name` 字段（而非 `skill`）
   - ✅ 使用 `description` 字段
   - ✅ 文件位于 `.cursor/skills/<skill-name>/SKILL.md`

2. **Hooks 规范**
   - ✅ 使用 `hooks.json` 配置文件
   - ✅ 脚本路径使用项目根目录相对路径（`.cursor/hooks/`）
   - ✅ 脚本文件具有可执行权限
   - ✅ 从 stdin 读取 JSON 输入
   - ✅ 输出 JSON 格式到 stdout
   - ✅ 正确处理退出码（0=成功，2=阻止）
   - ✅ 支持 timeout 和 loop_limit 配置

3. **Commands 规范**
   - ✅ 命令文件位于 `.cursor/commands/` 目录
   - ✅ 使用 Markdown 格式
   - ✅ 文件名简洁明了

4. **Rules 规范**
   - ✅ 使用 frontmatter 定义元数据
   - ✅ `alwaysApply` 用于始终应用的规则
   - ✅ `globs` 用于文件模式匹配
   - ✅ `description` 用于规则描述

### 兼容性

- ✅ 保持与 `.claude/commands/` 的兼容性（Claude Code）
- ✅ Cursor 优先使用 `.cursor/` 目录下的配置
- ✅ 两个工具共享相同的开发规范

### 参考文档

- [Cursor Rules 文档](https://cursor.com/cn/docs/context/rules)
- [Cursor Skills 文档](https://cursor.com/cn/docs/context/skills)
- [Cursor Commands 文档](https://cursor.com/cn/docs/context/commands)
- [Cursor Hooks 文档](https://cursor.com/cn/docs/agent/hooks)
