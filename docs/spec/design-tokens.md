# 设计令牌

参照 Clash Verge Rev **v2.5.9** `160bb4fe8eaef2bc626163d89622229e5d7c5b59`。

颜色默认值来自 `src/pages/_theme.tsx`。运行时由 `src/pages/_layout/hooks/use-custom-theme.ts` 的 `createTheme` 灌进 MUI palette，并写入一组 CSS 变量。用户在设置「主题设置」里改色后，主色 / 文字色 / 状态色可被覆盖；**背景纸色不读用户配置**，始终用当前模式的 `background_color`（`use-custom-theme.ts:166-169`）。

MUI 间距单位 8px。`borderRadius` 数字倍数基于默认 `shape.borderRadius = 4px`。

## 1. 调色板

### 1.1 主题默认色（`src/pages/_theme.tsx:5-32`）

| 令牌 | 亮色 | 暗色 | MUI palette |
| --- | --- | --- | --- |
| primary | `#007AFF` | `#0A84FF` | `palette.primary.main` |
| secondary | `#FC9B76` | `#FF9F0A` | `palette.secondary.main` |
| primary_text | `#000000` | `#FFFFFF` | `palette.text.primary` |
| secondary_text | `#3C3C4399` | `#EBEBF599` | `palette.text.secondary` |
| info | `#007AFF` | `#0A84FF` | `palette.info.main` |
| error | `#FF3B30` | `#FF453A` | `palette.error.main` |
| warning | `#FF9500` | `#FF9F0A` | `palette.warning.main` |
| success | `#06943D` | `#30D158` | `palette.success.main` |
| background | `#F5F5F5` | `#2E303D` | `palette.background.paper` 与 `background.default` |

`secondary_text` 带 alpha（亮色约 60% 黑，暗色约 60% 白）。

窗口原生背景与主题纸色一致：暗 `#2E303D`（`Color(46, 48, 61, 255)`），亮 `#F5F5F5`（`window.rs:11-14`）。主题未就绪时的占位底是亮 `#fff` / 暗 `#181a1b`（`_layout.tsx:96`）。

### 1.2 CSS 变量（`use-custom-theme.ts:202-242`，`:root` 初值在 `src/assets/styles/index.scss:15-23`）

| 变量 | 亮色 | 暗色 | 用途 |
| --- | --- | --- | --- |
| `--background-color` | `#ECECEC` | `defaultDarkTheme.background_color` = `#2E303D` | `body` 背景。注意与 palette 纸色 `#F5F5F5` 不同 |
| `--divider-color` | `rgba(0, 0, 0, 0.06)` | `rgba(255, 255, 255, 0.06)` | 页头、侧栏右边框 |
| `--selection-color` | `#f5f5f5` | `#3E3E3E` | 选中文字颜色（`::selection` 的 color） |
| `--scroller-color` | `#90939980` | `#555555` | 旧滚动条拇指（`index.scss`） |
| `--scrollbar-bg` | `#f1f1f1` | `#2E303D` | WebKit 滚动条轨道 |
| `--scrollbar-thumb` | `#c1c1c1` | `#555555` | WebKit 滚动条拇指 |
| 拇指 hover | `#a1a1a1` | `#666666` | `use-custom-theme.ts:271-273` |
| `--primary-main` | 当前 primary | 同左 | 选区背景 |
| `--background-color-alpha` | primary @ 10% | 同左 | `alpha(primary.main, 0.1)` |
| `--window-border-color` | `#cccccc` | `#1E1E1E` | `.MuiPaper-root` 边框 |
| `--border-radius` | `8px` | `8px` | `:root`，页面样式少量引用 |

选区：文字色 `--selection-color`，背景 `--primary-main`（`index.scss:25-28`）。

对话框纸色强制：亮 `#ffffff`，暗 `#2E303D`（`use-custom-theme.ts:299-301`）。

### 1.3 页面与卡片表面色（写死，不走 palette）

