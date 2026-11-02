# 首页文案对照（zh / en）

参照 Clash Verge Rev v2.5.9 `160bb4fe8eaef2bc626163d89622229e5d7c5b59`。

- 中文：`src/locales/zh/home.json`
- 英文：`src/locales/en/home.json`
- 加载：`src/services/i18n.ts` 以文件名作为顶层键，合并进 i18next 默认命名空间 `translation`。调用形式 `t('home.…')`。

首页。`HomePage` 与 `src/components/home/*`。

| key | 中文 | English |
| --- | --- | --- |
| `home.page.tooltips.lightweightMode` | 轻量模式 | Lightweight Mode |
| `home.page.tooltips.manual` | 使用手册 | Manual |
| `home.page.tooltips.settings` | 首页设置 | Home Settings |
| `home.page.cards.trafficStats` | 流量统计 | Traffic Stats |
| `home.page.cards.networkSettings` | 网络设置 | Network Settings |
| `home.page.cards.proxyMode` | 代理模式 | Proxy Mode |
| `home.page.settings.cards.profile` | 订阅卡 | Profile Card |
| `home.page.settings.cards.currentProxy` | 当前代理卡 | Current Proxy Card |
| `home.page.settings.cards.network` | 网络设置卡 | Network Settings Card |
| `home.page.settings.cards.proxyMode` | 代理模式卡 | Proxy Mode Card |
| `home.page.settings.cards.traffic` | 流量统计卡 | Traffic Stats Card |
| `home.page.settings.cards.tests` | 网站测试卡 | Website Tests Card |
| `home.page.settings.cards.ip` | IP 信息卡 | IP Information Card |
| `home.page.settings.cards.clashInfo` | Clash 信息卡 | Clash Info Cards |
| `home.page.settings.cards.systemInfo` | 系统信息卡 | System Info Cards |
| `home.page.settings.title` | 首页设置 | Home Settings |
| `home.page.title` | 首页 | Home |
| `home.components.proxyTun.status.systemProxyEnabled` | 系统代理已启用，您的应用将通过代理访问网络 | System Proxy Enabled |
| `home.components.proxyTun.status.systemProxyDisabled` | 系统代理已关闭，建议大多数用户打开此选项 | System Proxy Disabled |
| `home.components.proxyTun.status.tunModeServiceRequired` | TUN 模式需要服务模式，请先安装服务 | TUN Mode Service Required |
| `home.components.proxyTun.status.tunModeEnabled` | TUN 模式已启用，应用将通过虚拟网卡访问网络 | TUN Mode Enabled |
| `home.components.proxyTun.status.tunModeDisabled` | TUN 模式已关闭，适用于特殊应用 | TUN Mode Disabled |
| `home.components.proxyTun.tooltips.systemProxy` | 修改操作系统的代理设置，如果开启失败，可手动修改操作系统的代理设置 | System Proxy Info |
| `home.components.proxyTun.tooltips.tunMode` | TUN 模式可以接管所有应用流量，适用于特殊不遵循系统代理设置的应用 | TUN Mode Intercept Info |
| `home.components.clashInfo.title` | Clash 信息 | Clash Info |
| `home.components.clashInfo.fields.coreVersion` | 内核版本 | Core Version |
| `home.components.clashInfo.fields.systemProxyAddress` | 系统代理地址 | System Proxy Address |
| `home.components.clashInfo.fields.mixedPort` | 混合代理端口 | Mixed Port |
| `home.components.clashInfo.fields.uptime` | 运行时间 | Uptime |
| `home.components.clashInfo.fields.rulesCount` | 规则数量 | Rules Count |
| `home.components.systemInfo.title` | 系统信息 | System Info |
| `home.components.systemInfo.fields.osInfo` | 操作系统信息 | OS Info |
| `home.components.systemInfo.fields.autoLaunch` | 开机自启 | Auto Launch |
| `home.components.systemInfo.fields.runningMode` | 运行模式 | Running Mode |
| `home.components.systemInfo.fields.lastCheckUpdate` | 最后检查更新 | Last Check Update |
| `home.components.systemInfo.fields.vergeVersion` | Verge 版本 | Verge Version |
| `home.components.systemInfo.actions.settings` | 设置 | Settings |
| `home.components.systemInfo.badges.notRunning` | 内核未运行 | Core Not Running |
| `home.components.systemInfo.badges.adminMode` | 管理员模式 | Administrator Mode |
| `home.components.systemInfo.badges.serviceMode` | 服务模式 | Service Mode |
| `home.components.systemInfo.badges.sidecarMode` | 用户模式 | User Mode |
| `home.components.systemInfo.badges.adminServiceMode` | 管理员 + 服务模式 | Admin + Service Mode |
| `home.components.ipInfo.title` | IP 信息 | IP Information |
| `home.components.ipInfo.labels.ip` | IP | IP |
| `home.components.ipInfo.labels.asn` | 自治域 | ASN |
| `home.components.ipInfo.labels.isp` | 服务商 | ISP |
| `home.components.ipInfo.labels.org` | 组织 | ORG |
| `home.components.ipInfo.labels.location` | 位置 | Location |
| `home.components.ipInfo.labels.timezone` | 时区 | Timezone |
| `home.components.ipInfo.labels.autoRefresh` | 自动刷新 | Auto refresh |
| `home.components.ipInfo.labels.unknown` | 未知 | Unknown |
| `home.components.ipInfo.errors.load` | 获取 IP 信息失败 | Failed to get IP info |
| `home.components.ipInfo.errors.loadWithDetails` | 获取 IP 信息失败：{{message}} | Failed to get IP info: {{message}} |
| `home.components.ipInfo.errors.noServices` | 暂无可用的 IP 检测服务 | No IP detection services are available |
| `home.components.currentProxy.title` | 当前节点 | Current Node |
| `home.components.currentProxy.actions.refreshDelay` | 延迟测试 | Delay check |
| `home.components.currentProxy.labels.globalMode` | 全局模式 | Global Mode |
| `home.components.currentProxy.labels.directMode` | 直连模式 | Direct Mode |
| `home.components.currentProxy.labels.group` | 代理组 | Group |
| `home.components.currentProxy.labels.proxy` | 节点 | Proxy |
| `home.components.currentProxy.labels.noActiveNode` | 暂无激活的代理节点 | No active proxy node |
| `home.components.currentProxy.status.testing` | 测试中 | Testing |
| `home.components.currentProxy.status.untested` | 未测试 | Not tested |
| `home.components.currentProxy.status.error` | 错误 | Error |
| `home.components.currentProxy.status.timeout` | 超时 | Timeout |
| `home.components.currentProxy.status.latencyHigh` | 延迟较高 | High latency |
| `home.components.currentProxy.status.latencyMedium` | 延迟中等 | Moderate latency |
| `home.components.currentProxy.status.latencyGood` | 延迟良好 | Good latency |
| `home.components.currentProxy.status.latencyExcellent` | 延迟极佳 | Excellent latency |
| `home.components.currentProxy.status.uninitialized` | 未初始化 | Not initialized |
| `home.components.currentProxy.status.noProxyNode` | 无代理节点 | No proxy node |
| `home.components.tests.title` | 网站测试 | Website Tests |
| `home.components.traffic.metrics.uploadSpeed` | 上传速度 | Upload Speed |
| `home.components.traffic.metrics.downloadSpeed` | 下载速度 | Download Speed |
| `home.components.traffic.metrics.activeConnections` | 活跃连接 | Active Connections |
| `home.components.traffic.metrics.memoryUsage` | 内核占用 | Core Usage |
| `home.components.traffic.legends.upload` | 上传 | Upload |
| `home.components.traffic.legends.download` | 下载 | Download |
| `home.components.traffic.patterns.minutes` | {{time}} 分钟 | {{time}} Minutes |
| `home.components.traffic.unknownTime` | 未知时间 | Unknown time |
| `home.components.traffic.chartStyles.smooth` | 平滑 | Smooth |
| `home.components.traffic.chartStyles.linear` | 线性 | Linear |
| `home.components.traffic.diagnostics` | 数据点：{{points}} \| 压缩后：{{compressed}} \| 帧率：{{fps}} | Points: {{points}} \| Compressed: {{compressed}} \| FPS: {{fps}} |
| `home.components.clashMode.errors.communication` | 内核通信错误，请重启电脑再试 | Core communication error, please restart your computer and try again |
| `home.components.clashMode.labels.rule` | 规则 | Rule |
| `home.components.clashMode.labels.global` | 全局 | Global |
| `home.components.clashMode.labels.direct` | 直连 | Direct |
| `home.components.clashMode.descriptions.rule` | 基于预设规则智能判断流量走向，提供灵活的代理策略 | Automatically choose proxies according to the rule set. |
| `home.components.clashMode.descriptions.global` | 所有流量均通过代理服务器，适用于需要全局科学上网的场景 | Forward all network requests through the selected proxy. |
| `home.components.clashMode.descriptions.direct` | 所有流量不经过代理节点，但经过 Clash 内核转发连接目标服务器，适用于需要通过内核进行分流的特定场景 | Bypass the proxy and connect to the internet directly. |

