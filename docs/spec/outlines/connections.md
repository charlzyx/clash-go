# 连接（规格，P1）

路由 `/connections`。参照 v2.5.9 `160bb4fe`。页面 `src/pages/connections.tsx` 的 `ConnectionsPage`。组件目录 `src/components/connection/`。数据 `src/hooks/use-connection-data.ts`。布局设置 `src/hooks/use-connection-setting.ts`。文案对照见 `i18n/connections.md`。

> 本版修正旧大纲两处：右下角 Fab 不是「回到顶部」，而是已关闭列表的「清除」；默认布局是**表格**，不是列表。

## 1. 结构

```text
BasePage full  标题 connections.page.title（whiteSpace: nowrap）
│  contentStyle: height 100%、flex 列、overflow hidden、borderRadius 8px、minHeight 0
├── header（flex，gap 16px）
│   ├── 「下载量: 12.3MB」shared.labels.downloaded + parseTraffic(downTotal)   mx 8px
│   ├── 「上传量: 4.56MB」shared.labels.uploaded + parseTraffic(upTotal)       mx 8px
│   ├── IconButton small inherit：表格时显示 TableRowsRounded（title 列表视图），列表时显示 TableChartRounded（title 表格视图）
│   └── Button small contained「关闭全部」shared.actions.closeAll（nowrap）
├── 工具行（sticky top 0，zIndex 2，pt 8px，mb 4px，mx 10px，minHeight 36，gap 8px）
│   ├── ButtonGroup（mr 8px）：「活跃 N」|「已关闭 N」，size small，当前 contained、另一个 outlined
│   ├── BaseStyledSelect 排序（仅列表布局）：默认 / 上传速度 / 下载速度
│   ├── BaseSearchBox（flex 1，占满剩余宽度）
│   └── 列设置 IconButton small ViewColumnRounded（仅表格布局且有数据；Tooltip 列设置）
├── 主体
│   ├── 过滤结果为空 → BaseEmpty（InboxRounded 4em +「空空如也」）
│   ├── 表格布局 → ConnectionTable
│   └── 列表布局 → VirtualList（estimateSize 56，key=活跃/已关闭，换 tab 时重建）→ ConnectionRowItem
├── ConnectionDetail（Snackbar，右下）
└── Zoom Fab「清除」（仅已关闭 tab 且列表非空）
```

出处：`connections.tsx:172-342`。

## 2. 页头

### 2.1 累计流量

- 文案：`{t('shared.labels.downloaded')}: ` + `parseTraffic(traffic?.downTotal || 0)`。上传同理（`connections.tsx:190-196`）。
- `parseTraffic` 返回 `[数值, 单位]`，React 直接拼接，**数值与单位之间没有空格**：如 `下载量: 12.3MB`、`上传量: 0B`。冒号后有一个空格。
- 数值规则（`src/utils/parse-traffic.ts`）：1024 进制，单位 `B KB MB GB TB …`；取整后 ≥1000 用 `toFixed(0)`，否则 `toPrecision(3)`。所以 0 显示为 `0.00B`，页头初值是 `下载量: 0.00B`。
- 数据来自流量 WebSocket（`useTrafficData({ enabled: pageVisible })`），页面不可见时不订阅。

### 2.2 布局切换

- 单个 IconButton。当前是表格时显示「切到列表」图标 `TableRowsRounded`，`titleAccess` = `shared.actions.listView`；当前是列表时显示 `TableChartRounded`，`titleAccess` = `shared.actions.tableView`（`connections.tsx:197-213`）。
- 值存 `localStorage['connections-setting']`，JSON `{ layout: 'table' | 'list' }`，**默认 `table`**（`use-connection-setting.ts:3`）。

### 2.3 关闭全部

`Button size=small variant=contained`。调用内核 `closeAllConnections`，`errorNotice: false`，失败不弹提示（`connections.tsx:156-161`）。不弹确认框。

## 3. 工具行

### 3.1 活跃 / 已关闭

- 文案 `connections.components.actions.active|closed` + 空格 + 条数，如「活跃 12」「已关闭 0」（`connections.tsx:237-254`）。数据未到时条数为空（`undefined` 渲染为空，按钮只显示「活跃 」）。
- 切换时：关闭详情、关闭列设置对话框、切换列表（`connections.tsx:136-144`）。

### 3.2 排序（仅列表布局）

`BaseStyledSelect`：宽 120，高 33.375，`mr 8px`，亮色底 `#fff`，暗色透明（`base-styled-select.tsx`）。菜单项文字 14px。

