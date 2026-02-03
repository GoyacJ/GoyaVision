# 添加资产 - 流媒体接入 设计与实现

## 1. 目标与原则

在「添加资产」流程中增加**流媒体接入**能力，使用户可以将实时流（RTSP/RTMP/HLS 等）登记为媒体资产，供工作流、预览、录制、点播等使用。

**设计原则**（详见 **`docs/stream-asset-mediamtx-design.md`**）：

- **本平台与 MediaMTX 深度集成，作为 MediaMTX 的客户端**：创建流媒体资产时必须接入 MediaMTX，不提供“仅登记流地址、不接入 MediaMTX”的模式；所有流媒体资产均对应 MediaMTX 上的一个 path，具备完整流媒体能力。
- **引入 MediaSource 表**：媒体源（MediaSource）与 MediaMTX path 一一对应；流媒体资产通过 `source_id` 关联 MediaSource，`path` 存 path name。

## 2. 概念

- **流媒体资产**：`type=stream`、`source_type=live` 的 MediaAsset；**必须**关联 MediaSource（`source_id` 非空），`path` 存 MediaMTX path name（与 MediaSource.path_name 一致）。
- **媒体源（MediaSource）**：与 MediaMTX path 一一对应的实体；创建/更新/删除 MediaSource 时同步 MediaMTX AddPath / PatchPath / DeletePath；API 见 `docs/api.md` 媒体源（Sources）章节。
- **MediaMTX path**：MediaMTX 中每个“流”对应一个 path（如 `live/camera1`），可配置拉流地址、录制、多协议预览等。

## 3. 接入方式（两种入口）

### 3.1 新建流并创建资产（主入口）

| 项目 | 说明 |
|------|------|
| 入口 | 添加资产 → Tab「流媒体接入」→ 输入流地址、资产名称、标签 |
| 必填 | 资产名称、流地址 |
| 选填 | 标签 |
| 后端流程 | 1）按 **path_name 生成规则**（见 stream-asset-mediamtx-design 2.2.1）生成 path name；2）MediaMTX **AddPath**(name=path name, source=流地址)；3）创建 **MediaSource**(path_name, url, type=pull, …)；若写表失败则 **DeletePath 回滚**；4）创建 **MediaAsset**(type=stream, source_type=live, source_id=MediaSource.ID, path=path name, name, tags) |
| 校验 | name、流地址非空；path name 全局唯一（表唯一索引 + 生成规则含 short_uuid）；MediaMTX AddPath 成功后再写库 |

适用于：用户首次接入该流，需要完整预览/状态/录制/点播能力。

### 3.2 从已有媒体源创建资产

| 项目 | 说明 |
|------|------|
| 入口 | 添加资产 → 流媒体接入 → 选择「从已有媒体源」→ 选择 MediaSource、填写资产名称、标签 |
| 必填 | 资产名称、媒体源（source_id） |
| 选填 | 标签 |
| 后端流程 | 1）校验 MediaSource 存在；2）创建 **MediaAsset**(type=stream, source_type=live, source_id=MediaSource.ID, path=MediaSource.path_name, name, tags) |
| 校验 | name、source_id 非空；MediaSource 存在 |

适用于：同一路流需要多个资产记录（如不同标签、不同用途），或从媒体源管理页已有源创建资产。

## 4. 前端实现要点

- 添加资产对话框：Tab「流媒体接入」下提供两种子入口：
  - **新建流**：资产名称、流地址、标签；提交后调用“创建流媒体资产（新建流）”接口（如 POST /api/v1/assets 且 body 含 stream_url 等，或专用接口）。
  - **从已有媒体源**：媒体源下拉（GET /api/v1/sources）、资产名称、标签；提交后调用“创建流媒体资产（source_id）”接口（POST /api/v1/assets，body 含 source_id、name、type=stream、source_type=live、path=由后端或前端从 source 取 path_name）。
- 流地址提示：支持 RTSP、RTMP、HLS 等协议；创建后可在资产详情/媒体源详情中查看预览、状态、录制、点播。

## 5. 后端实现要点

- **MediaSource 实体与表**：见 `docs/stream-asset-mediamtx-design.md` 2.2、2.2.1 节；path_name 唯一索引；path_name 生成规则（slug + short_uuid，非法字符过滤）；创建/更新/删除 MediaSource 时同步 MediaMTX。
- **创建流媒体资产（新建流）**：按 2.2.1 生成 path name → MediaMTX AddPath → 创建 MediaSource（**写表失败则 DeletePath 回滚**）→ 创建 MediaAsset(source_id, path=path_name)。
- **创建流媒体资产（从已有源）**：校验 source_id → 创建 MediaAsset(source_id, path=MediaSource.path_name)。
- **删除流媒体资产**：仅删资产记录，不删 MediaSource 与 MediaMTX path。
- **删除媒体源**：**禁止删除有关联流媒体资产的媒体源**；仅当无 MediaAsset.source_id 指向该源时允许删除；顺序：MediaMTX DeletePath → 删 MediaSource；若有关联资产返回 409。
- **推流（push）**：创建/详情/预览 API 中 type=push 时返回 **push_url**（rtmp_address + path_name），便于 OBS 等配置。
- **错误码**：与设计 5.3 统一（503 MediaMTX 不可用、409 path 冲突或有关联资产、400 参数错误、404 不存在）。

## 6. 文档与进度

- **总体设计（结合 MediaMTX + MediaSource）**：`docs/stream-asset-mediamtx-design.md`（原则、MediaSource 表、path_name 规则 2.2.1、失败回滚、删除策略、录制/点播约定、错误码、分层与实现顺序、**开发准备清单 第 9 节**）。
- 需求与设计：`docs/requirements.md` 3.1.2 节、本文档。
- API：媒体源 CRUD 与预览/状态/录制/点播见 `docs/api.md` 媒体源（Sources）章节；资产创建在流媒体场景下见本设计 3.1、3.2。
- 开发进度：在 `docs/development-progress.md` 中更新「媒体源管理」「媒体资产管理」「添加资产 - 流媒体接入」状态。
