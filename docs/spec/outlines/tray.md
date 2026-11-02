# 托盘菜单（规格，P2）

参照 v2.5.9 `160bb4fe`。菜单项定义 `src-tauri/src/core/tray/menu_def.rs`，组装 `src-tauri/src/core/tray/mod.rs`（`create_tray_menu` 587–869 行），图标点击 `on_tray_icon_event`（871–916），菜单点击 `on_menu_event`（918 起）。macOS 状态栏速率 `src-tauri/src/core/tray/speed_task.rs` + `src-tauri/src/utils/tray_speed.rs` + `src-tauri/src/utils/speed.rs`。文案在 `crates/clash-verge-i18n/locales/{zh,en}.yml` 的 `tray.*`（不是前端 `src/locales`），对照表见 `i18n/tray.md`。

v2.5.9 已去掉托盘「测试延迟」。托盘测速与代理页不一致的问题是靠删入口修的，不要再加。

> 本版修正旧大纲：第 1 项文案是「仪表板 / Dashboard」，不是「打开面板」。

## 1. 菜单顺序

`menu_items` 拼装见 `mod.rs:831-866`。`—` 为原生分隔线。

```text
仪表板                         tray.dashboard          热键 open_or_close_dashboard
—
出站模式 (规则)  ▸              tray.outboundModes + " (" + tray.rule|global|direct + ")"
   ✓ 规则模式                   tray.ruleMode           热键 clash_mode_rule
     全局模式                   tray.globalMode         热键 clash_mode_global
     直连模式                   tray.directMode         热键 clash_mode_direct
—
订阅  ▸                         tray.profiles
   ✓ 当前订阅名 / 其他订阅名 …  （CheckMenuItem，当前项勾选）
代理  ▸                         tray.proxies（显示模式 default 且有组时）
   组名 ▸                       每个可见代理组一个子菜单
      ✓ 节点名   | 123ms        当前节点勾选
—
✓ 系统代理                      tray.systemProxy        热键 toggle_system_proxy
✓ TUN 模式                      tray.tunMode            热键 toggle_tun_mode；TUN 不可用时禁用
—
  轻量模式                      tray.lightweightMode    勾选态 = 是否处于轻量模式；热键 entry_lightweight_mode
打开目录  ▸                     tray.openDir
   配置目录 / 内核目录 / 日志目录 / 应用日志 / 内核日志
更多  ▸                         tray.more
   复制环境变量 / 关闭所有连接 / 重启 Clash 内核 / 重启应用 / Verge 版本 2.5.9
—
退出                            tray.exit               热键 quit；macOS 未设时用 Cmd+Q
```

### 1.1 出站模式

- 默认是子菜单，标签 `{tray.outboundModes} ({当前模式短名})`，如「出站模式 (规则)」「Outbound Modes (Rule)」（`mod.rs:692-711`）。当前模式未知时按「规则」显示。
- 设置 `tray_inline_outbound_modes = true`（设置页「将出站模式显示在托盘一级菜单」）时，三项直接平铺在第一条分隔线之后，不再有子菜单（`:833-841`）。
- 三项都是 CheckMenuItem，勾选 = 当前 `mode`（`:665-690`）。点击调用 `change_clash_mode`。

### 1.2 订阅

- 子菜单 `tray.profiles`，列出 `profiles.items` 里**同时有 uid 和 name** 的每一项，按配置文件顺序；当前项勾选（`config/profiles.rs:397-411`、`mod.rs:439-457`）。
- 点击切换到该订阅（`feat::toggle_proxy_profile`）。

### 1.3 代理

由设置 `tray_proxy_groups_display_mode`（「托盘代理组显示模式」）决定（`mod.rs:722-733、845-853`）：

| 值 | 设置页文案 | 托盘表现 |
| --- | --- | --- |
| `default`（缺省） | 默认 | 一级菜单一项「代理 ▸」，下面每组一个子菜单；没有组时整项不出现 |
| `inline` | 一级菜单 | 每个组的子菜单直接放在「订阅」下面 |
| `disable` | 禁用 | 不出现代理项，也不拉取代理数据 |

组与节点（`mod.rs:459-549`）：

