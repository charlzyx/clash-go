# 布局文案对照（zh / en）

参照 Clash Verge Rev v2.5.9 `160bb4fe8eaef2bc626163d89622229e5d7c5b59`。

- 中文：`src/locales/zh/layout.json`
- 英文：`src/locales/en/layout.json`
- 加载：`src/services/i18n.ts` 以文件名作为顶层键，合并进 i18next 默认命名空间 `translation`。调用形式 `t('layout.…')`。

导航与壳层文案。首页/设置会直接引用导航 tab 与服务弹窗 key。

| key | 中文 | English |
| --- | --- | --- |
| `layout.components.navigation.tabs.home` | 首 页 | Home |
| `layout.components.navigation.tabs.proxies` | 代 理 | Proxies |
| `layout.components.navigation.tabs.profiles` | 订 阅 | Profiles |
| `layout.components.navigation.tabs.connections` | 连 接 | Connections |
| `layout.components.navigation.tabs.rules` | 规 则 | Rules |
| `layout.components.navigation.tabs.logs` | 日 志 | Logs |
| `layout.components.navigation.tabs.unlock` | 测 试 | Test |
| `layout.components.navigation.tabs.settings` | 设 置 | Settings |
| `layout.components.navigation.menu.reorderMode` | 菜单排序模式 | Menu reorder mode |
| `layout.components.navigation.menu.restoreDefaultOrder` | 恢复默认排序 | Restore default order |
| `layout.components.navigation.menu.unlock` | 解锁菜单排序 | Unlock menu order |
| `layout.components.navigation.menu.lock` | 锁定菜单排序 | Lock menu order |
| `layout.components.navigation.menu.collapseNavBar` | 收起导航栏 | Collapse navigation bar |
| `layout.components.navigation.menu.expandNavBar` | 展开导航栏 | Expand navigation bar |
| `layout.components.serviceMigration.locationRefusedMessage` | 系统服务的安装路径或目录权限未通过安全检查，无法启动内核。请重新安装服务以恢复服务文件和权限，随后重试。Windows 可能会请求管理员授权。 | The service cannot start the core because its installation path or directory permissions failed a security check. Reinstall the service to restore its files and permissions, then try again. Windows may ask for administrator approval. |
| `layout.components.serviceMigration.notAutoStartedMessage` | Clash Verge 系统服务的启动类型已被改为手动或禁用（常见于优化软件的开机加速功能），开机后不会自动运行，本次已改用 Sidecar。请重新安装服务以恢复自动启动。Windows 可能会请求管理员授权。 | The Clash Verge service startup type was changed to Manual or Disabled (often by system optimizer tools), so it no longer starts with Windows; this session is using Sidecar. Reinstall the service to restore automatic startup. Windows may ask for administrator approval. |
| `layout.components.serviceMigration.title` | 需要更新系统服务 | Service update required |
| `layout.components.serviceMigration.message` | Clash Verge 系统服务版本与当前应用要求不符。请选择「重新安装服务」卸载现有服务并安装所需版本，也可以仅在本次会话继续使用 Sidecar。 | The Clash Verge service version does not match the version required by this app. Choose Reinstall Service to uninstall the existing service and install the required version, or continue with Sidecar for this session. |
| `layout.components.serviceMigration.checkingMessage` | 正在确认系统服务状态及当前运行模式，请稍候…… | Checking the system service status and current running mode. Please wait… |
| `layout.components.serviceMigration.unavailableMessage` | 所需的 Clash Verge 系统服务当前不可用。你可以安装或修复服务，或仅在本次会话继续使用 Sidecar。 | The required Clash Verge service is not available. Install or repair it, or continue with Sidecar for this session. |
| `layout.components.serviceMigration.reinstall` | 重新安装服务 | Reinstall service |
| `layout.components.serviceMigration.repair` | 修复服务 | Repair service |
| `layout.components.serviceMigration.continueSidecar` | 继续使用 Sidecar | Continue with Sidecar |
| `layout.components.serviceMigration.success` | Clash Verge 系统服务已更新 | The Clash Verge service was updated successfully |
| `layout.components.serviceMigration.errors.actionFailed` | 服务操作失败 | The service operation failed |
| `layout.components.serviceMigration.errors.stateRefreshFailed` | 刷新服务状态失败 | Failed to refresh the service state |
| `layout.components.serviceMigration.errors.restartFailed` | 重启内核失败 | Failed to restart the core |
| `layout.components.serviceMigration.errors.revalidationFailed` | 刷新运行状态失败 | Failed to refresh the running state |
| `layout.components.serviceMigration.errors.sidecarFailed` | 继续使用 Sidecar 失败 | Failed to continue with Sidecar |
| `layout.components.sysproxyPrivilege.title` | 无法设置系统代理 | Could not set the system proxy |
| `layout.components.sysproxyPrivilege.message` | macOS 以权限不足为由拒绝了这次修改。安装 Clash Verge 系统服务后，代理设置将由服务以管理员身份完成，之后不会再出现这个提示。 | macOS refused the change for lack of privileges. Installing the Clash Verge service lets it make proxy changes as an administrator, and this prompt will not come back. |
| `layout.components.sysproxyPrivilege.serviceReadyMessage` | Clash Verge 系统服务已就绪，但内核仍以 Sidecar 运行，因此设置系统代理时被拒绝。切换到服务模式即可解决，无需重新安装。 | The Clash Verge service is ready, but the core is still running as a sidecar, which is why setting the system proxy was refused. Switching to service mode fixes it — no reinstall needed. |
| `layout.components.sysproxyPrivilege.alternative` | 另一种办法：在系统设置 → 隐私与安全性 → 高级中，关闭「访问系统范围设置需要管理员密码」。 | Alternatively, turn off “require an administrator password to access system-wide settings” in System Settings → Privacy & Security → Advanced. |
| `layout.components.sysproxyPrivilege.later` | 稍后再说 | Later |
| `layout.components.sysproxyPrivilege.installing` | 正在安装系统服务，可能会请求管理员密码，请稍候…… | Installing the system service. This may ask for an administrator password. |
| `layout.components.sysproxyPrivilege.restarting` | 服务已安装，正在通过服务重启内核…… | Service installed. Restarting the core through it… |
| `layout.components.sysproxyPrivilege.applying` | 正在重新应用系统代理…… | Applying your system proxy change… |
| `layout.components.sysproxyPrivilege.tunTitle` | TUN 模式需要系统服务 | TUN mode needs the system service |
| `layout.components.sysproxyPrivilege.tunMessage` | TUN 模式需要 Clash Verge 系统服务才能接管网络流量。安装完成后会自动为你开启 TUN 模式。 | TUN mode needs the Clash Verge system service to take over network traffic. Once it is installed, TUN mode will be turned on for you. |
| `layout.feedback.errors.restartCoreStopFailed` | 无法安全停止内核，已取消重启 | The core could not be stopped safely, so the restart was cancelled |
| `layout.feedback.errors.quitCoreStopFailed` | 无法安全停止内核，已取消退出 | The core could not be stopped safely, so quitting was cancelled |

