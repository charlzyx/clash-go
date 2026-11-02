# 代理

优先级 **P0**。路由 `/proxies`。参照 v2.5.9 `160bb4fe`。

入口 `src/pages/proxies.tsx` 的 `ProxyPage`。列表 `src/components/proxy/proxy-groups.tsx` 的 `ProxyGroups`。

## 1. 结构

```text
BasePage full contentStyle.height=100%
├── title
│   ├── 普通：proxies.page.title.default（代理组）
│   └── 链式：proxies.page.title.chainMode + WarningRounded tooltip
├── header
│   ├── ProviderButton
│   ├── ButtonGroup 规则 / 全局 / 直连
│   └── Button 链式代理
└── 主体
    ├── mode=direct → BaseEmpty「直连模式」
    ├── 数据未到 → BaseLoading
    ├── 空态 → ProxyEmptyState
    └── 否则
        ├── 非链式：StickyVirtualList + ProxyGroupNavigator
        └── 链式：ProxyChain 面板 + 节点列表
```

`BasePage` 的 `full` 去掉内容区内边距（`base-page.tsx:35`，`proxies.tsx:147-149`）。

## 2. 页头

### 2.1 模式按钮

三个模式，顺序固定 `rule | global | direct`（`proxies.tsx:23、173-184`）。文案 `proxies.page.modes.{rule|global|direct}`。`ButtonGroup size="small"`。当前模式 `contained`，其余 `outlined`。`textTransform: capitalize`。

点击调用 `patchClashMode`。失败 `showNotice.error`，成功后刷新 clash 配置。未知模式字符串会自动改回 `rule`（`proxies.tsx:140-144`）。

与首页模式卡是同一后端动作，两处状态应一致。

### 2.2 链式代理开关

文案 `proxies.page.actions.toggleChain`。关：`outlined` + `LanOutlined`。开：`contained` + `LanRounded`。左边距 8px。

- 从关到开：不立即开启，先弹确认框（`proxies.tsx:95-98、210-224`）。
- 从开到关：清除运行时链式配置，失败则保持开启并报错；成功后写 `localStorage['proxy-chain-mode-enabled']='false'`。
- 开关状态持久化在 `localStorage['proxy-chain-mode-enabled']`，值为 `'true'` / 其他。

确认框 `BaseDialog`，宽 `xs 320 / sm 420`：

| 项 | key |
| --- | --- |
| 标题 | `proxies.page.actions.toggleChain` |
| 正文 | `proxies.page.chain.warning` |
| 确定 | `proxies.page.chain.confirmEnable` |
| 取消 | `proxies.page.chain.keepDisabled` |

确定后只把本地开关打开，不在这一步连接。

标题在链式模式下是「链式代理模式」加 `WarningRounded`（`color="warning"`），tooltip 同一段 warning 文案。

### 2.3 代理集合按钮

`ProviderButton`（`src/components/proxy/provider-button.tsx`）。图标 `StorageOutlined`。点击打开对话框，标题 `proxies.page.provider.title`。

- 列表项：名称、类型角标（secondary @ 50% 边 / 80% 字，10px）、流量、过期 `YYYY-MM-DD`（无则 `-`）。
- 单项「更新」`proxies.page.provider.actions.update`：该项 `LinearProgress`，成功 notice `proxies.feedback.notifications.provider.updateSuccess`（`{{name}}`），失败 `...updateFailed`（`{{name}}` `{{message}}`）。
- 「更新全部」：没有 provider 时 info `...none`；逐个更新，最后刷新并 `...allUpdated`。
- `providerState === 'unavailable'` 时列表不可用（按钮仍在，对话框内按该状态展示）。

## 3. 列表前置状态

`resolveProxyListState`（`proxy-empty-state-model.ts:33-61`），在 `ProxyGroups` 里分支（`proxy-groups.tsx:637-653`）。

