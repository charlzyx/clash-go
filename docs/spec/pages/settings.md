# 设置

优先级 **P0**。路由 `/settings`。参照 v2.5.9 `160bb4fe`。

入口 `src/pages/settings.tsx`。四块列表：`setting-system.tsx`、`setting-clash.tsx`、`setting-verge-basic.tsx`、`setting-verge-advanced.tsx`。行组件 `src/components/setting/mods/setting-comp.tsx`。对话框在 `src/components/setting/mods/`。

文案对照：`docs/spec/i18n/settings.md`（与本页同时落地；未生成前 key 以 `src/locales/{zh,en}/settings.json` 为准）。

## 1. 页面结构

`BasePage` 标题 `settings.page.title`。

页头 `ButtonGroup variant="contained"`，`aria-label` = `settings.page.actionsGroupLabel`。三个 `IconButton size="medium" color="inherit"`：

| 图标 | title key | 打开 |
| --- | --- | --- |
| `HelpOutlineRounded` | `settings.page.actions.manual` | `https://clash-verge-rev.github.io/index.html` |
| `Telegram` | `...telegram` | `https://t.me/clash_verge_re` |
| `GitHub` | `...github` | `https://github.com/clash-verge-rev/clash-verge-rev` |

失败 `showNotice.error`。

主体 `Grid spacing={1.5}`，列 `xs/sm` 6、`md` 12。两列各 `size={6}`，窄屏上下堆叠。每块卡片圆角 8px，底色亮 `#ffffff` / 暗 `#282a36`（`settings.tsx:77-116`）。

```text
左列
├── 系统设置 SettingSystem          marginBottom 12px
└── Clash 设置 SettingClash
右列
├── Verge 基础 SettingVergeBasic    marginBottom 12px
└── Verge 高级 SettingVergeAdvanced
```

## 2. 行控件

`SettingList`：分组标题 16px / 700，颜色 `text.primary`，背景透明，不吸顶。

`SettingItem`：

- 主文案 14px。`extra` 跟在文字后（通常是 tooltip 图标）。
- **不可点**（无 `onClick`）：左右排列，上下 padding 5px，右侧是 children（开关、下拉、只读框）。
- **可点**：整行 `ListItemButton`。同步点击立刻执行。异步点击后行 `disabled`，尾部 20px `CircularProgress`；平时尾部 `ChevronRightRounded`。

`GuardState`：先乐观改 UI，再 `patch`。失败交给 `onError` → 页级 `showNotice.error`，并按各字段自己的回滚逻辑恢复。开关的 `onFormat` 取 `checked`。

下拉统一 `size="small"`，内边距 `py: 7.5px`。

## 3. 系统设置

标题 `settings.sections.system.title`。`SysproxyViewer`、`TunViewer` 挂在列表上，由行内齿轮或首页网络卡打开。本列表里的系统代理 / TUN 两行复用 `ProxyControlSwitches`（与首页同一组件，见 `pages/home.md` 第 7 节）：

| 行 | 控件 | 说明 |
| --- | --- | --- |
| `settings.sections.system.toggles.tunMode` | 开关 + 齿轮开 TUN 对话框 | 无服务时打开会请求安装服务 |
| `...systemProxy` | 开关 + 齿轮开系统代理对话框 | 核心未运行时不允许打开，错误 `settings.feedback.errors.sysproxy.coreNotReady` |
| `...fields.autoLaunch` | 开关，默认关 | 先改 UI 再 `patchVerge`；失败把开关拨回去 |
| `...fields.silentStart` | 开关，默认关 | 旁有 tooltip `...tooltips.silentStart`，图标透明度 0.7 |

## 4. Clash 设置

标题 `settings.sections.clash.title`。行顺序（`setting-clash.tsx:158-351`）：