| value | key | 规则 |
| --- | --- | --- |
| `default` | `connections.components.order.default` | 按开始时间新→旧 |
| `uploadSpeed` | `connections.components.order.uploadSpeed` | 当前上传速度大→小 |
| `downloadSpeed` | `connections.components.order.downloadSpeed` | 当前下载速度大→小 |

排序选项只存在组件 state，不持久化，默认 `default`（`connections.tsx:90`）。

表格布局下：无搜索时**完全不排序**，保持内核返回顺序；有搜索时仍会套用当前（隐藏的）排序函数（`connections.tsx:114-128`）。表格自己的列头排序见 4.3。

### 3.3 搜索

`BaseSearchBox`（三开关：区分大小写 / 全字 / 正则，详见 `pages/proxies.md` §7）。有文字才过滤。匹配字段：`metadata.host`、`metadata.destinationIP`、`metadata.process`，任一命中即保留（`connections.tsx:120-125`）。不匹配链路、规则、端口。

## 4. 表格布局 `ConnectionTable`

自绘虚拟表格（`connection-table.tsx`），不是 MUI DataGrid。

### 4.1 尺寸

| 项 | 值 | 出处 |
| --- | --- | --- |
| 行高 | 40px，底边 1px `palette.divider` | `:33、271-283` |
| 单元格 | padding 8px，字号 13px，单行省略，右对齐列靠右 | `:286-302` |
| 表头 | sticky top 0，zIndex 2，底色 `palette.background.paper`（亮 `#F5F5F5` / 暗 `#2E303D`），底边 divider | `:780-793` |
| 表头文字 | 13px / 600，`text.secondary`，不可选中 | `:800-809` |
| 列宽拖拽柄 | 宽 6px，贴列右缘居中（`translateX(50%)`），`col-resize` 光标 | `:34、837-853` |
| 滚动容器 | 横纵都可滚，圆角 8px；表宽 = 可见列宽之和，最小 100% | `:762-779` |
| 虚拟化 | 上下各多画 6 行 | `:35` |

### 4.2 列

默认顺序、默认宽 / 最小宽（`connection-table.tsx:381-465`）。12 列默认**全部可见**（`visibility` 缺省为 true）。

| # | field | 表头 key | 默认宽 | 最小宽 | 对齐 | 单元格内容 |
| --- | --- | --- | --- | --- | --- | --- |
| 1 | `host` | `connections.components.fields.host` | 180 | 140 | 左 | `host‖destinationIP‖remoteDestination` + `:` + 端口 |
| 2 | `download` | `shared.labels.downloaded` | 76 | 60 | 右 | 累计下载，如 `1.23 MB`（**有空格**） |
| 3 | `upload` | `shared.labels.uploaded` | 76 | 60 | 右 | 累计上传 |
| 4 | `dlSpeed` | `connections.components.fields.dlSpeed` | 76 | 60 | 右 | 如 `0.00 B/s` |
| 5 | `ulSpeed` | `connections.components.fields.ulSpeed` | 76 | 60 | 右 | 同上 |
| 6 | `chains` | `connections.components.fields.chains` | 280 | 160 | 左 | chains **倒序**用 ` / ` 连接（出口在前… 见待确认 2） |
| 7 | `rule` | `connections.components.fields.rule` | 220 | 160 | 左 | `rule(rulePayload)`；无 payload 只写 rule |
| 8 | `process` | `connections.components.fields.process` | 180 | 140 | 左 | `process‖processPath` |
| 9 | `time` | `connections.components.fields.time` | 100 | 80 | 右 | 相对时间（dayjs `from`，每 5s 刷新） |
| 10 | `source` | `connections.components.fields.source` | 160 | 120 | 左 | `sourceIP:sourcePort` |
| 11 | `remoteDestination` | `connections.components.fields.destination` | 160 | 120 | 左 | `destinationIP‖remoteDestination` + `:` + 端口 |
| 12 | `type` | `connections.components.fields.type` | 120 | 80 | 左 | `type(network)`，如 `HTTPS(tcp)` |

流量格式 `formatConnectionTraffic`（`connection-row-view.ts:30-40`）：与 `parseTraffic` 同算法，但数值与单位之间**有空格**。

### 4.3 列头排序

点击表头循环：升序 `▲` → 降序 `▼` → 取消（无标记）（`:618-624、830-835`）。数值列按数值比，其余 `localeCompare`。排序状态不持久化；「重置为默认值」会清空。

### 4.4 列宽 / 顺序 / 可见性持久化

