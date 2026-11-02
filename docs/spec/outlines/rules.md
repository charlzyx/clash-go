# 规则（大纲，P1）

路由 `/rules`。页面 `src/pages/rules.tsx`。组件 `src/components/rule/rule-item.tsx`、`provider-button.tsx`。文案 `src/locales/{zh,en}/rules.json`（本轮未生成对照表）。

## 结构

```text
BasePage full 标题 rules.page.title
├── header：规则 Provider 按钮
└── 搜索行高 36（左右 10px）
    └── BaseSearchBox（匹配 payload）
VirtualList estimateSize 40
空列表不渲染列表（无匹配时区域为空，不是 BaseEmpty；待核对 rules.tsx:85 之后）
滚过 100px 显示 ScrollTopButton
```

## 已确认行为

- 进入页面以及页面重新可见时刷新规则和规则 provider（`rules.tsx:28-36`）。
- 行号是 UI 派生的 `index + 1`，不改内核数据。
- 搜索只对 `payload` 做 `BaseSearchBox` 匹配。
- `RuleItem`：序号、类型、payload、代理策略。样式见 `rule-item.tsx`（序号最小宽 30，payload 最小宽 120）。
- Provider 按钮样式与代理页同类：secondary 角标，列表 `minHeight 250`，可更新单个 / 全部。

待补：行内完整字段与颜色、provider 对话框字段、无规则时的空态是否真的空白、点击行有没有详情。