- `global` 模式只显示 `GLOBAL` 组；其他模式显示除 `GLOBAL` 外的全部组。`hidden` 组不显示。没有 `all` 成员或成员为空的组跳过。
- 组顺序按运行时配置 `proxy-groups` 的先后；不在配置里的组排在后面，保持原顺序。
- 节点项文字：`{节点名}   | {延迟}`（**节点名后 3 个空格**，竖线后 1 个空格）。延迟取该节点最后一次历史：`0` 或 ≥10000 显示 `-ms`，无历史 `-ms`，否则 `{n}ms`。
- 当前节点（组的 `now`）勾选。点击切换节点（`feat::switch_proxy_node`）。
- 拉取代理数据超时 1000ms，超时则本次不显示代理组（`:599-604`）。

### 1.4 系统代理 / TUN

- 系统代理 CheckMenuItem，勾选 = `enable_system_proxy`。
- TUN CheckMenuItem，勾选 = `enable_tun_mode`；`enabled` = 运行态 `tun_capable()`，服务不可用时灰显（`mod.rs:203-207、744-751`）。

### 1.5 打开目录 / 更多

| 菜单 | 子项（顺序） | 动作 |
| --- | --- | --- |
| 打开目录 `tray.openDir` | `tray.confDir` 配置目录、`tray.coreDir` 内核目录、`tray.logsDir` 日志目录、`tray.appLog` 应用日志、`tray.coreLog` 内核日志 | 用系统打开目录 / 最新日志文件 |
| 更多 `tray.more` | `tray.copyEnv` 复制环境变量、`tray.closeAllConnections` 关闭所有连接、`tray.restartClash` 重启 Clash 内核、`tray.restartApp` 重启应用、`tray.vergeVersion` + 空格 + 版本号 | 版本项可点击但无动作（`on_menu_event` 未处理） |

### 1.6 热键标注

菜单项右侧显示的快捷键来自设置 `hotkeys`（`"功能,键"` 列表）。含 `NUMPAD` 的键不标注；Windows 下含 `CMD` 的不标注（`mod.rs:416-437`）。macOS 退出项未设热键时固定标 `Cmd+Q`（`:822-827`）。

## 2. 图标点击

左键**按下**时按设置 `tray_event`（缺省 `main_window`）（`mod.rs:871-916`）：

| 值 | 行为 |
| --- | --- |
| `main_window` | 退出轻量模式或显示主窗口 |
| `tray_menu` | 左键直接弹菜单（`set_show_menu_on_left_click(true)`，`:169-186`） |
| `system_proxy` | 切换系统代理 |
| `tun_mode` | 切换 TUN |
| `disable` 及其他 | 无动作（记 warn 日志） |

右键始终弹菜单（系统默认）。托盘点击和菜单点击都经过限流器 `should_handle_tray_click`（`:402-408`），过快点击会被忽略。Linux 设置页不提供此项。

## 3. 图标

- macOS：设置 `tray_icon`（「托盘图标」）`monochrome`（缺省）/ `colorful`。单色时作为 **template image**，随菜单栏明暗自动反色（`mod.rs:96-112、260-272`）。
- 三种状态图标：常规、系统代理开、TUN 开（`IconKind::Common / SysProxy / Tun`）。可分别启用自定义图标（`common_tray_icon` 等，`find_target_icons`）。
- 资源（单色）：`src-tauri/icons/tray-icon-mono.ico`、`tray-icon-sys-mono-new.ico`、`tray-icon-tun-mono-new.ico`；彩色：`tray-icon.ico`、`tray-icon-sys.ico`、`tray-icon-tun.ico`。
- clash-go 品牌图标替换方案见待确认 1。

## 4. 悬停提示（tooltip）

`mod.rs:277-336`，固定 4 行：

```text
Clash Verge 2.5.9
系统代理: on
TUN: off
订阅: 当前订阅名
```

- 开关值固定英文 `on` / `off`，不翻译。TUN 显示 `enable_tun_mode && tun_capable()`。
- 无当前订阅时名称为 `None`。
- 标签 key：`tray.tooltip.systemProxy`、`tray.tooltip.tun`、`tray.tooltip.profile`。冒号为半角 `: `。

## 5. macOS 状态栏速率

设置 `enable_tray_speed`（「启用托盘速率」，仅 macOS，缺省关）（`mod.rs:338-352`、`layout-viewer.tsx`）。