| 行 | 控件 | 行为 |
| --- | --- | --- |
| `...form.fields.allowLan` | 开关，默认关 | 旁 `LanRounded`，title `...tooltips.networkInterface`，点击打开网卡对话框。开关写 clash `allow-lan` |
| `...fields.dnsOverwrite` | 开关 | 无当前订阅、正在提交、或确认框开着时 `disabled`。旁齿轮打开 DNS 对话框；旁 info `settings.modals.dns.dialog.profileScope`。按**当前订阅**记覆盖，不写全局开关当唯一来源 |
| `...fields.ipv6` | 开关 | 写 clash `ipv6` |
| `...fields.unifiedDelay` | 开关 | tooltip `...tooltips.unifiedDelay` |
| `...fields.logLevel` | 下拉宽 100 | tooltip `...tooltips.logLevel`。选项 debug / info / warning / error / silent，文案 `...options.logLevel.*`。内核若是 `warn`，界面显示 `warning` |
| `...fields.portConfig` | 只读框宽 100，显示混合端口，指针 | 点击打开端口对话框并 blur |
| `...fields.external` | 可点行 | 行点击打开外部控制器。旁齿轮 `stopPropagation` 打开 CORS，title `settings.sections.externalCors.tooltips.open` |
| `...fields.webUI` | 可点行 | Web UI 对话框 |
| `...fields.clashCore` | 右侧纯文本内核版本 | 齿轮打开内核对话框 |
| `...fields.openUwpTool` | 可点行 | **仅 Windows**。tooltip `...tooltips.openUwpTool` |
| `...fields.updateGeoData` | 可点行 | 成功 `settings.feedback.notifications.clash.geoDataUpdated` |
| `...fields.tunnels.title` | 可点行 | 隧道对话框 |

DNS 打开若后端返回 `confirmation_required`，弹出保护确认（`setting-clash.tsx:133-155`）：

- 标题 `settings.modals.dns.protection.title`
- 正文 `...protection.message`
- 确定 `...protection.enableAnyway`
- 取消 `...protection.keepDisabled`
- 宽 `xs 320 / sm 420`。提交中 loading，取消禁用。

## 5. Verge 基础

标题 `settings.components.verge.basic.title`。

| 行 | 控件 | 默认 | 说明 |
| --- | --- | --- | --- |
| `...fields.language` | 下拉宽 110 | `en` | 选项是语言自称，不走 i18n：English、Русский、简体中文、فارسی、Татар、Bahasa Indonesia、العربية、한국어、Türkçe、Deutsch、Español、日本語、繁體中文、Português。改完同时切 dayjs（`zh` → `zh-cn`） |
| `...fields.themeMode` | `ThemeModeSwitch` 三按钮 | | `light` / `dark` / `system`，文案 `settings.sections.appearance.{light\|dark\|system}`。当前 `contained`，其余 `outlined`。`system` 跟随 OS，并 `setTheme(null)` |
| `...fields.trayClickEvent` | 下拉宽 140 | `main_window` | **Linux 不渲染此行**。选项：显示主窗口 `main_window`、显示托盘菜单 `tray_menu`、系统代理、TUN、禁用 `disable`。文案 key 见 `...trayOptions.*` 与系统开关两条 |
| `...fields.copyEnvType` | 下拉宽 140 | Windows `powershell`，其他 `bash` | Bash / Fish / Nushell / CMD / PowerShell。旁复制图标，成功 `shared.feedback.notifications.common.copySuccess` 1 秒 |
| `...fields.startPage` | 下拉宽 140 | `/` | 八个导航项，文案用各自 tab key（带空格的「首 页」等） |
| `...fields.startupScript` | 只读 `Input` 宽 230，无下划线 | 空 | 「浏览」`...actions.browse` 选 `sh`/`bat`/`ps1`。已有路径时多一个「清除」`shared.actions.clear` |
| `...fields.themeSetting` | 可点行 | | 主题对话框 |
| `...fields.layoutSetting` | 可点行 | | 布局对话框 |
| `...fields.misc` | 可点行 | | 杂项，点保存才写入 |
| `...fields.hotkeySetting` | 可点行 | | 热键 |

## 6. Verge 高级

