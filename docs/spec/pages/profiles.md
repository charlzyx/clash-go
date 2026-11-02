# 订阅

优先级 **P0**。路由 `/profile`（不是 `/profiles`）。参照 v2.5.9 `160bb4fe`。

入口 `src/pages/profiles.tsx` 的 `ProfilePage`。卡片 `src/components/profile/profile-item.tsx`，壳 `profile-box.tsx`。

## 1. 结构

```text
BasePage full title=profiles.page.title  contentStyle.height=100%
├── header：普通工具 或 批量工具
├── 导入行（高 36，左右 10px）
│   ├── 链接输入
│   ├── 导入
│   └── 新建
└── 滚动区 height=calc(100% - 48px)
    ├── 订阅网格（可拖拽排序）
    ├── Divider
    └── 两列：全局扩展配置 Merge | 全局扩展脚本 Script
ProfileViewer
ConfigViewer（运行时配置，只读）
```

出处：`profiles.tsx:694-983`。

## 2. 页头

### 2.1 普通模式

按钮都是 `IconButton size="small"`，`gap` 8px。

| 按钮 | 色 | title key | 行为 |
| --- | --- | --- | --- |
| `CheckBoxOutlineBlankRounded` | inherit | `profiles.page.batch.title` | 进入批量 |
| `RefreshRounded` | inherit | `profiles.page.actions.updateAll` | 更新所有远程订阅 |
| `TextSnippetOutlined` | inherit | `...viewRuntimeConfig` | 打开只读运行时配置 `ConfigViewer` |
| `LocalFireDepartmentRounded` | primary | `...reactivate` | `onEnhance(true)` 重新激活 |
| `ClearRounded` | warning | `profiles.page.feedback.tooltips.forceRefreshStaleData` | **仅** `error` 或数据 `isStale` 时出现。2s 脉冲透明度 1↔0.5。强制刷新 |

### 2.2 批量模式

| 控件 | 行为 |
| --- | --- |
| 复选图标 | 全选 `CheckBoxRounded`；部分 `IndeterminateCheckBoxRounded`；全不选 `CheckBoxOutlineBlankRounded`。title 在全选时是 `...batch.actions.deselectAll`，否则 `selectAll` |
| `DeleteRounded` | `color="error"`。`selected` 为空时 `disabled`。title `...batch.actions.delete` |
| 「完成」 | `Button size="small" variant="outlined"`，`...batch.actions.done`，退出批量并清空选择 |
| 右侧计数 | `text.secondary`：`...batch.summary.selected` + 数字 + `...batch.summary.items` |

批量时卡片左侧多一枚复选框，点击复选不触发切换订阅。右键菜单里的「删除」在批量模式下改为切换勾选，不弹确认（`profile-item.tsx:512-521`）。

## 3. 导入行

`BaseStyledTextField`，placeholder `profiles.page.importForm.placeholder`。

- 输入为空：尾部粘贴 `ContentPasteRounded`，title `profiles.page.importForm.actions.paste`，读剪贴板写入输入。
- 输入非空：尾部 `ClearRounded`，title `shared.actions.clear`。
- Enter 且不是 IME 组合中：有 URL、未禁用、未在 loading 时导入。
- 「导入」`contained` `small` 圆角 6px，文案 `profiles.page.actions.import`。无 URL 或 `disabled` 时禁用；请求中 `loading`。
- 「新建」始终可点，圆角 6px，文案 `shared.actions.new`，打开 `ProfileViewer.create()`。

URL 必须以 `http://` 或 `https://` 开头，否则错误 `profiles.page.feedback.errors.invalidUrl`。导入失败会再试一次走 Clash 代理（notice `...importRetry` / `...importFail`）。成功 `...importSuccess`；导入了但可能要手动刷新时 `...importNeedsRefresh`。

