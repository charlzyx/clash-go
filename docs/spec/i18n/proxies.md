# 代理文案对照（zh / en）

参照 Clash Verge Rev v2.5.9 `160bb4fe8eaef2bc626163d89622229e5d7c5b59`。

- 中文：`src/locales/zh/proxies.json`
- 英文：`src/locales/en/proxies.json`
- 加载：`src/services/i18n.ts` 以文件名作为顶层键，合并进 i18next 默认命名空间 `translation`。调用形式 `t('proxies.…')`。

代理页。`src/pages/proxies.tsx` 与 `src/components/proxy/*`。

| key | 中文 | English |
| --- | --- | --- |
| `proxies.page.modes.rule` | 规则 | Rule |
| `proxies.page.modes.global` | 全局 | Global |
| `proxies.page.modes.direct` | 直连 | Direct |
| `proxies.page.actions.toggleChain` | 链式代理 | Chain Proxy |
| `proxies.page.actions.connect` | 连接 | Connect |
| `proxies.page.actions.disconnect` | 断开 | Disconnect |
| `proxies.page.actions.connecting` | 连接中... | Connecting... |
| `proxies.page.actions.clearChainConfig` | 删除链式配置 | Delete Chain Config |
| `proxies.page.provider.title` | 代理集合 | Proxy Provider |
| `proxies.page.provider.actions.updateAll` | 更新全部 | Update All |
| `proxies.page.provider.actions.update` | 更新 | Update |
| `proxies.page.rules.title` | 代理规则 | Proxy Rules |
| `proxies.page.rules.select` | 选择规则 | Select Rules |
| `proxies.page.labels.proxyCount` | 节点数量 | Proxy Count |
| `proxies.page.labels.delayCheckReset` | 进行延迟测试，以取消固定 | Delay check to cancel fixed |
| `proxies.page.labels.nodeCount` | {{count}} 个节点 | {{count}} nodes |
| `proxies.page.tooltips.locate` | 当前节点 | locate |
| `proxies.page.tooltips.delayCheck` | 延迟测试 | Delay check |
| `proxies.page.tooltips.sortDefault` | 默认排序 | Sort by default |
| `proxies.page.tooltips.sortDelay` | 按延迟排序 | Sort by delay |
| `proxies.page.tooltips.sortName` | 按名称排序 | Sort by name |
| `proxies.page.tooltips.delayCheckUrl` | 延迟测试链接 | Delay check URL |
| `proxies.page.tooltips.showBasic` | 隐藏节点细节 | Proxy basic |
| `proxies.page.tooltips.showDetail` | 展示节点细节 | Proxy detail |
| `proxies.page.tooltips.filter` | 过滤节点 | Filter |
| `proxies.page.placeholders.delayCheckUrl` | 延迟测试链接 | Delay check URL |
| `proxies.page.empty.noAvailableGroups` | 暂无可用代理组 | No proxy groups are available |
| `proxies.page.empty.noProxies` | 暂无代理 | No proxies |
| `proxies.page.empty.noSubscriptions.title` | 尚未导入订阅 | No subscription imported |
| `proxies.page.empty.noSubscriptions.description` | 代理页没有可展示的内容。请先前往“订阅”页面导入并启用订阅。 | There is nothing to display on the Proxies page. Go to Profiles to import and activate a subscription. |
| `proxies.page.empty.inactiveSubscription.title` | 尚未激活订阅 | No active subscription |
| `proxies.page.empty.inactiveSubscription.description` | 检测到已有订阅，但当前没有有效订阅被激活。请前往“订阅”页面选择并激活一个订阅。 | Subscriptions were found, but none is currently active. Go to Profiles and select one to activate. |
| `proxies.page.empty.coreUnavailable.title` | 内核暂不可用 | Core temporarily unavailable |
| `proxies.page.empty.coreUnavailable.description` | 无法从内核获取代理信息。请尝试重启内核；若问题仍然存在，请打开日志目录检查内核日志。 | Proxy information could not be retrieved from the core. Try restarting the core. If the problem persists, open the logs folder and check the core logs. |
| `proxies.page.empty.noProxyInfo.title` | 未获取到代理信息 | No proxy information available |
| `proxies.page.empty.noProxyInfo.description` | 当前没有可展示的代理组或节点。请检查当前订阅是否包含有效的代理组和节点，或尝试重启内核并查看日志。 | There are no proxy groups or nodes to display. Check that the active subscription contains valid proxy groups and nodes, or restart the core and review the logs. |
| `proxies.page.empty.actions.openProfiles` | 前往订阅 | Go to Profiles |
| `proxies.page.empty.actions.restartCore` | 重启内核 | Restart Core |
| `proxies.page.empty.actions.openLogs` | 打开日志目录 | Open Logs Folder |
| `proxies.page.chain.header` | 代理链配置 | Chain Proxy Config |
| `proxies.page.chain.warning` | 链式代理会显著降低网速。如果您没有明确的业务需求，或不了解此项配置，请保持关闭 | Chain proxy will significantly reduce network speed. If you do not have a clear business need, or do not understand this configuration, keep it disabled. |
| `proxies.page.chain.confirmEnable` | 我明白，开启 | I understand, enable |
| `proxies.page.chain.keepDisabled` | 保持关闭 | Keep disabled |
| `proxies.page.chain.empty` | 链式代理会显著降低网速。如果您没有明确的业务需求，或不了解此项配置，请保持关闭 | Chain proxy will significantly reduce network speed. If you do not have a clear business need, or do not understand this configuration, keep it disabled. |
| `proxies.page.chain.instruction` | 顺序点击节点添加到代理链中 | Click nodes in order to add to proxy chain |
| `proxies.page.chain.minimumNodes` | 链式代理至少需要 2 个节点 | Chain proxy requires at least 2 nodes |
| `proxies.page.chain.minimumNodesHint` | 链式代理至少需要 2 个节点，请再添加一个节点。 | Chain proxy requires at least 2 nodes. Please add one more node. |
| `proxies.page.chain.connectFailed` | 连接链式代理失败 | Failed to connect to proxy chain |
| `proxies.page.chain.disconnectFailed` | 断开链式代理失败 | Failed to disconnect from proxy chain |
| `proxies.page.chain.duplicateNode` | 该节点已在链式代理表中 | Proxy node already exists in chain |
| `proxies.page.chain.entryNode` | 入口 | Entry |
| `proxies.page.chain.exitNode` | 出口 | Exit |
| `proxies.page.messages.directMode` | 直连模式 | Direct Mode |
| `proxies.page.title.default` | 代理组 | Proxy Groups |
| `proxies.page.title.chainMode` | 链式代理模式 | Proxy Chain Mode |
| `proxies.feedback.warnings.invalidTestUrl` | Test URL must start with http:// or https:// | Test URL must start with http:// or https:// |
| `proxies.feedback.notifications.provider.updateSuccess` | {{name}} 更新成功 | {{name}} updated successfully |
| `proxies.feedback.notifications.provider.updateFailed` | {{name}} 更新失败: {{message}} | Failed to update {{name}}: {{message}} |
| `proxies.feedback.notifications.provider.genericError` | 更新失败: {{message}} | Update failed: {{message}} |
| `proxies.feedback.notifications.provider.none` | 没有可更新的 provider | No providers available to update |
| `proxies.feedback.notifications.provider.allUpdated` | 所有 provider 均已更新 | All providers updated successfully |
| `proxies.components.enums.strategies.select` | 手动选择代理 | Select proxy manually |
| `proxies.components.enums.strategies.url-test` | 根据 URL 测试延迟选择代理 | Select proxy based on URL test delay |
| `proxies.components.enums.strategies.fallback` | 不可用时切换到另一个代理 | Switch to another proxy on error |
| `proxies.components.enums.strategies.load-balance` | 根据负载均衡分配代理 | Distribute proxy based on load balancing |
| `proxies.components.enums.strategies.relay` | 根据定义的代理链传递 | Pass through the defined proxy chain |
| `proxies.components.enums.policies.DIRECT` | 直接连接 (DIRECT) | Data goes directly outbound (DIRECT) |
| `proxies.components.enums.policies.REJECT` | 拦截请求 (REJECT) | Intercepts requests (REJECT) |
| `proxies.components.enums.policies.REJECT-DROP` | 丢弃请求 (REJECT-DROP) | Discards requests (REJECT-DROP) |
| `proxies.components.enums.policies.PASS` | 跳过此项 (PASS) | Skips this rule when matched (PASS) |

