# 首页

优先级 **P0**。路由 `/`。参照 v2.5.9 `160bb4fe`。

入口：`src/pages/home.tsx` 的 `HomePage`。卡片壳：`src/components/home/enhanced-card.tsx` 的 `EnhancedCard`。页面壳：`BasePage`（`src/components/base/base-page.tsx`）。

## 1. 结构

```text
BasePage title=home.page.title
├── header
│   ├── IconButton 轻量模式 HistoryEduOutlined
│   ├── IconButton 手册 HelpOutlineRounded
│   └── IconButton 首页设置 SettingsOutlined
└── Grid spacing=1.5 columns={xs:6, sm:6, md:12}
    ├── [可关] profile     HomeProfileCard          size 6
    ├── [可关] proxy       CurrentProxyCard         size 6
    ├── [可关] network     NetworkSettingsCard      size 6
    ├── [可关] mode        ClashModeEnhancedCard    size 6
    ├── [可关] traffic     EnhancedTrafficStats     size 12
    ├── [可关] test        TestCard                 size 6  懒加载
    ├── [可关] ip          IpInfoCard               size 6  懒加载
    ├── [可关] clashinfo   ClashInfoCard            size 6  懒加载
    └── [可关] systeminfo  SystemInfoCard           size 6  懒加载
```

出处：`home.tsx:270-360`。`contentStyle` 为 `padding: 2`（`home.tsx:326`）。窄屏（小于 md 900px）每张默认卡占满 6 列，即单列；md 及以上两列。流量卡始终整行。

默认显示（`home.tsx:89-100`）：除 `info` 外全部为 true。`info` 没有对应卡片，勾选列表里也不出现。

懒加载卡片在 chunk 到达前显示 `Skeleton variant="rectangular" height={200}`（`home.tsx:298-318`）。

## 2. 卡片壳 `EnhancedCard`

`enhanced-card.tsx:42-121`。

- 外盒：高 100%，列方向，圆角 8px，底色亮 `#ffffff` / 暗 `#282a36`。
- 头：底边 `divider`。左：38×38 圆角 6px 的图标井，底色 `iconColor` @ 12%，图标色 `iconColor.main`。标题 18px / medium，超出省略，`title` 属性为全文。
- 右：可选 `action`，左边距 16px。
- 内容：默认 padding 16px。
- `iconColor` 取值：`primary | secondary | error | warning | info | success`。

## 3. 页头动作

| 按钮 | tooltip key | 行为 | 出处 |
| --- | --- | --- | --- |
| 轻量模式 | `home.page.tooltips.lightweightMode` | 调用 `entry_lightweight_mode()`。失败不弹 notice（`errorNotice: false`） | `home.tsx:329-342` |
| 手册 | `home.page.tooltips.manual` | 打开 `https://clash-verge-rev.github.io/index.html`，失败 `showNotice.error` | `:343-347、246-250` |
| 首页设置 | `home.page.tooltips.settings` | 打开 `HomeSettingsDialog` | `:348-352` |

三个按钮都是 `size="small"`、`color="inherit"`。Tooltip 带箭头。

## 4. 首页设置对话框

`home.tsx:115-230`。`Dialog` `maxWidth="xs"` `fullWidth`。标题 `home.page.settings.title`。

复选框（顺序固定）：

| 状态键 | 文案 key | 默认 |
| --- | --- | --- |
| `profile` | `home.page.settings.cards.profile` | 开 |
| `proxy` | `...currentProxy` | 开 |
| `network` | `...network` | 开 |
| `mode` | `...proxyMode` | 开 |
| `traffic` | `...traffic` | 开 |
| `test` | `...tests` | 开 |
| `ip` | `...ip` | 开 |
| `clashinfo` | `...clashInfo` | 开 |
| `systeminfo` | `...systemInfo` | 开 |

交互：勾选只改对话框本地状态。取消 `shared.actions.cancel` 丢弃。保存 `shared.actions.save` 调 `patchVerge({ home_cards })` 后关闭。对话框用卡片标志串作 `key`，外部配置变化时重置（`home.tsx:363-368`）。

