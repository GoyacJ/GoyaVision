# 提交前检查与提交信息

用于生成符合 Conventional Commits 的提交信息并进行简单检查（提交信息需使用中文描述）。

## 1. 检查分支命名

```bash
git branch --show-current
```
确认分支名符合 `feature/*`, `fix/*`, `hotfix/*`, `refactor/*` 或 `release/*` 规范。

## 2. 查看当前变更

```bash
git status -sb
```

## 2. 确认提交范围

```xml
<ask_followup_question>
  <question>请提供本次变更类型（feat/fix/docs/refactor/test/chore/perf/style）与 scope（如 asset/workflow/ui），以及简要描述（中文）。</question>
</ask_followup_question>
```

## 3. 生成提交命令

输出符合规范的提交信息，并提示是否需要更新文档。
