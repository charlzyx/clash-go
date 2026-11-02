# 测试文案对照（zh / en）

参照 Clash Verge Rev v2.5.9 `160bb4fe8eaef2bc626163d89622229e5d7c5b59`。

- 中文：`src/locales/zh/tests.json`
- 英文：`src/locales/en/tests.json`
- 加载：`src/services/i18n.ts` 以文件名作为顶层键，合并进 i18next 默认命名空间 `translation`。调用形式 `t('tests.…')`。

首页「网站测试」卡（`src/components/home/test-card.tsx`、`src/components/test/*`）与解锁测试页（`src/pages/unlock.tsx`）共用。下表按用途分两段。

## 首页网站测试卡 / 共用

| key | 中文 | English |
| --- | --- | --- |
| `tests.page.actions.testAll` | 测试全部 | Test All |
| `tests.page.title` | 测试 | Test |
| `tests.components.item.actions.test` | 测试 | Test |
| `tests.modals.test.title.create` | 新建测试 | Create Test |
| `tests.modals.test.title.edit` | 编辑测试 | Edit Test |
| `tests.modals.test.fields.url` | 测试地址 | Test URL |
| `tests.modals.test.errors.nameRequired` | 请输入名称 | Enter a name |
| `tests.modals.test.errors.urlRequired` | 请输入测试链接 | Enter a test URL |
| `tests.modals.test.errors.invalidSvg` | 图标不是有效的 SVG | The icon is not valid SVG |
| `tests.modals.test.errors.uidMissing` | 缺少测试项标识 | The test item identifier is missing |
| `tests.modals.test.errors.saveFailed` | 保存测试项失败 | Failed to save test item |

`tests.page.actions.testAll`、`tests.components.item.actions.test` 也被解锁测试页使用（「测试全部」按钮、单项按钮 Tooltip）。

## 解锁测试页

| key | 中文 | English |
| --- | --- | --- |
| `tests.statuses.test.pending` | 待检测 | Pending |
| `tests.statuses.test.yes` | 支持 | Yes |
| `tests.statuses.test.no` | 不支持 | No |
| `tests.statuses.test.failed` | 测试失败 | Failed |
| `tests.statuses.test.risky` | 存疑 | Risky |
| `tests.statuses.test.completed` | 检测完成 | Completed |
| `tests.statuses.test.disallowedIsp` | 不允许的 ISP | Disallowed ISP |
| `tests.statuses.test.originalsOnly` | 仅限原创 | Originals Only |
| `tests.statuses.test.noDisney` | 不支持（IP 被 Disney+ 禁止） | No (IP Banned By Disney+) |
| `tests.statuses.test.unsupportedRegion` | 不支持的国家/地区 | Unsupported Country/Region |
| `tests.statuses.test.failedNetwork` | 测试失败（网络连接问题） | Failed (Network Connection) |
| `tests.unlock.page.actions.testing` | 测试中... | Testing... |
| `tests.unlock.page.empty` | 暂无解锁测试项目 | No unlock test items |
| `tests.unlock.page.messages.detectionFailedWithName` | {{name}} 检测失败 | Detection failed for {{name}} |
| `tests.unlock.page.messages.detectionTimeout` | 检测超时或失败 | Detection timeout or failed |
| `tests.unlock.page.title` | 解锁测试 | Unlock Test |

状态值 `Soon`、`Blocked` 没有 key，页面直接显示英文原值。服务显示名（如「哔哩哔哩大陆」、`ChatGPT Web`）写死在 Rust，不翻译。

本页引用的 shared 文案（定义在 `i18n/shared.md`，此处只列出便于核对）：

| key | 中文 | English |
| --- | --- | --- |
| `shared.statuses.empty` | 空空如也 | Empty |
