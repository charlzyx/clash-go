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
    ├── 级别下拉：all / debug / info / warn / err
    └── BaseSearchBox
VirtualList → LogItem
空 → BaseEmpty
```

菜单 value 以 `logs.tsx` 为准：`all` / `debug` / `info` / `warn` / `err`，文案 ALL / DEBUG / INFO / WARN / ERROR。`type.includes` 下 `warn` 匹配 `warning`，`err` 匹配 `error`。初稿里的 warning / error / silent 不采用。

## 已确认行为

- 暂停写 `clashLog.enable`。暂停 title `shared.actions.pause`，继续 `shared.actions.resume`。
- 顺序默认 `asc`。降序把过滤结果反转。升序且视口靠近底部时，新日志 smooth 滚到最后一条。
- 搜索文本是 `时间 + 类型 + payload` 再转小写，再交给搜索框匹配器。级别过滤：`all` 或 `type.includes(logState)`。
- 清空调用 `refreshGetClashLog(true)`。

级别色与时间格式按 `log-item.tsx` / `use-log-data.ts`：`error` / `warning` / `info` 上色，时间 `MM-DD HH:mm:ss`，正文自动换行。等宽不采用，行组件用主题字体。

待补：搜索命中高亮；`clash-log` 跨启动持久化（当前只留在页面状态里，和连接布局一样）。升序贴底滚动用列表的跟随末尾，停住的阈值约半个像素，比 Verge 的 20px 更紧。
