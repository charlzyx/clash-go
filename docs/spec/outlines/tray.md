# 托盘菜单（大纲，P2）

菜单定义 `src-tauri/src/core/tray/menu_def.rs`，组装 `src-tauri/src/core/tray/mod.rs` 约 651–868 行。文案在 `crates/clash-verge-i18n/locales/{zh,en}.yml` 的 `tray.*`，不是前端 `src/locales`。本轮不生成托盘文案表。

v2.5.9 已去掉托盘「测试延迟」。托盘测速与代理页不一致的问题是靠删入口修的，不要再加。

## 顺序

1. 打开面板 `tray.dashboard`（热键 `open_or_close_dashboard`）
2. 分隔线
3. 出站模式：默认子菜单 `tray.outboundModes`，标题带当前模式括号（`tray.rule|global|direct`）。`tray_inline_outbound_modes` 为真时三项平铺。三项是可勾选：`tray.ruleMode` / `globalMode` / `directMode`，热键与设置热键相同。
4. 分隔线
5. 订阅子菜单 `tray.profiles`：当前配置列表，当前项勾选
6. 代理：`default` 为子菜单 `tray.proxies`，下面每个代理组一个子菜单；`inline` 把组直接铺在主菜单。其他显示模式不放代理项。
7. 分隔线
8. 系统代理 `tray.systemProxy`（勾选，热键）
9. TUN `tray.tunMode`（服务不可用时禁用，热键）
10. 分隔线
11. 轻量模式 `tray.lightweightMode`（勾选）
12. 打开目录子菜单 `tray.openDir`：配置目录、内核目录、日志目录、应用日志、内核日志
13. 更多 `tray.more`：复制环境变量、关闭全部连接、重启内核、重启应用、`tray.vergeVersion` + 版本号
14. 分隔线
15. 退出 `tray.exit`。macOS 无自定义热键时用 `Cmd+Q`

## 图标点击

左键按下时按设置里的 `tray_event`（`mod.rs:871` 起）：主窗口、托盘菜单、切换系统代理、切换 TUN、或 `disable` 什么都不做。Linux 设置页不提供这一项。

待补：各 `tray.*` 的中英文表、代理组子菜单里节点的勾选与排序、模板图标与速度文字（`tray_speed.rs`，macOS 状态栏宽度）。
