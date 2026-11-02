# 规则文案对照（zh / en）

参照 Clash Verge Rev v2.5.9 `160bb4fe8eaef2bc626163d89622229e5d7c5b59`。

- 中文：`src/locales/zh/rules.json`
- 英文：`src/locales/en/rules.json`
- 加载：`src/services/i18n.ts` 以文件名作为顶层键，合并进 i18next 默认命名空间 `translation`。调用形式 `t('rules.…')`。

规则页。`src/pages/rules.tsx` 与 `src/components/rule/*`。`rules.modals.editor.*`（共 40+ 条）属于订阅页的规则编辑器 `src/components/profile/rules-editor-viewer.tsx`，不在本表。

| key | 中文 | English |
| --- | --- | --- |
| `rules.page.provider.trigger` | 规则集合 | Rule Provider |
| `rules.page.provider.dialogTitle` | 规则集合 | Rule Providers |
| `rules.page.provider.actions.updateAll` | 更新全部 | Update All |
| `rules.page.provider.actions.update` | 更新 | Update |
| `rules.page.title` | 规则 | Rules |
| `rules.feedback.notifications.provider.updateSuccess` | {{name}} 更新成功 | {{name}} updated successfully |
| `rules.feedback.notifications.provider.updateFailed` | {{name}} 更新失败: {{message}} | Failed to update {{name}}: {{message}} |
| `rules.feedback.notifications.provider.genericError` | 更新失败: {{message}} | Update failed: {{message}} |
| `rules.feedback.notifications.provider.none` | 没有可更新的规则集合 | No providers available to update |
| `rules.feedback.notifications.provider.allUpdated` | 所有规则集合均已更新 | All providers updated successfully |

本页引用的 shared 文案（定义在 `i18n/shared.md`，此处只列出便于核对）：

| key | 中文 | English |
| --- | --- | --- |
| `shared.labels.updateAt` | 更新于 | Update At |
| `shared.actions.close` | 关闭 | Close |
| `shared.statuses.empty` | 空空如也 | Empty |
| `shared.placeholders.filter` | 过滤条件 | Filter conditions |

「更新于」后的相对时间同样用 dayjs `fromNow()`，如「3 小时前」/ `3 hours ago`。