| 条件 | 界面 |
| --- | --- |
| `mode === 'direct'` | `BaseEmpty`，图标 `InboxRounded` 4em，文案 `proxies.page.messages.directMode`。链式开关即使开着也不进链式列表 |
| 订阅、代理视图或运行模式仍在加载 | `BaseLoading` |
| 没有任何 `remote`/`local` 订阅 | `ProxyEmptyState` reason `no-subscriptions` |
| 有订阅但当前激活项不在其中 | reason `inactive-subscription` |
| 列表渲染后没有任何可见行，且核心未运行或代理视图报错 | reason `core-unavailable`（`resolveEmptyListReason`） |
| 列表渲染后为空，核心正常 | reason `no-proxy-info` |

空态是描边 `Alert severity="warning"`，最大宽 640，垂直水平居中，padding 16px（`proxy-empty-state.tsx:61-76`）。

| reason | 标题 key | 说明 key | 按钮 |
| --- | --- | --- | --- |
| `no-subscriptions` | `proxies.page.empty.noSubscriptions.title` | `...description` | 「前往订阅」`contained` → `/profile` |
| `inactive-subscription` | `...inactiveSubscription.*` | | 同上 |
| `core-unavailable` | `...coreUnavailable.*` | | 「重启内核」`contained`（loading 中禁用「打开日志」）+「打开日志目录」`outlined` |
| `no-proxy-info` | `...noProxyInfo.*` | | 三个按钮都有 |

重启中按钮 `loading`，失败由统一 mutate 报错。重启后刷新代理和系统状态。

## 4. 普通列表（规则 / 全局）

`NormalProxyGroups`（`proxy-groups.tsx:372`）。`StickyVirtualList`：分组头估计高 76，节点行估计高 64。滚动位置按模式存，节流 500ms，回来时最多 30 帧尝试恢复。

### 4.1 画哪些组

`use-render-list.ts:260-265`：

- `rule` 与 `global`：全部 `proxyView.groups`。
- `direct` 不会走到这里（上一节已空态）。若模式数据异常只剩 GLOBAL 组，则只画 `proxyView.global`，且节点区强制展开（`headState.open || !useRule`，`useRule` 仅 rule/global 为真，`:260、296`）。

隐藏组（`group.hidden`）的头不计入「有内容」，但其已展开的成员仍算内容（`hasRenderableItems`，`:57-59`）。

### 4.2 行类型

| type | 内容 | 组件 |
| --- | --- | --- |
| 0 | 分组头 | `ProxyRender` 头 |
| 1 | 仅非 rule/global 时的第二行工具条 | `ProxyHead` |
| 2 | 单列节点 | `ProxyItem`，高 40 |
| 3 | 组展开但过滤后为空 | `InboxRounded` 2.5em + `proxies.page.empty.noProxies` |
| 4 | 多列节点 | 网格高 56，`gap` 8px，列数 `repeat(n, 1fr)`，`ProxyItemMini` |

### 4.3 列数

`use-render-list.ts:91-98`。`verge.proxy_layout_column` 缺省按 6 处理。

- 配置为 1–5：固定该列数。
- 否则（含 6=自动，见设置杂项）：按 **body 宽度**：`>1920 → 5`，`>1450 → 4`，`>1024 → 3`，`>=600 → 2`，否则 1。

### 4.4 分组头

`proxy-render.tsx:212-243`。padding `4px 8px`。按钮高 100%，圆角 8px，底色亮 `#ffffff` / 暗 `#282A36`。点击切换 `open`。收起前若头已滚过视口顶，先把该头滚回顶部（`proxy-groups.tsx:528-546`）。

粘住且展开时阴影 `0 4px 8px rgba(0,0,0,0.2) !important`（全局 `box-shadow: none` 的例外）。

头内从左到右（工具在左时左右对调，`verge.proxy_group_tools_position === 'left'`）：

1. 组图标（`enable_group_icon` 默认 true）：URL / data / 内联 SVG，宽 32，右边距 12，URL 与 data 圆角 6px。
2. 组名 + 类型角标 + 当前节点名 `group.now`。
3. 工具条 `ProxyGroupTools`，高 36，`gap` 4px。
4. 节点数 Chip：primary @ 10% 底、primary 字，tooltip `proxies.page.labels.proxyCount`。工具在右时 Chip 在工具左侧。
5. 展开 `ExpandLessRounded` / 收起 `ExpandMoreRounded`。

分组头可拖拽排序（`useProxyGroupHeaderLayout`）。拖拽中元素用全局 `[data-dnd-dragging]` 阴影。

