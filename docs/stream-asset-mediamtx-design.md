# 结合 MediaMTX 的流媒体资产设计

## 0. 设计原则

- **本平台与 MediaMTX 深度集成，作为 MediaMTX 的客户端**：创建流媒体资产时必须接入 MediaMTX，不提供“仅登记流地址、不接入 MediaMTX”的模式。所有流媒体资产均对应 MediaMTX 上的一个 path，具备完整的预览、状态、录制、点播等能力。
- **引入 MediaSource 表**：媒体源（MediaSource）与 MediaMTX path 一一对应，作为“已接入的流”的独立资源；流媒体资产通过 `source_id` 关联 MediaSource，`path` 存 path name（便于兼容与查询）。

---

## 1. MediaMTX 在本项目中的角色

本项目使用 [MediaMTX](https://github.com/bluenviron/mediamtx) 作为流媒体服务器，负责：

| 能力 | 说明 | 本项目中的使用 |
|------|------|----------------|
| **路径（Path）** | 每个“流”对应一个 path，如 `live/camera1`；path 可配置拉流地址（source）或等待推流（publisher） | 每个 MediaSource 对应一个 MediaMTX path |
| **多协议分发** | 同一 path 自动提供 RTSP、RTMP、HLS、WebRTC 等协议 | 预览 URL 由 config 的 hls_address / rtsp_address 等 + path name 拼出 |
| **录制** | 按 path 录制，可配置 recordPath、recordFormat、segmentDuration | 录制启停、录制段列表、点播 URL 通过 MediaMTX API 获取 |
| **点播（Playback）** | 按 path + 时间戳提供点播 URL | 录制完成后可生成点播 URL 或点播资产 |
| **状态** | Path 有 ready、available、online、readers 等状态 | 媒体源/流媒体资产详情可展示“流是否在线”等 |

### 1.1 已有适配器能力

- **internal/adapter/mediamtx/client.go**：MediaMTX HTTP API 客户端
- **主要接口**：AddPath、GetPathConfig、PatchPath、DeletePath、ListPathConfigs、GetPath、ListPaths、IsPathReady、EnableRecording、DisableRecording、GetRecordings、ListRecordings、DeleteRecordingSegment 等
- **configs/config.<env>.yaml**：`mediamtx` 段配置 api_address、rtsp_address、rtmp_address、hls_address、webrtc_address、playback_address、record_path、record_format、segment_duration

### 1.2 MediaMTX Path 与 URL 对应关系

- **Path name**：如 `live/camera1`，在 API 中作为路径标识。
- **拉流源（source）**：PathConfig.Source 为拉流时填 RTSP/RTMP/HLS 等 URL；推流时为 `publisher`。
- **预览 URL**（以 config 为准）：HLS `{hls_address}/{path_name}/index.m3u8`，RTSP `{rtsp_address}/{path_name}`，RTMP `{rtmp_address}/{path_name}`，WebRTC 等同理。
- **点播 URL**：Playback 服务 `{playback_address}/{path_name}/...?start=...`。

---

## 2. MediaSource（媒体源）与 MediaMTX Path

### 2.1 一一对应关系

- **MediaSource** 是系统中“已接入的流”的抽象，与 **MediaMTX path** 一一对应。
- **创建**：生成 path name（见 2.2.1）→ 调用 MediaMTX **AddPath** → 写 MediaSource 表。若 AddPath 成功但写表失败，**必须**调用 MediaMTX **DeletePath** 回滚，避免 MediaMTX 残留 path。
- **更新**：若修改拉流地址等，调用 MediaMTX **PatchPath**，再更新表。需注意 MediaMTX PatchPath 对部分字段的兼容性（按实际 API 行为处理）。
- **删除**：**禁止删除有关联流媒体资产的媒体源**。仅当该 MediaSource 下无流媒体资产（无 MediaAsset.source_id 指向该源）时，才允许删除；删除顺序：先 MediaMTX **DeletePath**，再删 MediaSource 记录。若有关联资产则返回 409 并提示先解除关联。若需“解除关联再删源”，需先提供解除关联能力并定义解除后资产的处理方式（如资产标记为无效或清空 source_id），再由业务决定是否删除源。

### 2.2 MediaSource 表结构（建议）

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uuid | 主键 |
| name | varchar(255) | 显示名称 |
| path_name | varchar(255) | MediaMTX path name，**全局唯一**，如 `live/camera1-a1b2c3d4` |
| type | varchar(20) | 源类型：pull / push |
| url | varchar(1024) | 拉流地址（type=pull 时必填；push 可为空） |
| protocol | varchar(20) | 协议：rtsp、rtmp、hls、webrtc、srt 等（可由 url scheme 解析或用户选择） |
| enabled | boolean | 是否启用（可映射为 MediaMTX path 的启用/禁用，若 MediaMTX 支持） |
| record_enabled | boolean | 是否开启录制（可选，或由 PathConfig 推断） |
| created_at | timestamp | 创建时间 |
| updated_at | timestamp | 更新时间 |

- **path_name 唯一索引**：表上对 `path_name` 建唯一索引，保证全局唯一。
- 创建/更新/删除 MediaSource 时，通过 **internal/adapter/mediamtx** 的 Client 同步 MediaMTX。

#### 2.2.1 path_name 生成规则（必读）

- **规则**：`path_name = "live/" + slug(name) + "-" + short_uuid`。其中 `slug(name)` 为名称的合法化片段（仅保留字母、数字、连字符，空格转 `-`，长度上限由实现定）；`short_uuid` 为 UUID 前 8 位或等价唯一后缀。**禁止**仅用 `live/{slug(name)}`，否则多用户/多设备同名会导致 AddPath 冲突或覆盖。
- **非法字符过滤**：生成 slug 时过滤 MediaMTX path 不支持的字符（如空格、`/`、`#`、`?` 等），仅保留字母、数字、`-`、`_`，避免 AddPath 报错。
- **示例**：`name="摄像头1"` → slug 如 `camera1`，path_name 如 `live/camera1-a1b2c3d4`。

### 2.3 MediaSource API（与 docs/api.md 对齐）

- **GET /api/v1/sources**：列表，支持 type、enabled 等过滤；可带 with_status 从 MediaMTX GetPath 拉取实时状态。
- **POST /api/v1/sources**：创建；按 2.2.1 生成 path_name，调用 MediaMTX AddPath(source=url 或 publisher)，再写 MediaSource 表；**若写表失败则调用 DeletePath 回滚**。请求体需区分 type=pull（必填 url）与 type=push（url 可为空）。
- **GET /api/v1/sources/:id**：详情；可带 with_status 返回 MediaMTX Path 状态。**type=push 时响应中必须包含 `push_url`**（推流地址），格式为 `{rtmp_address}/{path_name}`，便于 OBS 等配置推流目标。
- **PUT /api/v1/sources/:id**：更新；若 url 等变更则 PatchPath，再更新表。
- **DELETE /api/v1/sources/:id**：删除；**仅当无流媒体资产关联该源时允许**；先 MediaMTX DeletePath，再删表记录；若有关联资产则返回 409 并提示先解除关联。
- **POST /api/v1/sources/:id/enable**、**POST /api/v1/sources/:id/disable**：启用/禁用（若业务需要）。
- **GET /api/v1/sources/:id/status**：实时状态（MediaMTX GetPath）。
- **GET /api/v1/sources/:id/preview**：各协议预览 URL（由 config + path_name 拼出）；**type=push 时同时返回 `push_url`**。
- **GET /api/v1/sources/:id/preview/ready**：是否就绪（MediaMTX IsPathReady）。
- **POST /api/v1/sources/:id/record/start**、**stop**，**GET /api/v1/sources/:id/record/status**：录制控制（MediaMTX EnableRecording/DisableRecording 等）；录制参数与 config 一致（见 2.4）。
- **GET /api/v1/sources/:id/record/sessions**、**GET /api/v1/sources/:id/record/files**：录制会话/文件列表（按 MediaMTX 能力封装）。
- **GET /api/v1/sources/:id/playback?start=...**、**GET /api/v1/sources/:id/playback/segments**：点播 URL 与录制段列表；点播 URL 格式与 config 一致（见 2.4）。

### 2.4 录制与点播参数约定（与 config 一致）

- **录制**：创建 path 时若开启录制，PatchPath 传 `RecordPath`、`RecordFormat`、`RecordSegmentDuration` 等与 **configs/config.<env>.yaml** 中 `mediamtx.record_path`、`mediamtx.record_format`、`mediamtx.segment_duration` 一致，或使用 MediaMTX 全局默认，避免与设计中的录制路径不一致。录制启停即 PatchPath 的 `record` 字段。
- **点播 URL**：按 MediaMTX Playback 文档与 config 的 `playback_address` 拼出，格式固定为 `{playback_address}/{path_name}/...?start={unix_ts}`（具体路径与参数以 MediaMTX 文档为准），在 API 文档中写明，实现时统一使用该格式。

---

## 3. 流媒体资产与 MediaSource

### 3.1 关系约定

- **流媒体资产**：`type=stream`、`source_type=live` 的 MediaAsset；**必须**关联 MediaSource（`source_id` 非空），`path` 存 MediaMTX path name（与 MediaSource.path_name 一致），便于兼容与按 path 调 MediaMTX。
- **不再支持**“仅登记流地址、不接入 MediaMTX”的流媒体资产：所有流媒体资产均通过 MediaSource 对应一个 MediaMTX path，具备完整流媒体能力。

### 3.2 创建流媒体资产的两种入口

| 入口 | 说明 | 后端流程 |
|------|------|----------|
| **新建流并创建资产** | 用户输入流地址、资产名称、标签等 | 1）按 2.2.1 生成 path_name；2）MediaMTX AddPath(source=流地址)；3）创建 MediaSource(path_name, url, type=pull, …)；若步骤 3 失败则 DeletePath 回滚；4）创建 MediaAsset(type=stream, source_type=live, source_id=MediaSource.ID, path=path_name, name, tags) |
| **从已有媒体源创建资产** | 用户选择已有 MediaSource，填写资产名称、标签 | 1）校验 MediaSource 存在；2）创建 MediaAsset(type=stream, source_type=live, source_id=MediaSource.ID, path=MediaSource.path_name, name, tags) |

两种入口下，流媒体资产的 `path` 均为 path name；预览、状态、录制、点播均可通过 `source_id` 找到 MediaSource，再通过 path_name 调 MediaMTX。

**API 约定（新建流并创建资产）**：前端可调用 **POST /api/v1/assets**，请求体含 `type=stream`、`source_type=live`、`name`、`stream_url`（流地址）、`tags` 等；后端根据 `stream_url` 执行上述四步，返回的资产中 `path` 为 path name、`source_id` 为新建的 MediaSource.ID。或由前端先 **POST /api/v1/sources**（创建媒体源），再 **POST /api/v1/assets**（body 含 `source_id`、`name`、`path=source.path_name`）；本设计推荐“新建流并创建资产”一步完成，减少前端两次请求。

### 3.3 删除与一致性

- **删除流媒体资产**：仅删资产记录，不删 MediaSource 与 MediaMTX path（同一媒体源可被多个资产引用）。
- **删除媒体源**：**禁止删除有关联流媒体资产的媒体源**。仅当无 MediaAsset.source_id 指向该源时，才允许删除；删除顺序：先 MediaMTX DeletePath，再删 MediaSource。若有关联资产则返回 409，提示先删除或解除关联相关资产。

---

## 4. 流媒体资产相关 API（基于 MediaSource + MediaMTX）

以下 API 针对流媒体资产（`type=stream`、`source_id` 非空）；通过 `source_id` 取 MediaSource.path_name，再调 MediaMTX。

### 4.1 预览 URL

- **GET /api/v1/assets/:id/preview**（或 **GET /api/v1/sources/:source_id/preview** 由前端用 asset.source_id 调）
- **逻辑**：校验资产为 stream 且 source_id 存在；取 MediaSource.path_name，按 config 拼各协议预览 URL 返回。
- **响应示例**：`{ "path_name": "live/camera1", "hls_url": "...", "rtsp_url": "...", "rtmp_url": "...", "webrtc_url": "..." }`

### 4.2 流状态

- **GET /api/v1/assets/:id/stream/status**（或通过 source_id 调 **GET /api/v1/sources/:id/status**）
- **逻辑**：asset.source_id → MediaSource.path_name → MediaMTX GetPath，返回 ready、available、online、readers 等。

### 4.3 录制控制

- **POST /api/v1/assets/:id/stream/record/start**、**stop**，**GET /api/v1/assets/:id/stream/record/status**
- **逻辑**：asset.source_id → MediaSource.path_name → MediaMTX EnableRecording / DisableRecording / GetRecordings 等。

### 4.4 录制段 / 点播

- **GET /api/v1/assets/:id/stream/recordings**：asset.source_id → path_name → MediaMTX GetRecordings。
- **GET /api/v1/assets/:id/stream/playback?start=...**：返回点播 URL（playback_address + path_name + start）。

也可统一为：前端用 asset.source_id 直接调 **/api/v1/sources/:id/*** 的预览/状态/录制/点播接口，资产侧仅保留“获取预览 URL”等便捷接口（内部转发到 source）。

---

## 5. 配置、安全与错误码

### 5.1 配置

- **configs/config.<env>.yaml** 中 `mediamtx` 段：api_address 用于后端调 MediaMTX API（建议仅内网可达）；hls_address、rtsp_address、rtmp_address、playback_address 等用于拼预览/推流/点播 URL。
- **录制与点播**：record_path、record_format、segment_duration 与 2.4 节约定一致，实现时与 MediaMTX 全局或 path 级配置对齐。

### 5.2 安全与鉴权

- **鉴权**：媒体源与流媒体资产相关 API（预览、状态、录制、点播）均受现有 JWT 与权限控制；若预览/录制/点播 URL 直接暴露给前端，需考虑 MediaMTX 对外网开放时通过反向代理或鉴权参数控制访问，避免未授权拉流。
- **api_address**：建议仅内网可达；hls_address 等对外暴露时通过网关/鉴权控制。

### 5.3 错误码约定（与 internal/api/errors 统一）

| 场景 | HTTP 状态码 | 说明 |
|------|-------------|------|
| MediaMTX 不可用（Ping 失败、AddPath/GetPath 等超时或 5xx） | 503 | 提示“流媒体服务暂不可用” |
| path 已存在（AddPath 返回 path 冲突） | 409 | 提示“路径已存在”或“请重试”（path_name 含 uuid 时少见） |
| 参数错误（缺少 name、stream_url、source_id 等） | 400 | 提示具体缺失或非法字段 |
| 删除媒体源时有关联流媒体资产 | 409 | 提示“存在关联的流媒体资产，请先删除或解除关联” |
| 媒体源/资产不存在 | 404 | 与现有 ErrNotFound 一致 |

---

## 6. 实现顺序与分层

### 6.1 分层约定

- **Domain**：MediaSource 实体放在 `internal/domain/media_source.go`，不依赖 adapter/config。
- **Port**：Repository 扩展 MediaSource CRUD；若需抽象“媒体源与 MediaMTX 同步”，可定义 MediaMTXClient 或 MediaSourceSync 类端口，由 adapter 实现。
- **App**：MediaSourceService 编排 Repository 与 MediaMTX 客户端；创建/更新/删除时调用 AddPath/PatchPath/DeletePath，并实现失败回滚（写表失败则 DeletePath）。
- **API**：Sources 相关 Handler、DTO、路由；错误码与 5.3 统一（503、409、400、404）。

### 6.2 实现顺序建议

| 步骤 | 内容 |
|------|------|
| 1 | **MediaSource 实体与表**：domain（media_source.go）、repository 接口与实现、migration；path_name 唯一索引；path_name 生成规则（2.2.1）与非法字符过滤 |
| 2 | **MediaSource 与 MediaMTX 同步**：创建时 AddPath → 写表，写表失败则 DeletePath 回滚；更新时 PatchPath → 更新表；删除时校验无关联资产 → DeletePath → 删表 |
| 3 | **Sources API**：CRUD、status、preview（含 push 时 push_url）、record、playback 等（见 2.3）；错误码 5.3 |
| 4 | **创建流媒体资产**：新建流时按 2.2.1 生成 path_name → AddPath → MediaSource → MediaAsset（含回滚）；从已有源创建时仅创建 MediaAsset(source_id, path=path_name) |
| 5 | **流媒体资产扩展 API**：预览/状态/录制/点播（可复用 sources 接口或资产侧封装） |
| 6 | **前端**：流媒体接入表单“必须接入 MediaMTX”；“从已有媒体源创建”入口；媒体源管理页；流媒体资产预览/状态/录制入口；push 类型展示推流地址 |

---

## 7. 与现有文档的对应

- **docs/asset-stream-ingestion.md**：更新为“创建流媒体资产必须接入 MediaMTX”，两种入口为“新建流（创建 MediaSource + 资产）”与“从已有媒体源创建资产”；移除“仅登记流地址”。
- **docs/requirements.md**：3.1.2 中“添加资产 - 流媒体接入”改为上述原则与两种入口；明确建 MediaSource 表。
- **docs/api.md**：媒体源（Sources）章节与本设计 2.3、4 节对齐；资产创建接口在流媒体场景下可接受“流地址 + 名称”或“source_id + 名称”。

---

## 8. 设计评估（业务 / 技术 / 规范）

对当前设计在**业务合理性**、**技术合理性**、**规范符合度**及**流媒体实际场景**下的评估结论与改进建议，见 **`docs/stream-asset-mediamtx-evaluation.md`**。本设计文档已按评估建议完成优化：path_name 生成规则（2.2.1）、创建失败回滚（2.1、2.3）、删除媒体源策略（2.1、3.3）、推流地址返回（2.3）、录制与点播参数约定（2.4）、错误码约定（5.3）、分层与实现顺序（6.1、6.2）。

---

## 9. 开发准备清单（实现前对照）

实现前请确认以下项已在设计或代码中落实：

| 项 | 设计位置 | 说明 |
|----|----------|------|
| path_name 生成规则 | 2.2.1 | `live/` + slug(name) + `-` + short_uuid；过滤非法字符；禁止仅用 slug |
| path_name 唯一索引 | 2.2 | 表上对 path_name 建唯一索引 |
| 创建失败回滚 | 2.1、2.3 | AddPath 成功但写 MediaSource 失败时调用 DeletePath |
| 删除媒体源策略 | 2.1、3.3 | 禁止删除有关联流媒体资产的媒体源；有关联时返回 409 |
| 推流地址 push_url | 2.3 | type=push 时创建/详情/预览响应含 push_url（rtmp_address + path_name） |
| 录制与点播参数 | 2.4、5.1 | recordPath/recordFormat/segmentDuration、playback URL 与 config 一致 |
| 错误码 | 5.3 | 503 MediaMTX 不可用、409 冲突/有关联资产、400 参数错误、404 不存在 |
| 分层 | 6.1 | domain/media_source.go、Repository、MediaSourceService、Handler/DTO |
| 拉流与推流区分 | 2.2、2.3 | type=pull 必填 url；type=push 时 Source=publisher，url 可为空 |
