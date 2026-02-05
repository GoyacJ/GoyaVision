---
name: goyavision-context
description: 获取 GoyaVision 项目架构、核心概念、API 端点与开发状态的完整上下文。用于新成员上手、修改核心模块或需要整体背景时。
---

# GoyaVision 项目上下文

提供架构分层、核心实体、API 端点与开发状态的速查信息。

## 主要内容
- 核心数据流：MediaSource → MediaAsset → Operator → Workflow → Task → Artifact。
- 分层架构：Domain / Port / App / Adapter / API。
- 关键实体与废弃概念。
- API 端点与配置项。
- 当前开发状态与里程碑。

## 使用方式
- 当你需要整体背景或不确定应修改哪一层时先调用此技能。
- 若仅阅读单文件，使用 read_file 工具即可。
