# 日志（大纲，P2）

路由 `/logs`。页面 `src/pages/logs.tsx`。行 `src/components/log/log-item.tsx`。文案 `src/locales/{zh,en}/logs.json`（本轮未生成对照表）。级别文案在 `shared.filters.logLevels.*`（已收入 `i18n/shared.md`）。

## 结构

```text
BasePage full 标题 logs.page.title
├── header
│   ├── 暂停/继续 PauseCircleOutlineRounded | PlayCircleOutlineRounded
│   ├── 顺序 SwapVertRounded（降序时 scaleY(-1)）
│   └── 清空 shared.actions.clear
└── 工具行高 39
    ├── 级别下拉：all / debug / info / warning / error / silent（以 logs.tsx 菜单为准）
    └── BaseSearchBox
VirtualList → LogItem
空 → BaseEmpty
```

## 已确认行为

- 暂停写 `clashLog.enable`。暂停 title `shared.actions.pause`，继续 `shared.actions.resume`。
- 顺序默认 `asc`。降序把过滤结果反转。升序且视口靠近底部时，新日志 smooth 滚到最后一条。
- 搜索文本是 `时间 + 类型 + payload` 再转小写，再交给搜索框匹配器。级别过滤：`all` 或 `type.includes(logState)`。
- 清空调用 `refreshGetClashLog(true)`。

待补：`LogItem` 的级别色、时间格式、等宽字体、自动换行；级别菜单的完整 value 列表。
