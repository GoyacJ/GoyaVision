---
name: commit
description: Git 提交规范指导（Conventional Commits）。用于提交前检查与撰写提交信息。
---

# Git 提交规范

用于按 Conventional Commits 规范完成提交。

## 何时使用
- 功能/修复完成准备提交
- 需要检查提交信息格式

## 提交格式
```
<type>(<scope>): <subject>
```

### 类型
- feat / fix / docs / refactor / test / chore / perf / style

### 范围（可选）
- asset / operator / workflow / task / auth / api / ui

## 提交前检查
- 代码已测试、已格式化
- 文档已同步（进度、CHANGELOG、API）
- 无调试代码
