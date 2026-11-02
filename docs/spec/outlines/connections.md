# 连接（大纲，P1）

路由 `/connections`。页面 `src/pages/connections.tsx`。组件目录 `src/components/connection/`。文案 `src/locales/{zh,en}/connections.json`（本轮未生成对照表）。

## 结构

```text
BasePage full 标题 connections.page.title（不换行）
├── header
│   ├── 累计下载 / 累计上传（shared.labels.downloaded|uploaded + parseTraffic）
│   ├── 列表/表格切换 TableRowsRounded | TableChartRounded
│   └── 「关闭全部」shared.actions.closeAll
└── 工具行（粘性，最小高 36）
    ├── ButtonGroup 活跃 | 已关闭（带条数）
    ├── 排序下拉（仅列表布局）
    ├── BaseSearchBox
    └── 表格布局时列管理 ViewColumnRounded
列表 VirtualList 或 ConnectionTable
点击行 → ConnectionDetail
空 → BaseEmpty
右下回到顶部 Fab（滚过 100px）
```

## 已确认行为

- 活跃 / 已关闭：`connections.components.actions.active|closed`，当前 `contained`。
- 列表排序三项：`connections.components.order.default`（开始时间新到旧）、`uploadSpeed`、`downloadSpeed`。表格布局隐藏这个下拉。
- 布局存在连接设置里，`layout: 'table' | 'list'`。
- 搜索用 `BaseSearchBox`（大小写 / 全字 / 正则），有文字才过滤。
- 「关闭全部」调用 `closeAllConnections`，失败不弹 notice（`errorNotice: false`）。
- 详情抽屉区分活跃与已关闭（`connectionsType === 'closed'`）。
- 行组件 `connection-row-item.tsx`，表格 `connection-table.tsx`，列管理 `connection-column-manager.tsx`，相对时间 `connection-relative-time.tsx`。

待补：每列字段与默认可见性、详情面板字段、行内关闭单条连接、空态文案、表格横向滚动与列宽。
