-- 目的：在算子重设计兼容迁移期，按 active_version 回填 operators 兼容字段，
-- 并在确认业务链路已切换至 active_version 后删除旧兼容列。
-- 说明：当前项目主要通过 AutoMigrate 管理结构，本脚本用于数据治理与手工执行。

BEGIN;

-- 1) 基于 active_version 回填 version / schema / config
UPDATE operators o
SET
  version = ov.version,
  input_schema = ov.input_schema,
  output_spec = ov.output_spec,
  config = ov.config
FROM operator_versions ov
WHERE o.active_version_id = ov.id;

-- 2) 仅当 active_version 为 HTTP 模式时回填 endpoint/method
UPDATE operators o
SET
  endpoint = COALESCE(ov.exec_config->'http'->>'endpoint', ''),
  method = COALESCE(NULLIF(UPPER(ov.exec_config->'http'->>'method'), ''), 'POST')
FROM operator_versions ov
WHERE o.active_version_id = ov.id
  AND ov.exec_mode = 'http';

-- 3) 非 HTTP 模式清空 endpoint/method（method 保持默认 POST，避免空值）
UPDATE operators o
SET
  endpoint = '',
  method = 'POST'
FROM operator_versions ov
WHERE o.active_version_id = ov.id
  AND ov.exec_mode <> 'http';

-- 4) 删除 operators 旧兼容执行字段（正式收口）
ALTER TABLE operators
  DROP COLUMN IF EXISTS version,
  DROP COLUMN IF EXISTS endpoint,
  DROP COLUMN IF EXISTS method,
  DROP COLUMN IF EXISTS input_schema,
  DROP COLUMN IF EXISTS output_spec,
  DROP COLUMN IF EXISTS config,
  DROP COLUMN IF EXISTS is_builtin;

COMMIT;