拖文件到窗口：只接受 YAML，否则 `...onlyYaml`。用文件名（路径最后一段，缺省 `New Profile`）创建本地订阅（`profiles.tsx:123-135`）。

## 4. 订阅网格

`profiles.tsx:887-939`。

```text
repeat(auto-fill, minmax(clamp(260px, calc((100% - 24px) / 4), 480px), 1fr))
```

注释：在列宽不超过 480px 前最多 4 列，24px 是 3 个 `gap: 1`。`px: 0.5`。`DragDropProvider` + PointerSensor + KeyboardSensor。拖完失败 notice `profiles.page.feedback.errors.reorderFailed`。

卡片 `ProfileItem`：

- `selected`：当前订阅，或正在切换的目标（`switchTarget ?? profiles.current`）。
- `activating`：激活列表里，或切换加载计时到点后的目标。切换时先等一会儿才出加载遮罩（`switchLoadingTimer`），避免闪一下。卸载页面取消尚未发出的切换。
- 点击卡片：非激活中则选中并切换。激活中点击被吃掉。
- 右键：在指针位置打开菜单。菜单上再次右键关闭。

### 4.1 卡片视觉 `ProfileBox`

`profile-box.tsx`。

- padding `8px 16px`，圆角 8px，左对齐，指针。
- 底色亮 `#ffffff` / 暗 `#282A36`。
- 选中：左边 3px `primary.main`，宽度 `calc(100% + 3px)`，左移 3px。标题 `h2` 变为 `primary.main`。未选中标题 `text.primary`。
- 正文色：亮 `text.secondary`；暗 `text.secondary` @ 65%。选中与否不改变正文色。

激活遮罩：绝对铺满，圆角继承，`zIndex 10`，`backdrop-filter: blur(2px)`，底 `rgba(0,0,0,0.1)`，居中 20px `CircularProgress`，1.5s pulse（`profile-item.tsx:660-681`）。

### 4.2 卡片内容

1. 第一行：批量复选（仅批量）+ `DragIndicatorRounded`（`cursor: move`，`text.primary`）+ 标题。标题 18px / 600 / 行高 26，单行省略，`title` 为全名。有 URL 时右上角刷新 `RefreshRounded`，绝对定位 `top: -1` `right: -5`，padding 3px。loading 时 1s 线性无限旋转并 `disabled`。激活中点击刷新无效。刷新走 `onUpdate(1)`（按订阅自身代理选项更新）。
2. 第二行高 26：有描述显示描述（14px）；否则有 URL 显示来源 host。右侧可点的相对时间（`dayjs.fromNow()`）或「下次更新」。hover 底边变 primary、字变 primary。点击在「上次」和「下次」之间切换。
   - tooltip：看下次时 `profiles.components.profileItem.tooltips.showLast`；看上次时是 `shared.labels.updateTime: {绝对时间}` 加换行再加 `...showNext`。
   - 下次更新文案：失败 `...status.lastUpdateFailed`；小于 1 分钟 `...nextUp <1m`；小于 1 小时 `...nextUp {n}m`；否则 `...nextUp {h}h {m}m`；无计划 `...noSchedule`；自动更新关 `...autoUpdateDisabled`；异常 `...unknown`。
3. 有流量 `extra`：左 `已用/总量`（tooltip `shared.labels.usedTotal`），右到期（tooltip `shared.labels.expireTime`），14px。无 extra：右对齐 12px 绝对更新时间（tooltip `shared.labels.updateTime`）。
4. `LinearProgress`。`total > 0` 才可见（否则 opacity 0，占位仍在）。值为用量百分比。

### 4.3 右键菜单

`minWidth 120`，dense，过渡 225ms。删除项 `error.main`。

远程（有 URL）顺序（`profile-item.tsx:440-524`）：