## 5. 订阅卡 `HomeProfileCard`

`src/components/home/home-profile-card.tsx`。图标 `CloudUploadOutlined`，`iconColor="info"`。

**有当前订阅**

- 标题：无 `home` 时纯文本 `current.name`；有 `home` 时标题是按钮链接，右侧 `LaunchOutlined`，点击用系统浏览器打开 `home`（`openProfileHome`，`:55-57、311-344`）。
- 右侧按钮：文案 `layout.components.navigation.tabs.profiles`（「订 阅」），`outlined` `small`，尾图标 `StorageOutlined`，点击 `navigate('/profile')`。
- 有 `url`：一行 `DnsOutlined` + `shared.labels.from` + URL 的 host；非法 URL 显示字面量 `local`（`:44-48`）。有 `home` 时 host 可点并带外链图标。
- 有 `updated`：`UpdateOutlined` + `shared.labels.updateTime` + `YYYY-MM-DD HH:mm`。图标和文字都可点，触发更新；更新中图标无限旋转 1.5s（`:37-41、162-179`）。
- 有 `extra`：`SpeedOutlined` + `shared.labels.usedTotal` + `已用 / 总量`（`upload+download` / `total`，经 `parseTraffic`）。`expire > 0` 时再一行 `EventOutlined` + `shared.labels.expireTime` + `YYYY-MM-DD`。
- 进度：百分比文字 + 高 8px、圆角 4px 的 `LinearProgress`。总量 ≤ 0 时百分比为 0。上限 100。轨道 primary @ 12%（`:210-227`）。

**空态** `EmptyProfile`（`:236-264`）：整块可点，hover `action.hover`，圆角 8px。60px `CloudUploadOutlined`。主文案是 `profiles.page.actions.import` 与 `profiles.page.title` 拼接（中文即「导入订阅」）。副文案 `profiles.components.card.labels.clickToImport`。点击去 `/profile`。

## 6. 当前节点卡 `CurrentProxyCard`

`src/components/home/current-proxy-card.tsx`。标题 `home.components.currentProxy.title`。

**图标井**：有节点时显示信号图标，颜色见设计令牌「信号格」；tooltip 为「状态文案: 延迟格式化值」。无节点时 `SignalNone` disabled，tooltip `...status.noProxyNode`（`:959-972`）。

**右侧动作**

| 控件 | 行为 |
| --- | --- |
| `NetworkCheckRounded` | tooltip / aria `...actions.refreshDelay`。直连模式或没有节点时禁用。点击对该组做延迟测试 |
| 排序图标 | 三态循环：默认 `SortRounded`、延迟 `AccessTimeRounded`、名称 `SortByAlphaRounded`。tooltip 用 `proxies.page.tooltips.sortDefault / sortDelay / sortName` |
| 「代 理」按钮 | `layout.components.navigation.tabs.proxies`，尾图标 `ChevronRight`，`navigate('/proxies')` |

**内容三态**

1. 内核数据未到：`py: 4` 的空盒，高 24px（`:1014-1015`）。
2. 无当前节点且（直连或没有选中组）：居中 `...labels.noActiveNode`（`:1131-1140`）。
3. 否则：摘要条 + 两个下拉。

摘要条（`:1018-1087`）：底 primary @ 5%，边 primary @ 10%，圆角 4px。

- 节点名；没有名字时 `...labels.noActiveNode`。
- 副行：节点 `type`；全局模式 Chip `...labels.globalMode`（primary）；直连模式 Chip `...labels.directMode`（success）；能力 Chip 描边：`UDP` `TFO` `XUDP` `MPTCP` `SMUX`（有则显示）。
- 右侧延迟 Chip：`formatDelay` + `convertDelayColor`。直连或没有当前节点时不显示。

下拉：

- 「代理组」`...labels.group`。全局或直连时禁用。选项为可选组名，当前项 `selected`。
- 「节点」`...labels.proxy`。直连时禁用。选项带延迟色 Chip（`PersistentProxySelect`，约 `:453-520`）。组被固定（`fixed`）时有固定提示，文案 `proxies.page.labels.delayCheckReset`（需在节点下拉里核对 `fixed` 分支，组件 `PersistentProxySelect`）。