### 4.5 工具条状态

`proxy-group-tools.tsx` 与直连用的 `proxy-head.tsx` 图标一致。每个按钮 `size="small"` `color="inherit"`，`preventDefault` + `stopPropagation`，避免触发展开。除过滤/URL 外，若组是收起的，先 `flushSync` 展开再执行。

| 按钮 | 图标 | title key | 行为 |
| --- | --- | --- | --- |
| 定位 | `MyLocationRounded` | `proxies.page.tooltips.locate` | 滚到当前节点，居中、smooth |
| 测速 | `NetworkCheckRounded` | `...delayCheck` | 测全组。自定义 URL 非空且当前不是过滤框时，展开 URL 输入。URL 非法则 warning `proxies.feedback.warnings.invalidTestUrl`（该条英文文案在 zh 文件里也是英文）且不测 |
| 排序 | 默认 `SortRounded` / 延迟 `AccessTimeRounded` / 名称 `SortByAlphaRounded` | `sortDefault` / `sortDelay` / `sortName` | `(sortType + 1) % 3` |
| 测速 URL | 开 `WifiTetheringRounded`，关 `WifiTetheringOffRounded` | `...delayCheckUrl` | 切换 `textState='url'`，并聚焦输入 |
| 细节 | 开 `VisibilityRounded`，关 `VisibilityOffRounded` | 开着 title 是 `showBasic`（隐藏细节），关着是 `showDetail` | 切换 `showType` |
| 过滤 | 开 `SearchOffRounded`，关 `SearchRounded` | `...filter` | 切换 `textState='filter'`。组收起时仅在将要打开过滤时才展开 |

`textState === 'filter'`：`BaseSearchBox` 占剩余宽度。输入防抖 **600ms** 才写入过滤条件；关掉过滤框或卸载时立刻 flush（`proxy-group-tools.tsx:75-93`）。`ProxyHead`（非 rule/global）不防抖。

`textState === 'url'`：`TextField` `size="small"`，placeholder `proxies.page.placeholders.delayCheckUrl`，输入 `py: 0.65` `px: 1`。实际测试 URL：输入去空格后的值，否则组自带 `testUrl`，否则 `verge.default_latency_test`，再否则 `http://cp.cloudflare.com/generate_204`。

工具条左右：过滤/URL 输入在 `side==='right'` 时位于按钮左侧，在 `'left'` 时位于按钮右侧。

头状态按**当前订阅 uid** 存在 `localStorage['proxy-head-state']`（`use-head-state.ts:21-31`）。默认：`open false`、`showType true`、`sortType 0`、过滤空、三个匹配开关 false、`textState null`、`testUrl ''`。

### 4.6 过滤与排序

`use-filter-sort.ts`。查询先 `trim`。空查询不过滤。

1. 正则 `delay([=<>])(\d+|timeout|error)` 不区分大小写。`error` 当成 `1e5`，`timeout` 当成 `3000`。未测（delay < 0）排除。`=error` 为 `>= 1e5`；`=timeout` 为 `< 1e5 且 >= 3000`；`=` 相等；`<` 实为 `<=`；`>` 实为 `>=`。
2. 否则 `type=(.*)`，类型子串、不区分大小写。
3. 否则按名称。`BaseSearchBox` 三个开关：区分大小写、全字、正则。正则非法则结果为空（`compiled.isValid === false`）。

排序：`0` 保持原序；`1` 按 `compareByDelay`（已测升序，然后超时、错误、测试中、未测，见 `src/utils/delay.ts:44-57`）；`2` 名称 `localeCompare`。超时阈值用 `verge.default_latency_timeout`，非法则 10000。

### 4.7 节点 `ProxyItem`（单列）

`proxy-item.tsx`。