| 表面 | 亮色 | 暗色 | 出处 |
| --- | --- | --- | --- |
| 页面内容区 `BasePage` section | `var(--background-color)` = `#ECECEC` | `#1e1f27` | `base-page.tsx:35-42` |
| 页面容器 `base-container` | `#ffffff` | `#1e1f27` | `base-page.tsx:37` |
| 首页卡片 / 设置分栏 | `#ffffff` | `#282a36` | `enhanced-card.tsx:49`、`settings.tsx:83` |
| 代理节点行底 | `#ffffff` | `#24252f` | `proxy-item.tsx:77` |
| 代理分组头底 | `#ffffff` | `#282A36` | `proxy-render.tsx:65` |
| 订阅卡片底 | `#ffffff` | `#282A36` | `profile-box.tsx:8` |
| 侧栏选中项文字 | `#1f1f1f` | `#ffffff` | `layout-item.tsx:59` |
| 侧栏选中底 | primary @ 15% | primary @ 35% | `layout-item.tsx:55-58` |
| 节点选中底 | primary @ 15% | primary @ 35% | `proxy-item.tsx:88-91` |
| 节点选中左边条 | `primary.main` | `primary.light` | `proxy-item.tsx:77` |

### 1.4 延迟色（与信号格不是同一套阈值）

`src/services/delay.ts:425-438`，`classifyDelay` 在 `src/utils/delay.ts:16-26`。默认超时 `DEFAULT_DELAY_TIMEOUT = 10000`。

| 状态 | 判定 | 节点行颜色 | 展示文案 |
| --- | --- | --- | --- |
| 未测 | `< 0` 且不是 `-2`，或非有限数 | 不着色 | `-` |
| 测试中 | `-2` | 不着色，转圈 | `testing`（源码英文字面量） |
| 错误 | `> 1e5` | `error.main` | `Error` |
| 超时 | `0` 或 `>= timeout` | `error.main` | `Timeout` |
| 偏高 | 实测 `>= 400` | `warning.main` | 毫秒数字 |
| 中等 | 实测 `>= 250` | `primary.main` | 毫秒数字 |
| 低 | 其余实测 | `success.main` | 毫秒数字 |

首页当前节点的信号格用另一套阈值（`current-proxy-card.tsx:136-190`）：

| 延迟 | 文案 key | 颜色 |
| --- | --- | --- |
| `-2` | `home.components.currentProxy.status.testing` | `text.secondary` |
| `-1` | `...untested` | `text.secondary` |
| `> 1e5` | `...error` | `error.main` |
| `0` 或 `>= 10000` | `...timeout` | `error.main` |
| `>= 500` | `...latencyHigh` | `error.main` |
| `>= 300` | `...latencyMedium` | `warning.main` |
| `>= 200` | `...latencyGood` | `info.main` |
| 其余 | `...latencyExcellent` | `success.main` |

源码注释写明两套阈值是故意分开的（`delay.ts:434`）。

## 2. 字体与字号

字体族（`_theme.tsx:15-17`）：

```text
-apple-system, BlinkMacSystemFont, "Microsoft YaHei UI", "Microsoft YaHei", Roboto, "Helvetica Neue", Arial, sans-serif, "Apple Color Emoji"
```

Windows 再追加 `twemoji mozilla`（`src/assets/styles/font.scss`，文件 `src/assets/fonts/Twemoji.Mozilla.ttf`）。用户自定义字体拼在默认字体族前面（`use-custom-theme.ts:173-175`）。`body` 再声明一次默认字体族（`use-custom-theme.ts:276-277`）。

抗锯齿：`-webkit-font-smoothing: antialiased`（`index.scss:7`）。正文默认 `user-select: none`。

| 角色 | 字号 / 字重 | 出处 |
| --- | --- | --- |
| 页面标题 | 20px / 700 | `base-page.tsx:26` |
| 卡片标题（`EnhancedCard` 的 `h6`） | 18px / medium，单行省略 | `enhanced-card.tsx:91-98` |
| 设置分组标题 | 16px / 700 | `setting-comp.tsx:75-80` |
| 设置行主文案 | 14px | `setting-comp.tsx:33` |
| 订阅卡片标题 | 18px / 600，行高 26px | `profile-item.tsx:721-727` |
| 订阅描述 / 流量 | 14px；无流量时更新时间 12px | `profile-item.tsx:769、825、832` |
| 代理节点名 | 14px | `proxy-item.tsx:106-109` |
| 协议角标（TypeBox） | 10px，行高 1.25，圆角 4，左右 padding 2px | `proxy-item.tsx:39-48` |
| 延迟胶囊 | 14px，padding `3px 6px`，圆角 4px | `proxy-item.tsx:33-36` |
| 侧栏流量单位 | 12px，`grey.500` | `layout-traffic.tsx:66-74` |
| IP 地址 | monospace 0.75rem | `ip-info-card.tsx:322-326` |
| 系统代理/TUN 开关标签 | 15px / 500 | `proxy-control-switches.tsx:105-109` |

