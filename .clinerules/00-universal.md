# GoyaVision 通用规则（始终生效）

## 项目目标
- GoyaVision 是基于 Go + Vue 3 的智能媒体处理平台。
- 核心理念：业务=配置、能力=插件、执行=引擎、资产驱动。

## 分层与依赖（严格遵守）
- Domain 不依赖任何其他层
- Port 可依赖 Domain
- App 仅依赖 Domain + Port（禁止依赖 Adapter）
- Adapter 实现 Port 接口，可依赖 Domain + Port
- API 可依赖 App + Port + Domain（禁止直接依赖 Adapter）

## 开发工作流要求（每次开发必须执行）
- 开始前查阅：docs/requirements.md、docs/architecture.md、docs/development-progress.md、docs/api.md、CHANGELOG.md。
- 开发中遵循本规则与对应的条件规则。
- 完成后必须更新：
  - docs/development-progress.md
  - CHANGELOG.md（[未发布] 章节，按新增/变更/修复/弃用/移除/安全分类）
  - 如 API 变更：docs/api.md
  - 如设计/需求变更：docs/architecture.md 或 docs/requirements.md
- 提交遵循 Conventional Commits：<type>(<scope>): <subject>（提交信息需使用中文描述）

## 信息完整性与提问规范
- 在执行用户请求前，请进行信息完整性检查：
  - 如果完成任务必须依赖的关键信息缺失，你必须先向用户提出澄清问题，再继续执行。
  - 如果存在多个合理解释或执行路径，你必须指出歧义并向用户询问偏好。
  - 如果继续执行会导致高风险错误或不可逆后果，你必须先确认用户意图。
- 提问规范：
  - 一次只问最少必要的问题（最多 3 个）。
  - 问题必须是具体、可操作、可回答的，不要泛问。
  - 不要重复已确认的信息。
- 禁止行为：
  - 在关键信息缺失时擅自假设。
  - 为了显得“有帮助”而继续胡乱生成结果。
- 当信息充分时，直接执行任务，不要再提问。

## 通用代码质量
- 不吞错误，必须返回或记录。
- 使用 context.Context 处理可取消操作。
- 遵循项目既有命名与结构，不引入新风格。