| 项 | key | 禁用 |
| --- | --- | --- |
| 首页 | `profiles.components.menu.home` | 仅当有 `home` 字段才出现 |
| 使用 | `...select` | |
| 分享二维码 | `...shareQrCode` | |
| 编辑信息 | `...editInfo` | |
| 编辑文件 | `...editFile` | |
| 编辑规则 | `...editRules` | 无 `option.rules` |
| 编辑节点 | `...editProxies` | 无 `option.proxies` |
| 编辑代理组 | `...editGroups` | 无 `option.groups` |
| 扩展覆写配置 | `...extendConfig` | 无 `option.merge` |
| 扩展脚本 | `...extendScript` | 无 `option.script` |
| 打开文件 | `...openFile` | |
| 更新 | `...update` | 不走代理（type 0） |
| 更新（代理） | `...updateViaProxy` | 有 self_proxy 用自身，否则系统代理（type 2） |
| 删除 | `shared.actions.delete` | 非批量时先确认 |

本地文件菜单（`:525-585`）没有首页、二维码、更新、更新（代理），其余相同，禁用规则相同。

「使用」强制选中。切换成功 `profiles.page.feedback.notifications.profileSwitched`；忙 `...switchBusy`；失败 `...errors.switchFailed`。重新激活成功 `...profileReactivated`。

删除确认 `BaseDialog`（`profile-item.tsx:959` 起）：标题 `profiles.modals.confirmDelete.title`，正文 `...message`（此操作不可逆），确定 `shared.actions.confirm`，取消 `shared.actions.cancel`，宽 `xs 320 / sm 420`，正文可选中。当前正在用的订阅被删后的行为以页面 `onDelete` 为准：删的是当前项会走增强/切换失败提示 `...deleteFailed`。

批量删除成功 `...batchDeleted`。

## 5. 全局扩展

分割线宽 `calc(100% - 32px)`，颜色 `divider`。下面 `Grid`：`xs` 各占 12，`sm` 起各占 6。

两个 `ProfileMore`：

| id | 标题 key | 角标 |
| --- | --- | --- |
| `Merge` | `profiles.components.more.global.merge` | `...chips.merge`（字面量 Merge） |
| `Script` | `...global.script` | `...chips.script` |

Script 卡片带 `chainLogs['Script']`，可看脚本日志。保存且内容变化时静默重新激活（`onEnhance(false)`）。

编辑器顶部 `EnhanceHint` 说明执行顺序，文案在 `profiles.modals.editor.enhance.*`：

1. 订阅原文 → 应用设置 → 全局扩展配置 → 全局扩展脚本 → 订阅扩展配置 → 订阅扩展脚本 → 应用设置回写。
2. 应用设置接管的字段以应用设置为准。扩展写入这些字段会被丢弃，notice `profiles.page.feedback.notices.discardedKeys`（`{{keys}}`）。

编辑器动作「格式化文档」`profiles.modals.editor.actions.format`。只读时 `profiles.modals.editor.messages.readOnly`。只读 Chip `shared.labels.readOnly`（运行时配置）。

## 6. 新建 / 编辑信息 `ProfileViewer`

`profile-viewer.tsx:257-472`。`BaseDialog` 宽 375，`maxHeight 80%`，底 padding 0。保存 `shared.actions.save`，取消 `shared.actions.cancel`。保存中 `loading`。标题：新建 `profiles.modals.profileForm.title.create`，编辑 `...edit`。