标题 `settings.components.verge.advanced.title`。全部是可点行或带复制图标的行。

| 行 | 行为 |
| --- | --- |
| `...fields.backupSetting` | 打开备份。tooltip `...tooltips.backupInfo` |
| `...fields.runtimeConfig` | 只读运行时配置，与订阅页同一 `ConfigViewer` |
| `...fields.openConfDir` | 打开配置目录。tooltip `...tooltips.openConfDir` |
| `...fields.openCoreDir` | 打开内核目录 |
| `...fields.openLogsDir` | 打开日志目录 |
| `...fields.checkUpdates` | 无更新：成功 notice `...notifications.latestVersion`。有更新：打开更新对话框。失败 warning |
| `...fields.openDevTools` | 打开开发者工具 |
| `...fields.liteModeSettings` | 轻量模式对话框。tooltip `...tooltips.liteMode` |
| `...fields.exit` | 退出应用 |
| `...fields.exportDiagnostics` | 行本身无点击；旁复制图标导出诊断信息，成功 copySuccess |
| `...fields.vergeVersion` | 右侧 `v{version}`，padding `7px 8px 7px 0`。复制图标 title `...actions.copyVersion`，成功 `...notifications.versionCopied` |

## 7. 对话框

通用 `BaseDialog`（`base-dialog.tsx`）：标题、内容、底栏。取消默认 `outlined`，确定 `contained` 且可 `loading`。`disableOk` / `disableCancel` / `disableFooter` 可去掉按钮。纸色亮白 / 暗 `#2E303D`。

### 7.1 系统代理 `SysproxyViewer`

标题 `settings.modals.sysproxy.title`。保存 / 取消。

- 分组「当前状态」：启用状态、服务器地址、PAC URL（只读展示）。
- 代理主机：可输入，placeholder `127.0.0.1`。
- 使用 PAC：开关。
- 代理守护：开关 + tooltip。守护时长。
- 绕过检查：开关。绕过列表；非法时提示 `...messages.invalidBypass`。Windows 默认绕过含 `*.local`（v2.5.9 行为，展示在绕过 chip 里）。
- PAC 脚本：查看/编辑，按钮 `...actions.editPac`。

字段 key：`settings.modals.sysproxy.fields.*` 与 `fieldsets.currentStatus`。

### 7.2 TUN `TunViewer`

标题 `settings.modals.tun.title`（对话框内是 `Typography h6`，不是一律走 BaseDialog 标题样式，以组件为准）。

| 字段 | 控件 |
| --- | --- |
| `...fields.stack` | `StackModeSwitch`：System / gVisor / Mixed / Mips，字面量不翻译 |
| `...fields.dnsMode` | `DnsModeSwitch`：Auto（空字符串）/ Disabled / Native / Hijack，字面量不翻译。tooltip `...tooltips.dnsMode` |
| `...fields.device` | 文本 |
| `...fields.autoRoute` | 开关 |
| `...fields.autoRedirect` | 开关 + tooltip |
| `...fields.strictRoute` | 开关 |
| `...fields.autoDetectInterface` | 开关 |
| `...fields.dnsHijack` | 文本，placeholder 用 tooltip 文案 |
| `...fields.mtu` | 数字 |
| `...fields.routeExcludeAddress` | 多值编辑。非法时 `...messages.invalidRouteExcludeAddress`，否则 hint |

### 7.3 网卡 `NetworkInterfaceViewer`

标题在组件内拼接。空列表 `BaseEmpty`。每项可选，展示 MAC `settings.modals.networkInterface.fields.macAddress`。

### 7.4 DNS `DnsViewer`

宽 550。标题 `settings.modals.dns.dialog.title`。确定 `shared.actions.save`。顶部说明 `...dialog.profileScope` 与 `...dialog.warning`。校验失败 `...errors.invalid`。

分组「常规」`...sections.general`：

