# 日志（规格，P2）

路由 `/logs`。参照 v2.5.9 `160bb4fe`。页面 `src/pages/logs.tsx` 的 `LogPage`。行 `src/components/log/log-item.tsx`。数据 `src/hooks/use-log-data.ts`，状态 `src/hooks/use-clash-log.ts`。文案：页面 `logs.*` 见 `i18n/logs.md`，级别 `shared.filters.logLevels.*` 见 `i18n/shared.md`。

> 本版修正旧大纲：级别下拉是 **5 项** `all / debug / info / warn / err`，没有 `silent`（`logs.tsx:164-175`）。`silent` 只出现在设置页「日志等级」。

## 1. 结构

```text
BasePage full  标题 logs.page.title
│  contentStyle: height 100%、flex 列、overflow auto
├── header（flex，gap 16px）
│   ├── IconButton small inherit 暂停/继续：开启时 PauseCircleOutlineRounded（title 暂停），暂停时 PlayCircleOutlineRounded（title 继续）
│   ├── IconButton small inherit 顺序：SwapVertRounded，降序时 scaleY(-1)，过渡 0.2s
│   └── Button small contained「清除」shared.actions.clear
├── 工具行（pt 8px，mb 4px，mx 10px，高 39，flex 居中）
│   ├── BaseStyledSelect 级别（宽 120，高 33.375，mr 8px）
│   └── BaseSearchBox
├── 有结果 → VirtualList（estimateSize 50，flex 1）→ LogItem
└── 无结果 → BaseEmpty
```

出处：`logs.tsx:90-204`。本页**没有**回到顶部按钮。

## 2. 页头按钮

| 按钮 | 行为 | title / aria-label |
| --- | --- | --- |
| 暂停 / 继续 | 切换 `clashLog.enable`。暂停后停止订阅日志 WebSocket，已显示的日志保留 | 开启中 `shared.actions.pause`；暂停中 `shared.actions.resume` |
| 顺序 | 切换 `clashLog.logOrder` `asc`↔`desc`；默认 `asc`（旧在上、新在下） | 当前降序时 `logs.actions.showAscending`（按时间正序）；当前升序时 `logs.actions.showDescending`（按时间倒序）——提示的是**点击后的效果** |
| 清除 | `refreshGetClashLog(true)`：清空前端列表，不清内核 | 按钮文字 `shared.actions.clear` |

状态持久化在 `localStorage['clash-log']`，默认 `{ enable: true, logLevel: 'INFO', logFilter: 'all', logOrder: 'asc' }`（`use-clash-log.ts:3-8`）。

## 3. 级别筛选

`BaseStyledSelect`，值写 `clashLog.logFilter`（持久化）。

| value | 文案 key | 显示 |
| --- | --- | --- |
| `all` | `shared.filters.logLevels.all` | ALL |
| `debug` | `shared.filters.logLevels.debug` | DEBUG |
| `info` | `shared.filters.logLevels.info` | INFO |
| `warn` | `shared.filters.logLevels.warn` | WARN |
| `err` | `shared.filters.logLevels.error` | ERROR |

过滤：`all` 全显示，否则 `log.type.includes(value)`（`logs.tsx:52-54`）。所以 `warn` 命中 `warning`，`err` 命中 `error`。

与「内核日志等级」的关系：WebSocket 订阅级别取 `clashLog.logLevel`（设置页「日志等级」修改时同步写入，`setting-clash.tsx:254-259`）。内核层只推送 ≥该等级的日志（`use-log-data.ts:15-21`）：

| 内核等级 | 前端接收的 type |
| --- | --- |
| DEBUG | debug / info / warning / error |
| INFO（默认） | info / warning / error |
| WARNING | warning / error |
| ERROR | error |
| SILENT | 无 |

内核等级变化时重新拉取初始日志并重连（`use-log-data.ts:136-151`）。

## 4. 搜索

`BaseSearchBox`。被匹配的文本是 `` `${time} ${type} ${payload}` `` **转小写后**的字符串（`logs.tsx:45-51`）。

命中的片段在行内高亮（`log-item.tsx`）：用同一组开关在**原文**上构造正则（非正则模式下转义；全字加 `\b`；大小写开关决定 `g` / `gi`；正则非法时退回转义文本）。高亮底色 亮 `#ffeb3b90` / 暗 `#ffeb3b40`，圆角 2，padding `0 2px`。

## 5. 日志行 `LogItem`

```text
04-21 13:05:22  INFO
[TCP] 127.0.0.1:52311 --> www.google.com:443 match DomainSuffix(google.com) using Proxy[HK 01]
```