| 字段 | key / 控件 | 何时出现 |
| --- | --- | --- |
| 类型 | `...fields.type`，Select：`remote` → `...types.remote`，`local` → `...types.local`。自动聚焦 | 始终 |
| 名称 | `shared.labels.name` | 始终 |
| 描述 | `...fields.description` | 始终 |
| 选择文件 | `profiles.components.fileInput.chooseFile`（`FileInput`） | 仅新建且 local。选中后若名称为空则用文件名 |
| 订阅链接 | `...fields.subscriptionUrl`，多行 | remote |
| User-Agent | `...fields.userAgent`，placeholder `clash-verge/v{version}` | remote |
| HTTP 超时 | `...fields.httpTimeout`，数字，placeholder `60`，后缀 `shared.units.seconds` | remote |
| 更新间隔 | `...fields.updateInterval`，数字，后缀 `shared.units.minutes` | remote |
| 使用系统代理 | `...fields.useSystemProxy`，Switch | remote |
| 使用内核代理 | `...fields.useClashProxy`，Switch | remote |
| 允许无效证书 | `...fields.acceptInvalidCerts`，Switch | remote |
| 允许自动更新 | `...fields.allowAutoUpdate`，Switch，缺省视为开 | remote |

校验：没类型 `...errors.typeRequired`；remote 没 URL `...errors.urlRequired`；编辑缺 uid `...errors.uidMissing`。保存失败 `...errors.saveFailed`。创建失败会再用 Clash 代理试一次（`...feedback.notifications.creationRetry` / `creationSuccess`）。

更新间隔 `> 0` 且 `< 1440` 分钟时，helper 为警告色：`...warnings.frequentUpdate`，`{{minutes}}` = 1440。不阻止保存（`profile-viewer.tsx:37-38、380-398`）。

关闭后约 500ms 重置表单。

## 7. 其余弹层

| 弹层 | 组件 | 要点 |
| --- | --- | --- |
| 编辑文件 | `EditorViewer` | YAML，路径 `profile:{uid}.yaml`。脏标记。保存失败重新加载 |
| 编辑规则 | `RulesEditorViewer` | 依赖 groups/merge uid 与 `option.rules` |
| 编辑节点 | `ProxiesEditorViewer` | 标题 `profiles.modals.proxiesEditor.title`。占位 `...placeholders.multiUri`。动作：前置 `...actions.prepend`、后置 `...actions.append` |
| 编辑代理组 | `GroupsEditorViewer` | 标题 `profiles.modals.groupsEditor.title`。字段见下表。名称空 `...errors.nameRequired`，重名 `...errors.nameExists`。前置/后置 `...actions.prepend/append` |
| 订阅扩展配置 | `EditorViewer` | 标题 `...menu.extendConfig`，YAML，`EnhanceHint stage=profileMerge` |
| 订阅扩展脚本 | `EditorViewer` | 标题 `...menu.extendScript`，JavaScript，`stage=profileScript` |
| 脚本日志 | `LogViewer` | 标题 `profiles.modals.logViewer.title` |
| 二维码 | `QrViewer` | 标题 `profiles.modals.qrViewer.title` |
| 运行时配置 | `ConfigViewer` | 与设置里同一个只读查看器，Chip `shared.labels.readOnly` |

代理组编辑字段（`profiles.json` 的 `modals.groupsEditor`，表单在 `groups-editor-viewer.tsx`）：

- 文本：类型、组名、图标、引入代理、引入代理集合、健康检查地址、期望状态码、间隔、最大失败次数、出站接口、路由标记、过滤、排除过滤、排除类型。
- 开关：`includeAll`、`includeAllProxies`、`includeAllProviders`、`lazy`、`disableUdp`、`hidden`。key 分别在 `fields.*` 与 `toggles.*`。

节点编辑器支持多行 URI，也接受 Base64（占位文案已说明）。

## 8. 反馈一览

错误 key 均在 `profiles.page.feedback.errors.*`：`invalidUrl`、`onlyYaml`、`createFailed`、`deleteFailed`、`enhanceFailed`、`importFailed`、`openFailed`、`readFailed`、`reorderFailed`、`switchFailed`、`updateFailed`。

通知：`importRetry`、`importFail`、`importNeedsRefresh`、`importSuccess`、`profileSwitched`、`profileReactivated`、`switchBusy`、`batchDeleted`。

强制刷新完成 `...notices.forceRefreshCompleted`；失败 `...emergencyRefreshFailed`（`{{message}}`）。
