# GoyaVision 文档总览

本文档用于说明 `docs/` 目录的当前结构、阅读顺序与维护约定。

## 一、优先阅读（核心文档）

1. `requirements.md`：需求边界与核心业务概念。
2. `architecture.md`：分层架构、依赖规则、关键模块设计。
3. `api.md`：对外 API 契约与请求/响应约定。
4. `development-progress.md`：当前阶段进展、风险与下一步计划。
5. `CHANGELOG.md`（仓库根目录）：版本变更记录。

## 二、专题设计文档

- `operator-redesign.md`：算子重设计目标方案。
- `operator-redesign-stage-report-2026-02-06.md`：算子重设计审计与阶段评估。
- `stream-asset-mediamtx-design.md`：流媒体资产与 MediaMTX 集成设计。
- `frontend-refactor-design.md`：前端重构设计方案。
- `design-tokens-migration-guide.md`：设计令牌迁移指南。
- `ui-design.md`：UI 设计规范。
- `DEPLOYMENT.md`：部署文档。

## 三、文档维护约定

- 发生功能改动时，至少同步更新：
  - `development-progress.md`
  - `CHANGELOG.md`
  - 若接口变更：`api.md`
  - 若设计/需求变更：`architecture.md` 或 `requirements.md`
- 已过时的阶段性总结文档应及时收敛到核心文档，避免多份文档描述同一事实。

最后更新：2026-02-07