排序作用于节点下拉选项，不改变代理页的排序存储。规则与代理页相同：0 默认序、1 按延迟（已测优先，其次超时、错误、测试中、未测）、2 按名称 `localeCompare`（`:198-211` 与 `use-filter-sort.ts:85-113`）。

## 7. 网络设置卡

壳在 `home.tsx:375-387`：标题 `home.page.cards.networkSettings`，图标 `DnsOutlined`，`iconColor="primary"`，无 action。

内容 `ProxyTunCard`（`proxy-tun-card.tsx`）。两个 Tab，选中态写入 `localStorage['clash-verge-proxy-active-tab']`，缺省 `'system'`。

| Tab | 图标 | 文案 key | 角标 |
| --- | --- | --- | --- |
| `system` | `ComputerRounded` | `settings.sections.system.toggles.systemProxy` | 系统代理真实开启时显示 |
| `tun` | `TroubleshootRounded` | `settings.sections.system.toggles.tunMode` | `enable_tun_mode && isTunModeAvailable` |

Tab 视觉与代理模式按钮同一套：选中 primary 底 + 白字 + 字重 600 + elevation 2，未选中纸色；hover 上移 1px；选中时底部 2×9px 的 primary 竖线。最大宽 160px，`flex: 1`。角标 8px 圆点，选中时白色，未选中 `success.main`，位于右上 8px（`:37-94`）。

描述条：宽 95% 居中，caption，secondary 色，1px primary 边框，纸色底，圆角 4px。右侧 14px `HelpOutlineRounded`（透明度 0.7）Tooltip：

| Tab | 描述 | Tooltip |
| --- | --- | --- |
| 系统代理开 | `home.components.proxyTun.status.systemProxyEnabled` | `...tooltips.systemProxy` |
| 系统代理关 | `...systemProxyDisabled` | 同上 |
| TUN 无服务 | `...tunModeServiceRequired` | `...tooltips.tunMode` |
| TUN 开 | `...tunModeEnabled` | 同上 |
| TUN 关 | `...tunModeDisabled` | 同上 |

描述以**系统代理的实际指示**为准，不是「请求打开」的配置（`proxy-tun-card.tsx:146-147`）。

下方开关区：primary @ 4% 底，圆角 8px，内放 `ProxyControlSwitches`，只渲染当前 Tab 对应的一行（`:224-241`）。

`ProxyControlSwitches`（`src/components/shared/proxy-control-switches.tsx`）一行结构：

- 左：播放/暂停图标（开 success，关 disabled）+ 15px 标签 + `SettingsRounded`（打开对应对话框）+ 可选额外图标。
- 右：`Switch`。乐观更新；失败或返回 `false` 时回滚。
- 开启中的行：底 `success.main` @ 7%。禁用时整体透明度 0.6。
- 系统代理：核心未运行时拒绝打开，notice `settings.feedback.errors.sysproxy.coreNotReady`（`:150-155`）。关闭不依赖核心。
- TUN：服务不可用时打开会走安装服务流程（`:158` 起），而不是直接写入。齿轮打开 `TunViewer` / `SysproxyViewer`，字段见 `pages/settings.md`。

## 8. 代理模式卡

壳 `home.tsx:390-402`：标题 `home.page.cards.proxyMode`，`RouterOutlined`，`iconColor="info"`。

`ClashModeCard`（`clash-mode-card.tsx`）三个按钮，从左到右：

| 模式 | 图标 | 文案 |
| --- | --- | --- |
| rule | `MultipleStopRounded` | `home.components.clashMode.labels.rule` |
| global | `LanguageRounded` | `...global` |
| direct | `DirectionsRounded` | `...direct` |

选中：primary 底、contrast 字、字重 600、elevation 2，底部 2×16px primary 竖线（比网络 Tab 更长）。未选中：纸色、regular。hover 上移 1px；active 下移 1px。点已选中的模式直接返回（`:89-90`）。