| 字段 | 控件 | placeholder |
| --- | --- | --- |
| enable | 开关 | |
| listen | 文本宽 150 | `:53` |
| enhancedMode | 下拉宽 150 | |
| fakeIpRange | 文本宽 150 | `198.18.0.1/16` |
| fakeIpRange6 | 文本宽 200 | `fdfe:dcba:9876::1/64` |
| fakeIpFilterMode | 下拉宽 150 | |
| ipv6 | 开关，带 description | |
| preferH3 | 开关，带 description | |
| respectRules | 开关 | |
| useHosts | 开关，带 description | |
| useSystemHosts | 开关 | |
| directPolicy | 文本 | |
| defaultNameserver | 文本 | `system,223.6.6.6, 8.8.8.8, 2400:3200::1, 2001:4860:4860::8888` |
| nameserver | 文本，带 description | `8.8.8.8, https://doh.pub/dns-query, https://dns.alidns.com/dns-query` |
| fallback | 文本，带 description | `https://dns.alidns.com/dns-query, https://dns.google/dns-query, https://cloudflare-dns.com/dns-query` |
| proxy | 文本，带 description | `https://doh.pub/dns-query, https://dns.alidns.com/dns-query` |
| directNameserver | 文本 | `system, 223.6.6.6` |
| fakeIpFilter | 文本 | `*.lan, *.local, localhost.ptlogin2.qq.com` |
| nameserverPolicy | 文本 | `+.arpa=10.0.0.1, rule-set:cn=https://doh.pub/dns-query;https://dns.alidns.com/dns-query` |

分组「回退过滤」`...sections.fallbackFilter`：geoip 开关、geoip code placeholder `CN`、IP CIDR placeholder `240.0.0.0/4, 127.0.0.1/8`、域名 placeholder `+.google.com, +.facebook.com, +.youtube.com`。

分组「Hosts」：文本 placeholder `*.clash.dev=127.0.0.1, alpha.clash.dev=::1, test.com=1.1.1.1;2.2.2.2, baidu.com=google.com`。

带 description 的字段，副文案用对应 `fields.*.description`。

### 7.5 端口 `ClashPortViewer`

宽 400。标题 `settings.modals.clashPort.title`。保存中确定按钮变为 20px 进度圈 + `shared.statuses.saving`。

| 行 | 说明 |
| --- | --- |
| 混合 `...fields.mixed` | 始终启用，无开关。若运行端口 ≠ 配置端口，副文案 `...messages.runningPort` |
| SOCKS `...fields.socks` | 开关 + 端口，默认端口 7898，默认关 |
| HTTP `...fields.http` | 同上 |
| Redir `...fields.redir` | **Windows 隐藏** |
| TProxy `...fields.tproxy` | **仅 Linux** |

每行可「随机端口」`...actions.random`，范围 1025–65535。端口必须 1–65535，否则 `...messages.invalidPort`。启用端口彼此重复 `...messages.duplicatePort`。被占用 `...messages.portInUse`（`{{port}}`）。成功 `...messages.saved`，失败 `...messages.saveFailed`。

### 7.6 外部控制器 `ControllerViewer`

宽 400。标题 `settings.sections.externalController.title`。

- 启用开关 `...fields.enable`。
- 地址 `...fields.address`，宽 175，placeholder `...placeholders.address`，旁复制，tooltip `...tooltips.copy`。
- 密钥 `...fields.secret`，同样。复制成功用 `...messages.secretCopied`（地址复制走同一 tooltip，成功文案以组件为准：密钥有专用文案）。

### 7.7 CORS `HeaderConfiguration`

宽 500。标题 `settings.sections.externalCors.title`。保存中确定文案变为 `shared.statuses.saving`。允许来源列表，placeholder `...placeholders.origin`，可增删。开关控制是否启用自定义 CORS。

### 7.8 Web UI `WebUIViewer`

宽 450，高 300，可滚动，文字可选。无确定按钮，取消文案 `shared.actions.close`。标题右侧「新建」`shared.actions.new`，编辑中禁用。