| 部位 | 样式 |
| --- | --- |
| 行容器 | padding `8px 0`，左右外边距 12px，行高 1.35，14px（0.875rem），主题字体（非等宽），底边 1px `palette.divider`，文字可选中 |
| 第一行 | 时间（`text.secondary`）+ 8px + 级别 |
| 级别 | inline-block，大写，600，圆角 2，居中。颜色：`error`/`err` → `error.main`；`warning`/`warn` → `warning.main`；`info`/`inf` → `info.main`；**`debug` 无专门颜色，继承 `text.primary`** |
| 第二行 | payload，`text.primary`，`overflow-wrap: anywhere`（长串强制断行） |

时间格式 `MM-DD HH:mm:ss`（`use-log-data.ts:101`、`services/cmds.ts:185`）。

## 6. 实时刷新与数量

- WebSocket `/logs?level=…`（`MihomoWebSocket.connect_logs`）。连接成功后先调 `getClashLogs()` 取内核已缓存的历史（只在列表为空时填充），再接收实时推送（`use-log-data.ts:115-129`）。
- 实时日志先进缓冲，**50ms** 合并一次写入列表；同一批次共用一个时间戳（批次第一条到达时的**本地时间**，不是内核时间）（`:62-82`）。
- 前端最多保留 **1000** 条，超出丢最旧（`:10、23-47`）。
- 断线：收到 `Websocket error` 文本后 1s 重连（`use-mihomo-ws-subscription.ts:11`）。
- 自动滚动：升序且视口距底 ≤20px 时，新日志到来后平滑滚到最后一条；降序时新日志在顶部，不自动滚（`logs.tsx:64-73、192-197`）。
- 日志订阅**不受页面可见性影响**（不像连接页），只受暂停开关控制。

## 7. 空态

过滤后为 0 条显示 `BaseEmpty`（收件箱 4em + `shared.statuses.empty`「空空如也」，次要色 75%）。暂停且已清除时同样是这个空态。

## 8. 亮 / 暗差异

| 项 | 亮 | 暗 |
| --- | --- | --- |
| 内容区 | `#ECECEC` | `#1e1f27` |
| 级别下拉底 | `#fff` | 透明 |
| ERROR 色 | `#FF3B30` | `#FF453A` |
| WARN 色 | `#FF9500` | `#FF9F0A` |
| INFO 色 | `#007AFF` | `#0A84FF` |
| 高亮底 | `#ffeb3b90` | `#ffeb3b40` |

## 9. 截图核对清单（`ui.NewTester`）

建议数据：6 条日志，info×3、warning×1、error×1、debug×1。

- [ ] `logs-light/dark`：标题「日志」；页头 暂停图标 → 上下箭头图标 →「清除」实心
- [ ] 暂停后图标变为播放（▶），tooltip「继续」
- [ ] 级别下拉默认 `ALL`，展开 5 项：ALL / DEBUG / INFO / WARN / ERROR，无 SILENT
- [ ] 工具行高 39，左右缩进 10px
- [ ] 每条两行：`MM-DD HH:mm:ss` 灰色 + 级别大写粗体；第二行 payload 自动换行
- [ ] 级别颜色：INFO 蓝、WARNING 橙、ERROR 红、DEBUG 正文色
- [ ] 选 WARN 只剩 warning；选 ERROR 只剩 error
- [ ] 搜索 `google`：命中片段黄色高亮
- [ ] 切倒序：箭头图标上下翻转，最新在最上
- [ ] 清除后：收件箱 +「空空如也」
- [ ] 英文：`Logs`；顺序按钮 tooltip `Newest first` / `Oldest first`；暂停 `Pause` / `Resume`；清除 `Clear`

## 10. 待确认

1. 搜索先把文本转小写再交给匹配器，所以开启「区分大小写」后输入含大写字母的关键词将**永远不匹配**，但高亮按原文区分大小写。这是原版 bug。**与 Verge 有意不同**：clash-go 不照搬，开启「区分大小写」后按原文大小写匹配，输入大写能匹配到大写内容（主程 2026-11-03 拍板，附测试）。
2. 级别下拉选 `debug` 但内核等级是 INFO 时，列表恒为空（内核不推 debug）。原版无提示，保留。
3. `getClashLogs` 历史行在 Service 模式与 Sidecar 模式来源不同（`core/manager/state.rs:153-158`），clash-go 内嵌 mihomo 只有一种来源，历史条数和格式需实测。
4. 级别标签圆角 2 但没有背景色，圆角肉眼不可见；保留即可。
