# GoyaVision 文档总览

本文档用于说明 `docs/` 目录的当前结构、阅读顺序与维护约定。V1.0 架构重构后，所有核心逻辑已收敛至本目录的四个核心文档中。

## 一、优先阅读（核心文档）

1. `requirements.md`：需求边界、核心业务概念（算子版本化、DAG 工作流）。
2. `architecture.md`：Clean Architecture 分层架构、CQRS 模式、执行器注册表设计。
3. `api.md`：对外 API 契约、多租户字段说明与 SSE 实时进度约定。
4. `development-progress.md`：当前阶段进展、技术债务与里程碑计划。
5. `CHANGELOG.md`（仓库根目录）：版本变更历史。

## 二、专题设计文档

- `operator-redesign.md`：算子系统重设计深度方案（版本化、多模式）。
- `operator-redesign-stage-report-2026-02-06.md`：算子系统迁移进度与审计。
- `stream-asset-mediamtx-design.md`：媒体源与 MediaMTX 深度集成方案。
- `ui-design.md`：前端 UI 规范与 Material Design 3 实践。
- `DEPLOYMENT.md`：基于 Docker Compose 的环境部署指南。

## 三、维护约定

- **同步更新**：代码发生功能改动或修复 BUG 后，必须同步更新 `development-progress.md` 和 `CHANGELOG.md`。
- **设计对齐**：若涉及 API 契约变更或系统架构调整，必须同步修正 `api.md` 和 `architecture.md`。
- **收敛原则**：所有已落地的阶段性设计文档应作为附录保留，但核心事实应收敛到上述核心文档中，确保信息源唯一。

最后更新：2026-02-08