切换：乐观改按钮，失败回滚并 `showNotice.error`，成功后写入缓存再刷新（`:89-110`）。

描述条样式同网络卡，无帮助图标。文案：

- 已解析模式：`home.components.clashMode.descriptions.{rule|global|direct}`
- 仍在加载：不换行空格 `\u00A0`
- 加载完仍无模式：`...errors.communication`

模式来源顺序：控制器配置 → 后端模式 → 运行时配置（`:77-81`）。未知模式字符串不当作当前模式。

## 9. 流量统计卡

壳：标题 `home.page.cards.trafficStats`，`SpeedOutlined`，`iconColor="secondary"`，整行（`home.tsx:285-295`）。

`EnhancedTrafficStats`（`enhanced-traffic-stats.tsx`）。

- 图表：`verge.traffic_graph !== false` 时显示（默认 true）。高 130px，圆角 8px，边框 divider @ 20%，点击在平滑 / 线性之间切换（`:181-199`，切换实现在 `EnhancedCanvasTrafficGraph.toggleStyle`）。图例文案 `home.components.traffic.legends.upload / download`。页面不可见时不画。
- 统计格：`Grid spacing={1}`，列 `xs/sm` 8、`md` 12，每格 `size={4}`（窄屏一行两个，宽屏一行三个）。

| 格 | 图标 | 文案 key | 色 | 单位 |
| --- | --- | --- | --- | --- |
| 上传速度 | `ArrowUpwardRounded` | `home.components.traffic.metrics.uploadSpeed` | secondary | `{unit}/s` |
| 下载速度 | `ArrowDownwardRounded` | `...downloadSpeed` | primary | `{unit}/s` |
| 活跃连接 | `LinkRounded` | `...activeConnections` | success | 空 |
| 累计上传 | `CloudUploadRounded` | `shared.labels.uploaded` | secondary | 流量单位 |
| 累计下载 | `CloudDownloadRounded` | `shared.labels.downloaded` | primary | 流量单位 |
| 内核内存 | `MemoryRounded` | `...memoryUsage` | error | 流量单位，仅 `enable_memory_usage !== false`（默认 true） |

格样式（`CompactStatCard`，`:60-80`）：无 elevation，圆角 8px，底色该色 @ 5%，边该色 @ 15%，padding 8px。图标井 32px 圆，底 @ 10%。标题 caption secondary，数值 body1 加粗。本卡这些格没有 `onClick`，hover 不变。

流量换算走 `parseTraffic`（B / KB / MB / GB…）。图表内部采样、FPS 诊断串 `home.components.traffic.diagnostics` 属于调试叠加，默认界面不要求展示；若实现画布右上角有样式切换提示，文案是 `home.components.traffic.chartStyles.smooth / linear`。

## 10. 网站测试卡 `TestCard`

`test-card.tsx`。标题 `home.components.tests.title`，图标 `NetworkCheck`。

动作：

- 「全部测试」tooltip `tests.page.actions.testAll`，发 `emitTestAll()`。
- 「新建」tooltip `tests.modals.test.title.create`，打开 `TestViewer` 的创建态。

列表：网格 `repeat(auto-fill, minmax(80px, 1fr))`，gap 8px，最大高度 180px，可纵向滚动。项可拖拽排序（`@dnd-kit` PointerSensor）。默认四项：Apple、GitHub、Google、YouTube，URL 与 SVG 见 `test-card.tsx:41-66`。首次 `verge.test_list` 为空时写入这份默认列表。

`TestItem`（`src/components/test/test-item.tsx`）展示图标、名称；单项测试、编辑、删除。编辑弹窗 `TestViewer`（`src/components/test/test-viewer.tsx`）字段以 `tests.modals.test` 为准：名称、URL、图标。删除立即 `patchVerge`。拖拽先改本地再在空闲时提交。

测试进行中的状态由 `TestItem` 内部维护。v2.5.9 修过「失败后一直显示测试中」，实现时失败必须回到可再次点击的终态，不能停在 loading。

