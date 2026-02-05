# 开发开始检查（GoyaVision）

该工作流用于在开始新需求或修复前快速完成项目必读检查。

## 1. 阅读当前进度与变更

我将读取以下文件，确认当前状态与近期变更：

```bash
sed -n '1,120p' docs/development-progress.md
```

```bash
sed -n '1,120p' CHANGELOG.md
```

## 2. 阅读需求与架构

```bash
sed -n '1,120p' docs/requirements.md
```

```bash
sed -n '1,120p' docs/architecture.md
```

## 3. 阅读 API 文档

```bash
sed -n '1,120p' docs/api.md
```

## 4. 确认任务范围

我会向你确认本次任务范围、目标与约束。

```xml
<ask_followup_question>
  <question>请确认本次任务的目标、涉及模块与交付物（例如后端/前端/文档/测试）。</question>
</ask_followup_question>
```
