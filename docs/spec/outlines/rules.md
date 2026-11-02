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

待补：点击行有没有详情（源码行本身不可点，没有详情面板）。

## 落地说明

对照 `rules.tsx` / `rule-item.tsx` / `provider-button.tsx`（`160bb4fe`）之后：

- 无匹配或没有规则时不画列表。Verge 这里渲染 `BaseEmpty`（收件箱 + 「空空如也」）。大纲写的是空白区域，M2 也没有共用的空态组件，所以这块保持空白，搜索行仍在。
- 搜索框只对 `payload` 做不区分大小写的包含匹配（`BaseSearchBox` 的默认行为，UI-R01）。区分大小写、全字、正则三个开关未做。
- 行号是过滤前的 `index+1`。策略颜色按 `rule-item.tsx`：`REJECT` / `REJECT-DROP` 用错误色，`DIRECT` 用正文色，其余按码点求和取色。
- 没有规则 provider 时不显示「规则集合」按钮，与 Verge 一致。有 provider 时对话框列表 `minHeight` 250，可更新单个或全部。
- 进入 `/rules` 刷新一次规则和规则 provider；窗口隐藏后再显示，并且仍停在该页时再刷新一次。Verge 在「已经可见」的挂载上会把两次刷新叠在同一次 effect 里，这里合并成一次。
- 滚过约 100px 显示回到顶部。mygo 的虚拟列表只给出可见行号，按行高 40 估算偏移（第 3 行及以后）。
