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

## 落地说明

对照 `unlock.tsx`（`160bb4fe`）之后：

- 路由 `/unlock`，侧栏「测 试」，标题「解锁测试」（`tests.unlock.page.title`）。页头按钮文案用大纲和首页已经落地的「全部测试」；i18n 表里同一 key 写成「测试全部」。测试进行中按钮换成 16px 进度圈和「测试中... n/N」，空列表或已有项目在测时禁用。
- 网格按窗口宽度：小于 650 一列，小于 900 两列，否则三列（`columns={{ xs:1, sm:2, md:3 }}`）。卡片左缘 4px 是状态色，其余是分割线。点卡片重测该项；正在测的卡片禁用。
- 状态文案、Chip 色、边框、图标按上面的表。`Soon`、`Blocked` 没有中文 key，Chip 上保留原文。地区是描边 info Chip；没有地区就不画。时间位是 `check_time`（`YYYY-MM-DD HH:MM:SS`），空显示 `-- --`。
- 空列表用跟连接页一样的空态（方框 + 「空空如也」）。Verge 这里给 BaseEmpty 传了 `tests.unlock.page.empty`（「暂无解锁测试项目」），大纲只写了 BaseEmpty，所以跟连接页用默认空文案。
- 结果写进设置里的 `prefs.clash_verge_unlock_results`，经现有 `SavePrefs` 进 `settings.json`，不另开存储。同名合并：忽略大小写和首尾空白，非 Pending 覆盖 Pending，目录外的名字丢掉，展示名用目录里的写法。进入页面时合并一次；离开再进来会再合并。停在本页时隐藏窗口不会重新加载。
- 探测本身不在这个页面里发请求。`CheckUnlock` 由调用方注入；测试塞假结果。没有注入时，单项失败 notice 是「{{name}} 检测失败」，全部测试失败 notice 是「检测超时或失败」，状态不改，也不写入缓存。真实的流媒体请求（Netflix、Disney+ 等，经混合端口）留到后面。
