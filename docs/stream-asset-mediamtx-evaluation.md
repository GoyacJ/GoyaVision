# 流媒体资产设计评估：业务、技术与规范

本文档从**业务合理性**、**技术合理性**、**规范符合度**三方面评估当前「结合 MediaMTX 的流媒体资产设计」，并结合实际流媒体场景给出结论与改进建议。设计文档见 `docs/stream-asset-mediamtx-design.md`、`docs/asset-stream-ingestion.md`。

---

## 一、业务合理性评估

### 1.1 典型流媒体业务场景对照

| 场景 | 用户诉求 | 当前设计是否覆盖 | 说明 |
|------|----------|------------------|------|
| **摄像头/设备拉流** | 输入 RTSP 地址，在平台内预览、录制、供工作流使用 | ✅ 覆盖 | 「新建流并创建资产」：输入流地址 → MediaMTX AddPath(pull) → MediaSource + MediaAsset；预览/录制/点播走 MediaMTX |
| **OBS/推流端推流** | 在 OBS 中推 RTMP 到平台，平台内预览、录制 | ✅ 覆盖 | MediaSource type=push，PathConfig.Source= publisher；需在创建媒体源时区分 pull/push，path 创建后告知用户推流地址（rtmp_address/path_name） |
| **同一路流多用途** | 同一摄像头既用于「监控」资产又用于「分析」资产（不同标签/工作流） | ✅ 覆盖 | 「从已有媒体源创建资产」：同一 MediaSource 可对应多个 MediaAsset，通过标签/名称区分 |
| **仅需登记流地址、不经过平台转发** | 工作流只引用外部 HLS URL，不需要平台预览/录制 | ❌ 不再支持 | 设计已明确：所有流媒体资产必须接入 MediaMTX。若业务确有「仅登记 URL」需求，需单独评估是否用普通 URL 资产（type=video + path=url）或放宽为可选 |
| **录制后点播/回放** | 对某路流录制后，按时间点播或生成点播资产 | ✅ 覆盖 | 录制启停、录制段列表、点播 URL 均通过 MediaMTX API；设计中有 record/playback 相关 API |
| **工作流以流为输入** | 工作流节点将流媒体资产作为输入，算子拉流或拉 HLS | ✅ 覆盖 | 流媒体资产具备 source_id、path（path name）；算子可通过预览 URL（由 path name 拼出）拉流，或通过资产 ID 解析出 path name 再拼 URL |

**结论（业务）**：典型场景中「拉流接入、推流接入、同源多资产、录制与点播、工作流输入」均能支撑。唯一被排除的是「仅登记流地址、不经过 MediaMTX」；在「平台与 MediaMTX 深度集成、作为其客户端」的前提下，该取舍**合理**。若后续确有仅登记 URL 的需求，可单独用非 stream 类型或扩展策略。

### 1.2 业务边界与策略需明确之处

| 点 | 当前设计 | 建议 |
|----|----------|------|
| **删除媒体源时有关联资产** | 文档写「需先解除关联或禁止删除（策略由业务定）」 | 建议明确：**禁止删除有关联流媒体资产的媒体源**，或提供「解除关联」使资产变为无效（需定义资产侧表现）；避免误删导致资产悬空 |
| **删除流媒体资产** | 仅删资产，不删 MediaSource/MediaMTX path | 合理；同一源可对应多资产。若希望「删掉最后一个关联资产时顺带删源」可作可选策略，不建议默认 |
| **path_name 唯一性** | 建议 `live/{slug(name)}` 或 `live/{short_uuid}` | 见下文技术部分：**必须保证全局唯一**，推荐带 uuid 或唯一后缀，避免重名冲突 |
| **同一物理流重复接入** | 未约束「同一 stream_url 是否允许多个 MediaSource」 | 业务上可允许（如不同命名、不同用途）；若需「同一 URL 只建一个源」可在创建前按 url 去重，由产品决定 |

---

## 二、技术合理性评估

### 2.1 MediaSource 与 MediaMTX 的同步

| 点 | 评估 | 建议 |
|----|------|------|
| **创建顺序** | 设计：生成 path_name → AddPath → 写 MediaSource → 创建 MediaAsset | 正确：先 MediaMTX 再 DB，避免 DB 有记录而 MediaMTX 无 path |
| **失败回滚** | AddPath 成功、写 MediaSource 失败时，MediaMTX 会残留 path | 建议：写 MediaSource 失败时调用 DeletePath 回滚；或采用「先写 MediaSource 状态=creating，AddPath 成功后更新为 active，失败则删 DB 并 DeletePath」等两阶段策略，减少残留 |
| **更新 MediaSource** | 修改 url 等 → PatchPath → 更新表 | 合理；需注意 PatchPath 部分字段时的兼容性（MediaMTX 部分字段可能不可 patch，需按实际 API 行为处理） |
| **删除顺序** | DeletePath → 删 MediaSource 记录 | 正确；若先删记录再 DeletePath，中间状态不一致且无法拿到 path_name 调 DeletePath |

