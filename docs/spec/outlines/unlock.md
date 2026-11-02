# 测试 · 解锁测试（规格，P2）

路由 `/unlock`。参照 v2.5.9 `160bb4fe`。侧栏文案「测 试」（`layout.components.navigation.tabs.unlock`），**页面标题是「解锁测试」**`tests.unlock.page.title`。页面 `src/pages/unlock.tsx` 的 `UnlockPage`。检测逻辑在 Rust：`src-tauri/src/cmd/media_unlock_checker/mod.rs` → `crates/clash-verge-media-unlock/`。文案对照见 `i18n/tests.md`（本版补全 unlock 部分）。

首页的「网站测试」卡（`src/components/test/*`、`pages/home.md` 第 10 节）不是这一页；`src/components/test/` 下的 `test-box / test-item / test-viewer` 只给首页卡用，本页不引用。

## 1. 结构

```text
BasePage（非 full，有默认内边距）标题 tests.unlock.page.title
├── header（flex，gap 8px）
│   └── Button contained small「测试全部」
│       startIcon：空闲 RefreshRounded；全部检测中 CircularProgress 16 inherit
│       文案：空闲 tests.page.actions.testAll；全部检测中 tests.unlock.page.actions.testing
│       禁用：列表为空 / 全部检测中 / 有任一单项在测
└── 主体
    ├── 列表为空 → Box（高 50%，居中）→ BaseEmpty textKey=tests.unlock.page.empty
    └── Grid container spacing 1.5（12px）columns { xs:1, sm:2, md:3 }，每卡 size 1
        └── Card × N
```

出处：`unlock.tsx:235-405`。列数按断点（design-tokens §7）：<650 一列，650–899 两列，≥900 三列。

## 2. 数据与顺序

- 项目清单来自 Rust `get_unlock_items()` → `default_unlock_items()`，初始全部 `Pending`（`crates/clash-verge-media-unlock/src/lib.rs:108-113`）。v2.5.9 共 **13** 项：

| Rust 常量名 | 显示名 |
| --- | --- |
| BilibiliChinaMainland | 哔哩哔哩大陆 |
| BilibiliHkMcTw | 哔哩哔哩港澳台 |
| ChatGptWeb | ChatGPT Web |
| Claude | Claude |
| Gemini | Gemini |
| MuseAi | Muse AI |
| YouTubePremium | YouTube Premium |
| BahamutAnime | Bahamut Anime |
| Netflix | Netflix |
| DisneyPlus | Disney+ |
| PrimeVideo | Prime Video |
| Spotify | Spotify |
| TikTok | TikTok |

- 显示名是写死字符串，**不随界面语言变化**（「哔哩哔哩大陆」在英文界面也是中文）。
- 页面挂载时：读 `localStorage['clash_verge_unlock_results']` 的上次结果，按名称（去空格、转小写）与默认清单合并，再按 `name.localeCompare` 排序（`unlock.tsx:109-140`）。同名重复项：非 Pending 优先，`region` / `check_time` 互补（`:52-90`）。
- 「测试全部」完成后按名称重新排序并写回 localStorage（`:143-169`）；单项检测完成也写回（`:172-193`）。

## 3. 卡片

| 部位 | 样式 | 出处 |
| --- | --- | --- |
| Card | `variant=outlined`，高 100%，圆角 8px，**左边框 4px 状态色**（见 §4），flex 列 | `:279-296` |
| 底色 | 亮 `#ffffff` / 暗 `#282a36`；hover 亮 `primary.light`@5% / 暗 `primary.dark`@5% | `:285-291` |
| 内容区 | padding 10.4px（`p: 1.3`），flex 1 | `:297` |
| 第一行 | 左：名称 subtitle1 16px / 600 `text.primary`；右：单项检测按钮。两端对齐 | `:298-346` |
| 单项按钮 | outlined primary small，32×32 圆形，图标 `RefreshRounded`；该项在测时图标 1s 匀速旋转；该项在测或全部检测中禁用。Tooltip「测试」`tests.components.item.actions.test` | `:315-345` |
| 第二行 | flex 换行 gap 8：状态 Chip（size small，带图标，Pending 字重 normal、其余 bold）+ 地区 Chip（有 region 才显示；outlined info small） | `:348-375` |
| 分隔 | 虚线 Divider，颜色 divider@20%，左右 8px | `:378-384` |
| 底行 | padding `1.6px 12px`；caption 0.7rem `text.secondary` 右对齐：`check_time`，无则 `-- --` | `:386-398` |

- `check_time` 由 Rust 生成，本地时间 `%Y-%m-%d %H:%M:%S`（`utils.rs:4-7`），如 `2026-11-03 02:41:09`。
- `region` 是 Rust 拼好的字符串：国旗 emoji + 两位国家码，如 `🇭🇰HK`（`types.rs:28-29`）。Disney+ 可能附加 `（即将上线）` 或 ` (from main page)`（`disney_plus.rs:93、181`）。
- 卡片本身不可点击；没有右键菜单、没有详情对话框。

## 4. 状态映射

文案（`unlock.tsx:38-50`）、Chip 色（`:196-211`）、图标（`:214-220`）、左边框色（`:223-231`）：

