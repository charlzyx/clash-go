# 托盘文案对照（zh / en）

参照 Clash Verge Rev v2.5.9 `160bb4fe8eaef2bc626163d89622229e5d7c5b59`。

- 中文：`crates/clash-verge-i18n/locales/zh.yml`
- 英文：`crates/clash-verge-i18n/locales/en.yml`
- 加载：Rust 宏 `clash_verge_i18n::t!("tray.…")`（`src-tauri/src/core/tray/menu_def.rs`），**不是**前端 `src/locales`。

托盘菜单与悬停提示。结构见 `outlines/tray.md`。

| key | 中文 | English |
| --- | --- | --- |
| `tray.dashboard` | 仪表板 | Dashboard |
| `tray.ruleMode` | 规则模式 | Rule Mode |
| `tray.globalMode` | 全局模式 | Global Mode |
| `tray.directMode` | 直连模式 | Direct Mode |
| `tray.outboundModes` | 出站模式 | Outbound Modes |
| `tray.rule` | 规则 | Rule |
| `tray.direct` | 直连 | Direct |
| `tray.global` | 全局 | Global |
| `tray.profiles` | 订阅 | Profiles |
| `tray.proxies` | 代理 | Proxies |
| `tray.systemProxy` | 系统代理 | System Proxy |
| `tray.tunMode` | TUN 模式 | TUN Mode |
| `tray.closeAllConnections` | 关闭所有连接 | Close All Connections |
| `tray.lightweightMode` | 轻量模式 | Lightweight Mode |
| `tray.copyEnv` | 复制环境变量 | Copy Environment Variables |
| `tray.confDir` | 配置目录 | Configuration Directory |
| `tray.coreDir` | 内核目录 | Core Directory |
| `tray.logsDir` | 日志目录 | Log Directory |
| `tray.openDir` | 打开目录 | Open Directory |
| `tray.appLog` | 应用日志 | Application Log |
| `tray.coreLog` | 内核日志 | Core Log |
| `tray.restartClash` | 重启 Clash 内核 | Restart Clash Core |
| `tray.restartApp` | 重启应用 | Restart Application |
| `tray.vergeVersion` | Verge 版本 | Verge Version |
| `tray.more` | 更多 | More |
| `tray.exit` | 退出 | Exit |
| `tray.tooltip.systemProxy` | 系统代理 | System Proxy |
| `tray.tooltip.tun` | TUN | TUN |
| `tray.tooltip.profile` | 订阅 | Profile |

拼接规则：

| 位置 | 格式 | 例（zh / en） |
| --- | --- | --- |
| 出站模式子菜单标题 | `{tray.outboundModes} ({tray.rule\|global\|direct})` | 出站模式 (规则) / Outbound Modes (Rule) |
| 版本项 | `{tray.vergeVersion} {版本号}` | Verge 版本 2.5.9 / Verge Version 2.5.9 |
| 代理节点项 | `{节点名}   \| {延迟}`（名后 3 空格） | `HK 01   \| 123ms`、`HK 02   \| -ms` |
| tooltip | `Clash Verge {版本}\n{tooltip.systemProxy}: on\|off\n{tooltip.tun}: on\|off\n{tooltip.profile}: {订阅名\|None}` | 见 `outlines/tray.md` §4 |

托盘相关的设置页文案（「托盘图标」「启用托盘速率」「托盘代理组显示模式」「将出站模式显示在托盘一级菜单」等）在 `i18n/settings.md` 的 `settings.components.verge.layout.*`。