### 2.2 path_name 生成与唯一性

| 点 | 评估 | 建议 |
|----|------|------|
| **仅 slug(name)** | 多用户/多设备同名（如「摄像头1」）会冲突 | **不推荐**仅用 slug；MediaMTX path name 全局唯一，冲突会导致 AddPath 失败或覆盖 |
| **推荐** | 使用带唯一后缀的 path name | **path_name = `live/` + slug(name) + `-` + short_uuid**，或 **`live/` + uuid 前 8 位**；既可读又唯一。需在 MediaSource 表对 path_name 建唯一索引 |
| **MediaMTX path 规范** | path name 通常支持字母数字、斜杠、连字符等 | 生成时过滤非法字符（如空格、特殊符），避免 AddPath 报错 |

### 2.3 拉流（pull）与推流（push）

| 点 | 评估 | 建议 |
|----|------|------|
| **pull** | PathConfig.Source = 流地址 URL；MediaSource.type=pull, url 必填 | 设计已区分；实现时 AddPath 传 Source=url |
| **push** | PathConfig.Source = `publisher`；MediaMTX 等待推流 | MediaSource.type=push 时 url 可为空；创建 path 后需把「推流地址」（如 rtmp_address/path_name）返回给用户，便于 OBS 等配置 |
| **protocol 字段** | 用于标识 rtsp/rtmp/hls 等 | 可由 url 解析（scheme 或后缀）或用户选择，便于前端展示与筛选 |

### 2.4 录制与点播

| 点 | 评估 | 建议 |
|----|------|------|
| **录制路径** | MediaMTX 使用 recordPath（含 %path 等变量） | 创建 path 时若需开启录制，PatchPath 传 RecordPath 等与 config 一致，或使用全局默认；避免与设计中的 record_path 不一致 |
| **录制启停** | EnableRecording / DisableRecording 按 path name 调用 | 设计正确；实现时注意 MediaMTX 的 record 是 path 级配置，启停即 PatchPath 的 record 字段 |
| **点播 URL 格式** | 依赖 MediaMTX Playback 文档 | 实现时按 MediaMTX 实际 playback API 拼 URL（path_name + start 等参数），并在文档中固定格式 |

### 2.5 事务与幂等

| 点 | 评估 | 建议 |
|----|------|------|
| **新建流并创建资产** | 四步：path_name → AddPath → MediaSource → MediaAsset | 若前端重试，可能重复 AddPath（同一 path_name 可能报错或覆盖）；建议：path_name 含 uuid，或先查 MediaSource 是否已存在（如按 stream_url 去重），再做幂等处理 |
| **从已有源创建资产** | 仅创建 MediaAsset | 天然幂等（同一 source_id + name 可视为允许重复，由业务定） |

**结论（技术）**：整体流程与 MediaMTX 能力匹配；需在实现时补齐**失败回滚**、**path_name 唯一性**、**push 场景推流地址返回**、**录制配置与 Playback URL 格式**等细节。

---

## 三、规范符合度评估

### 3.1 项目分层与依赖

| 点 | 评估 | 说明 |
|----|------|------|
| **Domain** | MediaSource 为新增实体，属领域层 | 建议放在 `internal/domain/media_source.go`，不依赖 adapter/config |
| **Port** | Repository 需扩展 MediaSource CRUD；可选 MediaMTXClient 端口 | 若已有 MediaMTX 客户端注入，可继续用 adapter 直接调；若希望抽象「媒体源与外部服务同步」，可定义 MediaSourceSync 类端口由 adapter 实现 |
| **App** | MediaSourceService 编排 Repository + MediaMTX 客户端 | 创建/更新/删除时调用 adapter 的 MediaMTX；符合「App 不直接依赖 adapter 实现」可透过 port 注入 |
| **API** | Sources 相关 Handler、DTO、路由 | 符合现有 API 层风格；与 docs/api.md 对齐即可 |

### 3.2 REST 与 API 约定

| 点 | 评估 | 建议 |
|----|------|------|
| **资源命名** | /api/v1/sources、/api/v1/assets | 符合 REST 资源命名；sources 与 assets 为不同资源，关系通过 source_id 表达，合理 |
| **创建流媒体资产** | POST /api/v1/assets 含 stream_url 或 source_id | 同一资源两种语义（新建流 vs 从已有源）可通过 body 区分，符合常规；若希望更显式，可增加子资源如 POST /api/v1/assets/from-stream（仅建议，非必须） |
| **错误码** | 需定义 MediaMTX 不可用、AddPath 失败、path_name 冲突等 | 建议 503（MediaMTX 不可用）、409（path 已存在）、400（参数错误），与现有 api/errors 统一 |

