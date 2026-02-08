# GoyaVision 文档总览

本文档说明 `docs/` 目录的阅读顺序与维护原则。当前项目已进入“统一上下文 + 算法库 + Agent 工程化”重构阶段，核心事实以本目录文档为准。

## 优先阅读（核心文档）

1. `requirements.md`：需求范围与核心业务对象。
2. `architecture.md`：Clean Architecture、CQRS、执行链路分层。
3. `api.md`：当前 REST API 契约（含算法库/上下文/Agent）。
4. `context-agent-algorithm-refactor-plan.md`：本轮重构的完整落地方案与阶段计划。
5. `development-progress.md`：最新实现状态、里程碑、待办项。
6. `DEPLOYMENT.md`：部署与运行说明。

## 专题文档

- `operator-redesign.md`：算子系统重设计
- `operator-redesign-stage-report-2026-02-06.md`：算子重构阶段报告
- `stream-asset-mediamtx-design.md`：流媒体资产方案
- `ui-design.md`：前端 UI 规范
- `refactoring-plan.md` / `refactoring-multitenancy-design.md`：阶段性重构设计
- `business-retrospective-2026-02-08.md`：业务复盘结论（价值/实用/可用/逻辑）
- `business-architecture-and-flow.md`：业务架构、模块关系与业务流转梳理
- `business-review-executive-summary.md`：面向管理层的复盘决策版
- `business-review-engineering-action-plan.md`：面向研发团队的落地执行版

## 维护约定

- 代码功能变更后，至少同步更新：`development-progress.md` 与相关专题文档。
- API 路由或数据结构变更后，必须同步更新：`api.md`。
- 架构边界或模块职责变更后，必须同步更新：`architecture.md`。
- 阶段设计文档可保留，但“当前事实”必须收敛回核心文档。

最后更新：2026-02-08