| 存储键 | 内容 |
| --- | --- |
| `localStorage['connection-table-widths']` | `{ field: px }`，拖拽实时写入，夹在 minWidth 内 |
| `localStorage['connection-table-visibility']` | `{ field: false }`，只记隐藏的列 |
| `localStorage['connection-table-order']` | field 数组；新增列补到末尾，未知列丢弃 |

### 4.5 列设置对话框 `ConnectionColumnManager`

- MUI `Dialog maxWidth="xs" fullWidth`，标题 `connections.components.columnManager.title`（`connection-column-manager.tsx:80-119`）。
- 每列一行：`ListItem` 圆角 4px、1px divider 边框、`px 8px py 4px`、行间距 8px。左 `Checkbox`，中列名（body2），右拖拽柄 `DragIndicatorRounded`（aria-label `...columnManager.dragHandle`），拖拽调整顺序。
- 至少保留 1 列：只剩一列可见时它的勾选框禁用。
- 底部：文字按钮「重置为默认值」`shared.actions.resetToDefault`（清可见性、顺序、列宽、排序）+ contained「关闭」`shared.actions.close`。`DialogActions` `px 24px pb 16px`。

### 4.6 交互

- 点击整行打开详情（§6）。表格行**没有**单条关闭按钮。
- 无右键菜单。
- 连接数变少导致滚动越界时，自动夹回最大 scrollTop（`:556-567`）。

## 5. 列表布局 `ConnectionRowItem`

`connection-row-item.tsx`。

| 项 | 值 |
| --- | --- |
| 行 | 最小高 56，flex 居中，gap 8，padding `6px 48px 6px 12px`，底边 1px `--divider-color` |
| 第一行 | host（同表格 host 列，含端口），14px，行高 1.4，单行省略 |
| 第二行标签 | 上边距 4，gap 4，可换行。每个标签：10px、行高 1.375、padding `0 4px`、1px `rgba(128,128,128,0.35)` 边框、圆角 4、单行省略 |
| 标签顺序 | network → type → process（有才显示）→ chains（有才显示）→ 相对时间 → `上传速度 / 下载速度`（任一 ≥100 B/s 才显示，如 `1.20 KB/s / 35.0 KB/s`） |
| 关闭按钮 | 仅活跃 tab。`IconButton small inherit` + `CloseRounded small`，绝对定位右 8 垂直居中，title / aria-label `connections.components.actions.closeConnection`；点击直接 `closeConnection(id)`，不确认 |
| 点击 | 点文字区打开详情；文字可选中 |

## 6. 详情 `ConnectionDetail`

不是对话框，是 MUI **Snackbar**，`anchorOrigin bottom/right`（`connection-detail.tsx:38-61`）。

- 内容框最大 520×480，超出纵向滚动；底色 `background.paper`，文字 `text.primary`。
- 已打开时再点别的行**不会**替换内容（`open` 时直接 return，`:29`）；点击外部（clickaway）或按 Esc 关闭（Snackbar 默认）。
- 打开时拍快照，**不随实时数据刷新**。
- 字段逐行 `<b>标签</b>: 值`，标签 `text.secondary` 粗体，值 `text.primary`，`word-break: break-all`（`:83-163`）：

| 顺序 | 标签 key | 值 |
| --- | --- | --- |
| 1 | `connections.components.fields.host` | `host:端口` |
| 2 | `shared.labels.downloaded` | `parseTraffic(download).join(' ')`，如 `1.23 MB` |
| 3 | `shared.labels.uploaded` | 同上 |
| 4 | `connections.components.fields.dlSpeed` | `… /s` 形式，如 `0.00 B/s` |
| 5 | `connections.components.fields.ulSpeed` | 同上 |
| 6 | `connections.components.fields.chains` | 倒序 ` / ` 连接 |
| 7 | `connections.components.fields.rule` | `rule(payload)` |
| 8 | `connections.components.fields.process` | `process(processPath)`；无 path 只写 process |
| 9 | `connections.components.fields.time` | `dayjs(start).fromNow()`（中文如「几秒前」，英文 `a few seconds ago`） |
| 10 | `connections.components.fields.source` | `sourceIP:sourcePort` |
| 11 | `connections.components.fields.destination` | destinationIP，没有则 remoteDestination（不带端口） |
| 12 | `connections.components.fields.destinationPort` | 端口 |
| 13 | `connections.components.fields.type` | `type(network)` |

- 活跃连接底部右对齐 contained 按钮「关闭连接」`...actions.closeConnection`，点击关闭连接并收起详情。已关闭连接不显示按钮。

