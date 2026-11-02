# 订阅文案对照（zh / en）

参照 Clash Verge Rev v2.5.9 `160bb4fe8eaef2bc626163d89622229e5d7c5b59`。

- 中文：`src/locales/zh/profiles.json`
- 英文：`src/locales/en/profiles.json`
- 加载：`src/services/i18n.ts` 以文件名作为顶层键，合并进 i18next 默认命名空间 `translation`。调用形式 `t('profiles.…')`。

订阅页。`src/pages/profiles.tsx` 与 `src/components/profile/*`。

| key | 中文 | English |
| --- | --- | --- |
| `profiles.page.actions.updateAll` | 更新所有订阅 | Update All Profiles |
| `profiles.page.actions.viewRuntimeConfig` | 查看运行时订阅 | View Runtime Config |
| `profiles.page.actions.reactivate` | 重新激活订阅 | Reactivate Profiles |
| `profiles.page.actions.import` | 导入 | Import |
| `profiles.page.batch.actions.delete` | 删除选中订阅 | Delete Selected Profiles |
| `profiles.page.batch.actions.selectAll` | 全选 | Select All |
| `profiles.page.batch.actions.deselectAll` | 取消全选 | Deselect All |
| `profiles.page.batch.actions.done` | 完成 | Done |
| `profiles.page.batch.summary.selected` | 已选中 | Selected |
| `profiles.page.batch.summary.items` | 项目 | items |
| `profiles.page.batch.title` | 批量操作 | Batch Operations |
| `profiles.page.importForm.placeholder` | 订阅文件链接 | Profile URL |
| `profiles.page.importForm.actions.paste` | 粘贴 | Paste |
| `profiles.page.feedback.errors.invalidUrl` | 无效的订阅链接，请输入以 http:// 或 https:// 开头的地址 | Invalid profile URL. Please enter a URL starting with http:// or https:// |
| `profiles.page.feedback.errors.onlyYaml` | 仅支持 YAML 文件 | Only YAML Files Supported |
| `profiles.page.feedback.errors.createFailed` | 创建订阅失败 | Failed to create profile |
| `profiles.page.feedback.errors.deleteFailed` | 删除订阅失败 | Failed to delete profile |
| `profiles.page.feedback.errors.enhanceFailed` | 重新激活订阅失败 | Failed to reactivate profiles |
| `profiles.page.feedback.errors.importFailed` | 导入订阅失败 | Failed to import profile |
| `profiles.page.feedback.errors.openFailed` | 打开订阅失败 | Failed to open profile |
| `profiles.page.feedback.errors.readFailed` | 读取订阅失败 | Failed to read profile |
| `profiles.page.feedback.errors.reorderFailed` | 调整订阅顺序失败 | Failed to reorder profiles |
| `profiles.page.feedback.errors.switchFailed` | 切换订阅失败 | Failed to switch profile |
| `profiles.page.feedback.errors.updateFailed` | 更新订阅失败 | Failed to update profile |
| `profiles.page.feedback.notifications.importRetry` | 订阅导入失败，尝试使用 Clash 代理导入 | Import failed, retrying with Clash proxy... |
| `profiles.page.feedback.notifications.importFail` | 使用 Clash 代理导入订阅失败 | Import failed even with Clash proxy |
| `profiles.page.feedback.notifications.importNeedsRefresh` | 订阅已导入，但可能需要手动刷新 | Profile imported but may need manual refresh |
| `profiles.page.feedback.notifications.importSuccess` | 订阅已成功导入，如未显示请重启应用 | Profile imported successfully, please restart if not visible |
| `profiles.page.feedback.notifications.profileSwitched` | 订阅已切换 | Profile Switched |
| `profiles.page.feedback.notifications.profileReactivated` | 订阅已激活 | Profile Reactivated |
| `profiles.page.feedback.notifications.switchBusy` | 正在切换订阅，请稍后重试 | Profile switch is already in progress |
| `profiles.page.feedback.notifications.batchDeleted` | 选中的订阅已成功删除 | Selected profiles deleted successfully |
| `profiles.page.feedback.notices.forceRefreshCompleted` | 数据已强制刷新 | Force refresh completed |
| `profiles.page.feedback.notices.discardedKeys` | 扩展配置或脚本写入了由应用设置接管的字段 {{keys}}，这些值已被丢弃 | Extensions wrote {{keys}}, which Settings manages, so those values were discarded |
| `profiles.page.feedback.notices.emergencyRefreshFailed` | 紧急刷新失败: {{message}} | Emergency refresh failed: {{message}} |
| `profiles.page.feedback.tooltips.forceRefreshStaleData` | 订阅数据异常，点击强制刷新 | Profile data is unavailable; click to force a refresh |
| `profiles.page.title` | 订阅 | Profiles |
| `profiles.components.card.labels.clickToImport` | 点击导入订阅 | Click to import subscription |
| `profiles.components.fileInput.chooseFile` | 选择文件 | Choose File |
| `profiles.components.menu.home` | 首 页 | Home |
| `profiles.components.menu.select` | 使用 | Select |
| `profiles.components.menu.shareQrCode` | 分享二维码 | Share QR Code |
| `profiles.components.menu.editInfo` | 编辑信息 | Edit Info |
| `profiles.components.menu.editFile` | 编辑文件 | Edit File |
| `profiles.components.menu.editRules` | 编辑规则 | Edit Rules |
| `profiles.components.menu.editProxies` | 编辑节点 | Edit Proxies |
| `profiles.components.menu.editGroups` | 编辑代理组 | Edit Proxy Groups |
| `profiles.components.menu.extendConfig` | 扩展覆写配置 | Extend Config |
| `profiles.components.menu.extendScript` | 扩展脚本 | Extend Script |
| `profiles.components.menu.openFile` | 打开文件 | Open File |
| `profiles.components.menu.update` | 更新 | Update |
| `profiles.components.menu.updateViaProxy` | 更新（代理） | Update via proxy |
| `profiles.components.more.global.merge` | 全局扩展覆写配置 | Global Extend Config |
| `profiles.components.more.global.script` | 全局扩展脚本 | Global Extend Script |
| `profiles.components.more.chips.merge` | Merge | Merge |
| `profiles.components.more.chips.script` | Script | Script |
| `profiles.components.profileItem.tooltips.showLast` | 点击查看上次更新时间 | Click to show last update time |
| `profiles.components.profileItem.tooltips.showNext` | 点击查看下一次更新 | Click to show next update |
| `profiles.components.profileItem.status.lastUpdateFailed` | 上次更新失败 | Last Update failed |
| `profiles.components.profileItem.status.nextUp` | 下次更新 | Next Up |
| `profiles.components.profileItem.status.noSchedule` | 没有计划 | No schedule |
| `profiles.components.profileItem.status.unknown` | 未知 | Unknown |
| `profiles.components.profileItem.status.autoUpdateDisabled` | 自动更新已禁用 | Auto update disabled |
| `profiles.modals.profileForm.title.create` | 新建配置 | Create Profile |
| `profiles.modals.profileForm.title.edit` | 编辑配置 | Edit Profile |
| `profiles.modals.profileForm.fields.type` | 类型 | Type |
| `profiles.modals.profileForm.fields.description` | 描述 | Descriptions |
| `profiles.modals.profileForm.fields.subscriptionUrl` | 订阅链接 | Subscription URL |
| `profiles.modals.profileForm.fields.userAgent` | 用户代理 | User Agent |
| `profiles.modals.profileForm.fields.httpTimeout` | HTTP 请求超时 | HTTP Request Timeout |
| `profiles.modals.profileForm.fields.updateInterval` | 更新间隔 | Update Interval |
| `profiles.modals.profileForm.fields.useSystemProxy` | 使用系统代理更新 | Use System Proxy |
| `profiles.modals.profileForm.fields.useClashProxy` | 使用内核代理更新 | Use Clash Proxy |
| `profiles.modals.profileForm.fields.acceptInvalidCerts` | 允许无效证书（危险） | Allows Invalid Certificates (Danger) |
| `profiles.modals.profileForm.fields.allowAutoUpdate` | 允许自动更新 | Allow Auto Update |
| `profiles.modals.profileForm.types.remote` | 远程 | Remote |
| `profiles.modals.profileForm.types.local` | 本地 | Local |
| `profiles.modals.profileForm.errors.typeRequired` | 请选择订阅类型 | Select a profile type |
| `profiles.modals.profileForm.errors.urlRequired` | 请输入订阅链接 | Enter a subscription URL |
| `profiles.modals.profileForm.errors.uidMissing` | 缺少订阅标识 | The profile identifier is missing |
| `profiles.modals.profileForm.errors.saveFailed` | 保存订阅失败 | Failed to save profile |
| `profiles.modals.profileForm.warnings.frequentUpdate` | 间隔过短会频繁触发内核重载并中断现有连接，建议使用默认值 {{minutes}} 分钟（24 小时） | A short interval reloads the core that often and drops existing connections. {{minutes}} minutes (24 hours) is recommended. |
| `profiles.modals.profileForm.feedback.notifications.creationRetry` | 订阅创建失败，尝试使用 Clash 代理创建 | Profile creation failed, retrying with Clash proxy... |
| `profiles.modals.profileForm.feedback.notifications.creationSuccess` | 使用 Clash 代理创建订阅成功 | Profile creation succeeded with Clash proxy |
| `profiles.modals.proxiesEditor.title` | 编辑节点 | Edit Proxies |
| `profiles.modals.proxiesEditor.placeholders.multiUri` | 多条 URI 请使用换行分隔（支持 Base64 编码） | Use newlines for multiple uri(Base64 encoding supported) |
| `profiles.modals.proxiesEditor.actions.prepend` | 添加前置代理节点 | Prepend Proxy |
| `profiles.modals.proxiesEditor.actions.append` | 添加后置代理节点 | Append Proxy |
| `profiles.modals.groupsEditor.title` | 编辑代理组 | Edit Proxy Groups |
| `profiles.modals.groupsEditor.errors.nameRequired` | 代理组名称不能为空 | Group Name Required |
| `profiles.modals.groupsEditor.errors.nameExists` | 代理组名称已存在 | Group Name Already Exists |
| `profiles.modals.groupsEditor.fields.type` | 代理组类型 | Group Type |
| `profiles.modals.groupsEditor.fields.name` | 代理组组名 | Group Name |
| `profiles.modals.groupsEditor.fields.icon` | 代理组图标 | Proxy Group Icon |
| `profiles.modals.groupsEditor.fields.proxies` | 引入代理 | Use Proxies |
| `profiles.modals.groupsEditor.fields.provider` | 引入代理集合 | Use Provider |
| `profiles.modals.groupsEditor.fields.healthCheckUrl` | 健康检查测试地址 | Health Check Url |
| `profiles.modals.groupsEditor.fields.expectedStatus` | 期望状态码 | Expected Status |
| `profiles.modals.groupsEditor.fields.interval` | 检查间隔 | Interval |
| `profiles.modals.groupsEditor.fields.maxFailedTimes` | 最大失败次数 | Max Failed Times |
| `profiles.modals.groupsEditor.fields.interfaceName` | 出站接口 | Interface Name |
| `profiles.modals.groupsEditor.fields.routingMark` | 路由标记 | Routing Mark |
| `profiles.modals.groupsEditor.fields.filter` | 过滤节点 | Filter |
| `profiles.modals.groupsEditor.fields.excludeFilter` | 排除节点 | Exclude Filter |
| `profiles.modals.groupsEditor.fields.excludeType` | 排除节点类型 | Exclude Type |
| `profiles.modals.groupsEditor.fields.includeAll` | 引入所有出站代理、代理集合 | Include All Proxies and Providers |
| `profiles.modals.groupsEditor.fields.includeAllProxies` | 引入所有出站代理 | Include All Proxies |
| `profiles.modals.groupsEditor.fields.includeAllProviders` | 引入所有代理集合 | Include All Providers |
| `profiles.modals.groupsEditor.toggles.lazy` | 懒惰状态 | Lazy |
| `profiles.modals.groupsEditor.toggles.disableUdp` | 禁用 UDP | Disable UDP |
| `profiles.modals.groupsEditor.toggles.hidden` | 隐藏代理组 | Hidden |
| `profiles.modals.groupsEditor.actions.prepend` | 添加前置代理组 | Prepend Group |
| `profiles.modals.groupsEditor.actions.append` | 添加后置代理组 | Append Group |
| `profiles.modals.editor.actions.format` | 格式化文档 | Format document |
| `profiles.modals.editor.messages.readOnly` | 无法在只读模式下编辑 | Cannot edit in read-only editor |
| `profiles.modals.editor.enhance.globalMerge` | 对所有订阅生效，在应用设置之后合并；后续的全局脚本、订阅扩展配置和订阅脚本仍可修改结果。 | Applies to all profiles, after app settings. The global script, profile config and profile script can still change the result. |
| `profiles.modals.editor.enhance.globalScript` | 对所有订阅生效，接收全局扩展配置合并后的结果；后续的订阅扩展配置和订阅脚本仍可修改结果。 | Applies to all profiles, using the result of the global config merge. The profile config and profile script can still change the result. |
| `profiles.modals.editor.enhance.profileMerge` | 仅对当前订阅生效，在全局扩展配置和全局脚本之后合并；后续的订阅脚本仍可修改结果。 | Applies only to this profile, after the global config and global script. The profile script can still change the result. |
| `profiles.modals.editor.enhance.profileScript` | 仅对当前订阅生效，接收订阅扩展配置合并后的结果，是应用设置回写前的最后一个扩展步骤。 | Applies only to this profile, using the result of the profile config merge. This is the last extension step before app settings are reapplied. |
| `profiles.modals.editor.enhance.order` | 执行顺序：订阅原文 → 应用设置 → <globalMerge>全局扩展配置</globalMerge> → <globalScript>全局扩展脚本</globalScript> → <profileMerge>订阅扩展配置</profileMerge> → <profileScript>订阅扩展脚本</profileScript> → 应用设置回写 | Order: Original profile → App settings → <globalMerge>Global extend config</globalMerge> → <globalScript>Global extend script</globalScript> → <profileMerge>Profile extend config</profileMerge> → <profileScript>Profile extend script</profileScript> → App settings reapplied |
| `profiles.modals.editor.enhance.settingsPriority` | 注意：由应用设置接管的字段最终以应用设置为准，在扩展配置或脚本中修改这些字段不会生效。 | Note: Fields managed by app settings always use those settings in the final config. Changing these fields in extensions has no effect. |
| `profiles.modals.confirmDelete.title` | 确认删除 | Confirm deletion |
| `profiles.modals.confirmDelete.message` | 此操作不可逆 | This operation is not reversible |
| `profiles.modals.logViewer.title` | 脚本控制台输出 | Script Console |
| `profiles.modals.qrViewer.title` | 订阅二维码 | Subscription QR Code |