- 行高 40，下边距 8，圆角 4px。底色亮 `#ffffff` / 暗 `#24252f`。
- `dense`。未解析成员 `disabled`，不可点。
- 选中：宽 `calc(100% + 3px)`，左移 3px，左边 3px 条（亮 `primary.main` / 暗 `primary.light`），底 primary @ 15%（亮）或 35%（暗）。
- 主文案 14px：`name`；`showType && now` 时追加 ` - ${now}`。然后类型角标；细节开时再追加有则显示的 `UDP` `XUDP` `TFO` `MPTCP` `SMUX`。角标见设计令牌。
- 右侧：
  - delay `-2`：`BaseLoading`。
  - 非测试中：hover 出「检测」`shared.actions.check`，默认 `display:none`，hover 行时显示（已有延迟数字时改显示数字、藏检测）。点击只测这一条，不选中。
  - delay `> 0`：彩色数字，颜色见设计令牌。hover 底 primary @ 15%。
  - 已选中且没有正延迟、也不是测试中：`CheckCircleOutlineRounded` 16px。hover 时让位给「检测」。
- 预设节点（`isPreset`）不显示右侧图标区。
- 点击行：仅组类型 `Selector` / `URLTest` / `Fallback` 且成员可交互时切换 `now`（`proxy-groups.tsx:475-478`）。失败打日志并刷新列表。

`ProxyItemMini` 是多列紧凑版，同样有选中左边条、延迟、类型；点选规则相同。

### 4.8 分组导航

`ProxyGroupNavigator`，只在 rule / global（`proxy-groups.tsx:611-618`）。竖向首字母，点击滚到该组头（`align: start`）。`enable_hover_jump_navigator` 默认 true 时，悬停 `hover_jump_navigator_delay` 毫秒后跳转，默认 `DEFAULT_HOVER_DELAY = 280`（`proxy-group-navigator.tsx:11`）。首字符取组名显示字符（`getGroupDisplayChar`）。

## 5. 链式模式

`ChainProxyGroups` + `ProxyChain`（`proxy-chain.tsx`）。开启后拉 `getRuntimeProxyChainConfig(exitNode)`。出口节点在 `localStorage['proxy-chain-exit-node']`，没有则配置数据为空字符串。

链表面板：

- 头：警告图标 + 标题区。已有节点时显示删除 `proxies.page.actions.clearChainConfig`（error 色，hover error.light 加 `20` 十六进制透明）。清空运行时配置并删 `proxy-chain-group` / `proxy-chain-exit-node` / `proxy-chain-items`。
- 连接按钮最小宽 90：未连接 `success` + `Link`，文案 `...actions.connect`；已连接 `error` + `LinkOff`，文案 `...disconnect`；进行中 `...connecting` 且禁用。
- 禁用条件：连接中；或未连接且（节点少于 2，或有节点没有 `recordId`，或 global 模式没有 global 组，或非 global 没有选中组）。少于 2 个时 title 为 `proxies.page.chain.minimumNodes`。
- 提示条：1 个节点时 `warning` + `...minimumNodesHint`，否则 `info` + `...instruction`。
- 空链：文案 `...chain.empty`，下面 `outlined`「关闭」`shared.actions.close`，点击退出链式模式（`proxy-chain.tsx:697-712`）。
- 非空：可拖拽排序。第一项标 `...chain.entryNode`，最后一项（且多于 1 个）标 `...chain.exitNode`。重复添加 notice `...duplicateNode`。连接失败 `...connectFailed`，断开失败 `...disconnectFailed`。

节点列表在链面板之外，只含 `Selector` 与 `URLTest` 组（`proxy-groups.tsx:201-207`）。rule 模式默认选第一组。点节点是加入链，不是切换 `now`。过滤条件在链式渲染里被传成空字符串、排序 0（`use-render-list.ts:234-239`），即链式节点列表不做组内过滤。

## 6. 延迟展示

见 `design-tokens.md` 第 1.4 节。测试中写 `-2`，批量并发最多 10，组内从第 3 个起随机错开 0–200ms（`delay.ts:364-384`）。单条失败记 `1e6`（错误色）。组内测速进行中不因结果回来而重排，避免鼠标下的行跳动（`delay.ts:371` 附近注释）。

## 7. 搜索框控件

`BaseSearchBox`（`base-search-box.tsx`）：

- 亮色输入底 `#fff`。
- 三个 SVG 开关：`match_case`、`match_whole_word`、`use_regular_expression`。激活 path 填充 `primary.light`，未激活 `#A7A7A7`，`aria-label` 为 `active` / `inactive`。
- 有文字时显示 `ClearRounded` 清空。
