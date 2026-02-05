# 开发完成收尾（GoyaVision）

该工作流用于开发完成后的文档同步与提交检查。

## 1. 更新开发进度

```xml
<ask_followup_question>
  <question>请提供本次完成/变更的模块与状态（用于更新 docs/development-progress.md）。</question>
</ask_followup_question>
```

## 2. 更新 CHANGELOG

在 `[未发布]` 下按类型记录变更（新增/变更/修复/弃用/移除/安全）。

## 3. 同步 API/架构/需求文档

如涉及 API/设计变更，请同步更新 `docs/api.md`、`docs/architecture.md`、`docs/requirements.md`。

## 4. 提交前检查

```bash
git status -sb
```

如需提交，请按 Conventional Commits 格式生成提交信息。
