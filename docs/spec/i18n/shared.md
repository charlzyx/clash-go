# 共用文案对照（zh / en）

参照 Clash Verge Rev v2.5.9 `160bb4fe8eaef2bc626163d89622229e5d7c5b59`。

- 中文：`src/locales/zh/shared.json`
- 英文：`src/locales/en/shared.json`
- 加载：`src/services/i18n.ts` 以文件名作为顶层键，合并进 i18next 默认命名空间 `translation`。调用形式 `t('shared.…')`。

四页共用：按钮、状态、流量标签、单位。来源 `shared.json`，调用 `t('shared....')`。

| key | 中文 | English |
| --- | --- | --- |
| `shared.actions.cancel` | 取消 | Cancel |
| `shared.actions.close` | 关闭 | Close |
| `shared.actions.confirm` | 确认 | Confirm |
| `shared.actions.save` | 保存 | Save |
| `shared.actions.delete` | 删除 | Delete |
| `shared.actions.edit` | 编辑 | Edit |
| `shared.actions.new` | 新建 | New |
| `shared.actions.enable` | 启用 | Enable |
| `shared.actions.upgrade` | 升级内核 | Upgrade |
| `shared.actions.restart` | 重启内核 | Restart |
| `shared.actions.resetToDefault` | 重置为默认值 | Reset to Default |
| `shared.actions.refresh` | 刷新 | Refresh |
| `shared.actions.retry` | 重试 | Retry |
| `shared.actions.refreshPage` | 刷新页面 | Refresh Page |
| `shared.actions.showDetails` | 显示详情 | Show Details |
| `shared.actions.hideDetails` | 隐藏详情 | Hide Details |
| `shared.actions.listView` | 列表视图 | List View |
| `shared.actions.tableView` | 表格视图 | Table View |
| `shared.actions.pause` | 暂停 | Pause |
| `shared.actions.resume` | 继续 | Resume |
| `shared.actions.closeAll` | 关闭全部 | Close All |
| `shared.actions.clear` | 清除 | Clear |
| `shared.actions.previous` | 上一页 | Previous |
| `shared.actions.next` | 下一页 | Next |
| `shared.actions.check` | 检测 | Check |
| `shared.labels.updateAt` | 更新于 | Update At |
| `shared.labels.timeout` | 超时 | Timeout |
| `shared.labels.icon` | 图标 | Icon |
| `shared.labels.name` | 名称 | Name |
| `shared.labels.readOnly` | 只读 | ReadOnly |
| `shared.labels.expireTime` | 到期时间 | Expire Time |
| `shared.labels.updateTime` | 更新时间 | Update Time |
| `shared.labels.usedTotal` | 已使用 / 总量 | Used / Total |
| `shared.labels.from` | 来自 | From |
| `shared.labels.password` | 密码 | Password |
| `shared.labels.retryAttempts` | 重试次数 | Retry attempts |
| `shared.labels.downloaded` | 下载量 | Downloaded |
| `shared.labels.uploaded` | 上传量 | Uploaded |
| `shared.statuses.enabled` | 已启用 | Enabled |
| `shared.statuses.disabled` | 未启用 | Disabled |
| `shared.statuses.saving` | 保存中... | Saving... |
| `shared.statuses.empty` | 空空如也 | Empty |
| `shared.statuses.loading` | 加载中... | Loading... |
| `shared.units.milliseconds` | 毫秒 | ms |
| `shared.units.seconds` | 秒 | seconds |
| `shared.units.minutes` | 分钟 | mins |
| `shared.units.hours` | 小时 | hrs |
| `shared.units.kilobytes` | KB | KB |
| `shared.units.files` | 文件 | Files |
| `shared.placeholders.resetInput` | 清空输入框 | Clear Input |
| `shared.placeholders.filter` | 过滤条件 | Filter conditions |
| `shared.placeholders.matchCase` | 区分大小写 | Match Case |
| `shared.placeholders.matchWholeWord` | 全字匹配 | Match Whole Word |
| `shared.placeholders.useRegex` | 使用正则表达式 | Use Regular Expression |
| `shared.validation.invalidRegex` | 无效的正则表达式 | Invalid regular expression |
| `shared.window.maximize` | 最大化 | Maximize |
| `shared.window.minimize` | 最小化 | Minimize |
| `shared.editorModes.visualization` | 可视化 | Visualization |
| `shared.editorModes.advanced` | 高级 | Advanced |
| `shared.feedback.errors.unexpected` | 出现了一些问题 | Something went wrong |
| `shared.feedback.errors.label` | 错误 | Error |
| `shared.feedback.errors.unknown` | 未知错误 | Unknown error |
| `shared.feedback.errors.details` | 错误详情 | Error details |
| `shared.feedback.errors.stack` | 错误堆栈 | Error stack |
| `shared.feedback.errors.componentStack` | 组件堆栈 | Component stack |
| `shared.feedback.errors.trafficUnavailable` | 流量数据暂不可用 | Traffic data unavailable |
| `shared.feedback.errors.operationFailed` | 操作失败 | The operation failed |
| `shared.feedback.errors.trafficStats` | 流量统计错误 | Traffic Statistics Error |
| `shared.feedback.errors.trafficStatsDescription` | 流量统计组件发生错误，为防止崩溃已暂时停用。 | The traffic statistics component encountered an error and has been disabled to prevent crashes. |
| `shared.feedback.notices.raw` | {{message}} | {{message}} |
| `shared.feedback.notices.prefixedRaw` | {{prefix}} {{message}} | {{prefix}} {{message}} |
| `shared.feedback.notifications.importSuccess` | 导入订阅成功 | Profile Imported Successfully |
| `shared.feedback.notifications.importSubscriptionSuccess` | 导入订阅成功 | Import subscription successful |
| `shared.feedback.notifications.importWithClashProxy` | 使用 Clash 代理导入订阅成功 | Profile Imported with Clash proxy |
| `shared.feedback.notifications.updateAvailable` | 有可用更新 | Update Available |
| `shared.feedback.notifications.saved` | 保存成功 | Saved successfully |
| `shared.feedback.notifications.common.copySuccess` | 复制成功 | Copy Success |
| `shared.feedback.notifications.common.saveSuccess` | 配置保存成功 | Configuration saved successfully |
| `shared.feedback.notifications.common.saveFailed` | 配置保存失败 | Failed to save configuration |
| `shared.feedback.notifications.common.refreshFailed` | 刷新失败 | Refresh failed; showing last loaded value |
| `shared.feedback.validation.config.failed` | 订阅配置校验失败，请检查订阅配置文件，变更已撤销，错误详情： | Subscription configuration validation failed. Please check the subscription configuration file; modifications have been rolled back. |
| `shared.feedback.validation.config.bootFailed` | 启动订阅配置校验失败，已使用默认配置启动；请检查订阅配置文件，错误详情： | Boot subscription configuration validation failed. Started with the default configuration; please check the subscription configuration file. |
| `shared.feedback.validation.config.processTerminated` | 验证进程被终止，请检查是否有安全软件阻止了内核运行，或电脑 CPU/运行环境是否支持该内核程序。 | The validation process was terminated. Please check whether security software is preventing the core from running, or whether your computer's CPU and runtime environment support the core executable. |
| `shared.feedback.validation.script.syntaxError` | 脚本语法错误，变更已撤销 | Script syntax error, changes reverted |
| `shared.feedback.validation.script.missingMain` | 脚本错误，变更已撤销 | Script error, changes reverted |
| `shared.feedback.validation.script.fileNotFound` | 文件丢失，变更已撤销 | File missing, changes reverted |
| `shared.feedback.validation.script.fileError` | 脚本文件错误，变更已撤销 | Script file error, changes reverted |
| `shared.feedback.validation.yaml.syntaxError` | YAML 语法错误，变更已撤销 | YAML syntax error, changes reverted |
| `shared.feedback.validation.yaml.readError` | YAML 读取错误，变更已撤销 | YAML read error, changes reverted |
| `shared.feedback.validation.yaml.mappingError` | YAML 映射错误，变更已撤销 | YAML mapping error, changes reverted |
| `shared.feedback.validation.merge.syntaxError` | 覆写文件语法错误，变更已撤销 | Merge file syntax error, changes reverted |
| `shared.feedback.validation.merge.mappingError` | 覆写文件映射错误，变更已撤销 | Merge file mapping error, changes reverted |
| `shared.filters.logLevels.all` | ALL | ALL |
| `shared.filters.logLevels.debug` | DEBUG | DEBUG |
| `shared.filters.logLevels.info` | INFO | INFO |
| `shared.filters.logLevels.warn` | WARN | WARN |
| `shared.filters.logLevels.error` | ERROR | ERROR |