## 7. 已关闭「清除」Fab

- 条件：已关闭 tab 且过滤结果非空；`Zoom` 进出（`connections.tsx:323-341`）。
- `Fab size=medium variant=extended color=primary`，绝对定位 `right 16`，`bottom` 表格 70 / 列表 16。图标 `DeleteForeverRounded small`（mr 8）+「清除」`shared.actions.clear`。
- 只清前端缓存的已关闭列表，不调内核。

## 8. 实时刷新

`use-connection-data.ts`：

- 订阅内核 `/connections` WebSocket；**仅页面可见时订阅**（`useVisibility`：document 可见且窗口可见），不可见即断开（`connections.tsx:85-98`）。
- 合帧节流 500ms：两次刷新至少间隔 500ms，只保留最新一帧（`:5、234-248`）。
- 断线 1s 后重连（`:6`）。
- 当前速度 = 本帧累计值 − 上一帧累计值（`:111-112`），首次出现的连接速度为 0。
- 从活跃消失的连接移入「已关闭」，最多保留 **500** 条，超出丢最旧（`:4、177-198`）。已关闭列表只在前端内存，换页后仍在（模块级变量），重启应用清空。
- 相对时间统一 5s 定时刷新（`connection-relative-time.tsx`）。

## 9. 亮 / 暗差异

| 项 | 亮 | 暗 |
| --- | --- | --- |
| 内容区 | `base-container` `#ffffff`，section `#ECECEC` | `#1e1f27` |
| 表头底 | `#F5F5F5` | `#2E303D` |
| 排序下拉底 | `#fff` | 透明 |
| 行分隔 | 6% 黑 | 6% 白 |
| 标签边框 | `rgba(128,128,128,0.35)` 两色相同 | 同左 |
| 详情底 | `#F5F5F5` | `#2E303D` |

## 10. 截图核对清单（`ui.NewTester`）

建议固定数据：3 条活跃连接（含一条无 host 只有 IP、一条有 process、一条速度 ≥100 B/s），2 条已关闭。

- [ ] `connections-table-light/dark`：标题「连接」；页头顺序 下载量 → 上传量 → 视图图标 → 「关闭全部」实心
- [ ] 下载量 / 上传量数值与单位无空格（`12.3MB`，零值 `0.00B`），冒号后一个空格
- [ ] 工具行：「活跃 3」实心、「已关闭 2」描边；表格布局**没有**排序下拉；搜索框占满；最右列设置图标
- [ ] 表头 12 列，顺序与 §4.2 一致；13px/600 次要色；表头底色是纸色而非白色
- [ ] 行高 40，单元格 13px；速度、流量、时间列右对齐；流量值数值与单位**有空格**
- [ ] 点表头一次出现 `▲`，两次 `▼`，三次消失
- [ ] 列设置对话框：12 行勾选 + 拖拽柄；只剩 1 列时勾选框灰；按钮「重置为默认值」「关闭」
- [ ] `connections-list-light/dark`：切换到列表后出现排序下拉（默认「默认」）；行高约 56；第二行小标签；右侧 × 按钮
- [ ] 速度标签只出现在速度 ≥100 B/s 的那一行
- [ ] 详情：右下浮层，13 个字段顺序正确，「关闭连接」右对齐；已关闭 tab 打开的详情无按钮
- [ ] 已关闭 tab：右下「清除」扩展 Fab（表格时离底 70，列表离底 16）；列表为空时 Fab 消失
- [ ] 空态：搜索无结果显示收件箱图标 +「空空如也」，列设置图标同时隐藏
- [ ] 英文：`Connections`、`Active 3`、`Closed 2`、`Close All`、表头 `Host / Downloaded / Uploaded / DL Speed / UL Speed / Chains / Rule / Process / Time / Source / Destination / Type`

## 11. 待确认

1. 链路列把 `chains` 倒序显示。mihomo 的 `chains` 顺序通常是「出口节点 … 策略组」，倒序后变成「策略组 / … / 节点」。需要真实连接核对显示方向。
2. 「当前速度」= 两帧累计差。mihomo `/connections` 推送间隔（默认约 1s）决定这是否等于「每秒」；前端节流 500ms 不改差值基准。原版未除以时间间隔。
3. Snackbar 详情的内边距、圆角用 MUI `SnackbarContent` 默认值（源码未覆盖），具体像素需对原版截图量取。
4. 本页源码里没有任何右键菜单。若 clash-go 要加，属于扩展，不在 1:1 范围。