未覆盖的 MUI variant 用库默认：`body1` 16px、`body2` 14px、`caption` 12px、`subtitle1` 16px。

## 3. 间距

| 场景 | 值 | 出处 |
| --- | --- | --- |
| 页面头高度 | 58px，左右 padding 20px | `page.scss:7-16` |
| 内容区内边距 | 上下 10px；`.base-content` 宽 `calc(100% - 20px)` | `page.scss:26-36` |
| 首页 / 设置网格 | `spacing={1.5}`（12px），列：`xs/sm` 6 列、`md` 12 列，每卡 `size={6}`（半宽）；流量统计卡 `size={12}` | `home.tsx:356`、`settings.tsx:77` |
| 首页内容额外 padding | `contentStyle={{ padding: 2 }}` | `home.tsx:326` |
| 卡片头 | `px: 2`（16px）、`py: 1`（8px），图标 38×38，右边距 12px | `enhanced-card.tsx:54-86` |
| 卡片内容 | `p: 2`；`noContentPadding` 时为 0 | `enhanced-card.tsx:109-114` |
| 侧栏展开宽 | `flex: 1 0 200px`，右边框 1px `--divider-color` | `layout.scss:15-28` |
| 侧栏收起宽 | 72px；项 52×52，圆角 12px；logo 区 56px | `layout.scss:265-346` |
| 侧栏流量区 | 高 60px，内边距左右 20px | `layout.scss:72-78` |
| 订阅工具行 | 高 36px，左右外边距 10px | `profiles.tsx:807-816` |
| 订阅网格 | `repeat(auto-fill, minmax(clamp(260px, calc((100% - 24px) / 4), 480px), 1fr))`，`gap: 1`（8px） | `profiles.tsx:898-901` |
| 订阅卡片内边距 | `8px 16px` | `profile-box.tsx:48` |
| 代理分组头 | 内边距 `4px 8px`，圆角 8px | `proxy-render.tsx:213-225` |
| 多列节点行 | 高 56px，`my: 0.5`，`gap: 1`，`px: 2` | `proxy-render.tsx:296-303` |
| 单列节点行 | 高 40px，下边距 8px | `proxy-item.tsx:94-95` |
| 设置行（非点击） | `pt/pb: 5px` | `setting-comp.tsx:63` |

## 4. 圆角

| 场景 | 值 |
| --- | --- |
| `:root --border-radius` | 8px |
| Linux 窗口 | 8px（`_layout.tsx:150`） |
| 卡片 / 订阅卡片 / 分组头 | 8px（`borderRadius: 2` 或 `'8px'`） |
| 模式按钮、代理模式按钮 | `borderRadius: 1.5` = 6px |
| 导入 / 新建按钮 | 6px（`profiles.tsx:864`） |
| 节点行 | `borderRadius: 1` = 4px |
| 协议角标、延迟胶囊 | 4px |
| 滚动条拇指 | 4px（主题注入）或 6px（`index.scss:36`） |
| 流量统计图标圆 | 50% |
| 拖拽中的元素 | 无额外圆角，阴影见下 |

## 5. 阴影

- `createTheme` 的 `shadows` 25 档全部是 `'none'`（`use-custom-theme.ts:171`）。失败回退的 theme 没有再清空 shadows。
- 注入的全局 CSS：`* { box-shadow: none !important }`（`use-custom-theme.ts:303-307`）。
- 例外：粘性且展开的分组头写了 `boxShadow: '0 4px 8px rgba(0, 0, 0, 0.2) !important'`（`proxy-render.tsx:217-220`）。
- 拖拽中：`[data-dnd-dragging] { box-shadow: 0 6px 18px rgba(0, 0, 0, 0.32) !important }`（`index.scss:60-62`）。
- 模式按钮源码声明 hover `boxShadow: 1`、选中 `elevation={2}`（`clash-mode-card.tsx:126-128、178`）。在全局 `none !important` 下，这些阴影默认不可见。截图以实际观感为准，不要为了「有阴影」额外加一层。

