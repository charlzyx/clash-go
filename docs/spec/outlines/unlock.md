# 测试（大纲，P2）

路由 `/unlock`。侧栏文案是「测 试」，页面标题 `tests.unlock.page.title`。页面 `src/pages/unlock.tsx`。

首页里的网站测速卡不是这一页，规格在 `pages/home.md` 第 10 节，文案表 `i18n/tests.md`。

## 结构

```text
BasePage 标题 tests.unlock.page.title
└── header
    └── 「全部测试」contained small
        测试中：16px CircularProgress + tests.unlock.page.actions.testing
        空列表或已有项在测时禁用
网格卡片（约一半高度，unlock.tsx:270 附近）
每项：名称、状态 Chip、地区、耗时
点击单项重测
空 → BaseEmpty
```

## 状态

结果缓存在 `localStorage['clash_verge_unlock_results']`。同名项合并，非 Pending 覆盖 Pending。

| status | Chip 色 | 边框 | 图标 |
| --- | --- | --- | --- |
| Pending | default | divider | `PendingOutlined` |
| Yes | success | success | `CheckCircleOutlined` |
| No | error | error | `CancelOutlined` |
| Soon | warning | warning | `AccessTimeOutlined` |
| Risky | warning | warning | `HelpOutlined`（图标函数在 Soon 之后落到 Help） |
| Failed* | error | error | Help |
| Completed | info | info | Help |
| Disallowed ISP / Blocked / Unsupported Country/Region | error | divider（边框函数未单列这些字符串，只有 Chip 色特殊） | Help |

文案 key 在 `tests.statuses.test.*`（`unlock.tsx:38-50`），已在 `i18n/tests.md`。v2.5.9 给 Muse AI 增加 Risky。

失败 notice：`tests.unlock.page.messages.detectionFailedWithName`（`{{name}}`）。

待补：卡片网格列数、单项按钮禁用态、地区/耗时字段格式、列表来源 `getUnlockItems` 的空态。