默认三条（`web-ui-viewer.tsx:15-19`）：metacubexd、yacd、zashboard，URL 含 `%host` `%port` `%secret`。打开时 `%host` 换成 `127.0.0.1`，端口缺省 `9097`。无 clash 信息 `...errors.clashInfoUnavailable`；server 不含冒号 `...errors.invalidServer`；打开失败 `...errors.openFailed`。

列表空且不在编辑：`BaseEmpty` + 12px 说明 `...messages.placeholderInstruction`。每行可改、删、打开。

### 7.9 内核 `ClashCoreViewer`

宽 400。标题为 `settings.sections.clash.form.fields.clashCore`。变体 Chip：`settings.modals.clashCore.variants.release` 与 `...alpha`。检查后：有新版本 `settings.feedback.notifications.clash.versionUpdated`，已最新 `...alreadyLatestVersion`。

### 7.10 隧道 `TunnelsViewer`

标题 `settings.sections.clash.form.fields.tunnels.title`。列表里可标默认 `...tunnels.default`。底部添加 `...tunnels.actions.add`。

### 7.11 主题 `ThemeViewer`

标题 `settings.components.verge.theme.title`。打开时拷贝当前 `theme_setting`。每行：色块（当前值或该模式默认值）+ 输入宽 135，placeholder 为默认色。Enter 保存。字段：

| key | 配置字段 |
| --- | --- |
| `...fields.primaryColor` | `primary_color` |
| `...secondaryColor` | `secondary_color` |
| `...primaryText` | `primary_text` |
| `...secondaryText` | `secondary_text` |
| `...infoColor` | `info_color` |
| `...warningColor` | `warning_color` |
| `...errorColor` | `error_color` |
| `...successColor` | `success_color` |

另两行：字体族 `...fields.fontFamily`（文本，拼在默认字体前）；CSS `...fields.cssInjection`，按钮 `...actions.editCss` 打开编辑器，标题 `...dialogs.editCssTitle`。CSS 注入被包进 `@scope`，且不能含 `@charset` / `@import` 等顶层 at-rule（`use-custom-theme.ts:16-27`）。背景图、混合模式、不透明度也在 `theme_setting` 里，若主题对话框还有对应输入，以 `theme-viewer.tsx` 字体行之后的控件为准；运行时变量见 `design-tokens.md`。

### 7.12 布局 `LayoutViewer`

标题 `settings.components.verge.layout.title`。开关默认：

| 字段 | 默认 |
| --- | --- |
| `...fields.preferSystemTitlebar` | 跟随系统装饰。Linux 默认无装饰 |
| `...fields.trafficGraph` | 开 |
| `...fields.memoryUsage` | 开 |
| `...fields.proxyGroupIcon` | 开 |
| `...fields.pauseRenderTrafficStatsOnBlur` | 开 |
| `...fields.hoverNavigator` | 开，tooltip |
| `...fields.collapseNavBar` | 关 |

下拉：

- Toast 位置默认 `top-right`，宽 180：`top-right` / `top-left` / `bottom-right` / `bottom-left`，文案 `...options.toastPosition.*`。
- 导航图标默认 `monochrome`，宽 140：`monochrome` / `colorful` / `disable`。
- macOS 另有托盘图标 `...fields.trayIcon`，同样三档，默认 `monochrome`。

悬停延迟：数字宽 120，后缀 `shared.units.milliseconds`，min 0 max 5000 step 20，默认 280。导航悬停关时输入禁用。防抖 400ms。

自定义托盘图标（普通 / TUN）用「浏览」选文件，key `...fields.tunTrayIcon` 及普通托盘图字段（`layout-viewer.tsx:552` 附近两处 browse）。

### 7.13 杂项 `MiscViewer`

标题 `settings.modals.misc.title`。点确定才保存（本地草稿）。