## 6. 滚动条

宽高 8px（`use-custom-theme.ts:262-266` 与 `index.scss:30-34`）。侧栏菜单隐藏滚动条（`layout.scss:64-69`）。首页网站测试列表自绘 6px 拇指，`text.primary` @ 20%，圆角 3px（`test-card.tsx:27-37`）。

## 7. 断点

`use-custom-theme.ts:151-153`：`xs 0 / sm 650 / md 900 / lg 1200 / xl 1536`。与 MUI 默认 sm=600 不同。首页和设置的双列在 `md`（900px）及以上才变成 12 列网格里的半宽。

## 8. 图标

两套来源：

1. **@mui/icons-material**（Material Icons Outlined / Rounded，按组件 import 名）。
2. **本地 SVG**：`src/assets/image/`。侧栏在 `menu_icon === 'colorful'` 时用 SVG，默认或 `monochrome` 用 MUI（`layout-item.tsx:71-84`）。收起导航且 `menu_icon === 'disable'` 时强制单色（`layout-item.tsx:29-30`）。

### 8.1 侧栏（P0 页面都看得到）

| 位置 | MUI（icon[0]） | 本地 SVG（icon[1]） | 出处 |
| --- | --- | --- | --- |
| 首页 | `HomeOutlined` | `src/assets/image/itemicon/home.svg` | `_navigation.tsx:40` |
| 代理 | `WifiOutlined` | `itemicon/proxies.svg` | `:45` |
| 订阅 | `DnsOutlined` | `itemicon/profiles.svg` | `:50` |
| 连接 | `LanguageOutlined` | `itemicon/connections.svg` | `:55` |
| 规则 | `ForkRightOutlined` | `itemicon/rules.svg` | `:60` |
| 日志 | `SubjectOutlined` | `itemicon/logs.svg` | `:65` |
| 测试 | `LockOpenOutlined` | `itemicon/unlock.svg` | `:70` |
| 设置 | `SettingsOutlined` | `itemicon/settings.svg` | `:75` |
| Logo 亮色 | — | `src/assets/image/icon_light.svg` + `logo.svg` | `layout-sidebar.tsx` |
| Logo 暗色 | — | `icon_dark.svg` + `logo.svg` | 同上 |

### 8.2 首页

| 控件 | 图标 | 出处 |
| --- | --- | --- |
| 轻量模式 | `HistoryEduOutlined` | `home.tsx:341` |
| 手册 | `HelpOutlineRounded` | `:345` |
| 首页设置 | `SettingsOutlined` | `:350` |
| 订阅卡 | `CloudUploadOutlined`；外链 `LaunchOutlined`；来源 `DnsOutlined`；更新 `UpdateOutlined`；流量 `SpeedOutlined`；到期 `EventOutlined`；跳转 `StorageOutlined` | `home-profile-card.tsx` |
| 当前节点 | 信号 `SignalNone` / `SignalError` / `SignalWeak` / `SignalMedium` / `SignalGood` / `SignalStrong`；测速 `NetworkCheckRounded`；排序 `SortRounded` / `AccessTimeRounded` / `SortByAlphaRounded`；跳转 `ChevronRight` | `current-proxy-card.tsx` |
| 网络设置卡 | 卡图标 `DnsOutlined`；系统代理 `ComputerRounded`；TUN `TroubleshootRounded`；说明 `HelpOutlineRounded` | `home.tsx:380`、`proxy-tun-card.tsx` |
| 开关行 | 开 `PlayCircleOutlineRounded`（success）；关 `PauseCircleOutlineRounded`；设置 `SettingsRounded` | `proxy-control-switches.tsx` |
| 代理模式 | 卡 `RouterOutlined`；规则 `MultipleStopRounded`；全局 `LanguageRounded`；直连 `DirectionsRounded` | `clash-mode-card.tsx` |
| 流量统计 | 卡 `SpeedOutlined`；上传速度 `ArrowUpwardRounded`；下载速度 `ArrowDownwardRounded`；连接 `LinkRounded`；累计上传 `CloudUploadRounded`；累计下载 `CloudDownloadRounded`；内存 `MemoryRounded` | `enhanced-traffic-stats.tsx` |
| 网站测试 | 卡与「全部测试」`NetworkCheck`；新增 `Add` | `test-card.tsx` |
| 测试站点图 | `src/assets/image/test/{apple,github,google,youtube}.svg` | `test-card.tsx:14-17` |
| IP 卡显隐 | `VisibilityOutlined` / `VisibilityOffOutlined` | `ip-info-card.tsx` |
| Clash 信息 | `DeveloperBoardOutlined` | `clash-info-card.tsx:100` |
| 系统信息 | 卡 `InfoOutlined`；设置 `SettingsOutlined`；未运行 `PauseCircleOutlined`；管理员 `AdminPanelSettingsOutlined`；服务 `DnsOutlined`；Sidecar `ExtensionOutlined` | `system-info-card.tsx` |

