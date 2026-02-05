#!/bin/bash

# GoyaVision Design Tokens 迁移检查脚本
# 用于检查代码中是否还有需要迁移的旧设计系统引用

set -e

echo "🔍 GoyaVision Design Tokens 迁移检查"
echo "======================================"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查项计数
ISSUES_FOUND=0

# 1. 检查 secondary 色系
echo "📌 检查 1/6: Secondary 色系引用"
SECONDARY_REFS=$(grep -r "secondary-" web/src/ --include="*.vue" --include="*.ts" --include="*.tsx" --include="*.css" 2>/dev/null | wc -l | xargs)
if [ "$SECONDARY_REFS" -gt 0 ]; then
  echo -e "${RED}❌ 发现 $SECONDARY_REFS 处 secondary- 引用${NC}"
  grep -rn "secondary-" web/src/ --include="*.vue" --include="*.ts" --include="*.tsx" --include="*.css" | head -5
  echo "   ... (显示前 5 个，完整列表请运行: grep -rn 'secondary-' web/src/)"
  ISSUES_FOUND=$((ISSUES_FOUND + SECONDARY_REFS))
else
  echo -e "${GREEN}✅ 无 secondary- 引用${NC}"
fi
echo ""

# 2. 检查彩色阴影
echo "📌 检查 2/6: 彩色阴影引用"
COLOR_SHADOWS=$(grep -rE "shadow-(primary|secondary|success|error|warning|info)" web/src/ --include="*.vue" --include="*.ts" --include="*.tsx" --include="*.css" 2>/dev/null | wc -l | xargs)
if [ "$COLOR_SHADOWS" -gt 0 ]; then
  echo -e "${RED}❌ 发现 $COLOR_SHADOWS 处彩色阴影引用${NC}"
  grep -rnE "shadow-(primary|secondary|success|error|warning|info)" web/src/ --include="*.vue" --include="*.ts" --include="*.tsx" --include="*.css" | head -5
  echo "   ... (显示前 5 个)"
  ISSUES_FOUND=$((ISSUES_FOUND + COLOR_SHADOWS))
else
  echo -e "${GREEN}✅ 无彩色阴影引用${NC}"
fi
echo ""

# 3. 检查渐变背景
echo "📌 检查 3/6: 渐变背景引用"
GRADIENTS=$(grep -r "bg-gradient-to-" web/src/ --include="*.vue" --include="*.ts" --include="*.tsx" --include="*.css" 2>/dev/null | wc -l | xargs)
if [ "$GRADIENTS" -gt 0 ]; then
  echo -e "${RED}❌ 发现 $GRADIENTS 处渐变背景引用${NC}"
  grep -rn "bg-gradient-to-" web/src/ --include="*.vue" --include="*.ts" --include="*.tsx" --include="*.css" | head -5
  echo "   ... (显示前 5 个)"
  ISSUES_FOUND=$((ISSUES_FOUND + GRADIENTS))
else
  echo -e "${GREEN}✅ 无渐变背景引用${NC}"
fi
echo ""

# 4. 检查毛玻璃效果
echo "📌 检查 4/6: 毛玻璃效果引用"
BACKDROP_FILTER=$(grep -r "backdrop-filter" web/src/ --include="*.vue" --include="*.ts" --include="*.tsx" --include="*.css" 2>/dev/null | wc -l | xargs)
if [ "$BACKDROP_FILTER" -gt 0 ]; then
  echo -e "${RED}❌ 发现 $BACKDROP_FILTER 处 backdrop-filter 引用${NC}"
  grep -rn "backdrop-filter" web/src/ --include="*.vue" --include="*.ts" --include="*.tsx" --include="*.css" | head -5
  echo "   ... (显示前 5 个)"
  ISSUES_FOUND=$((ISSUES_FOUND + BACKDROP_FILTER))
else
  echo -e "${GREEN}✅ 无 backdrop-filter 引用${NC}"
fi
echo ""

# 5. 检查过度动画
echo "📌 检查 5/6: 过度动画引用"
OVER_ANIMATIONS=$(grep -rE "hover:(scale|(-)?translate)" web/src/ --include="*.vue" --include="*.ts" --include="*.tsx" --include="*.css" 2>/dev/null | wc -l | xargs)
if [ "$OVER_ANIMATIONS" -gt 0 ]; then
  echo -e "${YELLOW}⚠️  发现 $OVER_ANIMATIONS 处过度动画引用（建议简化）${NC}"
  grep -rnE "hover:(scale|(-)?translate)" web/src/ --include="*.vue" --include="*.ts" --include="*.tsx" --include="*.css" | head -5
  echo "   ... (显示前 5 个)"
fi
echo ""

# 6. 检查旧动画时长
echo "📌 检查 6/6: 旧动画时长引用"
OLD_DURATIONS=$(grep -rE "duration-(short|medium|long|extra-long)" web/src/ --include="*.vue" --include="*.ts" --include="*.tsx" --include="*.css" 2>/dev/null | wc -l | xargs)
if [ "$OLD_DURATIONS" -gt 0 ]; then
  echo -e "${YELLOW}⚠️  发现 $OLD_DURATIONS 处旧动画时长引用（建议迁移）${NC}"
  grep -rnE "duration-(short|medium|long|extra-long)" web/src/ --include="*.vue" --include="*.ts" --include="*.tsx" --include="*.css" | head -5
  echo "   ... (显示前 5 个)"
fi
echo ""

# 总结
echo "======================================"
if [ "$ISSUES_FOUND" -eq 0 ]; then
  echo -e "${GREEN}🎉 恭喜！未发现必须修复的问题${NC}"
  echo ""
  echo "✅ Design Tokens 迁移完成"
else
  echo -e "${RED}⚠️  发现 $ISSUES_FOUND 个必须修复的问题${NC}"
  echo ""
  echo "请参考以下文档进行迁移："
  echo "  - docs/design-tokens-migration-guide.md"
  echo "  - docs/frontend-refactor-design.md"
fi
echo ""

exit 0
