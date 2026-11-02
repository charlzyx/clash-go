# 连接文案对照（zh / en）

参照 Clash Verge Rev v2.5.9 `160bb4fe8eaef2bc626163d89622229e5d7c5b59`。

- 中文：`src/locales/zh/connections.json`
- 英文：`src/locales/en/connections.json`
- 加载：`src/services/i18n.ts` 以文件名作为顶层键，合并进 i18next 默认命名空间 `translation`。调用形式 `t('connections.…')`。

连接页。`src/pages/connections.tsx` 与 `src/components/connection/*`。

| key | 中文 | English |
| --- | --- | --- |
| `connections.page.title` | 连接 | Connections |
| `connections.components.fields.host` | 主机 | Host |
| `connections.components.fields.dlSpeed` | 下载速度 | DL Speed |
| `connections.components.fields.ulSpeed` | 上传速度 | UL Speed |
| `connections.components.fields.chains` | 链路 | Chains |
| `connections.components.fields.rule` | 规则 | Rule |
| `connections.components.fields.process` | 进程 | Process |
| `connections.components.fields.time` | 连接时间 | Time |
| `connections.components.fields.source` | 源地址 | Source |
| `connections.components.fields.destination` | 目标地址 | Destination |
| `connections.components.fields.destinationPort` | 目标端口 | Destination Port |
| `connections.components.fields.type` | 类型 | Type |
| `connections.components.order.default` | 默认 | Default |
| `connections.components.order.uploadSpeed` | 上传速度 | Upload Speed |
| `connections.components.order.downloadSpeed` | 下载速度 | Download Speed |
| `connections.components.actions.active` | 活跃 | Active |
| `connections.components.actions.closed` | 已关闭 | Closed |
| `connections.components.actions.closeConnection` | 关闭连接 | Close Connection |
| `connections.components.columnManager.title` | 列设置 | Columns |
| `connections.components.columnManager.dragHandle` | 拖拽控件 | Drag handle |

本页引用的 shared 文案（定义在 `i18n/shared.md`，此处只列出便于核对）：

| key | 中文 | English |
| --- | --- | --- |
| `shared.labels.downloaded` | 下载量 | Downloaded |
| `shared.labels.uploaded` | 上传量 | Uploaded |
| `shared.actions.listView` | 列表视图 | List View |
| `shared.actions.tableView` | 表格视图 | Table View |
| `shared.actions.closeAll` | 关闭全部 | Close All |
| `shared.actions.clear` | 清除 | Clear |
| `shared.actions.resetToDefault` | 重置为默认值 | Reset to Default |
| `shared.actions.close` | 关闭 | Close |
| `shared.statuses.empty` | 空空如也 | Empty |
| `shared.placeholders.filter` | 过滤条件 | Filter conditions |

相对时间（连接时间列、详情「连接时间」）用 dayjs `relativeTime`，语言 `zh-cn` / `en`（`_layout.tsx:31-34、85`），如「几秒前」「1 分钟前」/ `a few seconds ago`、`a minute ago`。