国旗用字符，Windows 走 `twemoji mozilla`（`ip-info-card.tsx:279-282`）。

### 8.3 代理页

| 控件 | 图标 | 出处 |
| --- | --- | --- |
| 链式代理按钮 | 关 `LanOutlined`，开 `LanRounded` | `proxies.tsx:191-196` |
| 链式警告 | `WarningRounded` | `proxies.tsx:160` |
| 分组工具 | 定位 `MyLocationRounded`；测速 `NetworkCheckRounded`；排序三态同上；测速 URL 开 `WifiTetheringRounded` / 关 `WifiTetheringOffRounded`；细节开 `VisibilityRounded` / 关 `VisibilityOffRounded`；过滤开 `SearchOffRounded` / 关 `SearchRounded` | `proxy-group-tools.tsx` |
| 展开 | `ExpandLessRounded` / `ExpandMoreRounded` | `proxy-render.tsx:238` |
| 组内空 | `InboxRounded` | `proxy-render.tsx:285` |
| 整页空 | `InboxRounded`（`BaseEmpty`） | `base-empty.tsx:35` |
| 选中节点 | `CheckCircleOutlineRounded`（16px，无延迟时） | `proxy-item.tsx:183-188` |
| Provider | `StorageOutlined`、刷新 `RefreshRounded` | `provider-button.tsx:1` |
| 搜索三态 | 本地 SVG：`match_case.svg`、`match_whole_word.svg`、`use_regular_expression.svg` | `base-search-box.tsx:15-17` |
| 搜索清空 | `ClearRounded` | `base-search-box.tsx` |

### 8.4 订阅页

| 控件 | 图标 | 出处 |
| --- | --- | --- |
| 批量 | `CheckBoxOutlineBlankRounded` / `CheckBoxRounded` / `IndeterminateCheckBoxRounded` | `profiles.tsx` |
| 全部更新 | `RefreshRounded` | 同上 |
| 运行时配置 | `TextSnippetOutlined` | 同上 |
| 重新激活 | `LocalFireDepartmentRounded` | 同上 |
| 异常刷新 | `ClearRounded`（warning，脉冲） | `profiles.tsx:741-759` |
| 删除 | `DeleteRounded` | 同上 |
| 粘贴 | `ContentPasteRounded` | `profiles.tsx:844` |
| 卡片拖拽 | `DragIndicatorRounded` | `profile-item.tsx:711` |
| 卡片刷新 | `RefreshRounded` | `profile-item.tsx:758` |

### 8.5 设置页

| 控件 | 图标 | 出处 |
| --- | --- | --- |
| 页头 | `HelpOutlineRounded`、`Telegram`、`GitHub` | `settings.tsx:1、49-73` |
| 可点击行尾 | `ChevronRightRounded`；异步中 `CircularProgress` 20px | `setting-comp.tsx:55-59` |
| 复制 | `ContentCopyRounded` | `setting-verge-basic.tsx`、`setting-verge-advanced.tsx` |
| 允许局域网旁 | `LanRounded` | `setting-clash.tsx:161` |
| DNS / 外部控制器 / 内核旁 | `SettingsRounded` | `setting-clash.tsx` |
| 主题 CSS | `EditRounded` | `theme-viewer.tsx:1` |

其余设置对话框以文字行和开关为主，不另引图标。