| 项 | 值 | 出处 |
| --- | --- | --- |
| 文本 | 两行：上行上传、下行下载，各右对齐到 6 字符宽 | `tray_speed.rs:38-42` |
| 单位 | `B/s` `K/s` `M/s` `G/s` `T/s`，数值与单位**无空格** | `speed.rs:2-3` |
| 数值 | <1000 B/s 显示整数 `512B/s`；否则 1024 进制，<9.95 一位小数（`1.2K/s`），否则取整（`35K/s`）；取整后 ≥1000 进下一级 | `speed.rs:12-32` |
| 字体 | 系统等宽字体 9.5pt Regular，颜色 `controlTextColor` | `tray_speed.rs:17、46-47` |
| 行高 | 10pt，行距 0 | `tray_speed.rs:19-23` |
| 宽度 | 文本宽 + 30pt，最小 58pt（防止数值跳动） | `tray_speed.rs:25-28、127-128` |
| 数据 | 订阅流量流；断开 1s 重连；5s 无数据视为断流 | `speed_task.rs:11-12` |

速率显示在图标右侧（attributed title）；关闭时清空标题。

## 6. 菜单刷新

托盘菜单整份重建（`update_menu` → `create_tray_menu`），由模式、订阅、系统代理、TUN、节点切换等状态变化触发；图标与 tooltip 走 `update_part`。源码中调用点共 13 处，未逐一列出（见待确认 3）。

## 7. 亮 / 暗差异

- 菜单是系统原生菜单，明暗跟随 macOS 外观，应用主题不影响。
- 单色 template 图标随菜单栏自动反色；彩色图标不反色。
- 速率文字颜色 `controlTextColor`，随系统外观变化。

## 8. 截图核对清单

托盘是原生菜单，`ui.NewTester` 拍不到。验收用 macos-14 runner 上的 `screencapture` 或菜单结构快照（遍历菜单树导出文本）对照：

- [ ] 一级菜单顺序与 §1 一致，共 5 条分隔线（出站模式平铺时仍是 5 条）
- [ ] 首项「仪表板」；**没有**「测试延迟」
- [ ] 「出站模式 (规则)」子菜单三项，当前模式勾选；切到全局后标签变「出站模式 (全局)」
- [ ] 「订阅」子菜单当前项勾选
- [ ] 「代理」下组顺序与配置一致；节点项 `名称   | 123ms`，超时或未测为 `-ms`
- [ ] 全局模式下代理只剩 GLOBAL 组
- [ ] 系统代理 / TUN 勾选与实际一致；TUN 不可用时灰
- [ ] 「更多」最后一项「Verge 版本 x.y.z」（clash-go 文案见待确认 1）
- [ ] 「退出」右侧 `⌘Q`
- [ ] 悬停 tooltip 4 行，`on/off` 英文
- [ ] 单色图标在浅色、深色菜单栏下都清晰；开系统代理 / TUN 后图标切换
- [ ] 开启托盘速率：两行等宽小字，右对齐，宽度不随数值抖动
- [ ] 英文菜单：`Dashboard`、`Outbound Modes (Rule)`、`Profiles`、`Proxies`、`System Proxy`、`TUN Mode`、`Lightweight Mode`、`Open Directory`、`More`、`Exit`

## 9. 待确认

1. 品牌：tooltip 首行 `Clash Verge {version}` 与「Verge 版本 / Verge Version」菜单项是否改为 `Clash Go`，以及托盘图标用 mygo `favicon.svg` 生成单色 template 还是沿用 Verge 图标，需主程拍板。1:1 复刻与品牌替换冲突。
2. `profiles_preview` 只要求有 uid 和 name，`items` 里若含全局扩展（Merge / Script）类条目，也会出现在「订阅」子菜单里。需在原版带全局扩展的配置下核对。
3. 菜单重建的 13 个调用点与触发时机未逐一列出；节点延迟在测速后是否立即刷新到托盘，需实测。
4. 节点 id 用 `proxy_{组名}_{节点名}`，点击时按第一个 `_` 拆分（`mod.rs:993-1001`）。组名含 `_` 时会拆错。这是原版 bug，clash-go 应避免（建议用索引映射），但菜单外观不变。
5. 「Verge 版本」项可点但无动作，是否置灰由主程决定；原版未置灰。