### 3.3 安全与鉴权

| 点 | 评估 | 建议 |
|----|------|------|
| **预览/录制/点播 URL** | 由后端按 config 拼出，经 API 返回 | 需保证这些 API 受 JWT/权限控制；若 URL 直接暴露给前端，需考虑 MediaMTX 是否对外网开放、是否经反向代理鉴权或带 token，避免未授权拉流 |
| **MediaMTX API** | 后端调 api_address（内网） | 建议 api_address 仅内网可达；hls_address 等对外暴露时通过网关/鉴权控制 |

**结论（规范）**：设计与项目分层、REST 习惯一致；实现时注意错误码统一与预览/录制 URL 的鉴权与暴露范围。

---

## 四、流媒体场景下的设计合理性小结

### 4.1 与「实际场景」的匹配度

- **拉流接入**：用户输入 RTSP/RTMP/HLS 地址 → 平台作为 MediaMTX 客户端拉流并多协议分发 → 符合「平台即 MediaMTX 客户端」的定位，**合理**。
- **推流接入**：创建 push 类型媒体源 → MediaMTX path 为 publisher → 用户向 rtmp_address/path_name 推流；设计需在文档/API 中明确返回「推流地址」，**合理**，需补实现细节。
- **同源多资产**：通过「从已有媒体源创建资产」实现，**合理**。
- **录制与点播**：依赖 MediaMTX 的录制与 Playback 能力，API 按 path 封装，**合理**；录制路径、段时长、点播 URL 格式需与 config 和 MediaMTX 文档一致。
- **工作流消费**：流媒体资产带 source_id 与 path（path name），工作流引擎或算子可通过「资产 → MediaSource → path_name → 预览 URL」拉流，**合理**。

### 4.2 与「业务」的匹配度

- **必须接入 MediaMTX**：强化了「所有流媒体能力统一经平台与 MediaMTX」的边界，避免「部分流可预览/录制、部分仅登记」的混合态，运维与权限模型更清晰，**业务上合理**。
- **MediaSource 表**：将「物理流」与「资产」解耦，便于独立管理源（启用/禁用、修改拉流地址）、统计与权限按「源」维度控制，**业务上合理**。

### 4.3 与「技术」的匹配度

- **MediaMTX 能力**：Path、AddPath/PatchPath/DeletePath、GetPath、录制、GetRecordings、Playback 等与设计中的 Sources/资产 API 一一对应，**技术可行**。
- **一致性**：创建/更新/删除 MediaSource 与 MediaMTX 的同步顺序正确；需在实现中落实失败回滚与 path_name 唯一性，**技术合理**。

---

## 五、总体结论与改进建议

### 5.1 总体结论

| 维度 | 结论 |
|------|------|
| **业务合理性** | **合理**。典型场景（拉流/推流、同源多资产、录制/点播、工作流输入）均能支撑；「必须接入 MediaMTX」的取舍在深度集成前提下成立。 |
| **技术合理性** | **合理**。与 MediaMTX 能力匹配，同步顺序正确；需在实现中补齐 path_name 唯一性、失败回滚、push 推流地址返回、录制/点播参数与 URL 格式。 |
| **规范符合度** | **符合**。与项目分层、REST、安全预期一致；需统一错误码与预览/录制 URL 的鉴权与暴露策略。 |

### 5.2 建议在设计文档或实现中补充的内容

1. **path_name 生成规则**：明确 `path_name = live/{slug}-{short_uuid}` 或等价规则，并对 path_name 建唯一索引；过滤非法字符。
2. **创建失败回滚**：AddPath 成功但写 MediaSource 失败时，调用 DeletePath 清理 MediaMTX；或采用两阶段（creating → active）减少残留。
3. **删除媒体源策略**：明确「禁止删除有关联流媒体资产的媒体源」，或定义解除关联后资产的处理方式。
4. **推流（push）**：在 MediaSource 创建响应或 GET 详情中返回「推流地址」（如 rtmp_address + path_name），便于 OBS 等配置。
5. **录制与点播**：在配置与 API 文档中固定 recordPath、segmentDuration、playback URL 格式，与 MediaMTX 与 config 一致。
6. **错误码**：MediaMTX 不可用、path 冲突、参数错误等与现有 api/errors 统一（如 503、409、400）。

按上述建议微调后，当前设计在业务、技术与规范上均可视为**合理且可落地**。
