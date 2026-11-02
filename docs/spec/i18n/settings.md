# 设置文案对照（zh / en）

参照 Clash Verge Rev v2.5.9 `160bb4fe8eaef2bc626163d89622229e5d7c5b59`。

- 中文：`src/locales/zh/settings.json`
- 英文：`src/locales/en/settings.json`
- 加载：`src/services/i18n.ts` 以文件名作为顶层键，合并进 i18next 默认命名空间 `translation`。调用形式 `t('settings.…')`。

设置页及全部子对话框。

| key | 中文 | English |
| --- | --- | --- |
| `settings.page.actions.manual` | 使用手册 | Manual |
| `settings.page.actions.telegram` | Telegram 频道 | Telegram Channel |
| `settings.page.actions.github` | GitHub 项目地址 | Github Repo |
| `settings.page.actionsGroupLabel` | 设置相关链接 | Settings links |
| `settings.page.title` | 设置 | Settings |
| `settings.sections.system.title` | 系统设置 | System Setting |
| `settings.sections.system.toggles.tunMode` | 虚拟网卡模式 | Tun Mode |
| `settings.sections.system.toggles.systemProxy` | 系统代理 | System Proxy |
| `settings.sections.system.tooltips.silentStart` | 程序启动时以后台模式运行，不显示程序面板 | Start the program in background mode without displaying the panel |
| `settings.sections.system.fields.autoLaunch` | 开机自启 | Auto Launch |
| `settings.sections.system.fields.silentStart` | 静默启动 | Silent Start |
| `settings.sections.system.notifications.tunMode.autoDisabled` | 由于服务不可用，TUN 模式已自动关闭 | TUN Mode automatically disabled due to service unavailable |
| `settings.sections.system.notifications.tunMode.autoDisableFailed` | 自动关闭 TUN 模式失败 | Failed to disable TUN Mode automatically |
| `settings.sections.proxyControl.tooltips.systemProxy` | 修改操作系统的代理设置，如果开启失败，可手动修改操作系统的代理设置 | Enable to modify the operating system's proxy settings. If enabling fails, modify the operating system's proxy settings manually |
| `settings.sections.proxyControl.tooltips.tunMode` | TUN（虚拟网卡）模式接管系统所有流量，启用时无须打开系统代理 | Tun (Virtual NIC) mode: Captures all system traffic, when enabled, there is no need to enable system proxy. |
| `settings.sections.proxyControl.tooltips.tunUnavailable` | TUN 模式需要安装服务模式或管理员模式 | TUN requires Service Mode or Admin Mode |
| `settings.sections.proxyControl.actions.installService` | 安装服务 | Install Service |
| `settings.sections.proxyControl.actions.uninstallService` | 卸载服务 | Uninstall Service |
| `settings.sections.proxyControl.actions.switchToServiceMode` | 切换到服务模式 | Switch to service mode |
| `settings.sections.proxyControl.fields.systemProxy` | 系统代理 | System Proxy |
| `settings.sections.proxyControl.fields.tunMode` | 虚拟网卡模式 | Tun Mode |
| `settings.sections.proxyControl.messages.installedCheckProxy` | 服务已安装，内核已通过服务运行。请检查系统代理开关，若刚才的改动未生效请重新操作一次。 | The service is installed and the core is running through it. Check the system proxy switch and make your change again if it did not take. |
| `settings.sections.proxyControl.messages.installedCoreNotOnService` | 服务已安装，但内核未通过服务运行。请检查内核状态后重试。 | The service is installed, but the core is not running through it. Check the core and try again. |
| `settings.sections.proxyControl.messages.installedProxyRestored` | 服务已安装，内核已通过服务运行，刚才的系统代理改动已重新应用。 | The service is installed, the core is running through it, and your system proxy change has been applied. |
| `settings.sections.externalController.title` | 外部控制器监听地址 | External Controller |
| `settings.sections.externalController.fields.enable` | 启用外部控制器 | Enable External Controller |
| `settings.sections.externalController.fields.address` | 外部控制器监听地址 | External Controller |
| `settings.sections.externalController.fields.secret` | API 访问密钥 | Core Secret |
| `settings.sections.externalController.placeholders.address` | 必填 | Required |
| `settings.sections.externalController.placeholders.secret` | 建议设置 | Recommended |
| `settings.sections.externalController.tooltips.copy` | 复制到剪贴板 | Copy to clipboard |
| `settings.sections.externalController.messages.addressRequired` | 控制器地址不能为空 | Controller address cannot be empty |
| `settings.sections.externalController.messages.secretRequired` | 访问密钥不能为空 | Secret cannot be empty |
| `settings.sections.externalController.messages.copyFailed` | 复制失败 | Failed to copy |
| `settings.sections.externalController.messages.controllerCopied` | 控制器地址已复制到剪贴板 | Controller address copied to clipboard |
| `settings.sections.externalController.messages.secretCopied` | 访问密钥已复制到剪贴板 | Secret copied to clipboard |
| `settings.sections.externalCors.title` | 外部控制跨域设置 | External Cors Configuration |
| `settings.sections.externalCors.fields.allowPrivateNetwork` | 允许专用网络访问 | Allow private network access |
| `settings.sections.externalCors.fields.allowedOrigins` | 允许的来源 | Allowed Origins |
| `settings.sections.externalCors.placeholders.origin` | 请输入有效的网址 | Please enter a valid url |
| `settings.sections.externalCors.actions.add` | 添加 | Add |
| `settings.sections.externalCors.messages.alwaysIncluded` | 始终包含来源：{{urls}} | Always included origins: {{urls}} |
| `settings.sections.externalCors.tooltips.open` | 外部控制跨域设置 | External Cors Settings |
| `settings.sections.appearance.light` | 浅色 | Light |
| `settings.sections.appearance.dark` | 深色 | Dark |
| `settings.sections.appearance.system` | 系统 | System |
| `settings.sections.clash.title` | Clash 设置 | Clash Setting |
| `settings.sections.clash.form.fields.allowLan` | 局域网连接 | Allow LAN |
| `settings.sections.clash.form.fields.dnsOverwrite` | DNS 覆写 | DNS Overwrite |
| `settings.sections.clash.form.fields.ipv6` | IPv6 | IPv6 |
| `settings.sections.clash.form.fields.unifiedDelay` | 统一延迟 | Unified Delay |
| `settings.sections.clash.form.fields.logLevel` | 日志等级 | Log Level |
| `settings.sections.clash.form.fields.portConfig` | 端口设置 | Port Config |
| `settings.sections.clash.form.fields.external` | 外部控制 | External |
| `settings.sections.clash.form.fields.webUI` | 网页界面 | Web UI |
| `settings.sections.clash.form.fields.clashCore` | Clash 内核 | Clash Core |
| `settings.sections.clash.form.fields.openUwpTool` | UWP 工具 | Open UWP tool |
| `settings.sections.clash.form.fields.updateGeoData` | 更新 GeoData | Update GeoData |
| `settings.sections.clash.form.fields.tunnels.title` | 流量隧道管理 | Tunnel Management |
| `settings.sections.clash.form.fields.tunnels.localAddr` | 本地监听地址 | Local Listen Address |
| `settings.sections.clash.form.fields.tunnels.localPort` | 本地监听端口 | Local Listen Port |
| `settings.sections.clash.form.fields.tunnels.targetAddr` | 目标地址 | Target Address |
| `settings.sections.clash.form.fields.tunnels.targetPort` | 目标端口 | Target Port |
| `settings.sections.clash.form.fields.tunnels.proxyGroup` | 代理组 | Proxy Group |
| `settings.sections.clash.form.fields.tunnels.proxyNode` | 代理节点 | Proxy Node |
| `settings.sections.clash.form.fields.tunnels.protocols` | 协议 | Protocol |
| `settings.sections.clash.form.fields.tunnels.existing` | 已配置的隧道 | Existing Tunnels |
| `settings.sections.clash.form.fields.tunnels.default` | 跟随当前配置 | Follow Current Configuration |
| `settings.sections.clash.form.fields.tunnels.optional` | 可选 | Optional |
| `settings.sections.clash.form.fields.tunnels.messages.incomplete` | 请完整填写隧道信息 | Please fill in all required tunnel fields |
| `settings.sections.clash.form.fields.tunnels.messages.invalidLocalAddr` | 无效的本地监听地址 | Invalid local listen address |
| `settings.sections.clash.form.fields.tunnels.messages.invalidLocalPort` | 无效的本地监听端口 | Invalid local listen port |
| `settings.sections.clash.form.fields.tunnels.messages.invalidTargetAddr` | 无效的目标地址 | Invalid Target Address |
| `settings.sections.clash.form.fields.tunnels.messages.invalidTargetPort` | 无效的目标端口 | Invalid Target Port |
| `settings.sections.clash.form.fields.tunnels.actions.add` | 添加 | Add |
| `settings.sections.clash.form.fields.tunnels.actions.addNew` | 添加新隧道 | Add New Tunnel |
| `settings.sections.clash.form.tooltips.networkInterface` | 网络接口 | Network Interface |
| `settings.sections.clash.form.tooltips.unifiedDelay` | 开启统一延迟时，会进行两次延迟测试，以消除连接握手等带来的不同类型节点的延迟差异 | When unified delay is turned on, two delay tests will be performed to eliminate the delay differences between different types of nodes caused by connection handshakes, etc |
| `settings.sections.clash.form.tooltips.logLevel` | 仅对日志目录 Service 文件夹下的内核日志文件生效 | This parameter is valid only for kernel log files in the log directory Service folder |
| `settings.sections.clash.form.tooltips.openUwpTool` | Windows 8 开始限制 UWP 应用（如微软商店）直接访问本地主机的网络服务，使用此工具可绕过该限制 | Since Windows 8, UWP apps (such as Microsoft Store) are restricted from directly accessing local host network services, and this tool can be used to bypass this restriction |
| `settings.sections.clash.form.options.logLevel.debug` | Debug | Debug |
| `settings.sections.clash.form.options.logLevel.info` | Info | Info |
| `settings.sections.clash.form.options.logLevel.warning` | Warn | Warn |
| `settings.sections.clash.form.options.logLevel.error` | Error | Error |
| `settings.sections.clash.form.options.logLevel.silent` | Silent | Silent |
| `settings.components.verge.basic.title` | Verge 基础设置 | Verge Basic Setting |
| `settings.components.verge.basic.actions.browse` | 浏览 | Browse |
| `settings.components.verge.basic.trayOptions.showMainWindow` | 显示主窗口 | Show Main Window |
| `settings.components.verge.basic.trayOptions.showTrayMenu` | 显示托盘菜单 | Show Tray Menu |
| `settings.components.verge.basic.trayOptions.disable` | 禁用 | Disable |
| `settings.components.verge.basic.fields.language` | 语言设置 | Language |
| `settings.components.verge.basic.fields.themeMode` | 主题模式 | Theme Mode |
| `settings.components.verge.basic.fields.trayClickEvent` | 托盘点击事件 | Tray Click Event |
| `settings.components.verge.basic.fields.copyEnvType` | 复制环境变量类型 | Copy Env Type |
| `settings.components.verge.basic.fields.startPage` | 启动页面 | Start Page |
| `settings.components.verge.basic.fields.startupScript` | 启动脚本 | Startup Script |
| `settings.components.verge.basic.fields.themeSetting` | 主题设置 | Theme Setting |
| `settings.components.verge.basic.fields.layoutSetting` | 界面设置 | Layout Setting |
| `settings.components.verge.basic.fields.misc` | 杂项设置 | Miscellaneous |
| `settings.components.verge.basic.fields.hotkeySetting` | 热键设置 | Hotkey Setting |
| `settings.components.verge.advanced.title` | Verge 高级设置 | Verge Advanced Setting |
| `settings.components.verge.advanced.tooltips.backupInfo` | 支持本地或 WebDAV 方式备份配置文件 | Support local or WebDAV backup of configuration files |
| `settings.components.verge.advanced.tooltips.openConfDir` | 如果软件运行异常，!备份!并删除此文件夹下的所有文件，重启软件 | If the software runs abnormally, BACKUP and delete all files in this folder then restart the software |
| `settings.components.verge.advanced.tooltips.liteMode` | 关闭 GUI 界面，仅保留内核运行 | Close the GUI and keep only the kernel running |
| `settings.components.verge.advanced.actions.copyVersion` | 复制 Verge 版本号 | Copy Version |
| `settings.components.verge.advanced.notifications.latestVersion` | 当前已是最新版本 | Currently on the Latest Version |
| `settings.components.verge.advanced.notifications.versionCopied` | Verge 版本已复制到剪贴板 | Version copied to clipboard |
| `settings.components.verge.advanced.fields.backupSetting` | 备份设置 | Backup Setting |
| `settings.components.verge.advanced.fields.runtimeConfig` | 当前配置 | Runtime Config |
| `settings.components.verge.advanced.fields.openConfDir` | 配置目录 | Open Conf Dir |
| `settings.components.verge.advanced.fields.openCoreDir` | 内核目录 | Open Core Dir |
| `settings.components.verge.advanced.fields.openLogsDir` | 日志目录 | Open Logs Dir |
| `settings.components.verge.advanced.fields.checkUpdates` | 检查更新 | Check for Updates |
| `settings.components.verge.advanced.fields.openDevTools` | 开发者工具 | Dev Tools |
| `settings.components.verge.advanced.fields.liteModeSettings` | 轻量模式设置 | LightWeight Mode Settings |
| `settings.components.verge.advanced.fields.exit` | 退出 | Exit |
| `settings.components.verge.advanced.fields.exportDiagnostics` | 导出诊断信息 | Export Diagnostic Info |
| `settings.components.verge.advanced.fields.vergeVersion` | Verge 版本 | Verge Version |
| `settings.components.verge.theme.title` | 主题设置 | Theme Setting |
| `settings.components.verge.theme.fields.primaryColor` | 主要颜色 | Primary Color |
| `settings.components.verge.theme.fields.secondaryColor` | 次要颜色 | Secondary Color |
| `settings.components.verge.theme.fields.primaryText` | 文本主要颜色 | Primary Text |
| `settings.components.verge.theme.fields.secondaryText` | 文本次要颜色 | Secondary Text |
| `settings.components.verge.theme.fields.infoColor` | 信息颜色 | Info Color |
| `settings.components.verge.theme.fields.warningColor` | 警告颜色 | Warning Color |
| `settings.components.verge.theme.fields.errorColor` | 错误颜色 | Error Color |
| `settings.components.verge.theme.fields.successColor` | 成功颜色 | Success Color |
| `settings.components.verge.theme.fields.fontFamily` | 字体系列 | Font Family |
| `settings.components.verge.theme.fields.cssInjection` | CSS 注入 | CSS Injection |
| `settings.components.verge.theme.actions.editCss` | 编辑 CSS | Edit CSS |
| `settings.components.verge.theme.dialogs.editCssTitle` | 编辑 CSS | Edit CSS |
| `settings.components.verge.layout.title` | 界面设置 | Layout Setting |
| `settings.components.verge.layout.fields.preferSystemTitlebar` | 优先使用系统标题栏 | Prefer System Titlebar |
| `settings.components.verge.layout.fields.trafficGraph` | 流量图显 | Traffic Graph |
| `settings.components.verge.layout.fields.memoryUsage` | 内核占用 | Core Usage |
| `settings.components.verge.layout.fields.proxyGroupIcon` | 代理组图标 | Proxy Group Icon |
| `settings.components.verge.layout.fields.toastPosition` | 通知位置 | Toast Position |
| `settings.components.verge.layout.fields.hoverNavigator` | 悬浮跳转导航 | Hover Jump Navigator |
| `settings.components.verge.layout.fields.hoverNavigatorDelay` | 悬浮跳转导航延迟 | Hover Jump Navigator Delay |
| `settings.components.verge.layout.fields.navIcon` | 导航栏图标 | Nav Icon |
| `settings.components.verge.layout.fields.collapseNavBar` | 收起导航栏 | Collapse Navigation Bar |
| `settings.components.verge.layout.fields.trayIcon` | 托盘图标 | Tray Icon |
| `settings.components.verge.layout.fields.proxyGroupsDisplayMode` | 托盘代理组显示模式 | Proxy Groups Display Mode |
| `settings.components.verge.layout.fields.showOutboundModesInline` | 将出站模式显示在托盘一级菜单 | Show Outbound Modes Inline |
| `settings.components.verge.layout.fields.commonTrayIcon` | 常规托盘图标 | Common Tray Icon |
| `settings.components.verge.layout.fields.systemProxyTrayIcon` | 系统代理托盘图标 | System Proxy Tray Icon |
| `settings.components.verge.layout.fields.tunTrayIcon` | TUN 模式托盘图标 | Tun Tray Icon |
| `settings.components.verge.layout.fields.enableTrayIcon` | 启用托盘图标 | Enable Tray Icon |
| `settings.components.verge.layout.fields.enableTraySpeed` | 启用托盘速率 | Enable Tray Speed |
| `settings.components.verge.layout.fields.pauseRenderTrafficStatsOnBlur` | 失去焦点时暂停渲染流量统计 | Pause rendering traffic statistics when blurred |
| `settings.components.verge.layout.tooltips.hoverNavigator` | 鼠标悬停在字母上时自动滚动到对应代理组 | Automatically scroll to the corresponding proxy group when hovering over alphabet letters |
| `settings.components.verge.layout.tooltips.hoverNavigatorDelay` | 鼠标悬停后触发自动跳转前等待的毫秒数 | Delay before auto scrolling when hovering, in milliseconds |
| `settings.components.verge.layout.options.icon.monochrome` | 单色图标 | Monochrome |
| `settings.components.verge.layout.options.icon.colorful` | 彩色图标 | Colorful |
| `settings.components.verge.layout.options.icon.disable` | 禁用 | Disable |
| `settings.components.verge.layout.options.toastPosition.topLeft` | 左上角 | Top Left |
| `settings.components.verge.layout.options.toastPosition.topRight` | 右上角 | Top Right |
| `settings.components.verge.layout.options.toastPosition.bottomLeft` | 左下角 | Bottom Left |
| `settings.components.verge.layout.options.toastPosition.bottomRight` | 右下角 | Bottom Right |
| `settings.components.verge.layout.options.proxyGroupsDisplayMode.default` | 默认 | Default |
| `settings.components.verge.layout.options.proxyGroupsDisplayMode.inline` | 一级菜单 | Inline |
| `settings.components.verge.layout.options.proxyGroupsDisplayMode.disable` | 禁用 | Disable |
| `settings.modals.clashPort.title` | 端口设置 | Port Config |
| `settings.modals.clashPort.fields.mixed` | 混合代理端口 | Mixed Port |
| `settings.modals.clashPort.fields.socks` | SOCKS 代理端口 | Socks Port |
| `settings.modals.clashPort.fields.http` | HTTP(S) 代理端口 | Http(s) Port |
| `settings.modals.clashPort.fields.redir` | Redir 透明代理端口 | Redir Port |
| `settings.modals.clashPort.fields.tproxy` | Tproxy 透明代理端口 | Tproxy Port |
| `settings.modals.clashPort.actions.random` | 随机端口 | Random Port |
| `settings.modals.clashPort.messages.portInUse` | 端口 {{port}} 已被占用 | Port {{port}} is already in use |
| `settings.modals.clashPort.messages.runningPort` | 当前运行端口：{{port}} | Currently using port {{port}} |
| `settings.modals.clashPort.messages.duplicatePort` | 已启用的监听端口不能重复 | Enabled listeners must use different ports |
| `settings.modals.clashPort.messages.invalidPort` | 请输入 1 至 65535 之间的端口号 | Enter a port between 1 and 65535 |
| `settings.modals.clashPort.messages.saved` | 端口设置已保存 | Port settings saved |
| `settings.modals.clashPort.messages.saveFailed` | 端口设置保存失败 | Failed to save port settings |
| `settings.modals.clashPort.messages.automaticFallback` | 端口 {{oldPort}} 已被占用，本次运行改用 {{newPort}}。设置未更改，下次启动会重新尝试 {{oldPort}}。 | Port {{oldPort}} was in use, so this session runs on {{newPort}}. Your setting is unchanged and {{oldPort}} is retried on the next start. |
| `settings.modals.clashPort.messages.automaticFallbackFailed` | 混合代理端口自动切换失败：{{error}} | Automatic mixed proxy port fallback failed: {{error}} |
| `settings.modals.clashPort.messages.portTooLow` | 端口号不能小于 1000 | The port must be 1000 or higher |
| `settings.modals.clashPort.messages.portTooHigh` | 端口号不能大于 65535 | The port must be 65535 or lower |
| `settings.modals.clashCore.variants.release` | 正式版 | Release Version |
| `settings.modals.clashCore.variants.alpha` | 预览版 | Alpha Version |
| `settings.modals.liteMode.title` | 轻量模式设置 | LightWeight Mode Settings |
| `settings.modals.liteMode.actions.enterNow` | 立即进入轻量模式 | Enter LightWeight Mode Now |
| `settings.modals.liteMode.toggles.autoEnter` | 自动进入轻量模式 | Auto Enter LightWeight Mode |
| `settings.modals.liteMode.tooltips.autoEnter` | 启用后，将在窗口关闭一段时间后自动激活轻量模式 | Enable to automatically activate LightWeight Mode after the window is closed for a period of time |
| `settings.modals.liteMode.fields.delay` | 自动进入轻量模式延迟 | Auto Enter LightWeight Mode Delay |
| `settings.modals.liteMode.messages.autoEnterHint` | 关闭窗口后，轻量模式将在 {{n}} 分钟后自动激活 | When closing the window, LightWeight Mode will be automatically activated after {{n}} minutes |
| `settings.modals.backup.title` | 备份设置 | Backup Setting |
| `settings.modals.backup.tabs.local` | 本地备份 | Local backup |
| `settings.modals.backup.tabs.webdav` | WebDAV 备份 | WebDAV backup |
| `settings.modals.backup.actions.selectTarget` | 选择备份目标 | Select backup target |
| `settings.modals.backup.actions.backup` | 备份 | Backup |
| `settings.modals.backup.actions.export` | 导出 | Export |
| `settings.modals.backup.actions.exportBackup` | 导出备份 | Export Backup |
| `settings.modals.backup.actions.importBackup` | 导入备份 | Import Backup |
| `settings.modals.backup.actions.deleteBackup` | 删除备份 | Delete Backup |
| `settings.modals.backup.actions.restore` | 恢复 | Restore |
| `settings.modals.backup.actions.restoreBackup` | 恢复备份 | Restore Backup |
| `settings.modals.backup.actions.viewHistory` | 查看记录 | View history |
| `settings.modals.backup.fields.webdavUrl` | WebDAV 服务器地址 http(s):// | WebDAV Server URL |
| `settings.modals.backup.fields.username` | 用户名 | Username |
| `settings.modals.backup.fields.info` | 在应用数据目录中创建本地备份，您可以通过下方列表进行恢复或删除。 | Backups are stored locally in the application data directory. Use the list below to restore or delete backups. |
| `settings.modals.backup.messages.webdavUrlRequired` | WebDAV 服务器地址不能为空 | WebDAV URL cannot be empty |
| `settings.modals.backup.messages.invalidWebdavUrl` | 无效的 WebDAV 服务器地址格式 | Invalid WebDAV URL format |
| `settings.modals.backup.messages.usernameRequired` | 用户名不能为空 | Username cannot be empty |
| `settings.modals.backup.messages.passwordRequired` | 密码不能为空 | Password cannot be empty |
| `settings.modals.backup.messages.webdavConfigSaved` | WebDAV 配置保存成功 | WebDAV configuration saved successfully |
| `settings.modals.backup.messages.webdavConfigSaveFailed` | 保存 WebDAV 配置失败: {{error}} | Failed to save WebDAV configuration: {{error}} |
| `settings.modals.backup.messages.backupCreated` | 备份创建成功 | Backup created successfully |
| `settings.modals.backup.messages.backupFailed` | 备份失败: {{error}} | Backup failed: {{error}} |
| `settings.modals.backup.messages.localBackupCreated` | 本地备份创建成功 | Local backup created successfully |
| `settings.modals.backup.messages.localBackupFailed` | 本地备份失败 | Local backup failed |
| `settings.modals.backup.messages.restoreSuccess` | 恢复成功，应用将在 1 秒后重启 | Restore Success, App will restart in 1s |
| `settings.modals.backup.messages.localBackupExported` | 本地备份导出成功 | Local backup exported successfully |
| `settings.modals.backup.messages.localBackupExportFailed` | 本地备份导出失败 | Failed to export local backup |
| `settings.modals.backup.messages.localBackupImported` | 本地备份导入成功 | Local backup imported successfully |
| `settings.modals.backup.messages.localBackupImportFailed` | 本地备份导入失败：{{error}} | Failed to import local backup: {{error}} |
| `settings.modals.backup.messages.webdavRefreshSuccess` | WebDAV 刷新成功 | WebDAV refresh succeeded |
| `settings.modals.backup.messages.webdavRefreshFailed` | WebDAV 刷新失败: {{error}} | WebDAV refresh failed: {{error}} |
| `settings.modals.backup.messages.confirmDelete` | 确认删除此备份文件吗？ | Confirm to delete this backup file? |
| `settings.modals.backup.messages.confirmRestore` | 确认恢复此份文件吗？ | Confirm to restore this backup file? |
| `settings.modals.backup.auto.title` | 自动备份 | Automatic backup |
| `settings.modals.backup.auto.scheduleLabel` | 启用定时备份 | Enable scheduled backup |
| `settings.modals.backup.auto.scheduleHelper` | 按设定频率在后台创建本地备份文件。 | Create local backups in the background at the configured interval. |
| `settings.modals.backup.auto.intervalLabel` | 备份频率 | Backup frequency |
| `settings.modals.backup.auto.changeLabel` | 关键变更时自动备份 | Backup on critical changes |
| `settings.modals.backup.auto.changeHelper` | 全局扩展覆写配置/脚本更改后自动备份。 | Automatically backup when Global Extend Config/Script changes. |
| `settings.modals.backup.auto.options.hours` | 每 {{n}} 小时 | Every {{n}} hours |
| `settings.modals.backup.auto.options.days` | 每 {{n}} 天 | Every {{n}} days |
| `settings.modals.backup.manual.title` | 手动备份 | Manual backup |
| `settings.modals.backup.manual.local` | 在本设备的应用数据目录中创建备份。 | Creates a snapshot on this device, stored under the app data directory. |
| `settings.modals.backup.manual.webdav` | 配置 WebDAV 后，可直接上传备份到服务器。 | Upload a snapshot to your WebDAV server once credentials are set. |
| `settings.modals.backup.manual.configureWebdav` | 配置 WebDAV | Configure WebDAV |
| `settings.modals.backup.history.title` | 备份记录 | Backup history |
| `settings.modals.backup.history.summary` | 共 {{count}} 份备份 · 最近 {{recent}} | {{count}} backups • latest {{recent}} |
| `settings.modals.backup.history.empty` | 暂无备份记录 | No backups available |
| `settings.modals.backup.history.unknownPlatform` | 未知 | unknown |
| `settings.modals.backup.history.unknownTime` | 未知时间 | Unknown time |
| `settings.modals.backup.webdav.title` | WebDAV 设置 | WebDAV settings |
| `settings.modals.backup.table.filename` | 文件名称 | Filename |
| `settings.modals.backup.table.backupTime` | 备份时间 | Backup Time |
| `settings.modals.backup.table.actions` | 操作 | Actions |
| `settings.modals.backup.table.noBackups` | 暂无备份 | No backups available |
| `settings.modals.backup.table.rowsPerPage` | 每页行数 | Rows per page |
| `settings.modals.misc.title` | 杂项设置 | Miscellaneous |
| `settings.modals.misc.fields.appLogLevel` | 应用日志等级 | App Log Level |
| `settings.modals.misc.fields.appLogMaxSize` | 应用日志最大大小 | App Log Max Size |
| `settings.modals.misc.fields.appLogMaxCount` | 应用日志最大数量 | App Log Max Count |
| `settings.modals.misc.fields.autoCloseConnections` | 自动关闭连接 | Auto Close Connections |
| `settings.modals.misc.fields.autoCheckUpdate` | 自动检查更新 | Auto Check Update |
| `settings.modals.misc.fields.enableBuiltinEnhanced` | 内置增强功能 | Enable Builtin Enhanced |
| `settings.modals.misc.fields.proxyLayoutColumns` | 代理页布局列数 | Proxy Layout Columns |
| `settings.modals.misc.fields.autoLogClean` | 自动清理日志 | Auto Log Clean |
| `settings.modals.misc.fields.autoDelayDetection` | 自动延迟检测 | Auto Delay Detection |
| `settings.modals.misc.fields.autoDelayDetectionInterval` | 自动延迟检测间隔 | Auto Delay Detection Interval |
| `settings.modals.misc.fields.defaultLatencyTest` | 默认测试链接 | Default Latency Test |
| `settings.modals.misc.fields.defaultLatencyTimeout` | 测试超时时间 | Default Latency Timeout |
| `settings.modals.misc.tooltips.autoCloseConnections` | 当代理组选中节点或代理模式变动时，关闭已建立的连接 | Terminate established connections when the proxy group selection or proxy mode changes |
| `settings.modals.misc.tooltips.enableBuiltinEnhanced` | 配置文件的兼容性处理 | Compatibility handling for the configuration file |
| `settings.modals.misc.tooltips.autoDelayDetection` | 后台定时检测当前节点延迟 | Periodically test the current node latency in the background |
| `settings.modals.misc.tooltips.defaultLatencyTest` | 仅用于 HTTP 客户端请求测试，不会对配置文件产生影响 | Used for HTTP client request testing only and won't make a difference to the configuration file |
| `settings.modals.misc.options.proxyLayoutColumns.auto` | 自动列数 | Auto Columns |
| `settings.modals.misc.options.autoLogClean.never` | 不清理 | Never Clean |
| `settings.modals.misc.options.autoLogClean.retainDays` | 保留 {{n}} 天 | Retain {{n}} Days |
| `settings.modals.update.title` | 新版本 v{{version}} | New Version v{{version}} |
| `settings.modals.update.actions.goToRelease` | 前往发布页 | Go to Release Page |
| `settings.modals.update.actions.update` | 更新 | Update |
| `settings.modals.update.messages.breakChangeError` | 此版本为重大更新，不支持应用内更新，请卸载后手动下载安装 | This version is a major update and does not support in-app updates. Please uninstall it and manually download and install the new version |
| `settings.modals.update.messages.available` | 有新版本可用 | A new version is available |
| `settings.modals.update.alerts.note` | 备注 | Note |
| `settings.modals.update.alerts.tip` | 提示 | Tip |
| `settings.modals.update.alerts.important` | 重要 | Important |
| `settings.modals.update.alerts.warning` | 警告 | Warning |
| `settings.modals.update.alerts.caution` | 注意 | Caution |
| `settings.modals.sysproxy.title` | 系统代理设置 | System Proxy Setting |
| `settings.modals.sysproxy.fieldsets.currentStatus` | 当前系统代理 | Current System Proxy |
| `settings.modals.sysproxy.fields.enableStatus` | 开启状态： | Enable Status: |
| `settings.modals.sysproxy.fields.serverAddr` | 服务地址： | Server Addr:  |
| `settings.modals.sysproxy.fields.pacUrl` | PAC 地址： | PAC URL:  |
| `settings.modals.sysproxy.fields.proxyHost` | 代理主机 | Proxy Host |
| `settings.modals.sysproxy.fields.usePacMode` | 使用 PAC 模式 | Use PAC Mode |
| `settings.modals.sysproxy.fields.proxyGuard` | 系统代理守卫 | Proxy Guard |
| `settings.modals.sysproxy.fields.guardDuration` | 代理守卫间隔 | Guard Duration |
| `settings.modals.sysproxy.fields.alwaysUseDefaultBypass` | 始终使用默认绕过 | Always use Default Bypass |
| `settings.modals.sysproxy.fields.enableBypassCheck` | 验证代理绕过格式 | Validate Proxy Bypass Format |
| `settings.modals.sysproxy.fields.proxyBypass` | 代理绕过设置： | Proxy Bypass Settings:  |
| `settings.modals.sysproxy.fields.bypass` | 当前绕过： | Bypass:  |
| `settings.modals.sysproxy.fields.pacScriptContent` | PAC 脚本内容 | PAC Script Content |
| `settings.modals.sysproxy.tooltips.proxyGuard` | 开启以防止其他软件修改操作系统的代理设置 | Enable to prevent other software from modifying the operating system's proxy settings |
| `settings.modals.sysproxy.messages.durationTooShort` | 代理守护间隔时间不得低于 1 秒 | Proxy Daemon Duration Cannot be Less than 1 Second |
| `settings.modals.sysproxy.messages.invalidBypass` | 无效的代理绕过格式 | Invalid Bypass Format |
| `settings.modals.sysproxy.messages.invalidProxyHost` | 代理主机格式无效 | Invalid Proxy Host Format |
| `settings.modals.sysproxy.actions.editPac` | 编辑 PAC | Edit PAC |
| `settings.modals.tun.title` | 虚拟网卡模式 | Tun Mode |
| `settings.modals.tun.fields.stack` | TUN 模式堆栈 | Tun Stack |
| `settings.modals.tun.fields.device` | 虚拟网卡名称 | Device Name |
| `settings.modals.tun.fields.autoRoute` | 自动设置全局路由 | Auto Route |
| `settings.modals.tun.fields.routeExcludeAddress` | 排除自定义网段 | Route Exclude Address |
| `settings.modals.tun.fields.strictRoute` | 严格路由 | Strict Route |
| `settings.modals.tun.fields.autoDetectInterface` | 自动选择流量出口接口 | Auto Detect Interface |
| `settings.modals.tun.fields.dnsHijack` | DNS 劫持 | DNS Hijack |
| `settings.modals.tun.fields.dnsMode` | DNS 模式 | DNS Mode |
| `settings.modals.tun.fields.mtu` | 最大传输单元 | Max Transmission Unit |
| `settings.modals.tun.fields.autoRedirect` | 自动重定向 | Auto Redirect |
| `settings.modals.tun.tooltips.dnsHijack` | 多个 DNS 服务器请使用 , 分隔 | Please use , to separate multiple DNS servers |
| `settings.modals.tun.tooltips.dnsMode` | Auto 保持旧行为；hijack 需 Windows 开启严格路由或 Linux 开启自动重定向 | Auto keeps the legacy behavior; hijack requires Strict Route on Windows or Auto Redirect on Linux |
| `settings.modals.tun.tooltips.autoRedirect` | 自动配置 nftables/iptables 的 TCP 重定向 | Automatically configures nftables/iptables TCP redirects |
| `settings.modals.tun.messages.applied` | 设置已应用 | Settings Applied |
| `settings.modals.tun.messages.invalidRouteExcludeAddress` | 请输入有效的 CIDR 网段 | Please enter a valid CIDR block |
| `settings.modals.tun.messages.routeExcludeAddressHint` | 仅支持 IPv4/IPv6 CIDR，例如 192.168.0.0/16 或 fd00::/8 | Only IPv4/IPv6 CIDR is supported, such as 192.168.0.0/16 or fd00::/8 |
| `settings.modals.dns.dialog.title` | DNS 覆写 | DNS Overwrite |
| `settings.modals.dns.dialog.profileScope` | DNS 覆写开关按订阅独立保存，仅对对应订阅生效。 | The DNS override switch is saved separately for each profile and applies only to that profile. |
| `settings.modals.dns.dialog.warning` | 如果你不清楚这里的设置请不要修改，并保持 DNS 覆写关闭 | If you are not familiar with these settings, please do not modify them and keep DNS Overwrite disabled |
| `settings.modals.dns.protection.title` | 是否开启 DNS 覆写？ | Enable DNS Overwrite? |
| `settings.modals.dns.protection.message` | 当前订阅包含自定义 DNS 配置。开启 DNS 覆写可能覆盖这些配置，导致节点无法连接或部分网站无法访问。是否仍要开启？ | This subscription contains custom DNS settings. DNS Overwrite may replace them, preventing nodes from connecting or some websites from loading. Enable it anyway? |
| `settings.modals.dns.protection.enableAnyway` | 仍然开启 | Enable anyway |
| `settings.modals.dns.protection.keepDisabled` | 保持关闭 | Keep disabled |
| `settings.modals.dns.protection.autoDisabled` | 检测到当前订阅包含自定义 DNS 配置，已自动关闭 DNS 覆写。 | DNS Overwrite was disabled to preserve the custom DNS settings in the current subscription. |
| `settings.modals.dns.sections.general` | DNS 设置 | DNS Settings |
| `settings.modals.dns.sections.fallbackFilter` | 回退过滤设置 | Fallback Filter Settings |
| `settings.modals.dns.sections.hosts` | Hosts 设置 | Hosts Settings |
| `settings.modals.dns.fields.enable` | 启用 DNS | Enable DNS |
| `settings.modals.dns.fields.listen` | DNS 监听地址 | DNS Listen |
| `settings.modals.dns.fields.enhancedMode` | 增强模式 | Enhanced Mode |
| `settings.modals.dns.fields.fakeIpRange` | Fake IP 范围 | Fake IP Range |
| `settings.modals.dns.fields.fakeIpRange6` | Fake IP 范围 (IPv6) | Fake IP Range (IPv6) |
| `settings.modals.dns.fields.fakeIpFilterMode` | Fake IP 过滤模式 | Fake IP Filter Mode |
| `settings.modals.dns.fields.ipv6.label` | IPv6 | IPv6 |
| `settings.modals.dns.fields.ipv6.description` | 启用 IPv6 DNS 解析 | Enable IPv6 DNS resolution |
| `settings.modals.dns.fields.preferH3.label` | 优先使用 HTTP/3 | Prefer H3 |
| `settings.modals.dns.fields.preferH3.description` | DNS DOH 使用 HTTP/3 协议 | DNS DOH uses HTTP/3 |
| `settings.modals.dns.fields.respectRules.label` | 遵循路由规则 | Respect Rules |
| `settings.modals.dns.fields.respectRules.description` | DNS 连接遵循路由规则 | DNS connections follow routing rules |
| `settings.modals.dns.fields.useHosts.label` | 使用 Hosts | Use Hosts |
| `settings.modals.dns.fields.useHosts.description` | 启用通过 hosts 文件解析域名 | Enable to resolve hosts through hosts file |
| `settings.modals.dns.fields.useSystemHosts.label` | 使用系统 Hosts | Use System Hosts |
| `settings.modals.dns.fields.useSystemHosts.description` | 启用通过系统 hosts 文件解析域名 | Enable to resolve hosts through system hosts file |
| `settings.modals.dns.fields.directPolicy.label` | 直连域名服务器遵循策略 | Direct Nameserver Follow Policy |
| `settings.modals.dns.fields.directPolicy.description` | 是否遵循 nameserver-policy 设置 | Whether to follow nameserver policy |
| `settings.modals.dns.fields.defaultNameserver.label` | 默认域名服务器 | Default Nameserver |
| `settings.modals.dns.fields.defaultNameserver.description` | 用于解析 DNS 服务器的默认 DNS 服务器 | Default DNS servers used to resolve DNS servers |
| `settings.modals.dns.fields.nameserver.label` | 域名服务器 | Nameserver |
| `settings.modals.dns.fields.nameserver.description` | DNS 服务器列表，用逗号分隔 | List of DNS servers, comma separated |
| `settings.modals.dns.fields.fallback.label` | 回退服务器 | Fallback |
| `settings.modals.dns.fields.fallback.description` | 回退 DNS 服务器列表，用逗号分隔 | List of fallback DNS servers, comma separated |
| `settings.modals.dns.fields.proxy.label` | 代理节点DNS | Proxy Server Nameserver |
| `settings.modals.dns.fields.proxy.description` | 代理节点域名解析服务器，仅用于解析代理节点的域名，用逗号分隔 | DNS servers for proxy node domain resolution |
| `settings.modals.dns.fields.directNameserver.label` | 直连域名服务器 | Direct Nameserver |
| `settings.modals.dns.fields.directNameserver.description` | 直连出口域名解析服务器，支持 system 关键字，用逗号分隔 | DNS servers for direct exit domain resolution, supports 'system' keyword, comma separated |
| `settings.modals.dns.fields.fakeIpFilter.label` | Fake IP 过滤 | Fake IP Filter |
| `settings.modals.dns.fields.fakeIpFilter.description` | 跳过 Fake IP 解析的域名，用逗号分隔 | Domains that skip fake IP resolution, comma separated |
| `settings.modals.dns.fields.nameserverPolicy.label` | 域名服务器策略 | Nameserver Policy |
| `settings.modals.dns.fields.nameserverPolicy.description` | 特定域名的 DNS 服务器，多个服务器使用分号分隔，格式: domain=server1;server2 | Domain-specific DNS server, multiple servers separated by semicolons, format: domain=server1;server2 |
| `settings.modals.dns.fields.geoipFiltering.label` | GeoIP 过滤 | GeoIP Filtering |
| `settings.modals.dns.fields.geoipFiltering.description` | 启用 GeoIP 回退过滤 | Enable GeoIP filtering for fallback |
| `settings.modals.dns.fields.geoipCode` | GeoIP 国家代码 | GeoIP Code |
| `settings.modals.dns.fields.fallbackIpCidr.label` | 回退 IP CIDR | Fallback IP CIDR |
| `settings.modals.dns.fields.fallbackIpCidr.description` | 不使用回退服务器的 IP CIDR，用逗号分隔 | IP CIDRs not using fallback servers, comma separated |
| `settings.modals.dns.fields.fallbackDomain.label` | 回退域名 | Fallback Domain |
| `settings.modals.dns.fields.fallbackDomain.description` | 使用回退服务器的域名，用逗号分隔 | Domains using fallback servers, comma separated |
| `settings.modals.dns.fields.hosts.label` | Hosts | Hosts |
| `settings.modals.dns.fields.hosts.description` | 自定义域名到 IP 或域名的映射，用逗号分隔 | Custom domain to IP or domain mapping |
| `settings.modals.dns.messages.saved` | DNS 设置已保存 | DNS settings saved |
| `settings.modals.dns.messages.configError` | DNS 配置错误： | DNS configuration error: |
| `settings.modals.dns.errors.invalid` | 配置无效 | Invalid configuration |
| `settings.modals.dns.errors.invalidYaml` | YAML 格式无效 | Invalid YAML format |
| `settings.modals.webUI.actions.openUrl` | 打开链接 | Open URL |
| `settings.modals.webUI.title` | 网页界面 | Web UI |
| `settings.modals.webUI.messages.supportedPlaceholders` | 支持 %host, %port, %secret | Support %host, %port, %secret |
| `settings.modals.webUI.messages.placeholderInstruction` | 使用 %host, %port, %secret 表示 主机, 端口, 访问密钥 | Replace host, port, secret with %host, %port, %secret |
| `settings.modals.webUI.errors.clashInfoUnavailable` | 内核连接信息暂不可用 | Core connection information is unavailable |
| `settings.modals.webUI.errors.invalidServer` | 内核服务地址“{{server}}”无效 | The core server address "{{server}}" is invalid |
| `settings.modals.webUI.errors.openFailed` | 打开 Web UI 失败 | Failed to open Web UI |
| `settings.modals.hotkey.toggles.enableGlobal` | 启用全局热键 | Enable Global Hotkey |
| `settings.modals.hotkey.title` | 热键设置 | Hotkey Setting |
| `settings.modals.hotkey.functions.rule` | 规则模式 | Rule Mode |
| `settings.modals.hotkey.functions.global` | 全局模式 | Global Mode |
| `settings.modals.hotkey.functions.openOrCloseDashboard` | 打开/关闭面板 | Open/Close Dashboard |
| `settings.modals.hotkey.functions.toggleSystemProxy` | 打开/关闭系统代理 | Enable/Disable System Proxy |
| `settings.modals.hotkey.functions.toggleTunMode` | 打开/关闭 TUN 模式 | Enable/Disable Tun Mode |
| `settings.modals.hotkey.functions.entryLightweightMode` | 进入轻量模式 | Entry Lightweight Mode |
| `settings.modals.hotkey.functions.direct` | 直连模式 | Direct Mode |
| `settings.modals.hotkey.functions.reactivateProfiles` | 重新激活订阅 | Reactivate Profiles |
| `settings.modals.networkInterface.title` | 网络接口 | Network Interface |
| `settings.modals.networkInterface.fields.ipAddress` | IP 地址 | IP Address |
| `settings.modals.networkInterface.fields.macAddress` | MAC 地址 | MAC Address |
| `settings.feedback.errors.clash.configUpdateFailed` | 更新内核配置失败 | Failed to update core configuration |
| `settings.feedback.errors.clash.modeUpdateFailed` | 切换代理模式失败 | Failed to change proxy mode |
| `settings.feedback.errors.clash.restartFailed` | 重启内核失败 | Failed to restart the core |
| `settings.feedback.errors.clash.startFailed` | 启动内核失败 | Failed to start the core |
| `settings.feedback.errors.clash.serviceCoreStopped` | 服务运行的内核意外停止 | The service's core stopped unexpectedly |
| `settings.feedback.errors.clash.stopFailed` | 停止内核失败 | Failed to stop the core |
| `settings.feedback.errors.clashService.installFailed` | 安装服务失败 | Failed to install the service |
| `settings.feedback.errors.clashService.uninstallFailed` | 卸载服务失败 | Failed to uninstall the service |
| `settings.feedback.errors.clashService.reinstallFailed` | 重新安装服务失败 | Failed to reinstall the service |
| `settings.feedback.errors.clashService.repairFailed` | 修复服务失败 | Failed to repair the service |
| `settings.feedback.errors.clashService.sidecarFailed` | 继续使用 Sidecar 失败 | Failed to continue with Sidecar |
| `settings.feedback.errors.clashService.coreAlreadyRunning` | 无法启动内核：已有内核进程正在运行。请退出相关应用或清理残留进程后重试。若仍无法解决，请重启电脑。 | Cannot start the core: a core process is already running. Close the related app or stop the leftover process, then retry. If the problem persists, restart your computer. |
| `settings.feedback.errors.sysproxy.privilegeRequired` | 无法设置系统代理：macOS 以权限不足为由拒绝了这次修改。请安装 Clash Verge 服务；管理员也可以在系统设置 → 隐私与安全性的高级选项中，关闭访问系统范围设置需要管理员密码的要求。 | Cannot set the system proxy: macOS refused the change for lack of privileges. Install the Clash Verge service; an administrator can also turn off the administrator-password requirement for system-wide settings, in the advanced options of System Settings → Privacy & Security. |
| `settings.feedback.errors.sysproxy.directFallback` | 系统代理设置失败，已回退为直连。 | Failed to set the system proxy; traffic has been switched back to a direct connection. |
| `settings.feedback.errors.sysproxy.guardStopped` | 系统代理守卫连续多次失败后已停止，系统代理不再被自动维持。 | The system proxy guard stopped after failing several times running. The system proxy is no longer being kept in place. |
| `settings.feedback.errors.sysproxy.sidecarWhileServiceReady` | 系统代理设置失败：Clash Verge 服务已就绪，但内核仍以 Sidecar 运行。重启内核即可切换到服务模式。 | The system proxy could not be set because the core is still running as a Sidecar, although the Clash Verge service is ready. Restart the core to switch to service mode. |
| `settings.feedback.errors.sysproxy.coreNotReady` | 无法开启系统代理：内核尚未运行。请先启动内核（或安装/修复 Clash Verge 服务）后重试。 | The system proxy could not be enabled because the core is not running. Start the core (or install/repair the Clash Verge service) and try again. |
| `settings.feedback.errors.sysproxy.systemCallFailed` | 系统代理写入失败：Windows 系统调用未成功，通常是瞬时故障。请稍后重试。 | The system proxy could not be written because a Windows system call failed. This is usually transient, so try again. |
| `settings.feedback.notifications.clash.restartSuccess` | 已重启 Clash 内核 | Clash Core Restarted |
| `settings.feedback.notifications.clash.versionUpdated` | 内核版本已更新 | Core Version Updated |
| `settings.feedback.notifications.clash.alreadyLatestVersion` | 已经是最新内核版本 | Already Using Latest Core Version |
| `settings.feedback.notifications.clash.changeSuccess` | 内核切换成功 | Core changed successfully |
| `settings.feedback.notifications.clash.changeFailed` | 无法切换内核 | Failed to change core |
| `settings.feedback.notifications.clash.geoDataUpdated` | 已更新 GeoData | GeoData Updated |
| `settings.feedback.notifications.clashService.appDataNotOwned` | 服务拒绝启动内核：磁盘上应用数据目录的所有者不是当前账户（常见于用户目录位于外置磁盘）。修复服务无法解决此问题，请在终端运行以下命令，然后重启软件： | The service refused to start the core because, on disk, the app data folder is not owned by your account (common when the home folder is on an external drive). Repairing the service will not fix this. Run this command in a terminal, then restart the app: |
| `settings.feedback.notifications.clashService.installSuccess` | 已成功安装服务 | Service Installed Successfully |
| `settings.feedback.notifications.clashService.permissionFallback` | {{reason}} 已切换到普通模式（Sidecar），普通权限下无法使用 TUN。可按文档修复后重装服务；不想或不会修复时，可直接以管理员身份运行软件，无需安装服务。 | {{reason}} Switched to normal mode (Sidecar); TUN requires administrator privileges. Follow the guide to fix and reinstall the service. If you prefer not to fix it or do not know how, run the app as administrator without installing the service. |
| `settings.feedback.notifications.clashService.permissionWritableReason` | 内核文件或上级文件夹允许普通账户修改，服务为防止内核被替换而拒绝高权限运行。 | Ordinary accounts can modify the core files or a parent folder, so the service refused to run the core with elevated privileges. |
| `settings.feedback.notifications.clashService.permissionRejectedReason` | 服务无法确认内核文件及所在目录是否安全，已拒绝高权限运行。 | The service could not confirm that the core files and folders are safe to run with elevated privileges. |
| `settings.feedback.notifications.clashService.permissionRepairGuide` | 查看修复方法 | View repair instructions |
| `settings.feedback.notifications.clashService.sidecarFallback` | 服务无法启动内核，已降级为 Sidecar 模式运行，服务仍不可用。 | The service could not start the core. Switched to Sidecar mode; the service remains unavailable. |
| `settings.feedback.notifications.clashService.uninstallSuccess` | 已成功卸载服务 | Service Uninstalled Successfully |
| `settings.feedback.notifications.updater.withClashProxySuccess` | 使用 Clash 代理更新成功 | Update with Clash proxy successfully |
| `settings.feedback.notifications.updater.withClashProxyFailed` | 使用 Clash 代理更新失败 | Update failed even with Clash proxy |
| `settings.statuses.clash.stopping` | 停止内核中... | Stopping Core... |
| `settings.statuses.clash.restarting` | 重启内核中... | Restarting Core... |
| `settings.statuses.clashService.installing` | 安装服务中... | Installing Service... |
| `settings.statuses.clashService.uninstalling` | 卸载服务中... | Uninstalling Service... |

