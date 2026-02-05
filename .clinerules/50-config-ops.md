---
paths:
  - "configs/**"
  - "config/**"
  - "docker-compose.yml"
  - "Makefile"
  - "scripts/**"
  - "cmd/**"
---

# 配置与运维规则

## 配置
- 默认配置位于 configs/config.<env>.yaml，避免硬编码敏感信息。
- 支持 GOYAVISION_* 环境变量覆盖，变更需同步文档。

## 部署与脚本
- Docker Compose 变更需同步 docs/DEPLOYMENT.md。
- Makefile 或脚本变更需保持跨平台兼容（bash + macOS/zsh）。

## 命令入口
- cmd/ 仅作为应用入口或工具入口，保持简洁。
- 入口中依赖注入必须符合分层规则。
