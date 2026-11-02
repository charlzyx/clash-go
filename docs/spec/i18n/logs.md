# 日志文案对照（zh / en）

参照 Clash Verge Rev v2.5.9 `160bb4fe8eaef2bc626163d89622229e5d7c5b59`。

- 中文：`src/locales/zh/logs.json`
- 英文：`src/locales/en/logs.json`
- 加载：`src/services/i18n.ts` 以文件名作为顶层键，合并进 i18next 默认命名空间 `translation`。调用形式 `t('logs.…')`。

日志页。`src/pages/logs.tsx` 与 `src/components/log/log-item.tsx`。

| key | 中文 | English |
| --- | --- | --- |
| `logs.page.title` | 日志 | Logs |
| `logs.actions.showDescending` | 按时间倒序 | Newest first |
| `logs.actions.showAscending` | 按时间正序 | Oldest first |

本页引用的 shared 文案（定义在 `i18n/shared.md`，此处只列出便于核对）：

| key | 中文 | English |
| --- | --- | --- |
| `shared.actions.pause` | 暂停 | Pause |
| `shared.actions.resume` | 继续 | Resume |
| `shared.actions.clear` | 清除 | Clear |
| `shared.filters.logLevels.all` | ALL | ALL |
| `shared.filters.logLevels.debug` | DEBUG | DEBUG |
| `shared.filters.logLevels.info` | INFO | INFO |
| `shared.filters.logLevels.warn` | WARN | WARN |
| `shared.filters.logLevels.error` | ERROR | ERROR |
| `shared.statuses.empty` | 空空如也 | Empty |
| `shared.placeholders.filter` | 过滤条件 | Filter conditions |

日志级别标签本身显示内核原值并 `text-transform: uppercase`（`INFO`、`WARNING`、`ERROR`、`DEBUG`），不走翻译。