| 字段 | 控件 |
| --- | --- |
| `...fields.appLogLevel` | 下拉 |
| `...fields.appLogMaxSize` | 数字 |
| `...fields.appLogMaxCount` | 数字 |
| `...fields.autoCloseConnections` | 开关 + tooltip |
| `...fields.autoCheckUpdate` | 开关 |
| `...fields.enableBuiltinEnhanced` | 开关 + tooltip |
| `...fields.proxyLayoutColumns` | 下拉宽 160。`6` 显示 `...options.proxyLayoutColumns.auto`，另有 1–5 |
| `...fields.autoLogClean` | 下拉宽 160。`never` 与若干 `retainDays` |
| `...fields.autoDelayDetection` | 开关 + tooltip |
| `...fields.defaultLatencyTest` | 文本，placeholder `http://cp.cloudflare.com/generate_204`，tooltip |
| `...fields.defaultLatencyTimeout` | 数字，placeholder `10000` |

### 7.14 热键 `HotkeyViewer`

标题 `settings.modals.hotkey.title`。保存 / 取消。总开关 `...toggles.enableGlobal`，默认开。每行一个功能 + `HotkeyInput` 宽 230，可录组合键，删除 title `shared.actions.delete`。功能：

| 存储键 | 文案 key |
| --- | --- |
| `open_or_close_dashboard` | `settings.modals.hotkey.functions.openOrCloseDashboard` |
| `clash_mode_rule` | `...rule` |
| `clash_mode_global` | `...global` |
| `clash_mode_direct` | `...direct` |
| `toggle_system_proxy` | `...toggleSystemProxy` |
| `toggle_tun_mode` | `...toggleTunMode` |
| `entry_lightweight_mode` | `...entryLightweightMode` |
| `reactivate_profiles` | `...reactivateProfiles` |

托盘菜单会显示其中已设置的快捷键（见托盘大纲）。

### 7.15 备份 `BackupViewer`

宽 `xs 360 / sm 520`。标题 `settings.modals.backup.title`。取消文案 `shared.actions.close`（无单独确定，或确定被内容区按钮代替）。

两块描边区域：

1. 自动备份 `...auto.title`：`AutoBackupSettings`。
   - 定时开关 `...auto.scheduleLabel`，副文案 `...scheduleHelper`，默认关。
   - 间隔 `...auto.intervalLabel`，数字，后缀 `shared.units.hours`，默认 24，范围 1–168，失焦或 Enter 提交，四舍五入后夹紧。定时关时禁用。
   - 变更时备份 `...auto.changeLabel` / `...changeHelper`，默认开。`verge` 未加载时全部禁用。
2. 手动 `...manual.title`：
   - 本地：说明 `...manual.local`。备份、查看历史、导入备份。
   - WebDAV：说明 `...manual.webdav`。备份、查看历史、配置 `...manual.configureWebdav`。

WebDAV 对话框宽同备份主框，标题 `settings.modals.backup.webdav.title`。字段：`...fields.webdavUrl`、`...fields.username`、`shared.labels.password`。清空后仍必须能保存（v2.5.9 修复点）。

历史宽 520，内容最小高 320。标题 `...history.title`。Tab：本地 `...tabs.local`、WebDAV `...tabs.webdav`。恢复/删除前再确认，确定 `shared.actions.confirm`，确认框宽 `xs 320 / sm 420`。

### 7.16 轻量模式 `LiteModeViewer`

标题 `settings.modals.liteMode.title`。保存 / 取消。

- 立即进入 `...actions.enterNow`（行点击，不等人点保存）。
- 自动进入开关 `...toggles.autoEnter`，tooltip `...tooltips.autoEnter`。
- 延迟 `...fields.delay`。
- 说明 `...messages.autoEnterHint`。

### 7.17 更新 `UpdateViewer`

由「检查更新」在 `available` 时打开。正文为版本说明与下载/安装动作，文案在 `settings` 的 update 段。无更新不打开此框，只 toast。

### 7.18 运行时配置 `ConfigViewer`

只读编辑器，标题旁 Chip `shared.labels.readOnly`。订阅页「查看运行时订阅」与这里打开的是同一个组件。