本卡引用 `tests` 文案，对照表在 `i18n/tests.md`。测试页 `/unlock` 是另一页，见 `outlines/unlock.md`。

## 11. IP 信息卡

`ip-info-card.tsx`。标题 `home.components.ipInfo.title`。

- 加载失败：`...errors.load` 或 `...errors.loadWithDetails`（带 `{{message}}`）。无可用服务：`...errors.noServices`。
- 成功：左侧国旗（国家码转 emoji）+ 国家名，未知用 `...labels.unknown`。IP 行默认掩码 `••••••••••`，眼睛按钮切换明文（`VisibilityOutlined` / `VisibilityOffOutlined`）。ASN 显示 `AS{n}`，没有则 `N/A`。
- 右侧：ISP（`organization`）、组织（`asn_organization`）、位置（城市与地区用 `, ` 连接）、时区。
- 底栏：小开关「自动刷新」`...labels.autoRefresh`，状态在 `localStorage['ip-info-auto-refresh']`。开启时标签后追加倒计时 `: Ns` 或 `...`（`:394-397`）。右侧 caption：`国家码, 经度两位, 纬度两位`，缺省 `N/A`。
- 自动刷新关闭后仍可手动刷新（v2.5.9 新增，开关只停自动；手动入口若在卡片错误态的重试按钮上，以 `ip-info-card.tsx` 里触发 `refetch` 的按钮为准）。查询 `staleTime: Infinity`，不随窗口焦点刷新（`:421-428`）。

## 12. Clash 信息卡

`clash-info-card.tsx`。标题 `home.components.clashInfo.title`，`DeveloperBoardOutlined`，`iconColor="warning"`。`clashConfig` 为空时内容为空。

行（左 secondary body2，右 medium body2，行间 `Divider`，`Stack spacing={1.5}`）：

| 文案 | 值 |
| --- | --- |
| `...fields.coreVersion` | 内核版本，空则 `-` |
| `...fields.systemProxyAddress` | `systemProxyAddress` |
| `...fields.mixedPort` | 展示用混合端口（可能与配置端口不同） |
| `...fields.uptime` | `H:MM:SS`，时不补零，分秒补零（`:18-23`） |
| `...fields.rulesCount` | 规则条数 |

## 13. 系统信息卡

`system-info-card.tsx`。标题 `home.components.systemInfo.title`，`InfoOutlined`，`iconColor="error"`。`verge` 未到时不渲染。右上角齿轮 `navigate('/settings')`，title `...actions.settings`。

| 行 | 值 | 交互 |
| --- | --- | --- |
| `...fields.osInfo` | `系统名 + 版本`，版本若已带系统名前缀则去掉 | 无 |
| `...fields.autoLaunch` | Chip：开 `shared.statuses.enabled` success filled；关 `...disabled` outlined | 点击切换 `enable_auto_launch` |
| `...fields.runningMode` | 见下表 | Sidecar（含管理员+Sidecar）时可点，下划线，hover 透明度 0.7，点击安装服务并重启核心 |
| `...fields.lastCheckUpdate` | 本地化时间，没有则 `-` | 下划线可点。无更新成功 notice `settings.components.verge.advanced.notifications.latestVersion`；有更新 info `shared.feedback.notifications.updateAvailable` 并跳转设置 |
| `...fields.vergeVersion` | `v{package.json version}` | 无 |

运行模式文案与图标（`:139-200`）：

| 条件 | 文案 key | 图标 |
| --- | --- | --- |
| `NotRunning` | `...badges.notRunning` | `PauseCircleOutlined` secondary |
| 管理员且非 Sidecar | `...badges.adminServiceMode` | 管理员 primary + 服务 success |
| 仅管理员 | `...badges.adminMode` | `AdminPanelSettingsOutlined` |
| Sidecar | `...badges.sidecarMode` | `ExtensionOutlined` info |
| 其余（服务模式） | `...badges.serviceMode` | `DnsOutlined` success |

`auto_check_update` 为真且本地没有「上次检查」记录时，进入首页 5 秒后自动查一次（`:64-73`）。