| status 原值 | 文案 key | Chip color | 左边框 | 图标 |
| --- | --- | --- | --- | --- |
| `Pending` | `tests.statuses.test.pending` | default | divider | `PendingOutlined` |
| `Yes` | `tests.statuses.test.yes` | success | success.main | `CheckCircleOutlined` |
| `No` | `tests.statuses.test.no` | error | error.main | `CancelOutlined` |
| `Soon` | **无 key，原样显示 `Soon`** | warning | warning.main | `AccessTimeOutlined` |
| `Risky` | `tests.statuses.test.risky` | warning | warning.main | `HelpOutlined` |
| `Failed` | `tests.statuses.test.failed` | error | error.main | `HelpOutlined` |
| `Failed (Network Connection)` | `tests.statuses.test.failedNetwork` | error | error.main | `HelpOutlined` |
| `Completed` | `tests.statuses.test.completed` | info | info.main | `HelpOutlined` |
| `No (IP Banned By Disney+)` | `tests.statuses.test.noDisney` | default | divider | `HelpOutlined` |
| `Originals Only` | `tests.statuses.test.originalsOnly` | default | divider | `HelpOutlined` |
| `Unsupported Country/Region` | `tests.statuses.test.unsupportedRegion` | error | divider | `HelpOutlined` |
| `Disallowed ISP` | `tests.statuses.test.disallowedIsp` | error | divider | `HelpOutlined` |
| `Blocked` | 无 key，原样 | error | divider | `HelpOutlined` |
| 其他未知 | 原样 | default | divider | `HelpOutlined` |

注意：`No (IP Banned By Disney+)` 不等于 `No`，所以 Chip 是 default 灰、边框 divider；`Failed*` 用 `includes('Failed')`，所以带后缀的也是红。

## 5. 检测行为

- 「测试全部」：先把所有项标为在测，调用 `check_media_unlock(channel)`。Rust 端**最多 4 个并发**（`lib.rs:25`），**每项 15s 超时**，超时记 `Failed`（`lib.rs:99-104`）。每完成一项通过 Channel 回推，前端立即更新该卡并解除旋转（`unlock.tsx:144-151`）。
- 全部失败（命令报错）：`showNotice.error('tests.unlock.page.messages.detectionTimeout')`（`:162-164`）。
- 单项：`check_media_unlock_item(name)`；失败 `showNotice.error('tests.unlock.page.messages.detectionFailedWithName', { name })`（`:183-188`）。
- 检测请求经本机 mixed 端口 `http://127.0.0.1:{port}` 走代理（`media_unlock_checker/mod.rs:12-27`），即结果反映当前所选节点。
- 页面不自动检测、不轮询；进入页面只显示缓存结果。

## 6. 亮 / 暗差异

| 项 | 亮 | 暗 |
| --- | --- | --- |
| 页面 | 非 full：容器 `#ffffff`、section `#ECECEC`，有默认内边距 | `#1e1f27` |
| 卡片底 | `#ffffff` | `#282a36` |
| hover | `primary.light`@5% | `primary.dark`@5% |
| 状态色 | success `#06943D` / error `#FF3B30` / warning `#FF9500` / info `#007AFF` | `#30D158` / `#FF453A` / `#FF9F0A` / `#0A84FF` |

## 7. 截图核对清单（`ui.NewTester`）

建议数据：13 项，混合 Yes（带 region）、No、Failed、Risky、Pending、Soon 各至少一项。

- [ ] `test-light/dark`：页头标题「解锁测试」（不是「测试」）；右侧实心小按钮「测试全部」带刷新图标
- [ ] 默认窗 940 宽减侧栏后内容宽 <900 时为 2 列；拉宽到内容 ≥900 为 3 列；<650 为 1 列（以内容区宽度还是窗口宽度计，见待确认 1）
- [ ] 卡片左侧 4px 彩条颜色与状态一致（Yes 绿、No 红、Risky/Soon 橙、Pending 灰）
- [ ] 卡片名称 16px/600；右上 32px 圆形描边刷新按钮
- [ ] 状态 Chip 带图标；Pending 字重正常，其余粗体；region Chip 蓝色描边带国旗
- [ ] 底部虚线 + 右对齐小字时间 `YYYY-MM-DD HH:mm:ss`，未检测为 `-- --`
- [ ] 「测试全部」进行中：按钮文字「测试中...」+ 16px 转圈，所有单项按钮禁用
- [ ] 单项检测中：该卡刷新图标旋转，「测试全部」禁用
- [ ] 「哔哩哔哩大陆 / 港澳台」在英文界面仍为中文
- [ ] 英文：`Unlock Test`、`Test All`、`Testing...`、状态 `Yes / No / Failed / Risky / Pending`
- [ ] 空清单：收件箱 +「暂无解锁测试项目」，位于上半区居中

## 8. 待确认

1. MUI Grid 断点按**窗口宽度**（媒体查询）判定，不是内容区宽度。940×700 默认窗下窗口宽 940 ≥ md 900 → 3 列，即便内容区只有约 740px。需在原版截图确认实际列数。
2. 中文名与英文名混排时 `localeCompare` 的结果依赖 WebView 默认 locale（未传 locale 参数）。「哔哩哔哩…」排在拉丁字母前还是后需原版截图确认；clash-go 若用 Go `collate` 要对齐。
3. `Soon`、`Blocked` 没有翻译 key，原版直接显示英文原值。v2.5.9 Rust 端只有 Disney+ 会产出 `Soon`；`Blocked` 与 `Disallowed ISP` 在 Rust 端未见产出（grep 无结果），前端映射保留即可。
4. 卡片 hover 背景 5% 主色与卡片白底叠加的观感、`Card variant=outlined` 的 1px 外框颜色（MUI 默认 `divider`）需截图核对。
